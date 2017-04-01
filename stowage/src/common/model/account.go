package model

import "github.com/astaxie/beego/orm"

type Account struct {
	Id        int    `orm:"pk;auto"`
	Accountid string `orm:"unique"` //以手机号为id,目前等于userid
	Userid    string //个人和企业id
	UserType  int    //1 个人，2代理商，3企业
	Banlance  int64  //账户余额
	Topup     int64  //充值金额,消费者
	Spend     int64  //消费金额
	Withdraw  int64  //代理商提现
	Status    int    //1正常，2冻结
}

func AddAccount(a *Account) (err error) {
	id, err := orm.NewOrm().Insert(a)
	if err != nil {
		return
	}
	a.Id = int(id)
	return
}
func CheckAccountExist(uid string) (b bool) {
	count, err := orm.NewOrm().QueryTable("Account").Filter("Userid", uid).Count()
	if err == nil && count > 0 {
		b = true
	}
	b = false
	return
}

func GetAccountByAccountid(id string) (a *Account, err error) {
	a = &Account{Accountid: id}
	err = orm.NewOrm().Read(a, "Accountid")
	return
}

func UpdateAccount(a *Account, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(a, fields...)
	return
}
