package order

import (
	"common/lib/errcode"
	"common/model"
	"common/service"
	"time"

	"github.com/astaxie/beego"
)

func (c *Controller) PayOnline() {
	userid := int(c.UserID)
	//price, _ := c.GetInt64("price")
	orderType, _ := c.GetInt("order_type")
	if userid == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	user, err := model.GetUser(userid)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	orderStatus := new(model.OrderStatus)
	status := model.YiWaitPayCanCancel
	orderStatus.Msg = "请在30分钟内完成支付"
	orderStatus.Status = status
	orderStatus.Time = time.Now().Format(model.TimeFormat)

	order := &model.Order{
		User:      user,
		Orderid:   service.GetTradeNO(4, userid),
		OrderType: orderType,
		Time:      time.Now().Format(model.TimeFormat),
		Status:    status,
	}
	beego.Debug(order)
}
