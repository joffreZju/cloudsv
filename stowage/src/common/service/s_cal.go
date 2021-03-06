package service

import (
	"common/lib/errcode"
	"common/model"
	"common/service/mqdto"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

func InsertTpl(t *model.CalTemplate) (err error) {
	err = model.InsertTemplate(t)
	if err != nil {
		beego.Error(err)
	}
	return
}

func GetTpl(uid int) (t *model.CalTemplate, err error) {
	t, err = model.GetTemplate(uid)
	if err != nil {
		beego.Error(err)
		err = errcode.ErrNoTpl
	}
	return
}

const calPrice = 1000

func InsertCalToDbAndSendToMq(uid int, cars []*model.CarSummary, goods []*model.CalGoods, record *model.CalRecord) (err error) {
	o := orm.NewOrm()
	o.Begin()
	calType := record.CalType
	if recordInDb, e := model.GetCalRecord(record.CalNo); e == nil {
		record = recordInDb
		if time.Now().Sub(record.Ctt).Hours() >= 48 {
			return errors.New("重复计算时间已超过48小时")
		}
		record.CalType = calType
		record.CalTimes += 1
		record.Ltt = time.Now()
		err = model.UpdateCalRecord(o, record)
		if err != nil {
			o.Rollback()
			beego.Error(err)
			return
		}
	} else {
		// 新的一次计算需要建立新账单，如果重复计算那么不需要创建账单
		or := &model.Order{
			Status:    model.YiOrderCreate,
			SubType:   model.CStowage,
			Price:     calPrice,
			OrderType: model.OrderConsume,
			Uid:       uid,
			OrderNo:   GetTradeNO(model.OrderConsume, uid),
			Desc:      "算配载计算费用",
			Time:      time.Now().Format(model.TimeFormat),
		}
		err = CreateOrder(or)
		if err != nil {
			return
		}
		record.OrderId = or.Id
		record.PayStatus = model.YiOrderCreate
		record.CalTimes = 1
		record.CalNo = fmt.Sprintf("%d%d", uid, time.Now().UnixNano())
		record.CalType = calType
		record.UserId = uid
		record.Ctt = time.Now()
		record.Ltt = time.Now()
		err = model.InsertCalRecord(o, record)
		if err != nil {
			o.Rollback()
			beego.Error(err)
			return
		}
	}
	for _, v := range cars {
		v.CalRecordId = record.Id
		v.CalTimes = record.CalTimes
		v.UserId = uid
		v.Ctt = time.Now()
	}
	for _, v := range goods {
		v.CalRecordId = record.Id
		v.CalTimes = record.CalTimes
		v.Ctt = time.Now()
	}
	err = model.InsertCars(o, cars)
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return
	}
	err = model.InsertGoods(o, goods)
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return
	} else if e := o.Commit(); e != nil {
		beego.Error(e)
	}
	err = SendCalToMq(cars, goods, record)
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Info("发送MQ消息成功:", record.CalNo, record.Id)
	return
}

func UpdateCalResult(result *mqdto.MQRespDto) (err error) {
	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		beego.Error(err)
		return
	}
	//更新spz_cal_record
	sql := `update cal_record
		set last_result =  ?,
			utt = ?
		where id = ? and last_result < ?`
	if _, err = o.Raw(sql, result.Cal_times, time.Now(), result.Using_id, result.Cal_times).Exec(); err != nil {
		o.Rollback()
		beego.Error("update cal result error with spz_cal_record", result.Using_id, err)
		return
	}
	//更新spz_cal_carsummary
	sql = `update car_summary
		set total_money =?,
			total_weight = ?,
			total_volume= ?,
			utt = ?,
			stowage_ratio = (?/max_weight + ?/max_volume)/2
		where cal_record_id = ? and cal_times = ? and car_no = ?`
	for _, v := range result.Car_summary {
		_, err = o.Raw(sql, v.Total_money, v.Total_weight, v.Total_volume, time.Now(),
			v.Total_weight, v.Total_volume,
			result.Using_id, result.Cal_times, v.Car_no).Exec()
		if err != nil {
			o.Rollback()
			beego.Error("update cal result error with spz_cal_carsummary", result.Using_id, err)
			return
		}
	}
	//更新spz_cal_lib
	sql = `update cal_goods
		set cal_result = ?,
			utt  = ?
		where id = ?`
	for _, v := range result.Result {
		if _, err = o.Raw(sql, v.Result, time.Now(), v.Id).Exec(); err != nil {
			o.Rollback()
			beego.Error("update cal result error with spz_cal_lib", result.Using_id, err)
			return
		}
	}
	if err = o.Commit(); err != nil {
		beego.Critical("update calculate result error with database transaction commit:", result.Using_id, err)
		return
	}
	beego.Debug("update calculate result success, using_id:", result.Using_id)
	return nil
}

//获取计算结果的carSummary
func GetCarsResult(calNo string) (cs []*model.CarSummary, err error) {
	cs = []*model.CarSummary{}
	sql := `select car.* from car_summary as car
			inner join cal_record as cr
			on car.cal_record_id = cr.id and car.cal_times=cr.cal_times
		where cr.cal_no=? and cr.pay_status = ? and cr.cal_times = cr.last_result`
	o := orm.NewOrm()
	if _, err = o.Raw(sql, calNo, model.YiPaid).QueryRows(&cs); err != nil {
		return nil, err
	}
	return cs, nil
}

//获取运单的计算结果
func GetGoodsResult(calNo string) (wbs []*model.CalGoods, err error) {
	wbs = []*model.CalGoods{}
	sql := `select cg.waybill_number,
			sum(cg.actual_volume) as actual_volume,
			sum(cg.actual_weight) as actual_weight,
			sum(cg.freight_charges) as freight_charges,
			sum(cg.package_number) as package_number,
			max(cg.necessary) as necessary,
			max(cg.understowed) as understowed,
			max(cg.split) as split,
			max(cg.other_info) as other_info,
			cg.cal_result
		from cal_goods as cg
		inner join cal_record as cr
		on cg.cal_record_id = cr.id and cg.cal_times=cr.cal_times
		where cr.cal_no=? and cg.split_info != ? and cr.pay_status = ? and cr.cal_times = cr.last_result
		group by cg.waybill_number,cg.cal_result`
	o := orm.NewOrm()
	if _, err = o.Raw(sql, calNo, model.WAYBILL_SPLIT_FROM, model.YiPaid).QueryRows(&wbs); err != nil {
		return nil, err
	}
	return wbs, nil
}

//重新计算，获取最后一次计算记录中已经编辑过的车辆数据
func GetEditedCars(calNo string) (cs []*model.CarSummary, err error) {
	cs = []*model.CarSummary{}
	sql := `select car.* from car_summary as car
			inner join cal_record as cr
			on car.cal_record_id = cr.id and car.cal_times=cr.cal_times
		where cr.cal_no=?`
	o := orm.NewOrm()
	if _, err = o.Raw(sql, calNo).QueryRows(&cs); err != nil {
		return nil, err
	}
	return cs, nil
}

//重新计算，获取最后一次计算记录的已经编辑过的运单数据，拆单的要合起来
func GetEditedWaybills(calNo string) (goods []*model.CalGoods, err error) {
	goods = []*model.CalGoods{}
	sql := `select cg.* from cal_goods as cg
			inner join cal_record as cr
			on cg.cal_record_id = cr.id and cg.cal_times=cr.cal_times
		where cr.cal_no=? and cg.split_info != ?`
	o := orm.NewOrm()
	if _, err = o.Raw(sql, calNo, model.WAYBILL_SPLIT_TO).QueryRows(&goods); err != nil {
		return nil, err
	}
	return goods, nil
}

type CalHistory struct {
	CalNo    string
	Ctt      time.Time
	CalTimes int
	Money    int64
}

//获取计算的历史记录
func GetCalHistory(uid, pageNumber, pageLimit int) (calRecords []*CalHistory, maxCount int, err error) {
	o := orm.NewOrm()
	//计数所有使用记录
	sql := `select count(distinct(cal_no)) from cal_record
		where user_id = ? and cal_times = last_result`
	if err = o.Raw(sql, uid).QueryRow(&maxCount); err != nil {
		return nil, -1, err
	}
	//分页查询，按utt降序,计算所有车辆的平均配载率
	//sql = `select cr.cal_no,
	//		max(cr.ctt) as ctt,
	//		max(cr.cal_times) as cal_times,
	//		avg(car.stowage_ratio) as stowage_ratio
	//	from cal_record as cr
	//		inner join car_summary as car
	//		on car.cal_record_id = cr.id and car.cal_times = cr.cal_times
	//	where cr.user_id = ? and cr.cal_times = cr.last_result
	//	group by cr.cal_no
	//	order by ctt desc
	//	limit ?
	//	offset ?`
	sql = `select cr.cal_no,
			cr.ctt,
			cr.cal_times,
			t2.money
		from cal_record as cr
			left join bill as t2
			on cr.order_id = t2.order_id
		where cr.user_id = ? and cr.cal_times = cr.last_result
		order by ctt desc
		limit ?
		offset ?`

	calRecords = []*CalHistory{}
	if _, err = o.Raw(sql, uid, pageLimit, pageLimit*(pageNumber-1)).
		QueryRows(&calRecords); err != nil {
		return nil, -1, err
	}
	return calRecords, maxCount, nil
}
