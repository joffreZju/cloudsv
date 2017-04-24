package mqdto

//请求MQ计算
type ReqMQDto struct {
	Callback   string
	Cal_type   string
	Using_id   int
	Cal_times  int
	Car_info   []MqCarInfo
	Goods_list []MQWaybill
}
type MqCarInfo struct {
	Car_no string
	Cubage float64
	Load   float64
}

//MQ通信DTO，简化json，对应Goods，用于向MQ发送消息，缩小json体积
type MQWaybill struct {
	Id  int
	Aw  float64
	Av  float64
	Fc  float64
	Ne  string
	Uns string
}

//接收MQ计算结果Dto
type MQRespDto struct {
	MqMsg_id    string
	Error_code  int
	Using_id    int
	Cal_times   int
	Car_summary []MQCarSummary
	Result      []MQCalResult
}
type MQCarSummary struct {
	Car_no       string
	Total_weight float64
	Total_volume float64
	Total_money  float64
}
type MQCalResult struct {
	Id     int
	Result string
}
