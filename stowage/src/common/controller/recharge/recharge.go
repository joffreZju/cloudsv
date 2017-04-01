package recharge

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/model"
	"common/service"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

//导入代金券
func (c *Controller) RechargeCreate() {
	start, err := c.GetInt("start")
	if err != nil {
		beego.Error("parameters wrong", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	end, err := c.GetInt("end")
	if err != nil {
		beego.Error("parameters wrong", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	err = service.AddRecharges(start, end)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//向代理商发放代金券
func (c *Controller) GrantAgent() {
	agentid := c.GetString("agent")
	start, err := c.GetInt("start")
	if err != nil {
		beego.Error("parameters wrong", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	end, err := c.GetInt("end")
	if err != nil {
		beego.Error("parameters wrong", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	err = service.GrantAgent(start, end, agentid)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//管理员回收代金券
func (c *Controller) RechargeRecycle() {
	agentuser := c.GetString("agentid")
	err := model.UpdateRechargeByAgent(agentuser)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//用户消费代金券
func (c *Controller) RechargeUsing() {
	code := c.GetString("verify")
	num, _ := c.GetInt("number")
	uid := c.GetString("userid")
	err := service.UsingRecharge(num, uid, code)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//浏览代金券数据
func (c *Controller) RechargeInfo() {

}
