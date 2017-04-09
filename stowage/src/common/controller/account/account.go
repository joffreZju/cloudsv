package account

import (
	"common/controller/base"
	"common/service"
)

type Controller struct {
	base.Controller
}

//查看账户信息
func (c *Controller) AccountInfo() {
	uid := int(c.UserId)
	account, err := service.GetAccount(uid)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(account)
	return
}