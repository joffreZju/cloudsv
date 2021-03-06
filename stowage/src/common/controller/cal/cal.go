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

const CommonErr = 99999

func (c *Controller) GetTplAndFrequentCars() {
	uid := int(c.UserID)
	oneTpl, e := model.GetTemplate(uid)
	if e != nil {
		c.ReplyErr(errcode.ErrNoTpl)
		beego.Error(e)
		return
	}
	cars, e := model.GetFrequentCars(uid)
	if e != nil {
		c.ReplyErr(errcode.ErrNoFrequentCars)
		beego.Error(e)
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
		beego.Error(e)
		return
	}
	tpl.UserId = uid
	tpl.Ctt = time.Now()
	e = model.InsertTemplate(tpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
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
		c.ReplyErr(errcode.ErrWrongJson)
		beego.Error(e)
		return
	}
	e = json.Unmarshal([]byte(goodsString), &goods)
	if e != nil {
		c.ReplyErr(errcode.ErrWrongJson)
		beego.Error(e)
		return
	}
	cr.CalType, e = c.checkCarsAndGoods(cars, goods)
	if e != nil {
		c.ReplyErr(errcode.ErrWrongCarsGoods)
		beego.Error(e)
		return
	}
	splitGoods := c.splitWaybill(goods)
	//插入计算数据并发起计算请求，建立账单
	e = service.InsertCalToDbAndSendToMq(uid, cars, splitGoods, cr)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"CalNo": cr.CalNo,
	})
}

func (c *Controller) GetCalResult() {
	calNo := c.GetString("CalNo")
	uid := int(c.UserID)
	cr, e := model.GetCalRecord(calNo)
	if e != nil {
		c.ReplyErr(errcode.ErrWrongCalNo)
		beego.Error(e)
		return
	} else if cr.PayStatus == model.YiFailed {
		c.ReplyErr(errcode.ErrCalPayFailed)
		return
	} else if cr.PayStatus == model.YiOrderCreate {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		return
	} else if cr.UserId != uid {
		c.ReplyErr(errcode.ErrCalNoUserNoMatch)
		return
	}
	cars, e := service.GetCarsResult(calNo)
	if e != nil || len(cars) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		beego.Error(e)
		return
	}
	goods, e := service.GetGoodsResult(calNo)
	if e != nil || len(goods) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		beego.Error(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"CarsResult":  cars,
		"GoodsResult": goods,
	})
}

func (c *Controller) GetCalResultExcel() {
	calNo := c.GetString("CalNo")
	uid := int(c.UserID)
	cr, e := model.GetCalRecord(calNo)
	if e != nil {
		c.ReplyErr(errcode.ErrWrongCalNo)
		beego.Error(e)
		return
	} else if cr.PayStatus == model.YiFailed {
		c.ReplyErr(errcode.ErrCalPayFailed)
		return
	} else if cr.PayStatus == model.YiOrderCreate {
		c.ReplyErr(errcode.ErrCalResultIsNull)
	} else if cr.UserId != uid {
		c.ReplyErr(errcode.ErrCalNoUserNoMatch)
		return
	}
	cars, e := service.GetCarsResult(calNo)
	if e != nil || len(cars) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		beego.Error(e)
		return
	}
	goods, e := service.GetGoodsResult(calNo)
	if e != nil || len(goods) == 0 {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		beego.Error(e)
		return
	}
	databyte, e := service.GetWbResultExcil(goods, cars)
	if e != nil {
		c.ReplyErr(errcode.ErrCalResultIsNull)
		beego.Error(e)
		return
	}
	//c.replyFileStream(calNo+".xlsx", databyte)
	c.ReplyFile("application/octet-stream", calNo+".xlsx", databyte)

}

//func (c *Controller) replyFileStream(name string, data []byte) {
//	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+name)
//	c.Ctx.Output.Header("Content-Description", "File Transfer")
//	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
//	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")
//	c.Ctx.Output.Header("Expires", "0")
//	c.Ctx.Output.Header("Cache-Control", "must-revalidate")
//	c.Ctx.Output.Header("Pragma", "public")
//	c.Ctx.Output.Body(data)
//}

func (c *Controller) GetEditedWbs() {
	calNo := c.GetString("CalNo")
	uid := int(c.UserID)
	cr, e := model.GetCalRecord(calNo)
	if e != nil {
		c.ReplyErr(errcode.ErrWrongCalNo)
		beego.Error(e)
		return
	} else if cr.UserId != uid {
		c.ReplyErr(errcode.ErrCalNoUserNoMatch)
		return
	}
	cars, e := service.GetEditedCars(calNo)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	goods, e := service.GetEditedWaybills(calNo)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
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
	calHistory, maxCount, e := service.GetCalHistory(uid, pageNumber, pageLimit)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(map[string]interface{}{
		"MaxCount": maxCount,
		"Result":   calHistory,
	})
}
