package dto

import "time"

type BaseReq struct {
	ID            uint      `json:"id" `
	IsDeleted     int       `json:"isDeleted" doc:"是否删除【0：未删除 1：删除】"`
	CreateBy      int       `json:"createBy" doc:"创建人"`
	CreateOrgCode string    `json:"createOrgCode" doc:"创建组织"`
	CreateTime    time.Time `json:"createTime" doc:"创建时间"`
	UpdateBy      int       `json:"updateBy" doc:"更新人"`
	UpdateTime    time.Time `json:"updateTime" doc:"更新时间"`
	CreateName    string    `json:"createName" doc:"创建人名称"`
	UpdateName    string    `json:"updateName" doc:"更新人名称"`
}
type BaseResp struct {
	ID            uint      `json:"id" `
	IsDeleted     int       `json:"isDeleted" doc:"是否删除【0：未删除 1：删除】"`
	CreateBy      int       `json:"createBy" doc:"创建人"`
	CreateOrgCode string    `json:"createOrgCode" doc:"创建组织"`
	CreateTime    time.Time `json:"createTime" doc:"创建时间"`
	UpdateBy      int       `json:"updateBy" doc:"更新人"`
	UpdateTime    time.Time `json:"updateTime" doc:"更新时间"`
	CreateName    string    `json:"createName" doc:"创建人名称"`
	UpdateName    string    `json:"updateName" doc:"更新人名称"`
}
type BasePageResp struct {
	Total int64 `json:"total" doc:"总数"`
	Page  int   `json:"page" doc:"当前页数"`
}

type Result struct {
	Status  bool        `json:"status" doc:"状态"`
	Code    int         `json:"code" doc:"错误码"`
	Message string      `json:"message" doc:"错误信息"`
	Data    interface{} `json:"data" doc:"返回体"`
}

const NO_LOGIN = -1
const SERVER_ERROR = 500
const SUCCESS = 200

func Success(data interface{}) Result {
	result := Result{}
	result.Status = true
	result.Code = SUCCESS
	result.Message = ""
	result.Data = data
	return result
}
func SuccessMessage(data interface{}, message string) Result {
	result := Result{}
	result.Status = true
	result.Code = SUCCESS
	result.Message = message
	result.Data = data
	return result
}
func Failure(data interface{}) Result {
	result := Result{}
	result.Status = false
	result.Code = SERVER_ERROR
	result.Message = ""
	result.Data = data
	return result
}
func FailureMessage(data interface{}, message string) Result {
	result := Result{}
	result.Status = false
	result.Code = SERVER_ERROR
	result.Message = message
	result.Data = data
	return result
}
