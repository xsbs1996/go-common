package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
)

// 密码强度等级，D为最低
const (
	PasswordStrengthLevel1 = iota // 原始等级,空字符串
	PasswordStrengthLevel2        // 纯数字
	PasswordStrengthLevel3        // 数字+小写字母
	PasswordStrengthLevel4        // 数字+小写字母+大写字母
	PasswordStrengthLevel5        // 数字+小写字符+大写字符+符号
)

// 密码强度错误
var (
	PasswordStrengthCheckShortErr = errors.New("password is too short")          // 密码太短
	PasswordStrengthCheckLongErr  = errors.New("password is too long")           // 密码太长
	PasswordStrengthCheckLevelErr = errors.New("insufficient password strength") // 密码强度不足
)

// PasswordEncrypt 生成密码
func PasswordEncrypt(pwd string, salt string) (encodePwd string) {
	d := []byte(fmt.Sprintf("%s%s", salt, pwd))
	m := md5.New()
	m.Write(d)
	encodePwd = hex.EncodeToString(m.Sum(nil))
	return
}

// PasswordStrengthCheck 密码强度校验
func PasswordStrengthCheck(minLength, maxLength, minLevel int, pwd string) error {
	// 首先校验密码长度是否在范围内
	if len(pwd) < minLength {
		return PasswordStrengthCheckShortErr
	}
	if len(pwd) > maxLength {
		return PasswordStrengthCheckLongErr
	}

	// 初始化密码强度等级为1，利用正则校验密码强度，若匹配成功则强度自增1
	var level int = PasswordStrengthLevel1
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}

	// 如果最终密码强度低于要求的最低强度，返回并报错
	if level < minLevel {
		return PasswordStrengthCheckLevelErr
	}
	return nil
}
