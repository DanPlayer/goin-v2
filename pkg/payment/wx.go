package payment

import (
	"encoding/json"
	"errors"
	"flying-star/pkg/payment/store"
	"flying-star/utils"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
	"net/http"
	"strconv"
	"strings"
)

type WechatPay struct {
	name string
	client *wechat.Client
	Store *store.Store
	notifyUrl string
}

// WxOptions 微信支付初始化参数
type WxOptions struct {
	Store *store.Store
	AppId string 			// 应用ID
	MchId string			// 商户ID
	ApiKey string			// API密钥
	KeyFilePath string		// apiclient_key.pem 文件路径
	CertFilePath string		// apiclient_cert.pem 文件路径
	NotifyUrl string		// 支付回调地址
	IsProd bool				// 是否为正式环境
}

// NewWxPay 初始化微信支付组件
func NewWxPay(options WxOptions) (*WechatPay, error) {
	ins := WechatPay{
		name: "wx_pay",
		Store: options.Store,
		notifyUrl: options.NotifyUrl,
	}

	client := wechat.NewClient(options.AppId, options.MchId, options.ApiKey, options.IsProd)
	if err := client.AddCertPemFilePath(options.CertFilePath, options.KeyFilePath); err != nil {
		return nil, err
	}

	//设置国家
	client.SetCountry(wechat.China)

	ins.client = client
	return &ins, nil
}

func (w WechatPay) BindAccount(options BindAccountOptions) (err error) {
	return w.Store.BindingSet(options.UserID, options.OpenID, options.RealName, w.name)
}

func (w WechatPay) GetStore() *store.Store {
	return w.Store
}

func (w WechatPay) RechargeBegin(options RechargeOptions) (res RechargeSchema, err error) {
	outTradeNum := strings.ReplaceAll(utils.GetUid(), "-", "")

	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32)).
		Set("body", options.Remark).
		Set("out_trade_no", outTradeNum).
		Set("total_fee", options.Price).
		Set("spbill_create_ip", "127.0.0.1").
		Set("notify_url", w.notifyUrl).
		Set("trade_type", wechat.TradeType_Native).
		Set("device_info", "WEB").
		Set("sign_type", wechat.SignType_MD5).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.SetBodyMap("h5_info", func(bm gopay.BodyMap) {
				bm.Set("type", "Wap")
				bm.Set("wap_url", w.notifyUrl)
				bm.Set("wap_name", options.Remark)
			})
		})
	fmt.Println(w.notifyUrl)

	//请求支付下单，成功后得到结果
	wxRsp, err := w.client.UnifiedOrder(bm)

	if err != nil {
		return  res, err
	}

	if wxRsp.ReturnCode != "SUCCESS" || wxRsp.ResultCode != "SUCCESS" {
		return res, errors.New(wxRsp.ReturnMsg)
	}

	bData, _ := json.Marshal(wxRsp)
	fmt.Println(string(bData))

	ts, err := w.Store.RechargeRecordSet(w.name, store.RechargeRecordSetOptions{
		UserId:        options.UserId,
		CorpId:        options.CorpId,
		Price:         options.Price,
		Remark:        options.Remark,
		ChargeID:      outTradeNum,
	})

	if err == nil {
		ts.Commit()
	}

	res.Uid = outTradeNum
	res.Content = wxRsp.CodeUrl
	return res, nil
}

func (w WechatPay) RechargeEnd(req *http.Request) (string, error)  {
	notify, _ := wechat.ParseNotifyToBodyMap(req)

	bData, _ := json.Marshal(notify)
	fmt.Println(string(bData))

	//验签
	ok, err := wechat.VerifySign(w.client.ApiKey, wechat.SignType_MD5, notify)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("非法请求")
	}

	outTradeNo, ok := notify["out_trade_no"].(string)
	if !ok || outTradeNo == "" {
		return "", nil
	}

	price := 0
	if priceStr, ok := notify["cash_fee"].(string); ok {
		price, _ = strconv.Atoi(priceStr)
	}

	//更新订单信息
	ts, userId, err := w.Store.RechargeRecordUpdate(w.name, store.RechargeRecordUpdateOptions{
		ChargeID:      outTradeNo,
		TransactionID: notify["transaction_id"].(string),
		Price:         uint(price),
	})

	if err != nil {
		return "", err
	}

	if ts != nil {
		if err := ts.Commit().Error; err != nil {
			ts.Rollback()
			return "", err
		}
	}

	fmt.Printf("订单 %s 支付完成，支付金额：%v \n", notify["out_trade_no"], notify["total_fee"])
	return userId, nil
}

func (w WechatPay) Transfer(options TransferOptions) error {
	bindInfo, err := w.Store.BindingGet(options.PayeeId, w.name)
	if err != nil || bindInfo.OpenID == "" {
		return errors.New("当前用户暂未绑定OpenId")
	}

	fmt.Println(bindInfo.OpenID, bindInfo.RealName)

	//存储转账记录
	ts, err := w.Store.TransferRecordSet(w.name, store.TransferRecordSetOptions{
		UserId:  options.UserId,
		CorpId:  options.CorpId,
		Price:   options.Price,
		PayeeId: options.PayeeId,
		PayeeName: bindInfo.RealName,
		Module: options.Module,
		Remark:  options.Remark,
		ReferId: options.ReferId,
	})

	if err != nil {
		return err
	}
	fmt.Println(ts)

	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32)).
		Set("partner_trade_no", util.GetRandomString(32)).
		Set("openid", bindInfo.OpenID).
		Set("check_name", "NO_CHECK"). 		// NO_CHECK：不校验真实姓名 , FORCE_CHECK：强校验真实姓名
		Set("re_user_name", bindInfo.RealName).       	// 收款用户真实姓名。 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
		Set("amount", options.Price).              	// 企业付款金额，单位为分
		Set("desc", options.Remark).             		// 企业付款备注，必填。注意：备注中的敏感词会被转成字符*
		Set("spbill_create_ip", "127.0.0.1")

	// 企业向微信用户个人付款（不支持沙箱环境）
	wxRsp, err := w.client.Transfer(bm)
	if err != nil {
		ts.Rollback()
		return err
	}

	bData, _ := json.Marshal(wxRsp)
	fmt.Println(string(bData))

	if wxRsp.ReturnCode != "SUCCESS" || wxRsp.ResultCode != "SUCCESS" {
		ts.Rollback()
		return errors.New(wxRsp.ErrCodeDes)
	}
	ts.Commit()
	return nil
}

