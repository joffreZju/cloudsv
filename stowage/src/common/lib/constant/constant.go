package constant

import (
	"time"
)

//define for login wrong password expire time
const (
	USER_LOGIN_WRONG_PWD_EXPIRE_TIME = time.Minute * 30
	NUMBER_OF_HISTORY_CARS           = 5
)

//define waybills 打底，必装(true,false)
const (
	STRING_TRUE           = "true"
	STRING_FALSE          = "" //false
	WAYBILL_SPLIT_FROM    = "split_from"
	WAYBILL_SPLIT_TO      = "split_to"
	ORDER_CAL_TYPE_MONEY  = "moneyOpt"
	ORDER_CAL_TYPE_LOAD   = "fullLoad"
	WAYBILL_SPLIT_Vmax    = 0.5    // 立方米
	WAYBILL_SPLIT_Wmax    = 100.0  // kg
	WAYBILL_SPLIT_V_div_W = 0.0045 // 立方米/kg
)

//define user_role
const (
	USER_ROLE_NORMAL = 1
	USER_ROLE_AGENT  = 2
)

//define user status
const (
	USER_STATUS_ACCESS = 0
	USER_STATUS_DENY   = 1
)

//define user_code prefix
const (
	USER_CODE_PREFIX_NORMAL = "SPZ"
	USER_CODE_PREFIX_AGENT  = "AGENT"
	USER_CODE_PREFIX_API    = "API"
)

//前段加密 Key and nonce 期限
const (
	PASSWORD_APPEND              = "suanpeizai"
	SECRET_REQUEST_SIGN_KEY      = "suanpeizai"
	NONCE_EXPIRE_TIME            = 60 * 5       // seconds
	NONCE_REPEAT_CAL_EXPIRE_TIME = 60 * 60 * 48 // seconds
)

//define for sms
//define sms verify type
const (
	SEND_SMS_CODE_EXISTED_MOBILE     = "existed_mobile"
	SEND_SMS_CODE_NOT_EXISTED_MOBILE = "not_existed_mobile"
)

//define sms template code
const (
	SMS_CODE_EXPIRE_TIME = time.Minute * 10
	SMS_SIGN_NAME        = "壹算科技"
	SMS_TEMPLATE_WEBV2   = "SMS_58265055"
	SMS_TEMPLATE_TO_CRR  = "SMS_60060604"
)

//define for token
const (
	TOKEN_ISSUER      = "allsum"
	TOKEN_EXPIRE_TIME = 60 * 60 * 2 //seconds
	TOKEN_SIGN_KEY    = "www.suanpeizai.com"
)
