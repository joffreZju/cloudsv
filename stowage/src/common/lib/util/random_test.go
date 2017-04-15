package util

import (
	"testing"
)

func TestGetTradeNo(t *testing.T) {
	var tp = 1
	var uid = 123456
	v := GetTradeNo(tp, uid)
	t.Log(v)
}
