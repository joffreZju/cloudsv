package service

import (
	"common/lib/util"
	"common/model"
	"common/service/mqdto"
	"encoding/json"
)

func SendCalToMq(cars []*model.CarSummary, goods []*model.CalGoods, record *model.CalRecord) (err error) {
	mqCars := make([]mqdto.MqCarInfo, len(cars))
	mqGoods := make([]mqdto.MQWaybill, len(goods))
	for _, v := range cars {
		mqCars = append(mqCars, mqdto.MqCarInfo{
			Car_no: v.CarNo,
			Cubage: v.MaxVolume,
			Load:   v.MaxWeight,
		})
	}
	for _, v := range goods {
		if v.SplitInfo == "split_from" {
			continue
		}
		mqGoods = append(mqGoods, mqdto.MQWaybill{
			Id:  v.Id,
			Av:  v.ActualVolume,
			Aw:  v.ActualWeight,
			Fc:  v.FreightCharges,
			Ne:  v.Necessary,
			Uns: v.Understowed,
		})
	}
	mqReq := mqdto.ReqMQDto{
		//todo 回调地址和calType要修改
		Callback:   "106.14.38.116:8090",
		Cal_type:   record.CalType,
		Using_id:   record.Id,
		Cal_times:  record.CalTimes,
		Car_info:   mqCars,
		Goods_list: mqGoods,
	}
	var b []byte
	b, err = json.Marshal(mqReq)
	if err != nil {
		return
	}
	for i := 0; i < 3; i++ {
		if err = util.Producer(string(b)); err == nil {
			break
		}
	}
	return
}
