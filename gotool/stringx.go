package gotool

import (
	"math/rand"
	"time"
)

const (
	UppercaseCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowercaseCharset = "abcdefghijklmnopqrstuvwxyz"
	NumericCharset   = "0123456789"
	SpecialCharset   = "!@#$%^&*()_+-="
)

// RandString 随机字符串
func RandString(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}
