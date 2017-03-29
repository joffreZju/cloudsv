package main

import (
	"common/controller/common"
	"common/controller/user"
	"common/filter"

	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
	CommonPrefix string = "/v2/common"
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

	//通用功能
	beego.Router(CommonPrefix+"/upload_file", &common.Controller{}, "POST:UploadFile")    //文件上传
	beego.Router(CommonPrefix+"/download_file", &common.Controller{}, "GET:DownloadFile") //文件下载

	beego.Router(CommonPrefix+"/add_document", &common.Controller{}, "POST:AddDocument")      //文件上传
	beego.Router(CommonPrefix+"/update_document", &common.Controller{}, "GET:UpdateDocument") //文件下载

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
