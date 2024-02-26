package auth

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type JwtTokenPayload struct {
	Info map[string]interface{} `json:"info,omitempty"` // 用户信息存储
	jwt.RegisteredClaims
}

func GetJwtPayload(info map[string]interface{}) *JwtTokenPayload {
	return &JwtTokenPayload{
		Info: info,
	}
}

var (
	jwtIsEmptyErr = errors.New("the token string is empty") // 字符串为空
	jwtSignErr    = errors.New("unexpected signing method") // 错误的签名方式
	jwtParseErr   = errors.New("failed to parse jwt token") // 解析令牌错误
	jwtInvalidErr = errors.New("invalid token")             // 令牌无效
)

// JwtParseToken token验证
func JwtParseToken(key string, tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, jwtIsEmptyErr
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwtSignErr
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, jwtParseErr
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logrus.WithFields(logrus.Fields{"invalidToken": tokenString}).Warn("invalid token")
		return nil, jwtInvalidErr
	}

	return claims, nil
}

// JwtCreateToken 生成token
func JwtCreateToken(key string, obj *JwtTokenPayload, expiresAt time.Time) (string, error) {
	rand.Seed(time.Now().UnixNano())

	obj.ExpiresAt = jwt.NewNumericDate(expiresAt)
	obj.IssuedAt = jwt.NewNumericDate(time.Now())
	obj.ID = strconv.Itoa(rand.Intn(1000))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, obj)
	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}
