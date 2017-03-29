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
	ManagePrefix string = "/v2/back"
)

func LoadRouter() {
	// aliyu 健康检测
	//beego.Router("/health", &maincontroller.Controller{}, "*:Check")
	// user 相关
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "Get:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "POST:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "POST:UserLogin")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "Get:GetUserInfo")
	beego.Router(UserPrefix+"/passwd/reset", &user.Controller{}, "Post:Resetpwd")
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
		ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register", ExemptPrefix + "/user/login",
	}

	// add filter
	// 请求合法性验证 这个要放在第一个
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
