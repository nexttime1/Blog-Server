package validate

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net"
	"reflect"
	"strings"
)

var trans ut.Translator

func init() {
	// 创建翻译器
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")

	// 注册翻译器
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			label = field.Name
		}
		name := field.Tag.Get("json")
		return fmt.Sprintf("%s---%s", name, label)
	})

	v.RegisterValidation("fip", func(fl validator.FieldLevel) bool {
		// 打印调试信息
		fmt.Println("fl.Field(): ", fl.Field())                     // 获取当前字段的值
		fmt.Println("fl.FieldName(): ", fl.FieldName())             // 获取字段名称（JSON标签名）
		fmt.Println("fl.StructFieldName(): ", fl.StructFieldName()) // 获取结构体字段名
		fmt.Println("fl.Parent(): ", fl.Parent())                   // 获取父结构体
		fmt.Println("fl.Top(): ", fl.Top())                         // 获取顶级结构体
		fmt.Println("fl.Param(): ", fl.Param())                     // 获取验证标签参数（如max=10中的"10"）  默认值

		// 类型断言，检查字段是否为字符串类型且非空
		ip, ok := fl.Field().Interface().(string) //reflect.Value.Interface() 方法  变成原生go
		if ok && ip != "" {
			// 传了值就去校验是不是IP地址
			ipObj := net.ParseIP(ip)
			return ipObj != nil // 如果解析成功则返回true，否则false
		}
		return true // 如果字段为空或非字符串类型，返回true（忽略验证）
	})
}

/*
{
  "name": "name参数必填",
}
*/

func ValidateErr(err error) (data map[string]any, msg string) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		msg = err.Error()
		return
	}
	data = map[string]any{}
	var msgList []string
	for _, e := range errs {
		msg := e.Translate(trans)
		_list := strings.Split(msg, "---")
		if e.Tag() == "fip" {
			data[strings.Split(e.Field(), "---")[0]] = "该ip地址不符合要求"
			continue
		}
		data[_list[0]] = _list[1]
		msgList = append(msgList, _list[1])
	}
	msg = strings.Join(msgList, ";")
	return
}
