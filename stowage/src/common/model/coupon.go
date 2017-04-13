package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Coupon struct {
	Id           int       `orm:"auto;pk;column(id)" json:"-"`
	Number       int       `orm:"unique"`
	VerifyCode   string    //核销码
	Denomination int       //面额
	Status       int       //0 初始，1 发放，2 回收，3 已使用
	Agentid      int       `orm:"null" json:",omitempty"` //代理商id
	Agent        *Agent    `orm:"-" json:",omitempty"`
	Userid       int       `orm:"null" json:",omitempty"` //消费者id
	User         *User     `orm:"-" json:",omitempty"`
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UsedTime     time.Time `orm:"type(datetime);null" json:",omitempty"`
}

func InsertCouponMulti(a []Coupon) (err error) {
	_, err = orm.NewOrm().InsertMulti(len(a), a)
	return err
}

func UpdateCouponMulti(s, e int, aid int) (count int64, err error) {
	count, err = orm.NewOrm().QueryTable("coupon").Filter("Number_gte", s).Filter("Number_lte", e).Filter("Status_exact", 0).Update(orm.Params{"Agentid": aid, "Status": 1})
	return
}

func GetCoupon(num int) (r *Coupon, err error) {
	r = &Coupon{Number: num}
	err = NewOrm(ReadOnly).Read(r)
	return
}

func UpdateCoupon(r *Coupon, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(r, fields...)
	return
}
func UpdateCouponByAgent(aid int) (err error) {
	_, err = orm.NewOrm().QueryTable("coupon").Filter("Agentid", aid).Filter("Status_exact", 1).Update(orm.Params{"Status": 2})
	return
}
