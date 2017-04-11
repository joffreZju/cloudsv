package main

import (
	"common/controller/account"
	"common/controller/agent"
	"common/controller/bill"
	"common/controller/doc"
	"common/controller/order"
	"common/controller/recharge"
	"common/controller/user"
	"common/filter"

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
	//wxpay callback
	beego.Router("/notify/wxpay", &order.Controller{}, "Post:WxPay")

	// user 相关
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "Get:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "POST:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "POST:UserLogin")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "Get:GetUserInfo")
	beego.Router(UserPrefix+"/passwd/modify", &user.Controller{}, "POST:Resetpwd")
	beego.Router(UserPrefix+"/edit_profile", &user.Controller{}, "POST:EditProfile")

	//用户的账户
	beego.Router(UserPrefix+"/account/info", &account.Controller{}, "Get:AccountInfo")

	//代理商相关
	beego.Router(ManagePrefix+"/agent/add", &agent.Controller{}, "POST:AgentCreate")
	beego.Router(ManagePrefix+"/agent/modify", &agent.Controller{}, "POST:AgentModify")
	beego.Router(ManagePrefix+"/agent/info", &agent.Controller{}, "Get:AgentGetInfo")
	beego.Router(ManagePrefix+"/agent/list", &agent.Controller{}, "Get:AgentList")
	beego.Router(ManagePrefix+"/agent/clients", &agent.Controller{}, "Get:AgentClients")

	//代金券
	beego.Router(ManagePrefix+"/recharge/input", &recharge.Controller{}, "POST:RechargeCreate") //录入
	beego.Router(ManagePrefix+"/recharge/grant", &recharge.Controller{}, "POST:GrantAgent")     //发放
	beego.Router(ManagePrefix+"/recharge/recycle", &recharge.Controller{}, "POST:RechargeRecycle")
	beego.Router(UserPrefix+"/recharge/using", &recharge.Controller{}, "POST:RechargeUsing")
	beego.Router(ManagePrefix+"/recharge/info", &recharge.Controller{}, "Get:RechargeInfo")
	//beego.Router(ManagePrefix+"/recharge/", &recharge.Controller{}, "POST:RechargeInfo")

	//订单交易
	beego.Router(UserPrefix+"/order/pay", &order.Controller{}, "Get:PayOnline")
	beego.Router(ManagePrefix+"/order/info", &order.Controller{}, "Get:OrderInfo")
	beego.Router(ManagePrefix+"/order/list_day", &order.Controller{}, "Get:OrderDay")

	//账单
	beego.Router(UserPrefix+"/bill/info", &bill.Controller{}, "Get:BillInfo")

	//文档
	beego.Router(ManagePrefix+"/doc/add", &doc.Controller{}, "POST:AddDocument")        //文档上传
	beego.Router(UserPrefix+"/doc/view", &doc.Controller{}, "GET:GetDocUsing")          //文档查看
	beego.Router(ManagePrefix+"/doc/list", &doc.Controller{}, "GET:GetDocList")         //文档列表
	beego.Router(ManagePrefix+"/doc/set_status", &doc.Controller{}, "Get:SetDocStatus") //文档列表

	beego.Router(UserPrefix+"/doc/file_add", &doc.Controller{}, "POST:AddFile")      //文件上传
	beego.Router(UserPrefix+"/doc/file_down", &doc.Controller{}, "GET:FileDownload") //文件下载

	// 非登录态列表
	notNeedAuthList := []string{
		// aliyun check
		"/",
		// user
		ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register", ExemptPrefix + "/user/login",
	}

	// add filter
	// 请求合法性验证 这个要放在第一个
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
