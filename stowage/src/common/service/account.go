package service

import (
	"common/model"
	"errors"
)

//建立资金账户
func AccountCreate(a *model.Account) (err error) {
	return model.AddAccount(a)
}

//个人充值，代理商收入
func ChargeAccount(aid string, money int64) (err error) {
	a, err := model.GetAccountByAccountid(aid)
	if err != nil {
		return err
	}
	a.Banlance = a.Banlance + money
	a.Topup = a.Topup + money
	err = model.UpdateAccount(a, "Banlance", "Topup")
	return
}

//个人消费
func DecreseAccount(aid string, money int64) (err error) {
	a, err := model.GetAccountByAccountid(aid)
	if err != nil {
		return err
	}
	if a.Banlance < money {
		return errors.New("账户余额不足")
	}
	a.Spend = a.Spend + money
	a.Banlance = a.Banlance - money
	err = model.UpdateAccount(a, "Banlance", "Spend")
	return
}

//提现,线下
//TODO   加锁
func WithdrawDeposit(aid string, money int64) (err error) {
	a, err := model.GetAccountByAccountid(aid)
	if err != nil {
		return err
	}
	if a.Banlance < money {
		return errors.New("账户余额不足")
	}
	a.Banlance -= money
	a.Withdraw += money
	err = model.UpdateAccount(a, "Banlance", "Withdraw")
	return
}