package service

import (
	"common/model"

	"common/lib/errcode"
	"common/service/mqdto"
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

func InsertCalInfo(cars []*model.CarSummary, goods []*model.CalGoods, record *model.CalRecord) (err error) {
	o := orm.NewOrm()
	o.Begin()
	err = model.InsertOrUpdateRec(o, record)
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return
	}
	for _, v := range cars {
		v.CalRecordId = record.Id
		v.CalTimes = record.CalTimes
	}
	for _, v := range goods {
		v.CalRecordId = record.Id
		v.CalTimes = record.CalTimes
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
	}
	return
}

func UpdateCalResult(cr *mqdto.MQRespDto) (err error) {
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
	if _, err = o.Raw(sql, cr.Cal_times, time.Now(), cr.Using_id, cr.Cal_times).Exec(); err != nil {
		o.Rollback()
		beego.Error("update cal result error with spz_cal_record", cr.Using_id, err)
		return
	}
	//更新spz_cal_carsummary
	sql = `update car_summary
		set total_money =?,
			total_weight = ?,
			total_volume= ?,
			utt = ?,
			stowage_ratio = (?/max_weight + ?/max_volume)/2
		where using_id = ? and cal_times = ? and car_no = ?`
	for _, v := range cr.Car_summary {
		_, err = o.Raw(sql, v.Total_money, v.Total_weight, v.Total_volume, time.Now(),
			v.Total_weight, v.Total_volume,
			cr.Using_id, cr.Cal_times, v.Car_no).Exec()
		if err != nil {
			o.Rollback()
			beego.Error("update cal result error with spz_cal_carsummary", cr.Using_id, err)
			return
		}
	}
	//更新spz_cal_lib
	sql = `update cal_goods
		set cal_result = ?,
			utt  = ?
		where id = ?`
	for _, v := range cr.Result {
		if _, err = o.Raw(sql, v.Result, time.Now(), v.Id).Exec(); err != nil {
			o.Rollback()
			beego.Error("update cal result error with spz_cal_lib", cr.Using_id, err)
			return
		}
	}
	if err = o.Commit(); err != nil {
		beego.Critical("update calculate result error with database transaction commit:", cr.Using_id, err)
		return
	}
	beego.Debug("update calculate result success, using_id:", cr.Using_id)
	return nil
}
