package order

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/service"
)

type Controller struct {
	base.Controller
}

func (c *Controller) OrderInfo() {
	id := c.GetString("order_no")
	or, err := service.GetOrderByOrderId(id)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(or)
}
func (c *Controller) OrderDay() {

}

func (c *Controller) OrderListIncomplete() {
	tp, err := c.GetInt("tp")
	if err != nil {
		c.ReplyErr(errcode.ErrParams)
		return
	}

	page, _ := c.GetInt("page")
	ct, list, err := service.GetOrdersIncomplete(page, 20, tp)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"total":  ct,
		"orders": list,
	})
}

func (c *Controller) OrdersCoupon() {
	stp, err := c.GetInt("stp")
	if err != nil {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	page, _ := c.GetInt("page")
	ct, list, err := service.GetOrdersSubtp(page, 20, stp)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"total":  ct,
		"orders": list,
	})
}
