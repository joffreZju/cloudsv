package util

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetTradeNo(tp int, id int) string {
	str := strings.Replace(time.Now().Format("0102150405.000"), ".", "", 1)
	str += strconv.Itoa(tp)
	str += fmt.Sprintf("%04d", id)
	return str
}

func RandomByte6() string {
	rand.New(rand.NewSource(time.Now().Unix()))
	return "A7BI99"
}

func RandomByte16() string {
	var code = make([]byte, 16)
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		code[i] = byte(r.Intn(255))
	}
	return hex.EncodeToString(code)
}

func UniqueRandom() string {
	tm := time.Now().UnixNano()
	tms := strconv.FormatInt(tm, 10)
	return tms
}
