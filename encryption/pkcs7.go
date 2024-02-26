package encryption

import "bytes"

// Pkcs7Pad 填充
func Pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// Pkcs7UnPad 去除填充
func Pkcs7UnPad(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	if unPadding <= length {
		return src[:length-unPadding]
	}
	return src
}
