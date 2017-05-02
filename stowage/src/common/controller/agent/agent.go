package agent

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"common/lib/util"
	"common/model"
	"common/service"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

const (
	defaultPasswd = "123456"
)

//创建代理商
func (c *Controller) AgentCreate() {
	licenceFile := c.GetString("licenseFile")
	tel := c.GetString("tel")
	passwdc := keycrypt.Sha256Cal(defaultPasswd)
	name := c.GetString("name")
	u := model.User{
		Tel:      tel,
		Password: passwdc,
		UserType: 2,
	}
	a := model.Agent{
		LicenseFile: licenceFile,
		Name:        name,
		User:        &u,
		Status:      1,
		Discount:    50,
	}
	err := service.AgentCreate(&a)
	if err != nil {
		beego.Error("agent create failed", err)
		c.ReplyErr(errcode.ErrAgentCreatFailed)
		return
	}
	//同时建立代理商资金账户
	ac := model.Account{
		AccountNo: util.RandomByte16(),
		Userid:    a.User.Id,
		UserType:  2,
		Status:    1,
	}
	err = service.AccountCreate(&ac)
	if err != nil {
		beego.Error("create user account failed:", err)
	}

	c.ReplySucc("success")
}

//获取代理商列表
func (c *Controller) AgentList() {
	//page, _ := c.GetInt("page")
	list, err := service.AgentGetList()
	if err != nil {
		c.ReplyErr(err)
		return
	}

	c.ReplySucc(list)
}

//代理商获取客户列表
func (c *Controller) AgentClients() {
	tel := c.GetString("tel")
	users, err := service.AgentClients(tel)
	if err != nil {
		beego.Error("get agent clients failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(users)
}

//获取代理商信息
func (c *Controller) AgentGetInfo() {
	id, _ := c.GetInt("aid")
	a, err := model.GetAgentInfo(id)
	if err != nil {
		beego.Error("not find the agent,userid:", id, err)
		c.ReplyErr(errcode.ErrAgentNotExisted)
		return
	}
	c.ReplySucc(*a)
	return

}

//修改代理商信息
func (c *Controller) AgentModify() {
	id, err := c.GetInt("aid")
	if err != nil {
		beego.Error("parameters error:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	status, err := c.GetInt("status")
	if err != nil {
		beego.Error("parameters error:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	licenseFile := c.GetString("licenseFile")
	desc := c.GetString("desc")
	a := model.Agent{
		Uid:         id,
		LicenseFile: licenseFile,
		Status:      status,
		Desc:        desc,
	}
	err = model.AgentUpdate(&a)
	if err != nil {
		beego.Error("update agent fail:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("OK")
}
