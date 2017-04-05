package model

import "github.com/astaxie/beego/orm"

const (
	//PaidType
	PwxPay  = 1
	PaliPay = 2

	OrderWaitProcess = 0 // 待处理
	OrderProcessing  = 1 // 处理中
	OrderFinished    = 2 // 处理完成
	OrderCanceled    = 3 // 取消

	//status
	YiUserOrder        = iota // 0 创建订单
	YiCancel                  //取消
	YiWaitPay                 //等待支付完成
	YiWaitPayCanCancel        //等待支付，可被取消
	YiPaid                    //支付完成
	YiPayBack                 //退款
	YiPayBackFinish           //推款成功

	//order type
	OrderStowage = 1
	OrderTopup   = 2
)

type OrderStatus struct {
	Id     int `orm:"auto;pk;"`
	Status int
	Time   string
	Msg    string
	Order  *Order `orm:"rel(fk)"`
	User   *User  `orm:"null;rel(fk);"`
}

type Order struct {
	Id            int `orm:"auto;pk"`
	Orderid       string
	PayOrderid    string `json:,omitempty`
	CreateT       string `json:,omitempty`
	Status        int    `json:,omitempty`
	ProcessStatus int    `json:,omitempty`
	PaidType      int    `json:",omitempty`  //支付方式，渠道
	Price         int64  `json:,omitempty`   //用户支付金额
	AgentSharing  int64  `json:,omitempty`   //代理商分成
	Desc          string `json:",omitempty"` //备注信息
	Remark        string `json:",omitempty"` //附加信息
	PaidBankType  string `json:",omitempty"` // 银行卡类型, 微信支付有
	OrderType     int    `json:",omitempty"` //1.算配载 2.算路由 3.充值,
	User          *User  `orm:"-" json:",omitempty"`
	Uid           int    `json:"-"`
	Bill          *Bill  `orm:"reverse(one);column(bill_id)" json:",omitempty"`
	Time          string `json:",omitempty"` // 下单时间
}

func (o *Order) UpdateProcessStatus() {
	if o == nil {
		return
	}
	switch o.OrderType {
	case OrderTopup:
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
}

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
	err = orm.NewOrm().QueryTable("Order").Filter("OrderId", orderId).One(o)
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
