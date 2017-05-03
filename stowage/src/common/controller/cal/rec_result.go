package cal

import (
	"common/lib/push"
	"common/model"
	"common/service"
	"common/service/mqdto"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RecController struct {
	beego.Controller
}
type MqResp struct {
	ErrorCode int
	ErrorMsg  string
}

var fail = MqResp{
	ErrorCode: 1,
	ErrorMsg:  "fail",
}
var success = MqResp{
	ErrorCode: 0,
	ErrorMsg:  "success",
}

func (c *RecController) HandleCalResult() {
	start := time.Now()
	key := fmt.Sprintf("%x", sha256.Sum256([]byte("allsum_suanpeizai2.0")))
	if key != c.Ctx.Request.Header.Get("key") {
		c.Data["json"] = fail
		c.ServeJSON()
		beego.Error("invalid key")
		return
	}
	calResult := mqdto.MQRespDto{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &calResult); err != nil {
		c.Data["json"] = fail
		c.ServeJSON()
		beego.Error(err)
		return
	}
	c.Data["json"] = success
	c.ServeJSON()

	if calResult.Error_code != 0 {
		notice := fmt.Sprintf("msgId：%s，usingId：%d，calTimes：%d，",
			calResult.MqMsg_id, calResult.Using_id, calResult.Cal_times)
		push.SendErrorSms("13735544671", "计算引擎故障,"+notice+time.Now().Format("2006-01-02 15:04:05"))
		return
	}

	var err error
	for i := 0; i < 5; i++ {
		if err = service.UpdateCalResult(&calResult); err == nil {
			beego.Error("计算结果写入成功,using_id:", calResult.Using_id, "耗时:", time.Now().Sub(start))
			break
		}
	}
	if err != nil {
		beego.Error("计算结果写入失败,using_id:", calResult.Using_id, err)
		return
	}
	// 写入结果成功，进行扣费，并更新record status
	cr, err := model.GetCalRecordById(calResult.Using_id)
	if err != nil {
		beego.Error(err)
		return
	}
	//扣费
	err = model.TransFinance(cr.OrderId)
	if err != nil {
		beego.Error(err)
		cr.PayStatus = model.YiFailed
	} else {
		cr.PayStatus = model.YiPaid
	}
	if e := model.UpdateCalRecord(orm.NewOrm(), cr); e != nil {
		beego.Error(e)
		return
	}
}
