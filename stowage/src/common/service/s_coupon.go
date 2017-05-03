package service

import (
	"common/lib/errcode"
	"common/lib/util"
	"common/model"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	LeastNumber  int = 10000000
	DefaultPrice int = 20000 //分为单位
)

func UsingCoupon(num int, uid int, code string) (err error) {
	r, err := model.GetCoupon(num)
	if err != nil {
		beego.Error(err)
		return errcode.ErrCouponNotExist
	}
	if r.Status == 3 {
		return errcode.ErrCouponUsed
	}
	if r.Status != 1 {
		return errcode.ErrCouponIllegal
	}
	if r.VerifyCode != code {
		return errcode.ErrCouponVerify
	}
	r.Status = 3
	r.Userid = uid
	r.UsedTime = time.Now()

	agent, err := model.GetAgentFromClient(uid)
	if err != nil {
		beego.Info("no find agent:", err)
		agent = nil
	}
	or := new(model.Order)
	or.Status = int(model.YiOrderCreate)
	or.ProcessStatus = model.OrderWaitProcess
	or.SubType = model.PCoupon
	or.CouponId = r.Id
	or.OrderNo = util.GetTradeNo(model.PCoupon, uid)
	or.Price = int64(DefaultPrice)
	or.OrderType = model.OrderTopup
	if agent != nil {
		or.AgentSharing = or.Price * int64(agent.Discount) / 100
		or.Agent = agent
	}
	err = model.TransCouponUsing(or, r, agent)
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

/*
func UsingCoupon(num int, uid int, code string) (err error) {
	r, err := model.GetCoupon(num)
	if err != nil {
		beego.Error(err)
		return errcode.ErrCouponNotExist
	}
	if r.Status == 3 {
		return errcode.ErrCouponUsed
	}
	if r.Status != 1 {
		return errcode.ErrCouponIllegal
	}
	if r.VerifyCode != code {
		return errcode.ErrCouponVerify
	}
	r.Status = 3
	r.Userid = uid
	r.UsedTime = time.Now()

	var sharing bool = true
	agent, err := model.GetAgentFromClient(uid)
	if err != nil {
		beego.Info("no find agent:", err)
		sharing = false
	}
	//TODO   事物管理
	//创建交易订单
	or := new(model.Order)
	or.Status = int(model.YiUserOrder)
	or.ProcessStatus = model.OrderWaitProcess
	or.SubType = model.PCoupon
	or.OrderNo = util.GetTradeNo(model.PCoupon, uid)
	or.Price = int64(DefaultPrice)
	or.OrderType = model.OrderTopup
	if sharing {
		or.AgentSharing = or.Price * int64(agent.Discount) / 100
		or.Agent = agent
	}
	err = CreateOrder(or)
	if err != nil {
		beego.Error("create order failed", err)
		return
	}
	//核销代金券
	err = model.UpdateCoupon(r, "Status", "Userid", "UsedTime")
	if err != nil {
		beego.Error("using coupon failed:", err)
		return err
	}
	//更新订单
	or.Status = model.YiPaid
	UpdateOrder(or, "Status")
	a, _ := GetAccountByUser(uid)

	//创建账单
	b := new(model.Bill)
	b.Money = int64(DefaultPrice)
	b.Time = time.Now().Format(model.TimeFormat)
	b.Order = or
	b.Type = model.OrderTopup
	b.UserId = uid
	b.AccountId = a.Id
	CreateBill(b)
	//更新用户资金账户
	ChargeAccount(a.Id, b.Money)
	//更新代理商账户

	return
}*/

func AddCoupons(start, end int) (err error) {
	caps := end - start
	if caps < 1 || start < LeastNumber {
		return fmt.Errorf("coupon range wrong:%d:%d", start, end)
	}
	charges := make([]model.Coupon, caps)
	for i := 0; i < caps; i++ {
		charges[i].Number = start + i
		charges[i].VerifyCode = util.RandomByte6(i + start)
		charges[i].Denomination = DefaultPrice
		charges[i].Status = 0
	}
	return model.InsertCouponMulti(charges)
}

//aid == agent's userid
func GrantAgent(start, end int, aid int) (err error) {
	caps := end - start + 1
	if caps < 1 || start < LeastNumber {
		beego.Error("coupon range wrong:%d:%d", start, end)
		err = errcode.ErrParams
	}
	/*
		user, err := model.GetUser(aid)
		if err != nil {
			return errcode.ErrUserNotExisted
		}*/

	count, err := model.GrantCouponMulti(start, end, aid)
	if err != nil {
		return err
	}
	if int(count) != caps {
		//TODO   rowback
		beego.Error("发放号段有误，请仔细核对")
		err = errcode.ErrCouponNo
	}
	return
}

func GetCoupon(no int) (c *model.Coupon, err error) {
	c, err = model.GetCoupon(no)
	if err != nil {
		if err == orm.ErrNoRows {
			err = errcode.ErrCouponNotExist
		} else {
			beego.Error("get coupon error:", err)
		}
	}
	return
}

func GetCouponList(aid, page int) (count int, list []*model.Coupon, err error) {
	count, list, err = model.CouponList(aid, page)
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

//func couponRecycle() (err error) {
//
//}
