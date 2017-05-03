package model

import "github.com/astaxie/beego/orm"

type fundOrder struct {
	Id      int    `orm:"auto;pk"`
	OrderId string `orm:"unique"` //订单号
}

const (
	//checkStatus
	BillNotchecked = 1
	BillChecked    = 2

	//bill type
	//same as order type
)

//用户在本平台服务消费记录也存此
type Bill struct {
	Id          int      `orm:"auto;pk"`
	BillNo      string   `orm:"unique"`
	User        *User    `orm:"-" json:",omitempty"`
	UserId      int      `json:"-"`
	AccountId   int      `jons:"-"`
	Account     *Account `orm:"-" json:",omitempty"`
	Order       *Order   `orm:"null;rel(one)" json:",omitempty"`
	Type        int      //账单类型
	SubType     int      `orm:"null"` //子类型
	Time        string
	Money       int64
	ExtMsg      string //附加信息
	CheckStatus int    //通过流水核对结果
}

//获取bill
func GetBill(bid int) (bill *Bill, err error) {
	bill = &Bill{Id: bid}
	err = orm.NewOrm().Read(bill)
	bill.User = &User{Id: bill.UserId}
	return
}

//获取bills
func GetBillByIds(ids []int) (list []*Bill, err error) {
	if len(ids) == 0 {
		return
	}
	billIds := []interface{}{}
	for _, id := range ids {
		billIds = append(billIds, id)
	}
	_, err = orm.NewOrm().QueryTable("Bill").Filter("Id__in", billIds).All(&list)
	for _, b := range list {
		b.User = &User{Id: b.UserId}
	}
	return
}

//创建账单
func InsertBill(b *Bill) (err error) {
	if b != nil && b.User != nil && b.UserId == 0 {
		b.UserId = b.User.Id
	}
	id, err := orm.NewOrm().Insert(b)
	if err != nil {
		b.Id = int(id)
	}
	return
}

//更新账单状态
func UpdateBillStatus(bid int) (err error) {
	o := orm.NewOrm()
	bill := &Bill{
		Id:          bid,
		CheckStatus: BillNotchecked,
	}
	_, err = o.Update(bill, "CheckStatus")
	if err != nil {
		return err
	}
	return
}

//根据用户id获取账单列表,分页
func GetBillsByUser(id, page int) (list []*Bill, err error) {
	_, err = orm.NewOrm().QueryTable("Bill").Filter("UserId", id).
		OrderBy("-Time").Limit(30).Offset(page * 30).All(&list)
	return
}

func GetUserBillsByType(uid, tp int) (list []*Bill, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("Bill").Filter("type", tp).Filter("UserId", uid).All(&list)
	return
}

//根据账单类型获取列表，分页
func GetBillsByType(page, offset, tp int) (cn int64, list []*Bill, err error) {
	o := orm.NewOrm()
	if page == 0 {
		cn, err = o.QueryTable("Bill").Filter("type", tp).Count()
		if err != nil {
			return
		}
	}
	_, err = o.QueryTable("Bill").Filter("type", tp).
		OrderBy("-Id").Limit(offset).Offset(page * offset).All(&list)
	return
}

//根据子类型获取列表
func GetBillsBySubType(page, offset, stp int) (cn int64, list []*Bill, err error) {
	o := orm.NewOrm()
	if page == 0 {
		cn, err = o.QueryTable("Bill").Filter("SubType", stp).Count()
		if err != nil {
			return
		}
	}
	_, err = o.QueryTable("Bill").Filter("SubType", stp).
		OrderBy("-Id").Limit(offset).Offset(page * offset).All(&list)
	return

}

//更新bill
func UpdateBill(bill *Bill, fields ...string) (err error) {
	newFields := []string{}
	for _, f := range fields {
		if f == "User" {
			if bill != nil && bill.User != nil {
				bill.UserId = bill.User.Id
			}
			newFields = append(newFields, "UserId")
		}
		newFields = append(newFields, f)
	}
	_, err = orm.NewOrm().Update(bill, fields...)
	return
}
