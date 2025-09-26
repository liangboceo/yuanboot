package utils

import (
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode"
)

func PadLeft(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func UppercaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LowercaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 修改前的 Contains 函数
func Contains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// LikeContains 函数，添加了通配符 ** 支持
func LikeContains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			item := targetValue.Index(i).Interface()
			// 精确匹配
			if item == obj {
				return true
			}

			// 尝试路径通配符匹配
			if strObj, okObj := obj.(string); okObj {
				if strItem, okItem := item.(string); okItem {
					// 检查是否包含通配符 **
					if strings.Contains(strItem, "**") {
						// 处理通配符 ** 匹配
						parts := strings.Split(strItem, "**")
						if len(parts) == 2 {
							// 确保前缀匹配
							if strings.HasPrefix(strObj, parts[0]) {
								// 如果前缀后面没有其他内容，或者后缀匹配（如果有）
								if len(parts[1]) == 0 || strings.HasSuffix(strObj, parts[1]) {
									return true
								}
							}
						}
					}
				}
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func ContainsStr(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Substr字符串截取
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// Explode 将字符串按字符拆成数组
func Explode(delimiter, datastr string) (arr []string) {
	ret := strings.Split(datastr, delimiter)

	for _, item := range ret {
		if item != "" {
			arr = append(arr, item)
		}
	}

	return arr
}

func GetRandStr(n int) (randStr string) {
	// 默认去掉了容易混淆的字符oOLl和数字01，要添加请使用addChars参数
	chars := "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen := len(chars)
	if n > 10 {
		n = 10
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}

func NewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
