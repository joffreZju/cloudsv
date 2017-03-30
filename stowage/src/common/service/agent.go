package service

import (
	"common/lib/errcode"
	"common/model"
	"fmt"

	"github.com/astaxie/beego"
)

func AgentCreate(a *model.Agent) (err error) {
	err = model.CreateUserIfNotExist(a.User)
	fmt.Printf("=========%+v=====%v\n", a.User, err)
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
