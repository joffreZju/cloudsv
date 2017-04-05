package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Agent struct {
	Id          int `orm:"auto;pk;column(id)"` // 表内自增
	Uid         int
	LicenseFile string
	Status      int       `orm:"default(1)"` //1 正常，2 禁用
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"`
	Consumers   []*User   `orm:"-" json:",omitempty"`
	User        *User     `orm:"rel(fk)" json:",omitempty"`
}

func InsertAgent(a *Agent) (err error) {
	id, err := orm.NewOrm().Insert(a)
	if err != nil {
		return
	}
	a.Id = int(id)
	return
}

func SetAgentStatus(a *Agent) (err error) {
	_, err = orm.NewOrm().Update(a, "Status")
	return err
}

func AgentUpdate(a *Agent) (err error) {
	_, err = orm.NewOrm().Update(a, "Status", "LicenseFile")
	return err
}

func GetAgentInfo(id int) (a *Agent, err error) {
	o := NewOrm(ReadOnly)
	a = new(Agent)
	err = o.QueryTable("public.agent").Filter("id", id).One(&a)
	if err != nil {
		return nil, err
	}
	_, err = o.LoadRelated(a, "User")
	if err != nil {
		return nil, err
	}
	return
}

func GetAgentAll() (list []*Agent, err error) {
	o := NewOrm(ReadOnly)
	_, err = o.QueryTable("agent").All(&list)
	return

}