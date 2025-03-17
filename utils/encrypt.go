package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Md5ToLower md5大写
func Md5ToLower(str string) string {
	md5str := Md5String(str)
	md5str = strings.ToLower(md5str)
	return md5str
}

// Md5ToUpper md5小写
func Md5ToUpper(str string) string {
	md5str := Md5String(str)
	md5str = strings.ToUpper(md5str)
	return md5str
}

func Md5String(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// Sha256ToLower 方法接收一个字符串输入，生成其对应的 SHA-256 哈希值，并以十六进制字符串形式返回
func Sha256ToLower(input string) string {
	// 创建一个新的 SHA-256 哈希对象
	hash := sha256.New()
	// 向哈希对象中写入输入字符串的字节表示
	hash.Write([]byte(input))
	// 计算哈希值
	hashedBytes := hash.Sum(nil)
	// 将哈希值的字节切片转换为十六进制字符串
	return strings.ToLower(hex.EncodeToString(hashedBytes))
}

// Sha256ToUpper 方法接收一个字符串输入，生成其对应的 SHA-256 哈希值，并以十六进制字符串形式返回
func Sha256ToUpper(input string) string {
	// 创建一个新的 SHA-256 哈希对象
	hash := sha256.New()
	// 向哈希对象中写入输入字符串的字节表示
	hash.Write([]byte(input))
	// 计算哈希值
	hashedBytes := hash.Sum(nil)
	// 将哈希值的字节切片转换为十六进制字符串
	return strings.ToUpper(hex.EncodeToString(hashedBytes))
}
