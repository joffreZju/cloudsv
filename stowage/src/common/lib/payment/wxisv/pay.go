package wxisv

import (
	"encoding/json"
	"fmt"
)

type Pay struct {
	Client *Client
}

func NewPay(cli ...*Client) *Pay {
	c := DefaultClient
	if len(cli) > 0 {
		c = cli[0]
	}
	return &Pay{c}
}

type MicroPayReply struct {
	OpendId     string `xml:"openid"`         // 用户在商户appid 下的唯一标识
	IsSubscribe string `xml:"is_subscribe"`   // 用户是否关注公众账号
	BankType    string `xml:"bank_type"`      // 银行类型
	TradeNo     string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo  string `xml:"out_trade_no"`   // 商户的订单号
}

func (p *Pay) MicroPay(storeId, subMchId, orderNo, desc string, totalAmount int64, authCode string) (reply *MicroPayReply, err error) {
	mp := map[string]string{
		"device_info":      storeId,
		"out_trade_no":     orderNo,
		"total_fee":        fmt.Sprintf("%.d", totalAmount),
		"body":             desc,
		"auth_code":        authCode,
		"spbill_create_ip": p.Client.LocalIP,
	}
	if len(subMchId) > 0 {
		mp["sub_mch_id"] = subMchId
	}
	var resp struct {
		CommonReply
		MicroPayReply
	}
	err = p.Client.sendCommand("pay/micropay", mp, &resp)
	if err != nil {
		return
	}
	reply = new(MicroPayReply)
	*reply = resp.MicroPayReply
	return
}

type PayReply struct {
	AppId     string `json:"appId,omitempty"`
	PartnerId string `json:"partnerId,omitempty"`
	PrepayId  string `json:",omitempty"`
	Package   string `json:"package,omitempty"`
	NonceStr  string `json:"nonceStr,omitempty"`
	Timestamp string `json:"timeStamp,omitempty"`
	Sign      string `json:",omitempty"`
	SignType  string `json:"signType,omitempty"`
	PaySign   string `json:"paySign,omitempty"`
	CodeUrl   string `json:",omitempty"`
}

type UnifiedOrderResult struct {
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	TradeType  string `xml:"trade_type"`
	PrepayId   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
}

func (p *Pay) PaymentRequest(tradeType string, prepay UnifiedOrderResult) *PayReply {
	param := make(map[string]string)
	param["appId"] = p.Client.AppId
	param["nonceStr"] = newNonceString()
	param["timeStamp"] = newTimestampString()

	payRequest := &PayReply{
		AppId:     p.Client.AppId,
		NonceStr:  param["nonceStr"],
		Timestamp: param["timeStamp"],
	}

	switch tradeType {
	case "APP": // app
		param["partnerid"] = p.Client.MchId
		param["package"] = "Sign=WXPay"
		payRequest.PartnerId = param["partnerid"]
		payRequest.Package = param["package"]
		payRequest.PrepayId = prepay.PrepayId
		sign := signMd5(param, p.Client.Key)
		payRequest.Sign = sign
	case "JSAPI": // 公众号
		param["signType"] = "MD5"
		param["package"] = "prepay_id=" + prepay.PrepayId
		payRequest.SignType = param["signType"]
		payRequest.Package = param["package"]
		sign := signMd5(param, p.Client.Key)
		payRequest.PaySign = sign
	case "NATIVE": // 二维码
		payRequest.CodeUrl = prepay.CodeUrl
	}

	return payRequest
}

func (p *Pay) WeChatOfficialAccountsPay(orderNo, desc string, totalAmount int64, openId string) (reply *PayReply, err error) {
	return p.unifiedOrder("JSAPI", orderNo, desc, totalAmount, openId)
}

func (p *Pay) QrPay(orderNo, desc string, totalAmount int64) (reply *PayReply, err error) {
	return p.unifiedOrder("NATIVE", orderNo, desc, totalAmount, "")
}

func (p *Pay) AppPay(storeId, subMchId, orderNo, desc string, totalAmount int64) (reply *PayReply, err error) {
	return p.unifiedOrder("APP", orderNo, desc, totalAmount, "")
}

func (p *Pay) unifiedOrder(tradeType, orderNo, desc string, totalAmount int64, openId string) (reply *PayReply, err error) {
	mp := map[string]string{
		"device_info":      "WEB",
		"out_trade_no":     orderNo,
		"total_fee":        fmt.Sprintf("%d", totalAmount),
		"body":             desc,
		"trade_type":       tradeType,
		"notify_url":       p.Client.NotifyUrl,
		"spbill_create_ip": p.Client.LocalIP,
	}
	if len(openId) > 0 {
		mp["openid"] = openId
	}
	var resp struct {
		CommonReply
		UnifiedOrderResult
	}
	err = p.Client.sendCommand("pay/unifiedorder", mp, &resp)
	if err != nil {
		return
	}
	reply = p.PaymentRequest(tradeType, resp.UnifiedOrderResult)
	return
}

type QueryOrderReply struct {
	TradeState string `xml:"trade_state"`
	OpendId    string `xml:"openid"`
	TotalFee   int64  `xml:"total_fee"`
	BankType   string `xml:"bank_type"`
	TradeNo    string `xml:"transaction_id"`
	TradeType  string `xml:"trade_type"`
	Attach     string `xml:"attach"`
}

func (p *Pay) QueryOrder(subMchId, orderNo string) (reply *QueryOrderReply, err error) {
	mp := map[string]string{
		"out_trade_no": orderNo,
	}
	if len(subMchId) > 0 {
		mp["sub_mch_id"] = subMchId
	}
	var resp struct {
		CommonReply
		QueryOrderReply
	}
	err = p.Client.sendCommand("pay/orderquery", mp, &resp)
	if err != nil {
		return
	}
	reply = new(QueryOrderReply)
	*reply = resp.QueryOrderReply
	return
}

func (p *Pay) ReportTil(url string, data interface{}) (err error) {
	mp := map[string]string{
		"interface_url": url,
	}
	trades, _ := json.Marshal(data)
	mp["trades"] = string(trades)

	var resp struct {
		CommonReply
	}
	err = p.Client.sendCommand("payitil/report", mp, &resp)

	return
}

func init() {
	setSignType("pay/micropay", "MD5")
	setSignType("pay/orderquery", "MD5")
	setSignType("pay/unifiedorder", "MD5")
	setSignType("payitil/report", "MD5")
}
