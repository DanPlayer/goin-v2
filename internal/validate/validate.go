package validate

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
	"time"
)

var Validator *validator.Validate
var trans ut.Translator

// New 初始化验证器
func init() {
	Validator = validator.New()
	
	//注册翻译器
	trans, _ = ut.New(zh.New()).GetTranslator("zh")
	_ = zhTranslations.RegisterDefaultTranslations(Validator, trans)
	
	//注册自定义验证函数
	_ = Validator.RegisterValidation("timeUnix", timeUnix)
}

// ParseStruct 验证结构体
func ParseStruct(s interface{}) error {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}
	if list, ok := err.(validator.ValidationErrors); ok {
		msgList := make([]string, 0)
		for _, val := range list {
			msgList = append(msgList, val.Translate(trans))
		}
		return errors.New(strings.Join(msgList, ","))
	}
	return errors.New("参数校验异常，请联系开发人员")
}

//TimeUnix 校验时间戳
func timeUnix(f validator.FieldLevel) bool {
	unix := f.Field().Int()
	if unix == 0 {
		return true
	}
	return time.Unix(unix, 0).Year() > 1970
}