package encryption

import "encoding/base64"

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

// Base64Encode base64加密
func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

// Base64Decode base64解密
func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}
