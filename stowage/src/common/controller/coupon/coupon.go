package coupon

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/model"
	"common/service"
	"strings"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

//导入代金券
func (c *Controller) CouponCreate() {
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
	err = service.AddCoupons(start, end)
	if err != nil {
		beego.Error(err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.ReplyErr(errcode.ErrCouponExist)
		} else {
			c.ReplyErr(errcode.ErrParams)
		}
		return
	}
	c.ReplySucc("success")
	return
}

//向代理商发放代金券
func (c *Controller) GrantAgent() {
	agentid, _ := c.GetInt("agent")
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
func (c *Controller) CouponRecycle() {
	aid, _ := c.GetInt("agent")
	err := model.UpdateCouponByAgent(aid)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//用户消费代金券
func (c *Controller) CouponUsing() {
	code := c.GetString("verify")
	num, _ := c.GetInt("number")
	uid := int(c.UserID)
	err := service.UsingCoupon(num, uid, code)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

//浏览代金券数据
func (c *Controller) CouponInfo() {

}
