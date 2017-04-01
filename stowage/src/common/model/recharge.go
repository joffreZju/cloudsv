package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Recharge struct {
	Id           int       `orm:"auto;pk;column(id)" json:"-"`
	Number       int       `orm:"unique"`
	VerifyCode   string    //核销码
	Denomination int       //面额
	Status       int       //0 初始，1 发放，2 回收，3 已使用
	AgentUser    string    `orm:"null"` //代理商id,手机号
	User         string    `orm:"null"` //消费者id
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UsedTime     time.Time `orm:"type(datetime);null" json:",omitempty"`
}

func InsertRechargeMulti(a []Recharge) (err error) {
	_, err = orm.NewOrm().InsertMulti(len(a), a)
	return err
}

func UpdateRechargeMulti(s, e int, aid string) (count int64, err error) {
	count, err = orm.NewOrm().QueryTable("recharge").Filter("Number_gte", s).Filter("Number_lte", e).Filter("Status_exact", 0).Update(orm.Params{"AgentUser": aid, "Status": 1})
	return
}

func GetRecharge(num int) (r *Recharge, err error) {
	r = &Recharge{Number: num}
	err = NewOrm(ReadOnly).Read(r)
	return
}

func UpdateRecharge(r *Recharge, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(r, fields...)
	return
}
func UpdateRechargeByAgent(aid string) (err error) {
	_, err = orm.NewOrm().QueryTable("recharge").Filter("AgentUser", aid).Filter("Status_exact", 1).Update(orm.Params{"Status": 2})
	return
}
