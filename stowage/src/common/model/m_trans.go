package model

import (
	"time"

	"common/lib/errcode"
	"github.com/astaxie/beego/orm"
)

func getBillNo(no string) string {
	return "9" + no
}

//使用数值加减函数
func TransCouponUsing(or *Order, co *Coupon, ag *Agent) (err error) {
	o := orm.NewOrm()
	o.Begin()
	//更新代金券
	_, err = o.Update(co, "Status", "Userid", "UsedTime")
	if err != nil {
		o.Rollback()
		return
	}
	//更新订单
	or.Status = YiPaid
	_, err = o.Update(or, "Status")
	if err != nil {
		o.Rollback()
		return
	}
	a := new(Account)
	err = o.QueryTable("Account").Filter("Userid", or.Uid).One(a)
	if err != nil {
		o.Rollback()
		return
	}

	//创建账单
	b := new(Bill)
	b.Money = or.Price
	b.Time = time.Now().Format(TimeFormat)
	b.Order = or
	b.BillNo = getBillNo(or.OrderNo)
	b.Type = or.OrderType
	b.UserId = or.Uid
	b.AccountId = a.Id
	_, err = o.Insert(b)
	if err != nil {
		o.Rollback()
		return
	}
	//更新资金账户
	_, err = o.Raw("update account set banlance= banlance + ?,topup= topup + ? where id= ?", b.Money, b.Money, a.Id).Exec()
	//a.Banlance = a.Banlance + b.Money
	//a.Topup = a.Topup + b.Money
	//err = o.Update(a, "Banlance", "Topup")
	if err != nil {
		o.Rollback()
		return
	}

	//代理商分成
	if ag != nil {
		_, err = o.Raw("update account set banlance= banlance + ? where id= ?", or.AgentSharing, ag.Account.Id).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	}

	err = o.Commit()
	return
}

func TransPayOnline(or *Order) (err error) {
	o := orm.NewOrm()
	o.Begin()
	_, err = o.Update(or)
	if err != nil {
		o.Rollback()
		return
	}
	a := new(Account)
	err = o.QueryTable("Account").Filter("Userid", or.Uid).One(a)
	if err != nil {
		o.Rollback()
		return
	}

	//创建账单
	b := new(Bill)
	b.Money = or.Price
	b.Time = time.Now().Format(TimeFormat)
	b.Order = or
	b.BillNo = getBillNo(or.OrderNo)
	b.Type = or.OrderType
	b.SubType = or.SubType
	b.UserId = or.Uid
	b.AccountId = a.Id
	_, err = o.Insert(b)
	if err != nil {
		o.Rollback()
		return
	}
	//更新资金账户
	_, err = o.Raw("update account set banlance= banlance + ?,topup= topup + ? where id= ?", b.Money, b.Money, a.Id).Exec()
	if err != nil {
		o.Rollback()
		return
	}
	err = o.Commit()

	return
}

func TransFinance(orderId int) (err error) {
	//更新order
	o := orm.NewOrm()
	or := Order{Id: orderId}
	err = o.Read(&or)
	if err != nil {
		return err
	}
	a, err := GetAccountById(or.Uid)
	if err != nil {
		return err
	}
	o.Begin()

	b := new(Bill)
	b.Money = or.Price
	b.Order = &or
	b.Type = or.OrderType
	b.Time = time.Now().Format(TimeFormat)
	b.UserId = or.Uid
	b.AccountId = a.Id
	_, err = o.Insert(b)
	if err != nil {
		o.Rollback()
		return
	}
	//更新订单
	or.Status = YiPaid
	_, err = o.Update(&or, "Status")
	if err != nil {
		o.Rollback()
		return
	}
	//和上面代码重复，表名写错
	//_, err = o.QueryTable("Order").Filter("Id", orderId).Update(orm.Params{"Status": YiPaid})
	//if err != nil {
	//	o.Rollback()
	//	return
	//}

	if or.OrderType == OrderTopup {
		a.Banlance = a.Banlance + b.Money
		a.Topup = a.Topup + b.Money
	} else {
		if a.Banlance < b.Money {
			o.Rollback()
			return errcode.ErrAccountFundLow
		}
		a.Spend = a.Spend + b.Money
		a.Banlance = a.Banlance - b.Money
	}
	_, err = o.Update(a)
	if err != nil {
		o.Rollback()
		return
	}
	err = o.Commit()
	return
}
