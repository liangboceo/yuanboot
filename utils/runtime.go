package utils

import (
	"runtime"
	"strings"
)

func GoId() string {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}
