package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	//PaidType
	PwxPay  = 1
	PaliPay = 2
	PCoupon = 3

	//process status
	OrderWaitProcess = 0 // 待处理
	OrderProcessing  = 1 // 处理中
	OrderFinished    = 2 // 处理完成
	OrderCanceled    = 3 // 取消

	//order status
	YiUserOrder        = iota // 0 创建订单
	YiCancel                  //取消
	YiWaitPay                 //等待支付完成
	YiWaitPayCanCancel        //等待支付，可被取消
	YiPaid                    //支付完成
	YiPayBack                 //退款
	YiPayBackFinish           //推款成功

	//order type
	OrderTopup   = 1
	OrderConsume = 2
	//sub type
	ConsumeStowage = 1
)

type OrderStatus struct {
	Id     int `orm:"auto;pk;"`
	Status int
	Time   string
	Msg    string
	Order  *Order `orm:"rel(fk)"`
	User   *User  `orm:"null;rel(fk);"`
}

//用户平台内消费也存于此
type Order struct {
	Id            int    `orm:"auto;pk"`
	Orderid       string `orm:"unique"`
	PayOrderId    string `orm:"null" json:,omitempty`
	CreateTime    string `orm:"auto_now_add;type(datetime)"`
	Status        int    `json:,omitempty`
	ProcessStatus int    `json:,omitempty`
	PaidType      int    `json:",omitempty`             //支付方式，渠道
	Price         int64  `json:,omitempty`              //用户支付金额
	AgentSharing  int64  `orm:"null" json:,omitempty`   //代理商分成
	Desc          string `orm:"null" json:",omitempty"` //备注信息
	Remark        string `orm:"null" json:",omitempty"` //附加信息
	PaidBankType  string `orm:"null" json:",omitempty"` // 银行卡类型, 微信支付有
	OrderType     int    `json:",omitempty"`            //1.充值2.消费
	SubType       int    `orm:"null" json:",omitempty"` //1.算配载 2.算路由
	Time          string `json:"-"`                     // 下单时间
	User          *User  `orm:"-" json:",omitempty"`
	Uid           int    `json:"-"`
	Bill          *Bill  `orm:"reverse(one);column(bill_id)" json:",omitempty"`
	Agent         *Agent `orm:"rel(fk);null;column(agent_id)" json:",omitempty"`
}

func (u *Order) TableName() string {
	return "allsum_order"
}

func (o *Order) UpdateProcessStatus() {
	if o == nil {
		return
	}
	if o.Status == YiUserOrder {
		o.ProcessStatus = OrderWaitProcess
	} else if o.Status == YiPaid {
		o.ProcessStatus = OrderFinished
	} else if o.Status == YiCancel {
		o.ProcessStatus = OrderCanceled
	} else {
		o.ProcessStatus = OrderProcessing
	}

}

func GetPaidOrderOfToday(aid int) (list []*Order, err error) {
	day := time.Now().Format("2006-01-02")
	start := day + " 00:00:00"
	end := time.Now().Format("2006-01-02 15:04:05")
	var tlist []*Order
	o := NewOrm(ReadOnly)
	_, err = o.QueryTable("allsum_order").Filter("agent_id", aid).
		Filter("Status", YiPaid).
		Filter("Time__gte", start).
		Filter("Time__lte", end).
		OrderBy("Time").All(&tlist)
	for _, l := range list {
		l.User = &User{Id: l.Uid}
	}
	return
}

/*
func GetUnhandleOrderList() (list []*Order, err error) {
	var blist []*Order
	o := NewOrm(ReadOnly)
	_, err = o.QueryTable("Order").Filter("Status", YiUserOrder)

}*/

func CreateOrder(o *Order) (err error) {
	if o != nil && o.User != nil && o.Uid == 0 {
		o.Uid = o.User.Id
	}
	id, err := orm.NewOrm().Insert(o)
	if err == nil {
		o.Id = int(id)
	}
	return
}

func GetOrder(oid int) (o *Order, err error) {
	o = &Order{Id: oid}
	err = orm.NewOrm().Read(o)
	o.User = &User{Id: o.Uid}
	return
}
func UpdateOrder(o *Order, fields ...string) (err error) {
	_, err = orm.NewOrm().Update(o, fields...)
	return
}

func GetOrderByOrderId(orderId string) (o *Order, err error) {
	o = new(Order)
	err = orm.NewOrm().QueryTable("allsum_order").Filter("OrderId", orderId).One(o)
	if o != nil {
		o.User = &User{Id: o.Uid}
	}
	return
}

func CreateOrderStatus(o *OrderStatus) (err error) {
	id, err := orm.NewOrm().Insert(o)
	if err == nil {
		o.Id = int(id)
	}
	return
}
