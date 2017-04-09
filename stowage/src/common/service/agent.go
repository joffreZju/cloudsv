package service

import (
	"common/lib/errcode"
	"common/model"

	"github.com/astaxie/beego"
)

func AgentCreate(a *model.Agent) (err error) {
	err = model.CreateUserIfNotExist(a.User)
	a.Uid = a.User.Id
	err = model.InsertAgent(a)
	if err != nil {
		beego.Error("agent create failed:", err)
		err = errcode.ErrAgentCreatFailed
		return
	}
	return
}
func AgentSetStatus() {}

func AgentClients(tel string) (users []*model.User, err error) {
	users, err = model.GetUsersByReferer(tel)
	if err != nil {
		return
	}
	return
}
