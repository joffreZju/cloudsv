package push

import (
	"common/lib/constant"
	"encoding/json"
	"fmt"
	"github.com/GiterLab/aliyun-sms-go-sdk/sms"
	"github.com/astaxie/beego"
)

func Init() (err error) {
	//alidayu.AppKey = "LTAIwxFn7egYfvra"
	//alidayu.AppKey = "23441572"
	//alidayu.AppSecret = "nBfpqo4StRZv9JreRsLQpFaZKKUT1h"
	//alidayu.AppSecret = "8d890383b59c43ccdeb50fe9d0074087"
	return
}

//
//func SendMsgWithDayuToMobile(mobile, code string, product string) bool {
//	sucess, response := alidayu.SendSMS(mobile, "登录验证", "SMS_13400735", fmt.Sprintf(`{"code":"%v","product":"%v"}`, code, product))
//	beego.Debug(response)
//	return sucess
//}
//
//func SendSMSWithDayu(mobile, name, tplId string, params map[string]string) bool {
//	args, _ := json.Marshal(params)
//	sucess, response := alidayu.SendSMS(mobile, name, tplId, string(args))
//	beego.Debug(response)
//	fmt.Println(response, sucess)
//	return sucess
//}

func SendSmsCodeToMobile(mobile, code string) error {
	param := make(map[string]string)
	param["smscode"] = code
	c := sms.New(constant.ALI_ACCESS_KEY_ID, constant.ALI_ACCESS_KEY_SECRET)
	str, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("send smscode failed,%v", err)
	}
	e, err := c.SendOne(mobile, constant.SMS_SIGN_NAME, constant.SMS_TEMPLATE_WEB, string(str))
	if err != nil {
		return fmt.Errorf("send sms failed,%v,%v", err, e.Error())
	}
	return nil
}

func SendErrorSms(mobile, content string) error {
	param := make(map[string]string)
	param["content"] = content
	c := sms.New(constant.ALI_ACCESS_KEY_ID, constant.ALI_ACCESS_KEY_SECRET)
	str, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("send sms failed,%v", err)
	}
	e, err := c.SendOne(mobile, constant.SMS_SIGN_NAME, constant.SMS_TEMPLATE_WHEN_ERROR, string(str))
	if err != nil {
		return fmt.Errorf("send sms failed,%v,%v", err, e.Error())
	}
	beego.Info("calculate failed:", content)
	return nil
}
