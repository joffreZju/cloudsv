package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"web/cons"
)

//producer And consumer use MqSign
func MqSigh(signStr, accessKey string) string {
	mac := hmac.New(sha1.New, []byte(accessKey))
	mac.Write([]byte(signStr))
	s := base64.StdEncoding.EncodeToString([]byte(mac.Sum(nil)))
	strings.TrimRight(s, " ")
	return s
}

//define MqMsg for MQ consumer to unMarshal json from MQ
type MqMsg struct {
	Body      string
	MsgHandle string
	MsgId     string
}

func Producer(bodyStr string) error {
	Topic := cons.MQ_TOPIC_PRODUCER
	URL := cons.MQ_URL
	AccessID := cons.ALI_ACCESS_KEY_ID
	AccessKEY := cons.ALI_ACCESS_KEY_SECRET
	ProducerID := cons.MQ_PRODUCER_ID
	newline := "\n"
	content := Md5Cal2String([]byte(bodyStr))
	date := fmt.Sprintf("%d", time.Now().UnixNano())[0:13]
	signStr := Topic + newline + ProducerID + newline + content + newline + date
	sign := MqSigh(signStr, AccessKEY)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL+"/message/?topic="+Topic+"&time="+date+"&tag=http"+"&key=http", strings.NewReader(bodyStr))
	if err != nil {
		return fmt.Errorf("MQ Producer error: %v", err)
	}

	req.Header.Set("Signature", sign)
	req.Header.Set("AccessKey", AccessID)
	req.Header.Set("ProducerID", ProducerID)
	req.Header.Set("Content-Type", "text/html;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("MQ Producer error: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Debug("read MQ response body error: ", err)
	}
	var respMsg MqMsg
	json.Unmarshal(body, &respMsg)
	beego.Debug("MQ producer status", respMsg.MsgId, resp.Status)

	if resp.StatusCode == 201 {
		return nil
	} else {
		return fmt.Errorf("MQ Producer error: %v", resp.Status)
	}
}

func Consumer(msgChan chan string, sig chan bool) {
	Topic := cons.MQ_TOPIC_CONSUMER
	URL := cons.MQ_URL
	AccessID := cons.ALI_ACCESS_KEY_ID
	AccessKEY := cons.ALI_ACCESS_KEY_SECRET
	ConsumerID := cons.MQ_CONSUMER_ID
	newline := "\n"
	receiveNumbers := "5"
	start := time.Now()
	timeLimit := 180.0 // seconds
	//如果一直没有计算结果那么 consumer 会陷入死循环,所以轮询三分钟
	for time.Now().Sub(start).Seconds() <= timeLimit {
		date := fmt.Sprintf("%d", time.Now().UnixNano())[0:13]
		signStr := Topic + newline + ConsumerID + newline + date
		sign := MqSigh(signStr, AccessKEY)

		client := &http.Client{}
		req, err := http.NewRequest("GET", URL+"/message/?topic="+Topic+"&time="+date+"&num="+receiveNumbers, nil)
		if err != nil {
			beego.Debug("MQ consumer get error: ", err)
		}
		req.Header.Set("Signature", sign)
		req.Header.Set("AccessKey", AccessID)
		req.Header.Set("ConsumerID", ConsumerID)
		resp, err := client.Do(req)

		if err != nil {
			beego.Debug("MQ consumer get error: ", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			beego.Debug("MQ consumer get error: ", err)
		}
		var msgs []MqMsg
		json.Unmarshal(body, &msgs)
		if len(msgs) == 0 {
			continue
		} else {
			for _, msg := range msgs {
				msgChan <- msg.Body
				beego.Debug("get MQ calculate result message and send to chan:", msg.MsgId)
			}
			close(msgChan)
		}
		//只有数据库写入结果没有出错的情况下，才会删除消息，如果出错那么不删除消息，继续轮询接收计算结果。
		//如果删除消息失败，那么继续轮询接受消息，如果再次受到相同的消息，再次update数据库中，结果不变。
		//未删除的消息扰乱了消息时序，因此不管数据库是否成功，都删除消息
		//dbDone := <-sig
		//if !dbDone {
		//	continue
		//}
		for _, msg := range msgs {
			date := fmt.Sprintf("%d", time.Now().UnixNano())[0:13]
			delUrl := URL + "/message/?msgHandle=" + msg.MsgHandle + "&topic=" + Topic + "&time=" + date
			signStr := Topic + newline + ConsumerID + newline + msg.MsgHandle + newline + date
			sign := MqSigh(signStr, AccessKEY)
			req, err := http.NewRequest(http.MethodDelete, delUrl, nil)
			if err != nil {
				beego.Debug(msg.MsgId, "MQ consumer delete error: ", err)
				continue
			}
			req.Header.Set("Signature", sign)
			req.Header.Set("AccessKey", AccessID)
			req.Header.Set("ConsumerID", ConsumerID)

			resp, err := client.Do(req)
			if err != nil {
				beego.Debug(msg.MsgId, "MQ consumer delete error: ", resp.Status, err)
			} else {
				beego.Debug(msg.MsgId, "MQ consumer delete success: ", resp.Status, err)
			}
		}
		sig <- true
		break
	}
}
