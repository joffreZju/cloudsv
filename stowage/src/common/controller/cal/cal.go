package cal

import (
	"common/controller/base"
	"common/lib/errcode"
	"common/model"
	"common/service"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type Controller struct {
	base.Controller
}

func (c *Controller) GetTplAndFrequentCars() {
	uid := int(c.UserID)
	oneTpl, e := model.GetTemplate(uid)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	cars, e := model.GetFrequentCars(uid)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"Cars": cars,
		"Tpl":  oneTpl,
	})
}

func (c *Controller) StoreTpl() {
	uid := int(c.UserID)
	s := c.GetString("Tpl")
	tpl := &model.CalTemplate{}
	e := json.Unmarshal([]byte(s), tpl)
	if e != nil || len(tpl.WaybillNumber) == 0 {
		c.ReplyErr(errcode.ErrTplIsNull)
		return
	}
	tpl.UserId = uid
	tpl.Ctt = time.Now()
	e = model.InsertTemplate(tpl)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	c.ReplySucc("success")
}

func (c *Controller) Calculate() {
	uid := int(c.UserID)
	calNo := c.GetString("CalNo")
	carString := c.GetString("Cars")
	goodsString := c.GetString("GoodsList")
	cars := []*model.CarSummary{}
	goods := []*model.CalGoods{}
	cr := &model.CalRecord{
		CalNo: calNo,
	}
	e := json.Unmarshal([]byte(carString), &cars)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
		return
	}
	e = json.Unmarshal([]byte(goodsString), &goods)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
		return
	}
	cr.CalType, e = c.checkCarsAndGoods(cars, goods)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
		return
	}
	splitGoods := c.splitWaybill(goods)
	//插入计算数据并发起计算请求，建立账单
	e = service.InsertCalToDbAndSendToMq(uid, cars, splitGoods, cr)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"CalNo": cr.CalNo,
	})
}

func (c *Controller) GetCalResult() {
	calNo := c.GetString("CalNo")
	cars, e := service.GetCarsResult(calNo)
	if e != nil || len(cars) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		return
	}
	goods, e := service.GetGoodsResult(calNo)
	if e != nil || len(goods) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"CarsResult":  cars,
		"GoodsResult": goods,
	})
}

func (c *Controller) GetEditedWbs() {
	calNo := c.GetString("CalNo")
	cars, e := service.GetEditedCars(calNo)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	goods, e := service.GetEditedWaybills(calNo)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"Cars":      cars,
		"GoodsList": goods,
	})
}

func (c *Controller) GetCalHistory() {
	uid := int(c.UserID)
	pageNumber, _ := c.GetInt("PageNumber")
	pageLimit, _ := c.GetInt("PageLimit")
	calRecords, maxCount, e := service.GetCalHistory(uid, pageNumber, pageLimit)
	if e != nil {
		c.ReplyErr(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"MaxCount": maxCount,
		"Result":   calRecords,
	})
}
