package order

import (
	"common/lib/errcode"
	"common/model"
	"common/service"
	"time"
)

func (c *Controller) PayOnline() {
	userid := int(c.UserId)
	price, _ := c.GetInt64("price")
	orderType, _ := c.GetInt("order_type")
	if userId == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	user, err := model.GetUser(userId)
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
		OrderId:   service.GetTradeNO(4, userId),
		OrderType: orderType,
		Time:      time.Now().Format(model.TimeFormat),
		Status:    status,
	}
}
