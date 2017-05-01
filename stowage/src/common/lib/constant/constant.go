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

//define system all kinds of error
//var (
//	SUCCESS_REQUEST = dto.ResponseError{
//		ErrorCode: 0,
//		ErrorMsg:  "请求成功",
//	}
//	ERR_SERVER = dto.ResponseError{
//		ErrorCode: 40000,
//		ErrorMsg:  "服务器内部错误",
//	}
//	ERR_DATABASE = dto.ResponseError{
//		ErrorCode: 40001,
//		ErrorMsg:  "数据库错误",
//	}
//	ERR_SIGN = dto.ResponseError{
//		ErrorCode: 10001,
//		ErrorMsg:  "签名错误",
//	}
//	ERR_NONCE = dto.ResponseError{
//		ErrorCode: 10002,
//		ErrorMsg:  "Nonce超时",
//	}
//	ERR_USER_NOT_SIGNUP = dto.ResponseError{
//		ErrorCode: 10003,
//		ErrorMsg:  "用户尚未注册",
//	}
//	ERR_USER_EXISTED = dto.ResponseError{
//		ErrorCode: 10008,
//		ErrorMsg:  "用户已经存在",
//	}
//	ERR_SEND_CODE_FAIL = dto.ResponseError{
//		ErrorCode: 10004,
//		ErrorMsg:  "发送验证码失败，请重新尝试",
//	}
//	ERR_SMS_CODE = dto.ResponseError{
//		ErrorCode: 10005,
//		ErrorMsg:  "验证码错误",
//	}
//	ERR_AGENT_CODE = dto.ResponseError{
//		ErrorCode: 10006,
//		ErrorMsg:  "推荐码错误",
//	}
//	ERR_TOKEN_GENERATE = dto.ResponseError{
//		ErrorCode: 10007,
//		ErrorMsg:  "生成token错误",
//	}
//	ERR_TOKEN_PARSE = dto.ResponseError{
//		ErrorCode: 10009,
//		ErrorMsg:  "解析token错误，请重新登录",
//	}
//	ERR_PASSWORD_WRONG = dto.ResponseError{
//		ErrorCode: 10010,
//		ErrorMsg:  "密码错误",
//	}
//	ERR_PASSWORD_WRONG_THREE = dto.ResponseError{
//		ErrorCode: 10011,
//		ErrorMsg:  "密码错误超过三次，密码登录方式被锁定三十分钟",
//	}
//	ERR_UPDATE_PASSWORD = dto.ResponseError{
//		ErrorCode: 10012,
//		ErrorMsg:  "密码更新失败",
//	}
//	ERR_PASSWORD_LENGTH = dto.ResponseError{
//		ErrorCode: 10013,
//		ErrorMsg:  "密码长度错误(32)",
//	}
//	ERR_TOKEN = dto.ResponseError{
//		ErrorCode: 10014,
//		ErrorMsg:  "token失效",
//	}
//	ERR_REQUEST_BODY = dto.ResponseError{
//		ErrorCode: 10015,
//		ErrorMsg:  "请求数据格式错误",
//	}
//	ERR_ACCOUNT_DENIED = dto.ResponseError{
//		ErrorCode: 10016,
//		ErrorMsg:  "账户余额不足或被禁用",
//	}
//	ERR_ORDERS_INFO = dto.ResponseError{
//		ErrorCode: 10017,
//		ErrorMsg:  "订单信息存在错误",
//	}
//	ERR_CAL_FAILED = dto.ResponseError{
//		ErrorCode: 10018,
//		ErrorMsg:  "计算失败，请重新请求计算",
//	}
//	ERR_GET_CAL_RESULT = dto.ResponseError{
//		ErrorCode: 10019,
//		ErrorMsg:  "获取计算结果失败，请重新请求",
//	}
//	ERR_GET_EDITED_WBS = dto.ResponseError{
//		ErrorCode: 10020,
//		ErrorMsg:  "获取运单数据失败，请重试",
//	}
//	ERR_GET_CAL_HISTORY = dto.ResponseError{
//		ErrorCode: 10021,
//		ErrorMsg:  "获取历史记录失败，请重试",
//	}
//	ERR_CAL_RESULT_NOT_EXISTED = dto.ResponseError{
//		ErrorCode: 10022,
//		ErrorMsg:  "计算结果还未返回，请耐心等待",
//	}
//	ERR_DATA_TO_EXCEL = dto.ResponseError{
//		ErrorCode: 10023,
//		ErrorMsg:  "转excel文件失败",
//	}
//)
