package order

import (
	"common/lib/errcode"
	"common/model"
	"common/service"
	"common/service/wxisv"
	"time"

	"github.com/astaxie/beego"
)

func (c *Controller) PayOnline() {
	pro := c.GetString("provider")
	money, _ := c.GetInt64("money")
	if pro != "wx" || money <= 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	uid := int(c.UserID)
	user, err := service.GetUserInfo(uid)
	if err != nil {
		c.ReplyErr(err)
		return
	}

	if pro == "wx" {
		//微信下单接口
		or := new(model.Order)
		or.Status = model.YiUserOrder
		or.SubType = model.PwxPay
		or.Price = money * 10
		or.OrderType = model.OrderTopup
		or.Uid = uid
		or.OrderNo = service.GetTradeNO(or.OrderType, or.Uid)
		or.Desc = "AllSum账户充值"
		or.User = user
		or.Time = time.Now().Format(model.TimeFormat)
		beego.Debug(or)

		err = service.CreateOrder(or)
		if err != nil {
			c.ReplyErr(err)
			return
		}
		reply, err := wxisv.Pay.QrPay(or.OrderNo, or.Desc, or.Price)
		if err != nil || len(reply.CodeUrl) == 0 {
			beego.Error("weixin qrpay failed:", err)
			c.ReplyErr(errcode.ErrWXPay)
			return
		}
		c.ReplySucc(reply.CodeUrl)

	} else {
		//ali支付接口
	}
	return
}
