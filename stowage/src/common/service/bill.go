package service

import (
	"common/lib/errcode"
	"common/model"

	"github.com/astaxie/beego"
)

func GetBillsConsume(page int) ([]*model.Bill, error) {
	return model.GetBillsByType(page, model.BillConsume)
}
func GetBillsTopup(page int) ([]*model.Bill, error) {
	return model.GetBillsByType(page, model.BillTopUp)
}

func GetBill(bid int) (bill *model.Bill, err error) {
	bill, err = model.GetBill(bid)
	if err != nil {
		beego.Error("GetBill error:", err)
		err = errcode.ErrGetBillFailed
	}
	return
}

func GetBillByIds(ids []int) (list []*model.Bill, err error) {
	list, err = model.GetBillByIds(ids)
	if err != nil {
		beego.Error("GetBillByIds error:", err)
	}
	return
}
