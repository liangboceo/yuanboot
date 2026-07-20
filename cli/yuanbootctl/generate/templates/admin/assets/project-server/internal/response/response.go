package response

import (
	"github.com/bytedance/sonic"
	"github.com/liangboceo/yuanboot/web/actionresult"
	"sendex-server/internal/dto"
)

func Success(data interface{}) actionresult.IActionResult {
	res, _ := sonic.Marshal(dto.Success(data))
	return actionresult.Data{ContentType: "application/json; charset=utf-8",
		Data: res}

}
func SuccessMessage(data interface{}, message string) actionresult.IActionResult {
	res, _ := sonic.Marshal(dto.SuccessMessage(data, message))
	return actionresult.Data{ContentType: "application/json; charset=utf-8",
		Data: res}
}
func Failure(data interface{}) actionresult.IActionResult {
	res, _ := sonic.Marshal(dto.Failure(data))
	return actionresult.Data{ContentType: "application/json; charset=utf-8",
		Data: res}
}
func FailureMessage(data interface{}, message string) actionresult.IActionResult {
	res, _ := sonic.Marshal(dto.FailureMessage(data, message))
	return actionresult.Data{ContentType: "application/json; charset=utf-8",
		Data: res}
}

func Result(data interface{}) actionresult.IActionResult {
	res, _ := sonic.Marshal(data)
	return actionresult.Data{ContentType: "application/json; charset=utf-8",
		Data: res}
}
