package main

import (
	"common/controller/user"
	"common/filter"

	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
)

func LoadRouter() {
	// aliyu 健康检测
	//beego.Router("/health", &maincontroller.Controller{}, "*:Check")
	// user 相关
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "POST:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "POST:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "POST:UserLogin")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "POST:GetUserInfo")
	beego.Router(UserPrefix+"/edit_profile", &user.Controller{}, "POST:EditProfile")

	// 非登录态列表
	notNeedAuthList := []string{
		// aliyun check
		"/",
		// user
		"/v2/app/init",
		"/v2/user/getcode", "/v2/user/login", "/v2/user/index",
	}

	// add filter
	// 请求合法性验证 这个要放在第一个
	filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("api_user", notNeedAuthList))
	beego.InsertFilter("/user/*", beego.BeforeRouter, filter.GetAuthFilter())
}
