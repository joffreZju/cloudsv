package cal

import (
	"common/service"
	"common/service/mqdto"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
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
	// todo wjf 这里修改的话计算引擎也要改，这里要不要签名,回调结果要不要签名等等。。
	start := time.Now()
	key := fmt.Sprintf("%x", sha256.Sum256([]byte("allsum_suanpeizai2.0")))
	if key != c.Ctx.Request.Header.Get("key") {
		c.Data["json"] = fail
		c.ServeJSON()
		beego.Debug("invalid key")
		c.StopRun()
	}
	calResult := mqdto.MQRespDto{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &calResult); err != nil {
		c.Data["json"] = fail
		c.ServeJSON()
		beego.Debug(err)
		c.StopRun()
	}
	c.Data["json"] = success
	c.ServeJSON()
	//if calResult.Error_code != 0 {
	//	notice := fmt.Sprintf("msgId：%s，usingId：%d，calTimes：%d，",
	//		calResult.MqMsg_id, calResult.Using_id, calResult.Cal_times)
	//	utils.SendSmsWhenError("13735544671", notice+time.Now().Format("2006-01-02 15:04:05"))
	//	c.StopRun()
	//}

	var err error
	for i := 0; i < 5; i++ {
		if err = service.UpdateCalResult(&calResult); err == nil {
			beego.Debug("计算结果写入成功,using_id:", calResult.Using_id, "耗时:", time.Now().Sub(start))
			break
		}
	}
	if err != nil {
		beego.Debug("计算结果写入失败,using_id:", calResult.Using_id, err)
	}
}
