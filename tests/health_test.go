package tests

import (
	"github.com/bytedance/sonic"
	abshealth "github.com/liangboceo/yuanboot/abstractions/health"
	"github.com/stretchr/testify/assert"
	"testing"
)

type downHealth struct{}

func (u downHealth) Health() abshealth.ComponentStatus {
	return abshealth.Down("downHealth").WithDetail("reason", "error:down")
}

type upHealth struct{}

func (u upHealth) Health() abshealth.ComponentStatus {
	return abshealth.Up("UpHealth").
		WithDetail("total", 1024).
		WithDetail("current", 50)
}

func TestHealth(t *testing.T) {
	var indicatorList []abshealth.Indicator
	indicatorList = append(indicatorList, downHealth{}, upHealth{}, abshealth.NewDiskHealthIndicator())
	builder := abshealth.NewHealthIndicator(indicatorList)
	m := builder.Build()
	bytes, _ := sonic.Marshal(m)
	jsonstr := string(bytes)
	assert.NotNil(t, jsonstr)
}
