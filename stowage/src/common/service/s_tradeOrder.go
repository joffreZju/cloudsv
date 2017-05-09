package service

import (
	"common/lib/errcode"
	"common/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

func GetTradeNO(orderType int, userId int) string {
	ntime := time.Now()
	str := strings.Replace(ntime.Format("0102150405.000"), ".", "", 1)
	str = strconv.Itoa(orderType) + str
	str = str + fmt.Sprintf("%04d", userId)
	return str
}

func GetOrderByOrderId(orderId string) (order *model.Order, err error) {
	order, err = model.GetOrderByOrderId(orderId)
	if err != nil {
		beego.Error("get order failed:", orderId, err)
		err = errcode.ErrNoOrder
	}
	return
}

func CreateOrder(o *model.Order) (err error) {
	err = model.CreateOrder(o)
	if err != nil {
		beego.Error("create Order failed:", err)
		err = errcode.ErrCreateOrderFailed
	}
	return
}

func CreateOrderStatus(o *model.OrderStatus) (err error) {
	err = model.CreateOrderStatus(o)
	if err != nil {
		beego.Error("CreateOrderStatus error", err)
		err = errcode.ErrCreateOrderFailed
	}
	return
}

func UpdateOrder(o *model.Order, fields ...string) (err error) {
	err = model.UpdateOrder(o, fields...)
	return
}

func GetOrdersSubtp(page, limit, stp int) (ct int64, list []*model.Order, err error) {
	ct, list, err = model.GetPaidOrdersStp(page, limit, stp)
	if err != nil {
		beego.Error(err)
	}
	return
}

func GetOrdersIncomplete(page, limit, tp int) (ct int64, list []*model.Order, err error) {
	ct, list, err = model.GetOrderTypeIncomplete(page, limit, tp)
	if err != nil {
		beego.Error(err)
	}
	return
}

func GetOrderMoreDetail(order *model.Order, fields ...string) (err error) {
	newFields := []string{}
	for _, f := range fields {
		if f == "User" {
			uid := order.Uid
			if uid == 0 {
				uid = order.User.Id
			}
			order.User, _ = model.GetUser(uid)
		} else {
			newFields = append(newFields, f)
		}
	}
	err = model.LoadRelations(order, newFields...)
	if err != nil {
		beego.Error("GetOrderMoreDetail error", err)
	}
	return
}
