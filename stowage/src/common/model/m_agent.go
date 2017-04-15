package model

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Agent struct {
	Id          int       `orm:"auto;pk;column(id)" json:"-"` // 表内自增
	Uid         int       `orm:"unique"`
	Name        string    `orm:"size(64)"`
	LicenseFile string    `orm:"size(64)"`
	Status      int       `orm:"default(1)"` //1 正常，2 禁用
	Desc        string    `orm:"null"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"`
	Consumers   []*User   `orm:"-" json:",omitempty"`
	ConsumNo    int       `orm:"-"`
	User        *User     `orm:"rel(fk)" json:",omitempty"`
	Account     *Account  `orm:"-" json:",omitempty"`
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
	if a.Status != 1 || a.Status != 2 {
		return fmt.Errorf("status param wrong")
	}
	_, err = orm.NewOrm().Update(a, "Status")
	return err
}

func AgentUpdate(a *Agent) (err error) {
	_, err = orm.NewOrm().QueryTable("Agent").Filter("Uid", a.Uid).Update(orm.Params{
		"Status":      a.Status,
		"LicenseFIle": a.LicenseFile,
		"Desc":        a.Desc,
	})
	return err
}

func GetAgentInfo(uid int) (a *Agent, err error) {
	o := NewOrm(ReadOnly)
	a = new(Agent)
	err = o.QueryTable("agent").Filter("Uid", uid).One(a)
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