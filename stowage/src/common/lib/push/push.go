package push

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/ltt1987/alidayu"
)

func Init() (err error) {
	alidayu.AppKey = "23441572"
	alidayu.AppSecret = "8d890383b59c43ccdeb50fe9d0074087"
	return
}

func SendMsgWithDayuToMobile(mobile, code string, product string) bool {
	sucess, response := alidayu.SendSMS(mobile, "登录验证", "SMS_13400735", fmt.Sprintf(`{"code":"%v","product":"%v"}`, code, product))
	beego.Debug(response)
	return sucess
}

func SendSMSWithDayu(mobile, name, tplId string, params map[string]string) bool {
	args, _ := json.Marshal(params)
	sucess, response := alidayu.SendSMS(mobile, name, tplId, string(args))
	beego.Debug(response)
	fmt.Println(response, sucess)
	return sucess
}
