package account

import (
	"common/controller/base"
	"common/service"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

//查看账户信息
func (c *Controller) AccountInfo() {
	uid := int(c.UserID)
	account, err := service.GetAccount(uid)
	if err != nil {
		beego.Error("get account failed:", err, uid)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(account)
	return
}
