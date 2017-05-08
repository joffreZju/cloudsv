package service

import (
	"common/lib/util"
	"common/model"
	"encoding/json"
	"github.com/tealeg/xlsx"
	"strconv"
)

const WB_HEADER_1 = "货运单号"
const WB_HEADER_2 = "实际重量（公斤）"
const WB_HEADER_3 = "实际体积（立方）"
const WB_HEADER_4 = "运费金额"
const WB_HEADER_5 = "件数"
const WB_HEADER_6 = "是否必装"
const WB_HEADER_7 = "是否打底"
const WB_HEADER_8 = "可否拆单"
const WB_HEADER_9 = "所属车辆"

const WB_SHEET1 = "结果"
const WB_SHEET2 = "滞留清单"

const CAS_EXCEL_NAME = "计算结果"
const WB_TRUE = "是"
const WB_FALSE = "否"

//	header := []string{"货运单号", "实际重量（公斤", "实际体积（立方）", "运费金额", "件数", "是否必装", "是否打底", "可否拆单"}

func returnBoolCNResult(value string) (res string) {
	res = WB_FALSE
	if value == "true" {
		res = WB_TRUE
	}
	return
}

func GetWbResultExcil(waybills []*model.CalGoods, carSummays []*model.CarSummary) (data []byte, err error) {
	file := util.MkExcelFile()

	header := []string{WB_HEADER_1, WB_HEADER_2, WB_HEADER_3, WB_HEADER_4, WB_HEADER_5, WB_HEADER_6, WB_HEADER_7, WB_HEADER_8, WB_HEADER_9}

	otherinfoMapTemp := make(map[string]string, 0)
	err = json.Unmarshal([]byte(waybills[0].OtherInfo), &otherinfoMapTemp)
	if err != nil {
		return
	}
	otherinfoheader := []string{}
	for k, _ := range otherinfoMapTemp {
		otherinfoheader = append(otherinfoheader, k)
		header = append(header, k)
	}
	//结果清单
	sheetMap := make(map[string]*xlsx.Sheet, 0)

	for _, carSummay := range carSummays {
		sheet, err := util.MkExcelSheet(carSummay.CarNo, file)
		if err != nil {
			return data, err
		}
		err = util.FillExcelSheet(header, sheet, true, false)
		if err != nil {
			return data, err
		}
		sheetMap[carSummay.CarNo] = sheet
	}
	//滞留清单
	sheet_sunk, err := util.MkExcelSheet(WB_SHEET2, file)
	if err != nil {
		return
	}
	sheetMap["sunk"] = sheet_sunk
	err = util.FillExcelSheet(header, sheet_sunk, true, false)
	if err != nil {
		return data, err
	}
	//充数据
	mapdata := make(map[string][][]string, 0)
	for _, waybill := range waybills {
		otherinfoMap := make(map[string]string, 0)
		err = json.Unmarshal([]byte(waybill.OtherInfo), &otherinfoMap)
		if err != nil {
			return data, err
		}
		dataline := []string{
			waybill.WaybillNumber,
			strconv.FormatFloat(waybill.ActualWeight, 'f', -1, 32),
			strconv.FormatFloat(waybill.ActualVolume, 'f', -1, 32),
			strconv.FormatFloat(waybill.FreightCharges, 'f', -1, 32),
			strconv.FormatInt(int64(waybill.PackageNumber), 10),
			returnBoolCNResult(waybill.Necessary),
			returnBoolCNResult(waybill.Understowed),
			returnBoolCNResult(waybill.Split),
			waybill.CalResult,
		}
		for _, headerkey := range otherinfoheader {
			dataline = append(dataline, otherinfoMap[headerkey])
		}
		if waybill.CalResult == "" {
			mapdata["sunk"] = append(mapdata["sunk"], dataline)
		} else {
			mapdata[waybill.CalResult] = append(mapdata[waybill.CalResult], dataline)
		}
	}
	for k, md := range mapdata {
		for i, line := range md {
			if i%2 == 0 {
				err = util.FillExcelSheet(line, sheetMap[k], false, true)
			} else {
				err = util.FillExcelSheet(line, sheetMap[k], false, false)
			}
			if err != nil {
				return data, err
			}
		}
	}
	data, err = util.GetExcelData(file)
	return
}
