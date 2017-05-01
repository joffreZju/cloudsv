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
	t, e := model.GetTemplate(uid)
	if e != nil {
		c.ReplyErr(e)
	}
	cars, e := model.GetFrequentCars(uid)
	if e != nil {
		c.ReplyErr(e)
	}
	c.ReplySucc(map[string]interface{}{
		"Cars": cars,
		"Tpl":  t,
	})
}

func (c *Controller) StoreTpl() {
	uid := int(c.UserID)
	s := c.GetString("Tpl")
	tpls := []*model.CalTemplate{}
	e := json.Unmarshal([]byte(s), tpls)
	if len(tpls) <= 0 || e != nil {
		c.ReplyErr(errcode.ErrTplIsNull)
	}
	tpls[0].UserId = uid
	tpls[0].Ctt = time.Now()
	e = model.InsertTemplate(tpls[0])
	if e != nil {
		c.ReplyErr(e)
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
	e := json.Unmarshal([]byte(carString), cars)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
	}
	e = json.Unmarshal([]byte(goodsString), goods)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
	}
	cr.CalType, e = c.checkCarsAndGoods(cars, goods)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
	}
	splitGoods := c.splitWaybill(goods)
	e = service.InsertCalToDbAndSendToMq(uid, cars, splitGoods, cr)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(e)
	}
	c.ReplySucc("sucess")
	//todo 建立账单，收到计算结果之后扣费
}

func (c *Controller) GetCalResult() {
	calNo := c.GetString("CalNo")
	cars, e := service.GetCarsResult(calNo)
	if e != nil {
		c.ReplyErr(errcode.ErrCalResultIsNull)
	}
	goods, e := service.GetGoodsResult(calNo)
	if e != nil {
		c.ReplyErr(errcode.ErrCalResultIsNull)
	}
	c.ReplySucc(map[string]interface{}{
		"CarsResult":  cars,
		"GoodsResult": goods,
		//"CalNo": calNo,
	})
}

func (c *Controller) GetEditedWbs() {
	calNo := c.GetString("CalNo")
	cars, e := service.GetEditedCars(calNo)
	if e != nil {
		c.ReplyErr(e)
	}
	goods, e := service.GetEditedWaybills(calNo)
	if e != nil {
		c.ReplyErr(e)
	}
	c.ReplySucc(map[string]interface{}{
		"Cars":      cars,
		"GoodsList": goods,
		//"CalNo": calNo,
	})
}

func (c *Controller) GetCalHistory() {
	uid := int(c.UserID)
	pageNumber, _ := c.GetInt("PageNumber")
	pageLimit, _ := c.GetInt("PageLimit")
	calRecords, maxCount, e := service.GetCalHistory(uid, pageNumber, pageLimit)
	if e != nil {
		c.ReplyErr(e)
	}
	c.ReplySucc(map[string]interface{}{
		"MaxCount": maxCount,
		"Result":   calRecords,
	})
}
