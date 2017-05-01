package constant

import (
	"github.com/astaxie/beego"
)

var ALI_ACCESS_KEY_ID string
var ALI_ACCESS_KEY_SECRET string
var MQ_URL string
var MQ_TOPIC_PRODUCER string
var MQ_TOPIC_CONSUMER string
var MQ_PRODUCER_ID string
var MQ_CONSUMER_ID string

func init() {
	ALI_ACCESS_KEY_ID = beego.AppConfig.String("ALI_ACCESS_KEY_ID")
	ALI_ACCESS_KEY_SECRET = beego.AppConfig.String("ALI_ACCESS_KEY_SECRET")
	MQ_URL = beego.AppConfig.String("MQ_URL")
	MQ_TOPIC_PRODUCER = beego.AppConfig.String("MQ_TOPIC_PRODUCER")
	MQ_TOPIC_CONSUMER = beego.AppConfig.String("MQ_TOPIC_CONSUMER")
	MQ_PRODUCER_ID = beego.AppConfig.String("MQ_PRODUCER_ID")
	MQ_CONSUMER_ID = beego.AppConfig.String("MQ_CONSUMER_ID")
}
