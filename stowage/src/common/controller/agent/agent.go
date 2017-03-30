package agent

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/lib/keycrypt"
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
	u := model.User{
		Tel:      tel,
		Password: passwdc,
		UserType: 2,
	}
	a := model.Agent{
		LicenseFile: licenceFile,
		User:        &u,
		Status:      1,
	}
	err := service.AgentCreate(&a)
	if err != nil {
		beego.Error("agent create failed", err)
		c.ReplyErr(errcode.ErrAgentCreatFailed)
		return
	}
	c.ReplySucc("success")

}

//获取代理商列表
func (c *Controller) AgentList() {
	list, err := model.GetAgentAll()
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(list)
}

//获取代理商信息
func (c *Controller) AgentGetInfo() {
	id, _ := c.GetInt("id")
	a, err := model.GetAgentInfo(id)
	if err != nil {
		beego.Error("not find the agent,id:", id, err)
		c.ReplyErr(errcode.ErrAgentNotExisted)
		return
	}
	c.ReplySucc(a)
	return

}

//修改代理商信息
func (c *Controller) AgentModify() {
	id, _ := c.GetInt("id")
	status, _ := c.GetInt("status")
	licenseFile := c.GetString("licenseFile")
	a := model.Agent{
		Id:          id,
		LicenseFile: licenseFile,
		Status:      status,
	}
	err := model.AgentUpdate(&a)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("OK")
}
