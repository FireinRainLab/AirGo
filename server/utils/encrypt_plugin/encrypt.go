package encrypt_plugin

import (
	"crypto/md5"
	crypt_rand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	PEM_BEGIN_RSA_PUBLIC  = "-----BEGIN RSA PUBLIC KEY-----\n"
	PEM_END_RSA_PUBLIC    = "\n-----END RSA PUBLIC KEY-----"
	PEM_BEGIN_RSA_PRIVATE = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END_RSA_PPRIVATE  = "\n-----END RSA PRIVATE KEY-----"
)

func FormatPublicKey(key string) string {
	if !strings.HasPrefix(key, PEM_BEGIN_RSA_PUBLIC) {
		key = PEM_BEGIN_RSA_PUBLIC + key
	}
	if !strings.HasSuffix(key, PEM_END_RSA_PUBLIC) {
		key = key + PEM_END_RSA_PUBLIC
	}
	return key
}
func FormatPrivateKey(key string) string {
	if !strings.HasPrefix(key, PEM_BEGIN_RSA_PRIVATE) {
		key = PEM_BEGIN_RSA_PRIVATE + key
	}
	if !strings.HasSuffix(key, PEM_END_RSA_PPRIVATE) {
		key = key + PEM_END_RSA_PPRIVATE
	}
	return key
}

func RSAEnCrypt(data, publicKey string) (string, error) {
	key := FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(key))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	encryptedData, err := rsa.EncryptPKCS1v15(crypt_rand.Reader, pubKey.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

func RSADecrypt(encryptedData, privateKey string) (string, error) {
	key := FormatPrivateKey(privateKey)

	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(key))
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	originalData, err := rsa.DecryptPKCS1v15(crypt_rand.Reader, priKey.(*rsa.PrivateKey), encryptedDecodeBytes)
	return string(originalData), err
}

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptEncode(password string) string {
	// Go 中的 bcrypt.DefaultCost 是 10
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptDecode(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// md5 encode
func Md5Encode(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	//转大写，strings.ToUpper(res)
	return res
}

// sha256 encode
func Sha256Encode(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

// 随机数，n为 位数
func RandomString(n int) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]rune, n)
	for i := range randomStr {
		randomStr[i] = defaultLetters[rand.Intn(len(defaultLetters))]
	}
	return string(randomStr)
}

// 随机数，n为 位数,去除大写字母和0
func RandomString2(n int) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyz123456789")
	rand.Seed(time.Now().UnixNano())
	randomStr := make([]rune, n)
	for i := range randomStr {
		randomStr[i] = defaultLetters[rand.Intn(len(defaultLetters))]
	}
	return string(randomStr)
}
