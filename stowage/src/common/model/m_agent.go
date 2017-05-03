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
	Discount    int       `orm:"default(50)"` //折扣率，百分比
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

func GetAgentFromClient(uid int) (a *Agent, err error) {
	o := NewOrm(ReadOnly)
	u := &User{Id: uid}
	err = o.Read(u)
	if err != nil || u.AgentUid == 0 {
		return nil, fmt.Errorf("no agent of user:%d", uid)
	}
	a, err = GetAgentInfo(u.AgentUid)
	return
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

func GetAgentList() (list []*Agent, err error) {
	o := NewOrm(ReadOnly)
	//if page == 0 {
	//	var ct int64
	//	ct, err = o.QueryTable("agent").Count()
	//	if err != nil {
	//		return
	//	}
	//	total = int(ct)
	//}
	_, err = o.QueryTable("agent").OrderBy("Id").All(&list)
	return
}
