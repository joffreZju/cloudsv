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
	userid, _ := c.GetInt("agent")
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
	err = service.GrantAgent(start, end, userid)
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
	err := model.RecycleCouponByAgent(aid)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("success")
	return
}

/*
func (c *Controller) CouponRecycleRange() {
	s, _ := c.GetInt("start")
	e, _ := c.GetInt("end")

	err = service.RecycleCouponRange(s,e)
}*/

//代理商卡券，分页
func (c *Controller) CouponList() {
	page, _ := c.GetInt("page")
	aid, _ := c.GetInt("agent")
	count, list, err := service.GetCouponList(aid, page)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"total":   count,
		"coupons": list,
	})
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

//查看某张代金券
func (c *Controller) CouponInfo() {
	no, _ := c.GetInt("number")
	co, err := service.GetCoupon(no)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(co)
	return
}
