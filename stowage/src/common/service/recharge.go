package service

import (
	"common/lib/util"
	"common/model"
	"errors"
	"time"
)

const (
	LeastNumber  int = 10000000
	DefaultPrice int = 200
)

func UsingRecharge(num int, uid string, code string) (err error) {
	r, err := model.GetRecharge(num)
	if err != nil {
		return err
	}
	if r.Status != 1 {
		return errors.New("代金券编号错误")
	}
	if r.VerifyCode != code {
		return errors.New("核销码错误")
	}
	r.Status = 3
	r.User = uid
	r.UsedTime = time.Now()
	return model.UpdateRecharge(r, "Status", "User", "UsedTime")
}

func AddRecharges(start, end int) (err error) {
	caps := end - start
	if caps < 1 || start < LeastNumber {
		return errors.New("recharge range wrong")
	}
	charges := make([]model.Recharge, caps)
	for i := 0; i < caps; i++ {
		charges[i].Number = start + i
		charges[i].VerifyCode = util.RandomByte6()
		charges[i].Denomination = DefaultPrice
		charges[i].Status = 0
	}
	return model.InsertRechargeMulti(charges)
}

func GrantAgent(start, end int, aid string) (err error) {
	caps := end - start
	if caps < 1 || start < LeastNumber {
		return errors.New("recharge range wrong")
	}
	count, err := model.UpdateRechargeMulti(start, end, aid)
	if err != nil {
		return err
	}
	if int(count) != caps {
		//TODO   rowback
		return errors.New("发放号段有误，请仔细核对")
	}
	return
}

//func RechargeRecycle() (err error) {
//
//}
