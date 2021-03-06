package main

import (
	"common/controller/doc"
	"common/filter"
	"s4s/common/controller/user"
	"stowage/src/common/controller/account"
	"stowage/src/common/controller/agent"

	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
	ManagePrefix string = "/v2/admin"
	TradePrefix  string = "/v2/trade"
)

func LoadRouter() {
	// aliyu 健康检测
	//beego.Router("/health", &maincontroller.Controller{}, "*:Check")
	// user 相关
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "*:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "*:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "*:UserLogin")
	beego.Router(ExemptPrefix+"/user/login_phone", &user.Controller{}, "*:UserLoginPhoneCode")
	beego.Router(UserPrefix+"/login_out", &user.Controller{}, "*:LoginOut")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "*:GetUserInfo")
	beego.Router(UserPrefix+"/passwd/modify", &user.Controller{}, "*:Resetpwd")
	beego.Router(UserPrefix+"/edit_profile", &user.Controller{}, "*:EditProfile")
	beego.Router(ExemptPrefix+"/passwd/retrieve", &user.Controller{}, "Post:Retrievepwd")

	//用户的账户
	beego.Router(UserPrefix+"/account/info", &account.Controller{}, "Get:AccountInfo")

	//代理商相关
	beego.Router(ManagePrefix+"/agent/add", &agent.Controller{}, "POST:AgentCreate")
	beego.Router(ManagePrefix+"/agent/modify", &agent.Controller{}, "POST:AgentModify")
	beego.Router(ManagePrefix+"/agent/info", &agent.Controller{}, "Get:AgentGetInfo")
	beego.Router(ManagePrefix+"/agent/list", &agent.Controller{}, "Get:AgentList")
	beego.Router(ManagePrefix+"/agent/clients", &agent.Controller{}, "Get:AgentClients")

	//文档
	beego.Router(ManagePrefix+"/doc/add", &doc.Controller{}, "POST:AddDocument")         //文档上传
	beego.Router(UserPrefix+"/doc/view", &doc.Controller{}, "GET:GetDocUsing")           //文档查看
	beego.Router(ManagePrefix+"/doc/list", &doc.Controller{}, "GET:GetDocList")          //文档列表
	beego.Router(ManagePrefix+"/doc/set_status", &doc.Controller{}, "Post:SetDocStatus") //文档列表

	beego.Router(UserPrefix+"/doc/file_add", &doc.Controller{}, "POST:AddFile")      //文件上传
	beego.Router(UserPrefix+"/doc/file_down", &doc.Controller{}, "GET:FileDownload") //文件下载

	// 非登录态列表
	notNeedAuthList := []string{
		// aliyun check
		"/",
		// user
		ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register", ExemptPrefix + "/user/login",
		UserPrefix + "/doc/file_down",
	}

	// 请求合法性验证 这个要放在第一个
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("allsum_account", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
