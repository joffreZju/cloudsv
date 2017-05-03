package cal

import (
	"common/model"
	"errors"
	mapset "github.com/deckarep/golang-set"
)

//检测车辆信息和运单信息是否正确，检测打底车辆是否存在
func (c *Controller) checkCarsAndGoods(cars []*model.CarSummary, goods []*model.CalGoods) (string, error) {
	calType := model.ORDER_CAL_TYPE_LOAD
	if len(cars) == 0 || len(goods) == 0 {
		return calType, errors.New("车辆、货物字段不能为空")
	}
	//检测车辆的基本信息是否正确，车辆是否重复
	carSet := mapset.NewSet()
	for _, v := range cars {
		if v.MaxVolume <= 0 || v.MaxWeight <= 0 {
			return calType, errors.New("车辆体积重量不能小于等于0")
		}
		carSet.Add(v.CarNo)
	}
	if carSet.Cardinality() != len(cars) {
		return calType, errors.New("车辆不能重复")
	}
	//检测运单信息是否正确，运单是否重复
	waybillSet := mapset.NewSet()
	for _, v := range goods {
		if v.FreightCharges > 0 {
			calType = model.ORDER_CAL_TYPE_MONEY
		}
		if v.ActualVolume < 0 || v.ActualWeight < 0 || v.FreightCharges < 0 {
			return calType, errors.New("运单重量体积金额不能小于0")
		}
		if v.Understowed != "" && !carSet.Contains(v.Understowed) {
			return calType, errors.New("运单打底车辆不正确")
		}
		waybillSet.Add(v.WaybillNumber)
	}
	if waybillSet.Cardinality() != len(goods) {
		return calType, errors.New("运单号不能重复")
	}
	return calType, nil
}

//拆单
func (c *Controller) splitWaybill(goods []*model.CalGoods) (wbs []*model.CalGoods) {
	primes := []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41,
		43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127,
		131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199}
	Vmax := model.WAYBILL_SPLIT_Vmax
	Wmax := model.WAYBILL_SPLIT_Wmax
	V_div_W := model.WAYBILL_SPLIT_V_div_W
	split_from := model.WAYBILL_SPLIT_FROM
	split_to := model.WAYBILL_SPLIT_TO
	for _, v := range goods {
		//不需要拆单
		if v.Split != model.STRING_TRUE || v.PackageNumber <= 1 ||
			(v.ActualVolume <= Vmax && v.ActualWeight <= Wmax) {
			wbs = append(wbs, v)
			continue
		}
		//原运单保留，但是 SplitInfo = "split_from",不写入计算结果
		tmpGood := v
		tmpGood.SplitInfo = split_from
		wbs = append(wbs, tmpGood)
		//计算重抛比较大，按照体积拆单
		if (float64(v.ActualVolume))/(float64(v.ActualWeight)) > V_div_W {
			//按照cutV拆分
			primeFlag := 1.0
			cutV := v.ActualVolume / float64(v.PackageNumber)
			for i := 0; i < len(primes) && cutV < Vmax; i++ {
				cutV = primes[i] * cutV / primeFlag
				primeFlag = primes[i]
			}
			cutW := primeFlag * v.ActualWeight / float64(v.PackageNumber)
			cutCharges := primeFlag * v.FreightCharges / float64(v.PackageNumber)
			//将原运单拆分，拆分后的运单 SplitInfo = "split_to"
			for v.ActualVolume-cutV > Vmax {
				wbs = append(wbs, &model.CalGoods{
					WaybillNumber:  v.WaybillNumber,
					ActualWeight:   cutW,
					ActualVolume:   cutV,
					FreightCharges: cutCharges,
					PackageNumber:  int(primeFlag),
					Necessary:      v.Necessary,
					Understowed:    v.Understowed,
					OtherInfo:      v.OtherInfo,
					Split:          v.Split,
					SplitInfo:      split_to,
				})
				v.ActualVolume -= cutV
				v.ActualWeight -= cutW
				v.FreightCharges -= cutCharges
				v.PackageNumber -= int(primeFlag)
			}
			if v.ActualVolume > 0 {
				wbs = append(wbs, &model.CalGoods{
					WaybillNumber:  v.WaybillNumber,
					ActualWeight:   v.ActualWeight,
					ActualVolume:   v.ActualVolume,
					FreightCharges: v.FreightCharges,
					PackageNumber:  v.PackageNumber,
					Necessary:      v.Necessary,
					Understowed:    v.Understowed,
					OtherInfo:      v.OtherInfo,
					Split:          v.Split,
					SplitInfo:      split_to,
				})
			}
		} else {
			//重抛比较小，按照重量拆分，按照cutW拆分
			primeFlag := 1.0
			cutW := v.ActualWeight / float64(v.PackageNumber)
			for i := 0; i < len(primes) && cutW < Wmax; i++ {
				cutW = primes[i] * cutW / primeFlag
				primeFlag = primes[i]
			}
			cutV := primeFlag * v.ActualVolume / float64(v.PackageNumber)
			cutCharges := primeFlag * v.FreightCharges / float64(v.PackageNumber)
			//将原运单拆分，拆分后的运单 SplitInfo = split_to
			for v.ActualWeight-cutW > Vmax {
				wbs = append(wbs, &model.CalGoods{
					WaybillNumber:  v.WaybillNumber,
					ActualWeight:   cutW,
					ActualVolume:   cutV,
					FreightCharges: cutCharges,
					PackageNumber:  int(primeFlag),
					Necessary:      v.Necessary,
					Understowed:    v.Understowed,
					OtherInfo:      v.OtherInfo,
					Split:          v.Split,
					SplitInfo:      split_to,
				})
				v.ActualVolume -= cutV
				v.ActualWeight -= cutW
				v.FreightCharges -= cutCharges
				v.PackageNumber -= int(primeFlag)
			}
			if v.ActualWeight > 0 {
				wbs = append(wbs, &model.CalGoods{
					WaybillNumber:  v.WaybillNumber,
					ActualWeight:   v.ActualWeight,
					ActualVolume:   v.ActualVolume,
					FreightCharges: v.FreightCharges,
					PackageNumber:  v.PackageNumber,
					Necessary:      v.Necessary,
					Understowed:    v.Understowed,
					OtherInfo:      v.OtherInfo,
					Split:          v.Split,
					SplitInfo:      split_to,
				})
			}
		}
	}
	return wbs
}
