package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Coupon struct {
	Id           int       `orm:"auto;pk;column(id)" json:"-"`
	Number       int       `orm:"unique"`   //编码
	VerifyCode   string    `orm:"size(16)"` //核销码
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

//发放代金券
func GrantCouponMulti(s, e int, aid int) (count int64, err error) {
	count, err = orm.NewOrm().QueryTable("Coupon").Filter("Number__gte", s).Filter("Number__lte", e).Filter("Status__exact", 0).Update(orm.Params{"Agentid": aid, "Status": 1})
	return
}

func GetCoupon(num int) (r *Coupon, err error) {
	r = &Coupon{Number: num}
	err = NewOrm(ReadOnly).Read(r)
	return
}
func VerifyCoupon(uid int, num int) (err error) {
	_, err = orm.NewOrm().QueryTable("Coupon").Filter("Number", num).Filter("Status", 1).Update(orm.Params{
		"Status":   3,
		"UsedTime": time.Now(),
		"Userid":   uid,
	})
	return
}
func UpdateCoupon(r *Coupon, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(r, fields...)
	return
}

func RecycleCouponByAgent(aid int) (err error) {
	_, err = orm.NewOrm().QueryTable("Coupon").Filter("Agentid", aid).Filter("Status__exact", 1).Update(orm.Params{"Status": 2})
	return
}

func CouponList(aid int, page int) (count int, list []*Coupon, err error) {
	o := orm.NewOrm()
	if page == 1 {
		var ct int64
		ct, err = o.QueryTable("Coupon").Filter("Agentid", aid).Count()
		if err != nil {
			return
		}
		count = int(ct)
	}
	_, err = orm.NewOrm().QueryTable("Coupon").Filter("Agentid", aid).Limit(20).OrderBy("Id").Limit(20).Offset(page * 20).All(&list)
	return
}
