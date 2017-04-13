package service

import (
	"common/lib/util"
	"common/model"
	"errors"
	"fmt"
	"time"
)

const (
	LeastNumber  int = 10000000
	DefaultPrice int = 200
)

func UsingCoupon(num int, uid int, code string) (err error) {
	r, err := model.GetCoupon(num)
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
	r.Userid = uid
	r.UsedTime = time.Now()
	return model.UpdateCoupon(r, "Status", "Userid", "UsedTime")
}

func AddCoupons(start, end int) (err error) {
	caps := end - start
	if caps < 1 || start < LeastNumber {
		return fmt.Errorf("coupon range wrong:%d:%d", start, end)
	}
	charges := make([]model.Coupon, caps)
	for i := 0; i < caps; i++ {
		charges[i].Number = start + i
		charges[i].VerifyCode = util.RandomByte6()
		charges[i].Denomination = DefaultPrice
		charges[i].Status = 0
	}
	return model.InsertCouponMulti(charges)
}

func GrantAgent(start, end int, aid int) (err error) {
	caps := end - start
	if caps < 1 || start < LeastNumber {
		return fmt.Errorf("coupon range wrong:%d:%d", start, end)
	}
	/*
		user, err := model.GetUser(aid)
		if err != nil {
			return errcode.ErrUserNotExisted
		}*/

	count, err := model.UpdateCouponMulti(start, end, aid)
	if err != nil {
		return err
	}
	if int(count) != caps {
		//TODO   rowback
		return errors.New("发放号段有误，请仔细核对")
	}
	return
}

//func couponRecycle() (err error) {
//
//}
