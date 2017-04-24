package order

import (
	"common/lib/payment/wxisv"
	"common/model"
	"common/service"
	"encoding/xml"
	"fmt"

	"github.com/astaxie/beego"
)

func (c *Controller) WxPayback() {
	body := c.Ctx.Input.RequestBody
	if len(body) == 0 {
		body = c.Ctx.Input.CopyBody(beego.BConfig.MaxMemory)
	}
	resp := wxisv.NotifyResponse{}
	err := xml.Unmarshal(body, &resp)
	if err != nil {
		beego.Error("WxPay xml.Unmarshal error:", err)
		c.Ctx.WriteString("error")
		return
	} else if resp.ReturnCode != "SUCCESS" {
		beego.Error("WxPay ReturnCode:", resp.ReturnCode)
		c.Ctx.WriteString("error")
		return
	}

	// 校验
	//xmlMap, err := wxisv.Xml2Map(resp)
	//if err != nil {
	//	beego.Error("WxPay Xml2Map convert error:", err)
	//	c.Ctx.WriteString("error")
	//	return
	//}

	//sign := wxisv.Sign(xmlMap, "9ot34qkz0o9qxo4tvdjp9g98um4zw9wx")
	//if resp.Sign != sign {
	//	fmt.Println("error, 签名校验失败")
	//	c.Ctx.WriteString("error")
	//	return
	//} else {
	trade_id := resp.OutTradeNo
	payOrderId := resp.TransactionId
	bankType := resp.BankType
	err = c.orderPayCallback(trade_id, payOrderId, model.PwxPay, bankType, resp.OpenId, false)
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	//}
	c.Ctx.WriteString("success")
}

func (c *Controller) orderPayCallback(orderId, payOrderId string, payType int, bankType string, payUserId string, splitAccount bool) (err error) {
	order, err := service.GetOrderByOrderId(orderId)
	if err != nil {
		return
	}

	beego.Debug(fmt.Sprintf("Order : %+v", *order))

	if order.Status == model.YiPaid {
		return nil
	}
	order.PayOrderId = payOrderId
	order.Status = model.YiPaid
	//order.UpdateProcessStatus()
	order.SubType = payType
	order.PaidBankType = bankType
	err = model.TransPayOnline(order)
	if err != nil {
		return
	}

	//service.GetOrderMoreDetail(order, "Bill")
	// 创建支付账单
	//if order.Bill == nil {
	//创建账单

	// 更新定单账单信息

	// 更新用户账户余额

	//} else {
	// TODO 已经有支付账单了，有问题
	//}
	// 插入到订单信息表
	//orderStatus := new(model.OrderStatus)
	//orderStatus.Status = model.YiPaid
	//orderStatus.Time = time.Now().Format(model.TimeFormat)
	//orderStatus.Order = order
	//service.CreateOrderStatus(orderStatus)
	beego.Info("success: ", order.OrderNo)

	return nil
}
