package constant

import (
	"github.com/astaxie/beego"
)

var ALI_ACCESS_KEY_ID string = "LTAIwxFn7egYfvra"
var ALI_ACCESS_KEY_SECRET string = "nBfpqo4StRZv9JreRsLQpFaZKKUT1h"
var MQ_URL string
var MQ_TOPIC_PRODUCER string
var MQ_TOPIC_CONSUMER string
var MQ_PRODUCER_ID string
var MQ_CONSUMER_ID string
var CAL_CALLBACK_URL string //计算引擎回调地址

func Init() error {
	MQ_URL = beego.AppConfig.String("MQ_URL")
	MQ_TOPIC_PRODUCER = beego.AppConfig.String("MQ_TOPIC_PRODUCER")
	MQ_TOPIC_CONSUMER = beego.AppConfig.String("MQ_TOPIC_CONSUMER")
	MQ_PRODUCER_ID = beego.AppConfig.String("MQ_PRODUCER_ID")
	MQ_CONSUMER_ID = beego.AppConfig.String("MQ_CONSUMER_ID")
	CAL_CALLBACK_URL = beego.AppConfig.String("cal_callback")

	return nil
}
