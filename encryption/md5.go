package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Encrypt md5加密
func Md5Encrypt(data []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	return encryptedData
}
