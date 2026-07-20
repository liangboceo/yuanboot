package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/google/uuid"
	context2 "github.com/liangboceo/yuanboot/web/context"
	"io"
	"net"
	"strings"
)

func GetGroupId(prefix string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return prefix + "_" + uuid.New().String()
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if addrs, err := iface.Addrs(); err == nil && len(addrs) > 0 {
				if mac := iface.HardwareAddr.String(); mac != "" {
					return prefix + "_" + strings.ReplaceAll(mac, ":", "")
				}
			}
		}
	}
	return prefix + "_" + uuid.New().String()
}
func GenSecret() string {
	return uuid.New().String()
}

func GetUserName(ctx *context2.HttpContext) string {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	if ctx.GetUser() == nil {
		return ""
	}
	return (ctx.GetUser())["username"].(string)
}

func GetUserId(ctx *context2.HttpContext) int {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	if ctx.GetUser() == nil {
		return 0
	}
	return int((ctx.GetUser())["userid"].(float64))
}

// AesKey AES加解密密钥（32字节 = AES-256）
const AesKey = "iot-link-bridge-dynamic-api-aes-k"

func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充数据到块大小
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// 生成随机 IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func AesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 生成随机 IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pkcs7RemovePad(ciphertext)
	return ciphertext, nil
}

// PKCS7 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去填充
func pkcs7RemovePad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
