package util

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func RandomByte6() string {
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
