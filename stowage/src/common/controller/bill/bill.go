package bill

import (
	"common/controller/base"
	"common/lib/errcode"
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

func (c *Controller) GetBillsType() {
	tp, err := c.GetInt("tp")
	if err != nil {
		c.ReplyErr(errcode.ErrParams)
		return
	}

	limit := 30
	page, _ := c.GetInt("page")
	ct, list, err := service.GetBillsByType(page, limit, tp)
	if err != nil {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"total": ct,
		"bills": list,
	})
}
