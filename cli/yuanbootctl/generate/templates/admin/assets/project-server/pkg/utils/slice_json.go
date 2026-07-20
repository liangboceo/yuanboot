package utils

import (
	"strings"

	"github.com/bytedance/sonic"
)

func MarshalUintSlice(items []uint) string {
	if len(items) == 0 {
		return "[]"
	}
	data, _ := sonic.Marshal(items)
	return string(data)
}

func UnmarshalUintSlice(data string) []uint {
	var items []uint
	if strings.TrimSpace(data) != "" {
		_ = sonic.Unmarshal([]byte(data), &items)
	}
	if items == nil {
		return []uint{}
	}
	return items
}

func MarshalStringSlice(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	data, _ := sonic.Marshal(items)
	return string(data)
}

func UnmarshalStringSlice(data string) []string {
	var items []string
	if strings.TrimSpace(data) != "" {
		_ = sonic.Unmarshal([]byte(data), &items)
	}
	if items == nil {
		return []string{}
	}
	return items
}
