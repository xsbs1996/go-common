package gotool

import (
	"regexp"
)

// PhpPhoneRegex 菲律宾电话验证
func PhpPhoneRegex(phoneNumber string) bool {
	re := regexp.MustCompile(`^(\+63|0)?\d{9,12}$`)
	return re.MatchString(phoneNumber)
}

// EmailRegex 邮箱验证
func EmailRegex(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
