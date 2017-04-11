package bill

import (
	"common/controller/base"
	"common/service"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) BillInfo() {
	bid, _ := c.GetInt("bid")
	b, err := service.GetBill(bid)
	if err != nil {
		beego.Error("get bill failed:", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(*b)
}
