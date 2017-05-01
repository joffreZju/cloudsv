package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type CalTemplate struct {
	Id             int       `orm:"column(id);auto;pk" json:"-"`
	UserId         int       `orm:"column(user_id);" json:"-"`
	WaybillNumber  string    `orm:"column(waybill_number);null"`
	ActualWeight   string    `orm:"column(actual_weight);null"`
	ActualVolume   string    `orm:"column(actual_volume);null"`
	FreightCharges string    `orm:"column(freight_charges);null"`
	PackageNumber  string    `orm:"column(package_number);null"`
	Ctt            time.Time `orm:"column(ctt);type(timestamp with time zone);null" json:"-"`
}

func InsertTemplate(t *CalTemplate) (err error) {
	_, err = orm.NewOrm().Insert(t)
	if err != nil {
		return
	}
	return
}

func GetTemplate(uid int) (t *CalTemplate, err error) {
	t = new(CalTemplate)
	err = orm.NewOrm().QueryTable("CalTemplate").Filter("UserId", uid).OrderBy("-Ctt").Limit(1).One(t)
	return

}

type CalRecord struct {
	Id        int    `orm:"column(id);auto;pk"`
	UserId    int    `orm:"column(user_id);"`
	AccountId int64  `orm:"column(account_id);null"`
	CalNo     string `orm:"column(cal_no);unique"`
	//UsingCost  float64   `orm:"column(using_cost);null"`
	CalTimes   int `orm:"column(cal_times);"`
	LastResult int `orm:"column(last_result);null"`
	//UserType   string    `orm:"column(user_type);null"`
	CalType string    `orm:"column(cal_type);null"` //计算类型，金额；
	Ctt     time.Time `orm:"column(ctt);type(timestamp with time zone);"`
	Ltt     time.Time `orm:"type(timestamp with time zone);"` //最后一次计算的createTime时间
	Utt     time.Time `orm:"column(utt);type(timestamp with time zone);null"`
}

type CarSummary struct {
	Id           int       `orm:"column(id);auto;pk" json:"-"`
	CalRecordId  int       `orm:"column(cal_record_id);null" json:"-"`
	CalTimes     int       `orm:"column(cal_times);null" json:"-"`
	UserId       int       `orm:"column(user_id);null" json:"-"`
	CarNo        string    `orm:"column(car_no);null"`
	MaxVolume    float64   `orm:"column(max_volume);null"`
	MaxWeight    float64   `orm:"column(max_weight);null"`
	TotalMoney   float64   `orm:"column(total_money);null"`
	TotalWeight  float64   `orm:"column(total_weight);null"`
	TotalVolume  float64   `orm:"column(total_volume);null"`
	StowageRatio float64   `orm:"column(stowage_ratio);null"`
	Ctt          time.Time `orm:"column(ctt);type(timestamp with time zone);null" json:"-"`
	Utt          time.Time `orm:"column(utt);type(timestamp with time zone);null" json:"-"`
}

type CalGoods struct {
	Id             int       `orm:"column(id);auto;pk"  json:"-"`
	CalRecordId    int       `orm:"column(cal_record_id);null" json:"-"`
	CalTimes       int       `orm:"column(cal_times);null" json:"-"` //计算次数
	WaybillNumber  string    `orm:"column(waybill_number);null"`     //运单号
	ActualWeight   float64   `orm:"column(actual_weight);null"`      //重量
	ActualVolume   float64   `orm:"column(actual_volume);null"`      //体积
	FreightCharges float64   `orm:"column(freight_charges);null"`    //金额
	PackageNumber  int       `orm:"column(package_number);null"`     //包裹
	Necessary      string    `orm:"column(necessary);null"`
	Understowed    string    `orm:"column(understowed);null"` //打底
	OtherInfo      string    `orm:"column(other_info);null"`
	Split          string    `orm:"column(split);null"` //拆分
	SplitInfo      string    `orm:"column(split_info);null" json:"-"`
	CalResult      string    `orm:"column(cal_result);null"` //保存计算结果,车牌号
	Ctt            time.Time `orm:"column(ctt);type(timestamp with time zone);null" json:"-"`
	Utt            time.Time `orm:"column(utt);type(timestamp with time zone);null" json:"-"`
}

func GetCalRecord(calNo string) (cr *CalRecord, err error) {
	cr = new(CalRecord)
	err = orm.NewOrm().QueryTable("CalTRecord").Filter("CalNo", calNo).One(cr)
	return
}

func UpdateCalRecord(o orm.Ormer, cr *CalRecord) (err error) {
	_, err = o.Update(cr)
	return
}

func InsertCalRecord(o orm.Ormer, cr *CalRecord) (err error) {
	var id int64
	id, err = o.Insert(cr)
	if err != nil {
		return
	}
	cr.Id = int(id)
	return
}

func UpsertCalRecord(o orm.Ormer, r *CalRecord) (err error) {
	//o := orm.NewOrm()
	err = o.QueryTable("CalRecord").Filter("CalNo", r.CalNo).One(r)
	if err == orm.ErrNoRows {
		var id int64
		id, err = o.Insert(r)
		if err != nil {
			return
		}
		r.Id = int(id)
		return
	} else if err == nil {
		_, err = o.Update(r)
	}
	return
}

//查找用户的历史使用频率最高的车辆
func GetFrequentCars(uid int) (cars []*CarSummary, err error) {
	//Limit := cons.NUMBER_OF_HISTORY_CARS
	sql := `select car_no,
			max(max_weight) as max_weight,
			max(max_volume) as max_volume,
			count(*) as co
			from car_summary
		where user_id = ?
		group by car_no
		order by co DESC`
	o := orm.NewOrm()
	cars = []*CarSummary{}
	if _, err = o.Raw(sql, uid).QueryRows(&cars); err != nil {
		return nil, err
	}
	return cars, nil
}

func InsertCars(o orm.Ormer, cs []*CarSummary) (err error) {
	//o := orm.NewOrm()
	for _, v := range cs {
		var id int64
		id, err = o.Insert(v)
		if err != nil {
			return
		}
		v.Id = int(id)
	}
	return
}

func InsertGoods(o orm.Ormer, gs []*CalGoods) (err error) {
	//o := orm.NewOrm()
	for _, v := range gs {
		var id int64
		id, err = o.Insert(v)
		if err != nil {
			return
		}
		v.Id = int(id)
	}
	return
}

func (t *CarSummary) TableName() string {
	return "car_summary"
}
