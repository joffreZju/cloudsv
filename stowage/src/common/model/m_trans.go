package model

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

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

	_, err = o.QueryTable("Order").Filter("Id", orderId).Update(orm.Params{"Status": YiPaid})
	if err != nil {
		o.Rollback()
		return
	}

	if or.OrderType == OrderTopup {
		a.Banlance = a.Banlance + b.Money
		a.Topup = a.Topup + b.Money
	} else {
		if a.Banlance < b.Money {
			o.Rollback()
			return errors.New("账户余额不足")
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
