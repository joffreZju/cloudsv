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
	return
}

func CreateOrder(o *model.Order) (err error) {
	err = model.CreateOrder(o)
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

func GetOrderMoreDetail(order *model.Order, fields ...string) (err error) {
	newFields := []string{}
	for _, f := range fields {
		if f == "User" {
			uid := order.UId
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
