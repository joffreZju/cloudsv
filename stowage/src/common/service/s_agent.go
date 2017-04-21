package service

import (
	"common/lib/errcode"
	"common/model"

	"github.com/astaxie/beego"
)

func AgentCreate(a *model.Agent) (err error) {
	err = model.CreateOrUpdateAuser(a.User)
	if err != nil {
		beego.Error(err)
		return
	}
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

func AgentGetList(page int) (total int, list []*model.Agent, err error) {
	total, list, err = model.GetAgentList(page)
	if err != nil {
		beego.Error(err)
		return
	}

	for _, v := range list {
		v.Account, _ = model.GetAccountByUserId(v.Uid)
		v.User, _ = model.GetUser(v.Uid)
		ct, _ := model.GetUserCountsByReferer(v.User.Tel)
		v.ConsumNo = int(ct)
	}
	return
}
