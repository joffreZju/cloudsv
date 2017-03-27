package main

import (
	"common/controller/user"
	"common/filter"

	"github.com/astaxie/beego"
)

func LoadRouter() {
	// aliyu 健康检测
	//beego.Router("/health", &maincontroller.Controller{}, "*:Check")
	// user 相关
	beego.Router("/v2/user/getcode", &user.Controller{}, "POST:GetCode")
	beego.Router("/v2/user/login", &user.Controller{}, "POST:LoginUser")
	beego.Router("/v2/user/upload_pic", &user.Controller{}, "POST:UploadPic")
	beego.Router("/v2/user/upload_report_pic", &user.Controller{}, "POST:UploadReportPic")
	beego.Router("/v2/user/index", &user.Controller{}, "POST:UserIndex")
	beego.Router("/user/index", &user.Controller{}, "*:UserIndex") //压测使用
	beego.Router("/v2/user/info", &user.Controller{}, "POST:GetUserInfo")
	beego.Router("/v2/user/edit_profile", &user.Controller{}, "POST:EditProfile")

	// 静态文件
	beego.SetStaticPath("/static", "../static")

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
