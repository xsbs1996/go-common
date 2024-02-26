package encryption

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"git.ejxcgit.com/ejhycommon/go-common/encryption/gorsa"
	"io"
)

var (
	RsaPrivateKeyIsNotSet = errors.New("private key is not set")
	RsaPublicKeyIsNotSet  = errors.New("public key is not set")
)

type CommonRsa struct {
	rsaPublicKey  *rsa.PublicKey  // 解析后公钥
	rsaPrivateKey *rsa.PrivateKey // 解析后私钥
}

func NewCommonRsa(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *CommonRsa {
	return &CommonRsa{
		rsaPublicKey:  publicKey,
		rsaPrivateKey: privateKey,
	}
}

// Encrypt 公钥加密
func (r *CommonRsa) Encrypt(data []byte) ([]byte, error) {
	if r.rsaPublicKey == nil {
		return nil, RsaPublicKeyIsNotSet
	}
	// blockLength = 密钥长度 = 一次能加密的明文长度。 "/8" 将bit转为bytes。"-11" 为 PKCS#1 建议的 padding 占用了 11 个字节
	blockLength := r.rsaPublicKey.N.BitLen()/8 - 11

	// 如果明文长度不大于密钥长度，可以直接加密
	if len(data) <= blockLength {
		return rsa.EncryptPKCS1v15(rand.Reader, r.rsaPublicKey, []byte(data))
	}

	// 明文长度大于密钥长度，创建一个新的缓冲区
	buffer := bytes.NewBufferString("")
	pages := len(data) / blockLength //切分为多少块
	//循环加密
	for i := 0; i <= pages; i++ {
		start := i * blockLength
		end := (i + 1) * blockLength
		if i == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}
		//分段加密
		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, r.rsaPublicKey, data[start:end])
		if err != nil {
			return nil, err
		}
		//写入缓冲区
		buffer.Write(chunk)
	}
	//读取缓冲区内容并返回，即返回加密结果
	return buffer.Bytes(), nil
}

// Decrypt 私钥解密
func (r *CommonRsa) Decrypt(data []byte) ([]byte, error) {
	if r.rsaPrivateKey == nil {
		return nil, RsaPrivateKeyIsNotSet
	}
	// 获取密钥长度。加密后的密文长度=密钥长度。如果密文长度大于密钥长度，说明密文非一次加密形成
	blockLength := r.rsaPublicKey.N.BitLen() / 8

	// 一次形成的密文直接解密
	if len(data) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, r.rsaPrivateKey, data)
	}

	// 循环解密
	buffer := bytes.NewBufferString("")
	pages := len(data) / blockLength
	for i := 0; i <= pages; i++ {
		start := i * blockLength
		end := (i + 1) * blockLength
		if i == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}
		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, r.rsaPrivateKey, data[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

// PriEncrypt 私钥加密
func (r *CommonRsa) PriEncrypt(data []byte) ([]byte, error) {
	if r.rsaPrivateKey == nil {
		return nil, RsaPrivateKeyIsNotSet
	}

	output := bytes.NewBuffer(nil)
	err := gorsa.PriKeyIO(r.rsaPrivateKey, bytes.NewReader(data), output, true)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// PubDecrypt 公钥解密
func (r *CommonRsa) PubDecrypt(data []byte) ([]byte, error) {
	if r.rsaPublicKey == nil {
		return nil, RsaPublicKeyIsNotSet
	}

	output := bytes.NewBuffer(nil)
	err := gorsa.PubKeyIO(r.rsaPublicKey, bytes.NewReader(data), output, false)
	if err != nil {
		return []byte(""), err
	}
	return io.ReadAll(output)
}

// Sign 私钥签名
func (r *CommonRsa) Sign(data []byte, sHash crypto.Hash) ([]byte, error) {
	hash := sHash.New()
	hash.Write(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.rsaPrivateKey, sHash, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, nil
}

// Verify 公钥验签
func (r *CommonRsa) Verify(data []byte, sign []byte, sHash crypto.Hash) bool {
	h := sHash.New()
	h.Write(data)
	return rsa.VerifyPKCS1v15(r.rsaPublicKey, sHash, h.Sum(nil), sign) == nil
}

// CreateKeys 生成pkcs1 格式的公钥私钥
func (r *CommonRsa) CreateKeys(keyLength int) (privateKey, publicKey string) {
	//根据 随机源 与 指定位数，生成密钥对。rand.Reader = 密码强大的伪随机生成器的全球共享实例
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}
	//编码私钥
	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY", //自定义类型
		Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
	}))
	//编码公钥
	objPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}
	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: objPkix,
	}))
	return
}

// CreatePkcs8Keys 生成pkcs8 格式公钥私钥
func (r *CommonRsa) CreatePkcs8Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	objPkcs8, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: objPkcs8,
	}))

	objPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: objPkix,
	}))
	return
}

// Pkcs1ToPkcs8 将pkcs1 转到 pkcs8 自定义
func (r *CommonRsa) Pkcs1ToPkcs8(key []byte) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = key
	k, _ := asn1.Marshal(info)
	return k
}
