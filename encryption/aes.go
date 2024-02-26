package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

// AesEncryptCBC 使用AES-128 CBC模式加密数据
func AesEncryptCBC(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS7填充
	plaintext = Pkcs7Pad(plaintext, block.BlockSize())

	// 创建加密模式为CBC的块模式
	iv := make([]byte, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密操作
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// AesDecryptCBC 使用AES-128 CBC模式解密数据
func AesDecryptCBC(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建加密模式为CBC的块模式
	iv := make([]byte, aes.BlockSize)
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密操作
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS7去除填充
	plaintext = Pkcs7UnPad(plaintext)

	return plaintext, nil
}

// AesEncryptECB 使用AES-128 ECB模式加密数据
func AesEncryptECB(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS7填充
	plaintext = Pkcs7Pad(plaintext, block.BlockSize())

	// 创建加密模式为CBC的块模式
	iv := make([]byte, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	mode := NewECBEncrypt(block, iv)

	// 加密操作
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// AesDecryptECB 使用AES-128 ECB模式解密数据
func AesDecryptECB(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建加密模式为CBC的块模式
	iv := make([]byte, aes.BlockSize)
	plaintext := make([]byte, len(ciphertext))
	mode := NewECBDecrypt(block, iv)

	// 解密操作
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS7去除填充
	plaintext = Pkcs7UnPad(plaintext)

	return plaintext, nil
}

// 自定义ECB加密器
type ecbEncrypt struct {
	b         cipher.Block
	blockSize int
}

func (x *ecbEncrypt) BlockSize() int { return x.blockSize }

func (x *ecbEncrypt) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// NewECBEncrypt 创建ECB加密器
func NewECBEncrypt(block cipher.Block, iv []byte) cipher.BlockMode {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		panic("cipher.NewECBEncrypt: IV length must equal block size")
	}
	return &ecbEncrypt{block, blockSize}
}

// 自定义ECB解密器
type ecbDecrypt struct {
	b         cipher.Block
	blockSize int
}

func (x *ecbDecrypt) BlockSize() int { return x.blockSize }

func (x *ecbDecrypt) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// NewECBDecrypt 创建ECB解密器
func NewECBDecrypt(block cipher.Block, iv []byte) cipher.BlockMode {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		panic("cipher.NewECBDecrypt: IV length must equal block size")
	}
	return &ecbDecrypt{block, blockSize}
}
