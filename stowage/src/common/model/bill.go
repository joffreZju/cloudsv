package model

type fundOrder struct {
	Id      int    `orm:"auto;pk"`
	OrderId string `orm:"unique"` //订单号
}
