package payment

import (
	"flying-star/pkg/payment/store"
	"net/http"
)

// BindAccountOptions 绑定账户参数
type BindAccountOptions struct {
	// 用户ID
	UserID string `json:"userId"`
	// 绑定ID
	OpenID string `json:"openId"`
	//用户真实姓名
	RealName string `json:"realName"`
}

// RechargeOptions 充值参数
type RechargeOptions struct {
	//付款用户
	UserId string `json:"userId"`
	//付款企业
	CorpId string `json:"corpId"`
	//付款金额
	Price uint `json:"price"`
	//支付备注
	Remark string `json:"remark"`
}

// RechargeSchema 充值结果
type RechargeSchema struct {
	//系统订单ID
	Uid string `json:"uid"`
	//支付内容
	Content string `json:"content"`
}

// TransferOptions 转账参数
type TransferOptions struct {
	//付款用户
	UserId string `json:"userId"`
	//付款企业
	CorpId string `json:"corpId"`
	//付款金额
	Price uint `json:"price"`
	//收款用户ID
	PayeeId string `json:"PayeeId"`
	//所属业务模块
	Module string `json:"module"`
	//支付备注
	Remark string `json:"remark"`
	//外部关联ID
	ReferId string `json:"referId"`
}

type Plugin interface {
	// BindAccount 绑定三方账户信息
	BindAccount(options BindAccountOptions) error
	// RechargeBegin 开始充值
	RechargeBegin(options RechargeOptions) (RechargeSchema, error)
	// RechargeEnd 充值完成
	RechargeEnd(c *http.Request) (userId string, err error)
	// Transfer 转账操作
	Transfer(options TransferOptions) error
	// GetStore 获取存储实例
	GetStore() *store.Store
}
