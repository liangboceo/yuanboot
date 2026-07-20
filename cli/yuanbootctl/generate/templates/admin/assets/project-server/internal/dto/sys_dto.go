package dto

import (
	"github.com/liangboceo/yuanboot/web/mvc"
	"time"
)

// ==================== 登录请求 ====================
type LoginReq struct {
	mvc.RequestBody `route:"/sys/Login" doc:"登录"`
	Username        string `json:"username" doc:"用户名" binding:"required"`
	Password        string `json:"password" doc:"密码" binding:"required"`
	CaptchaId       string `json:"captchaId" doc:"验证码流水ID" binding:"required"`
	VerifyCode      string `json:"verifyCode" doc:"图形验证码" binding:"required"`
}

type CaptchaReq struct {
	mvc.RequestBody `route:"/sys/captcha" doc:"获取图形验证码"`
}

type CaptchaResp struct {
	CaptchaId string `json:"captchaId" doc:"验证码流水ID"`
	Image     string `json:"image" doc:"Base64图片"`
}

type GetUserInfoReq struct {
	mvc.RequestBody `route:"/sys/GetUserInfo" doc:"获取用户信息"`
}

type GetRoutesReq struct {
	mvc.RequestBody `route:"/sys/GetRoutes" doc:"获取路由"`
}

type GetHomeStatsReq struct {
	mvc.RequestBody `route:"/sys/GetHomeStats" doc:"获取首页统计"`
}

type HomeStatsResp struct {
	ModelTotal        int64 `json:"modelTotal" doc:"模型总数"`
	EnabledModelTotal int64 `json:"enabledModelTotal" doc:"启用模型数"`
	UserTotal         int64 `json:"userTotal" doc:"用户总数"`
	EnabledUserTotal  int64 `json:"enabledUserTotal" doc:"正常用户数"`
	SkillTotal        int64 `json:"skillTotal" doc:"技能总数"`
	EnabledSkillTotal int64 `json:"enabledSkillTotal" doc:"启用技能数"`
	SessionTotal      int64 `json:"sessionTotal" doc:"会话总数"`
	MessageTotal      int64 `json:"messageTotal" doc:"消息总数"`
	AgentTotal        int64 `json:"agentTotal" doc:"Agent总数"`
	EnabledAgentTotal int64 `json:"enabledAgentTotal" doc:"启用Agent数"`
}

// ==================== 登录响应 ====================
type LoginResp struct {
	AccessToken  string   `json:"accessToken" doc:"访问令牌"`
	Expires      int64    `json:"expires" doc:"过期时间戳"`
	RefreshToken string   `json:"refreshToken" doc:"刷新令牌"`
	Avatar       string   `json:"avatar" doc:"头像"`
	Username     string   `json:"username" doc:"用户名"`
	Nickname     string   `json:"nickname" doc:"昵称"`
	Roles        []string `json:"roles" doc:"角色列表"`
	Permissions  []string `json:"permissions" doc:"权限列表"`
}

// ==================== 用户管理请求 ====================
type SysUserListReq struct {
	mvc.RequestBody `route:"/sys/user" doc:"获取用户列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Username        string `json:"username" doc:"用户名"`
	Nickname        string `json:"nickname" doc:"昵称"`
	Phone           string `json:"phone" doc:"手机号"`
	Status          string `json:"status" doc:"状态"`
	DeptId          string `json:"deptId" doc:"部门ID"`
}

type SysUserSaveReq struct {
	mvc.RequestBody `route:"/sys/user/save" doc:"保存用户"`
	ID              uint   `json:"id" doc:"用户ID（更新时必填）"`
	Username        string `json:"username" doc:"用户名" binding:"required"`
	Password        string `json:"password" doc:"密码"`
	Nickname        string `json:"nickname" doc:"昵称"`
	Email           string `json:"email" doc:"邮箱"`
	Phone           string `json:"phone" doc:"手机号"`
	Avatar          string `json:"avatar" doc:"头像"`
	DeptId          uint   `json:"deptId" doc:"部门ID"`
	Status          int    `json:"status" doc:"状态" default:"1"`
	RoleIds         []uint `json:"roleIds" doc:"角色ID列表"`
}

type SysUserDetailReq struct {
	mvc.RequestBody `route:"/sys/user/detail" doc:"获取用户详情"`
	ID              uint `json:"id" doc:"用户ID"`
}

type SysUserDelReq struct {
	mvc.RequestBody `route:"/sys/user/del" doc:"删除用户"`
	ID              uint `json:"id" doc:"用户ID" binding:"required"`
}

type SysUserResetPwdReq struct {
	mvc.RequestBody `route:"/sys/user/resetPwd" doc:"重置密码"`
	ID              uint   `json:"id" doc:"用户ID" binding:"required"`
	Password        string `json:"password" doc:"新密码" binding:"required"`
}

type SysUserChangePwdReq struct {
	mvc.RequestBody `route:"/sys/changePwd" doc:"修改当前用户密码"`
	OldPassword     string `json:"oldPassword" doc:"旧密码" binding:"required"`
	NewPassword     string `json:"newPassword" doc:"新密码" binding:"required"`
}

// ==================== 用户管理响应 ====================
type SysUserResp struct {
	BaseResp
	Username  string    `json:"username" doc:"用户名"`
	Nickname  string    `json:"nickname" doc:"昵称"`
	Email     string    `json:"email" doc:"邮箱"`
	Phone     string    `json:"phone" doc:"手机号"`
	Avatar    string    `json:"avatar" doc:"头像"`
	DeptId    uint      `json:"deptId" doc:"部门ID"`
	DeptName  string    `json:"deptName" doc:"部门名称"`
	Status    int       `json:"status" doc:"状态"`
	IsAdmin   int       `json:"isAdmin" doc:"是否管理员"`
	LoginIp   string    `json:"loginIp" doc:"最后登录IP"`
	LoginDate time.Time `json:"loginDate" doc:"最后登录时间"`
	RoleIds   []uint    `json:"roleIds" doc:"角色ID列表"`
	RoleNames string    `json:"roleNames" doc:"角色名称列表"`
}

type SysUserListResp struct {
	BasePageResp
	Data []SysUserResp `json:"list" doc:"用户列表"`
}

// ==================== 角色管理请求 ====================
type SysRoleListReq struct {
	mvc.RequestBody `route:"/sys/role" doc:"获取角色列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Name            string `json:"name" doc:"角色名称"`
	Code            string `json:"code" doc:"角色编码"`
	Status          string `json:"status" doc:"状态"`
}

type SysAllRoleListReq struct {
	mvc.RequestBody `route:"/sys/role/all" doc:"获取所有角色"`
}

type SysRoleSaveReq struct {
	mvc.RequestBody   `route:"/sys/role/save" doc:"保存角色"`
	ID                uint   `json:"id" doc:"角色ID（更新时必填）"`
	Name              string `json:"name" doc:"角色名称" binding:"required"`
	Code              string `json:"code" doc:"角色编码" binding:"required"`
	Sort              int    `json:"sort" doc:"排序" default:"1"`
	DataScope         int    `json:"dataScope" doc:"数据权限" default:"1"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly" doc:"菜单树选择范围是否关联显示" default:"true"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly" doc:"部门树选择范围是否关联显示" default:"true"`
	Status            int    `json:"status" doc:"状态" default:"1"`
	Remark            string `json:"remark" doc:"备注"`
	MenuIds           []uint `json:"menuIds" doc:"菜单ID列表"`
	DeptIds           []uint `json:"deptIds" doc:"部门ID列表"`
}

type SysRoleDetailReq struct {
	mvc.RequestBody `route:"/sys/role/detail" doc:"获取角色详情"`
	ID              uint `json:"id" doc:"角色ID"`
}

type SysRoleDelReq struct {
	mvc.RequestBody `route:"/sys/role/del" doc:"删除角色"`
	ID              uint `json:"id" doc:"角色ID" binding:"required"`
}

type SysRoleMenuReq struct {
	mvc.RequestBody `route:"/sys/role/menu" doc:"获取角色菜单权限"`
	RoleId          uint `json:"roleId" doc:"角色ID"`
}

type SysRoleMenuResp struct {
	CheckedKeys []uint        `json:"checkedKeys" doc:"选中的菜单ID列表"`
	MenuList    []SysMenuResp `json:"menuList" doc:"菜单树列表"`
}

// ==================== 角色管理响应 ====================
type SysRoleResp struct {
	BaseResp
	Name              string `json:"name" doc:"角色名称"`
	Code              string `json:"code" doc:"角色编码"`
	Sort              int    `json:"sort" doc:"排序"`
	DataScope         int    `json:"dataScope" doc:"数据权限"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly" doc:"菜单树选择范围是否关联显示"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly" doc:"部门树选择范围是否关联显示"`
	Status            int    `json:"status" doc:"状态"`
	Remark            string `json:"remark" doc:"备注"`
}

type SysRoleListResp struct {
	BasePageResp
	Data []SysRoleResp `json:"list" doc:"角色列表"`
}

// ==================== 菜单管理请求 ====================
type SysMenuListReq struct {
	mvc.RequestBody `route:"/sys/menu" doc:"获取菜单列表"`
	MenuName        string `json:"menuName" doc:"菜单名称"`
	Status          int    `json:"status" doc:"状态"`
}

type SysMenuTreeReq struct {
	mvc.RequestBody `route:"/sys/menu/tree" doc:"获取菜单树"`
}

type SysMenuSaveReq struct {
	mvc.RequestBody `route:"/sys/menu/save" doc:"保存菜单"`
	ID              uint   `json:"id" doc:"菜单ID（更新时必填）"`
	Name            string `json:"name" doc:"菜单名称" binding:"required"`
	ParentID        uint   `json:"parentId" doc:"父菜单ID" default:"0"`
	OrderNum        int    `json:"orderNum" doc:"显示顺序" default:"1"`
	Path            string `json:"path" doc:"路由地址"`
	Component       string `json:"component" doc:"组件路径"`
	MenuType        int    `json:"menuType" doc:"菜单类型【0：目录 1：菜单 2：按钮 3：内页】" default:"1"`
	Visible         int    `json:"visible" doc:"显示状态【0：显示 1：隐藏】" default:"0"`
	Status          int    `json:"status" doc:"状态【0：禁用 1：正常】" default:"1"`
	Icon            string `json:"icon" doc:"菜单图标"`
	IsFrame         int    `json:"isFrame" doc:"是否外链【0：否 1：是】" default:"0"`
	IsCache         int    `json:"isCache" doc:"是否缓存【0：否 1：是】" default:"0"`
	Permission      string `json:"permission" doc:"权限标识"`
	Query           string `json:"query" doc:"路由参数"`
	Perms           string `json:"perms" doc:"权限字符"`
}

type SysMenuDetailReq struct {
	mvc.RequestBody `route:"/sys/menu/detail" doc:"获取菜单详情"`
	ID              uint `json:"id" doc:"菜单ID"`
}

type SysMenuDelReq struct {
	mvc.RequestBody `route:"/sys/menu/del" doc:"删除菜单"`
	ID              uint `json:"id" doc:"菜单ID" binding:"required"`
}

// ==================== 菜单管理响应 ====================
type SysMenuResp struct {
	ID         uint          `json:"id" doc:"菜单ID"`
	ParentID   uint          `json:"parentId" doc:"父菜单ID"`
	Name       string        `json:"name" doc:"菜单名称"`
	OrderNum   int           `json:"orderNum" doc:"显示顺序"`
	Path       string        `json:"path" doc:"路由地址"`
	Component  string        `json:"component" doc:"组件路径"`
	MenuType   int           `json:"menuType" doc:"菜单类型"`
	Visible    int           `json:"visible" doc:"显示状态"`
	Status     int           `json:"status" doc:"状态"`
	Icon       string        `json:"icon" doc:"菜单图标"`
	IsFrame    int           `json:"isFrame" doc:"是否外链"`
	IsCache    int           `json:"isCache" doc:"是否缓存"`
	Permission string        `json:"permission" doc:"权限标识"`
	Query      string        `json:"query" doc:"路由参数"`
	Perms      string        `json:"perms" doc:"权限字符"`
	Children   []SysMenuResp `json:"children" doc:"子菜单"`
}

type SysMenuTreeResp struct {
	ID       uint              `json:"id" doc:"菜单ID"`
	ParentID uint              `json:"parentId" doc:"父菜单ID"`
	Name     string            `json:"label" doc:"菜单名称"`
	OrderNum int               `json:"orderNum" doc:"显示顺序"`
	MenuType int               `json:"menuType" doc:"菜单类型"`
	Icon     string            `json:"icon" doc:"菜单图标"`
	Children []SysMenuTreeResp `json:"children" doc:"子菜单"`
}

// ==================== 模型管理请求 ====================
type SysModelListReq struct {
	mvc.RequestBody `route:"/sys/GetModelConfigList" doc:"获取模型配置列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Name            string `json:"name" doc:"配置名称"`
	Provider        string `json:"provider" doc:"提供商"`
	ModelName       string `json:"modelName" doc:"模型名称"`
	Status          string `json:"status" doc:"状态"`
}

type SysModelSaveReq struct {
	mvc.RequestBody `route:"/sys/SaveModelConfig" doc:"保存模型配置"`
	ID              uint   `json:"id" doc:"模型配置ID（更新时必填）"`
	Name            string `json:"name" doc:"配置名称" binding:"required"`
	Provider        string `json:"provider" doc:"提供商" binding:"required"`
	ModelName       string `json:"modelName" doc:"模型名称" binding:"required"`
	BaseURL         string `json:"baseUrl" doc:"接口地址"`
	APIKey          string `json:"apiKey" doc:"API Key"`
	Status          int    `json:"status" doc:"状态" default:"1"`
	Remark          string `json:"remark" doc:"备注"`
}

type SysModelDelReq struct {
	mvc.RequestBody `route:"/sys/DelModelConfig" doc:"删除模型配置"`
	ID              uint `json:"id" doc:"模型配置ID" binding:"required"`
}

type SysModelTestReq struct {
	mvc.RequestBody `route:"/sys/TestModelConfig" doc:"测试模型配置连接"`
	Provider        string `json:"provider" doc:"提供商" binding:"required"`
	ModelName       string `json:"modelName" doc:"模型名称" binding:"required"`
	BaseURL         string `json:"baseUrl" doc:"接口地址"`
	APIKey          string `json:"apiKey" doc:"API Key"`
}

type SysModelResp struct {
	BaseResp
	Name      string `json:"name" doc:"配置名称"`
	Provider  string `json:"provider" doc:"提供商"`
	ModelName string `json:"modelName" doc:"模型名称"`
	BaseURL   string `json:"baseUrl" doc:"接口地址"`
	APIKey    string `json:"apiKey" doc:"API Key"`
	Status    int    `json:"status" doc:"状态"`
	Remark    string `json:"remark" doc:"备注"`
}

type SysModelListResp struct {
	BasePageResp
	Data []SysModelResp `json:"list" doc:"模型配置列表"`
}

type SysModelDiscoverReq struct {
	mvc.RequestBody `route:"/sys/ModelDiscover" doc:"模型发现"`
	ID              uint `json:"id" doc:"模型配置ID" binding:"required"`
}

type SysModelDiscoverResp struct {
	Discovered int `json:"discovered" doc:"发现模型数"`
	Added      int `json:"added" doc:"新增数量"`
	Updated    int `json:"updated" doc:"更新数量"`
	Deleted    int `json:"deleted" doc:"删除数量"`
}

type SysSkillListReq struct {
	mvc.RequestBody `route:"/sys/GetSkillList" doc:"获取技能列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Keyword         string `json:"keyword" doc:"关键词"`
	Category        string `json:"category" doc:"技能分类"`
	Status          string `json:"status" doc:"状态"`
}

type SysSkillSaveReq struct {
	mvc.RequestBody `route:"/sys/SaveSkill" doc:"保存技能"`
	ID              uint               `json:"id" doc:"技能ID"`
	Name            string             `json:"name" doc:"技能名称" binding:"required"`
	Code            string             `json:"code" doc:"技能编码" binding:"required"`
	Category        string             `json:"category" doc:"技能分类"`
	Icon            string             `json:"icon" doc:"技能图标"`
	Description     string             `json:"description" doc:"技能描述"`
	Content         string             `json:"content" doc:"技能内容"`
	Files           []SysSkillFileResp `json:"files" doc:"技能文件列表"`
	Status          int                `json:"status" doc:"状态" default:"1"`
	Remark          string             `json:"remark" doc:"备注"`
	SkipCompile     bool               `json:"skipCompile" doc:"是否跳过编译"`
}

type SysSkillDetailReq struct {
	mvc.RequestBody `route:"/sys/GetSkillDetail" doc:"获取技能详情"`
	ID              uint `json:"id" doc:"技能ID" binding:"required"`
}

type SysSkillDelReq struct {
	mvc.RequestBody `route:"/sys/DelSkill" doc:"删除技能"`
	ID              uint `json:"id" doc:"技能ID" binding:"required"`
}

type SysSkillGenerateReq struct {
	mvc.RequestBody `route:"/sys/GenerateSkill" doc:"生成技能"`
	ModelID         uint   `json:"modelId" doc:"模型配置ID" binding:"required"`
	Prompt          string `json:"prompt" doc:"业务描述" binding:"required"`
	SkillID         uint   `json:"skillId" doc:"技能ID"`
	Content         string `json:"content" doc:"当前技能草稿"`
}

type SysSkillFileResp struct {
	Path    string `json:"path" doc:"文件路径"`
	Name    string `json:"name" doc:"文件名称"`
	Content string `json:"content" doc:"文件内容"`
}

type SysSkillResp struct {
	BaseResp
	Name        string             `json:"name" doc:"技能名称"`
	Code        string             `json:"code" doc:"技能编码"`
	Category    string             `json:"category" doc:"技能分类"`
	Icon        string             `json:"icon" doc:"技能图标"`
	Description string             `json:"description" doc:"技能描述"`
	Content     string             `json:"content" doc:"技能内容"`
	Files       []SysSkillFileResp `json:"files" doc:"技能文件列表"`
	Status      int                `json:"status" doc:"状态"`
	Remark      string             `json:"remark" doc:"备注"`
}

type SysSkillListResp struct {
	BasePageResp
	Data []SysSkillResp `json:"list" doc:"技能列表"`
}

type SysSkillCategoryReq struct {
	mvc.RequestBody `route:"/sys/GetSkillCategories" doc:"获取技能分类"`
}

type SysSkillCategoryListReq struct {
	mvc.RequestBody `route:"/sys/GetSkillCategoryList" doc:"获取技能分类列表（含ID）"`
	Keyword         string `json:"keyword" doc:"关键词"`
}

type SysSkillCategorySaveReq struct {
	mvc.RequestBody `route:"/sys/SaveSkillCategory" doc:"保存技能分类"`
	ID              uint   `json:"id" doc:"分类ID（更新必填）"`
	Name            string `json:"name" doc:"分类名称" binding:"required"`
	Sort            int    `json:"sort" doc:"排序" default:"1"`
	Status          int    `json:"status" doc:"状态" default:"1"`
	Remark          string `json:"remark" doc:"备注"`
}

type SysSkillCategoryDelReq struct {
	mvc.RequestBody `route:"/sys/DelSkillCategory" doc:"删除技能分类"`
	ID              uint `json:"id" doc:"分类ID" binding:"required"`
}

type SysSkillCategoryResp struct {
	BaseResp
	Name   string `json:"name" doc:"分类名称"`
	Sort   int    `json:"sort" doc:"排序"`
	Status int    `json:"status" doc:"状态"`
	Remark string `json:"remark" doc:"备注"`
}

// ==================== Agent 管理请求 ====================
type SysAgentListReq struct {
	mvc.RequestBody `route:"/sys/GetAgentList" doc:"获取Agent列表"`
	PageNum         int    `json:"pageNum" doc:"分页页码" default:"1"`
	PageSize        int    `json:"pageSize" doc:"每页数量" default:"20"`
	Keyword         string `json:"keyword" doc:"关键词"`
	ModelMode       string `json:"modelMode" doc:"模型模式"`
	Status          string `json:"status" doc:"状态"`
}

type SysAgentDetailReq struct {
	mvc.RequestBody `route:"/sys/GetAgentDetail" doc:"获取Agent详情"`
	ID              uint `json:"id" doc:"Agent ID" binding:"required"`
}

type SysAgentSaveReq struct {
	mvc.RequestBody      `route:"/sys/SaveAgent" doc:"保存Agent"`
	ID                   uint     `json:"id" doc:"Agent ID"`
	Name                 string   `json:"name" doc:"Agent名称" binding:"required"`
	RoleName             string   `json:"roleName" doc:"自定义角色名称" binding:"required"`
	Avatar               string   `json:"avatar" doc:"头像类型"`
	Description          string   `json:"description" doc:"Agent描述"`
	Prompt               string   `json:"prompt" doc:"Agent提示词" binding:"required"`
	PromptEngineering    string   `json:"promptEngineering" doc:"提示词工程"`
	ModelMode            string   `json:"modelMode" doc:"模型模式【auto/model】"`
	ModelID              uint     `json:"modelId" doc:"模型配置ID"`
	SkillIDs             []uint   `json:"skillIds" doc:"Skill ID列表"`
	ToolCodes            []string `json:"toolCodes" doc:"系统默认Tool编码列表"`
	KnowledgeCategoryIDs []uint   `json:"knowledgeCategoryIds" doc:"知识库分类ID列表"`
	KnowledgeDocIDs      []uint   `json:"knowledgeDocIds" doc:"知识库文档ID列表"`
	WebhookIDs           []uint   `json:"webhookIds" doc:"Webhook ID列表"`
	GraphConfig          string   `json:"graphConfig" doc:"流程配置"`
	Status               int      `json:"status" doc:"状态" default:"1"`
	Remark               string   `json:"remark" doc:"备注"`
}

type SysAgentDelReq struct {
	mvc.RequestBody `route:"/sys/DelAgent" doc:"删除Agent"`
	ID              uint `json:"id" doc:"Agent ID" binding:"required"`
}

type SysAgentDebugReq struct {
	mvc.RequestBody `route:"/sys/GetAgentDebugInfo" doc:"获取Agent调试信息"`
	ID              uint `json:"id" doc:"Agent ID" binding:"required"`
}

type SysAgentToolResp struct {
	Code        string `json:"code" doc:"工具编码"`
	Name        string `json:"name" doc:"工具名称"`
	Description string `json:"description" doc:"工具描述"`
}

type KnowledgeDocBriefResp struct {
	ID       uint   `json:"id" doc:"文档ID"`
	Name     string `json:"name" doc:"文档名称"`
	FileName string `json:"fileName" doc:"文件名"`
	FileType string `json:"fileType" doc:"文件类型"`
}

type SysAgentResp struct {
	BaseResp
	Name                 string                  `json:"name" doc:"Agent名称"`
	RoleName             string                  `json:"roleName" doc:"自定义角色名称"`
	Avatar               string                  `json:"avatar" doc:"头像类型"`
	Description          string                  `json:"description" doc:"Agent描述"`
	Prompt               string                  `json:"prompt" doc:"Agent提示词"`
	PromptEngineering    string                  `json:"promptEngineering" doc:"提示词工程"`
	ModelMode            string                  `json:"modelMode" doc:"模型模式"`
	ModelID              uint                    `json:"modelId" doc:"模型配置ID"`
	ModelName            string                  `json:"modelName" doc:"模型名称"`
	SkillIDs             []uint                  `json:"skillIds" doc:"Skill ID列表"`
	Skills               []SysSkillResp          `json:"skills" doc:"已选Skill"`
	ToolCodes            []string                `json:"toolCodes" doc:"系统默认Tool编码列表"`
	Tools                []SysAgentToolResp      `json:"tools" doc:"已选系统Tool"`
	KnowledgeCategoryIDs []uint                  `json:"knowledgeCategoryIds" doc:"知识库分类ID列表"`
	KnowledgeDocIDs      []uint                  `json:"knowledgeDocIds" doc:"知识库文档ID列表"`
	KnowledgeDocs        []KnowledgeDocBriefResp `json:"knowledgeDocs" doc:"已选知识库文档"`
	WebhookIDs           []uint                  `json:"webhookIds" doc:"Webhook ID列表"`
	GraphConfig          string                  `json:"graphConfig" doc:"流程配置"`
	Status               int                     `json:"status" doc:"状态"`
	Remark               string                  `json:"remark" doc:"备注"`
}

type SysAgentListResp struct {
	BasePageResp
	Data []SysAgentResp `json:"list" doc:"Agent列表"`
}

type SysAgentOptionsReq struct {
	mvc.RequestBody `route:"/sys/GetAgentOptions" doc:"获取Agent配置选项"`
}

type SysAgentOptionsResp struct {
	Models   []SysModelResp       `json:"models" doc:"模型列表"`
	Skills   []SysSkillResp       `json:"skills" doc:"Skill列表"`
	Tools    []SysAgentToolResp   `json:"tools" doc:"系统默认Tool列表"`
	Webhooks []SysWebhookItemResp `json:"webhooks" doc:"可用的Webhook列表"`
}

type SysWebhookItemResp struct {
	ID   uint   `json:"id" doc:"Webhook ID"`
	Name string `json:"name" doc:"Webhook名称"`
	Type string `json:"type" doc:"Webhook类型"`
}

type SysSkillExportReq struct {
	mvc.RequestBody `route:"/sys/ExportSkill" doc:"导出单个技能为ZIP"`
	ID              uint `json:"id" doc:"技能ID" binding:"required"`
}

type SysSkillBatchExportReq struct {
	mvc.RequestBody `route:"/sys/BatchExportSkill" doc:"批量导出技能为ZIP"`
	IDs             []uint `json:"ids" doc:"技能ID列表" binding:"required"`
}

// SysDeptListReq SysDeptListReq SysDeptListReq SysDeptListReq SysDeptListReq SysDeptListReq ==================== 部门管理请求 ====================
type SysDeptListReq struct {
	mvc.RequestBody `route:"/sys/dept" doc:"获取部门列表"`
	DeptName        string `json:"deptName" doc:"部门名称"`
	Status          string `json:"status" doc:"状态"`
}

type SysDeptTreeReq struct {
	mvc.RequestBody `route:"/sys/dept/tree" doc:"获取部门树"`
}

type SysDeptSaveReq struct {
	mvc.RequestBody `route:"/sys/dept/save" doc:"保存部门"`
	ID              uint   `json:"id" doc:"部门ID（更新时必填）"`
	ParentID        uint   `json:"parentId" doc:"父部门ID" default:"0"`
	Name            string `json:"name" doc:"部门名称" binding:"required"`
	Sort            int    `json:"sort" doc:"显示顺序" default:"1"`
	Leader          string `json:"leader" doc:"负责人"`
	Phone           string `json:"phone" doc:"联系电话"`
	Email           string `json:"email" doc:"邮箱"`
	Status          int    `json:"status" doc:"状态" default:"1"`
}

type SysDeptDetailReq struct {
	mvc.RequestBody `route:"/sys/dept/detail" doc:"获取部门详情"`
	ID              uint `json:"id" doc:"部门ID"`
}

type SysDeptDelReq struct {
	mvc.RequestBody `route:"/sys/dept/del" doc:"删除部门"`
	ID              uint `json:"id" doc:"部门ID" binding:"required"`
}

// ==================== 部门管理响应 ====================
type SysDeptResp struct {
	BaseResp
	ParentID uint          `json:"parentId" doc:"父部门ID"`
	Name     string        `json:"name" doc:"部门名称"`
	Path     string        `json:"path" doc:"部门路径"`
	Sort     int           `json:"sort" doc:"显示顺序"`
	Leader   string        `json:"leader" doc:"负责人"`
	Phone    string        `json:"phone" doc:"联系电话"`
	Email    string        `json:"email" doc:"邮箱"`
	Status   int           `json:"status" doc:"状态"`
	Children []SysDeptResp `json:"children" doc:"子部门"`
}

type SysDeptTreeResp struct {
	ID       uint              `json:"id" doc:"部门ID"`
	ParentID uint              `json:"parentId" doc:"父部门ID"`
	Name     string            `json:"label" doc:"部门名称"`
	Children []SysDeptTreeResp `json:"children" doc:"子部门"`
}

// ==================== 获取用户信息响应 ====================
type UserInfoResp struct {
	ID          uint     `json:"id" doc:"用户ID"`
	Username    string   `json:"username" doc:"用户名"`
	Nickname    string   `json:"nickname" doc:"昵称"`
	Avatar      string   `json:"avatar" doc:"头像"`
	Email       string   `json:"email" doc:"邮箱"`
	Phone       string   `json:"phone" doc:"手机号"`
	DeptId      uint     `json:"deptId" doc:"部门ID"`
	DeptName    string   `json:"deptName" doc:"部门名称"`
	Roles       []string `json:"roles" doc:"角色列表"`
	Permissions []string `json:"permissions" doc:"权限列表"`
}

// ==================== 获取路由响应 ====================
type RouteResp struct {
	Path      string      `json:"path" doc:"路由路径"`
	Name      string      `json:"name" doc:"路由名称"`
	Component string      `json:"component" doc:"组件路径"`
	Redirect  string      `json:"redirect" doc:"重定向路径"`
	Meta      RouteMeta   `json:"meta" doc:"路由元信息"`
	Children  []RouteResp `json:"children" doc:"子路由"`
}

type RouteMeta struct {
	Title       string   `json:"title" doc:"菜单标题"`
	Icon        string   `json:"icon" doc:"菜单图标"`
	IsLink      string   `json:"isLink" doc:"外链地址"`
	IsHide      bool     `json:"isHide" doc:"是否隐藏"`
	IsFull      bool     `json:"isFull" doc:"是否全屏"`
	IsAffix     bool     `json:"isAffix" doc:"是否固定"`
	IsKeepAlive bool     `json:"isKeepAlive" doc:"是否缓存"`
	Auths       []string `json:"auths" doc:"权限列表"`
	MenuType    int      `json:"menuType" doc:"菜单类型"`
}

// ==================== 系统配置管理 ====================

type SystemConfigGetReq struct {
	mvc.RequestBody `route:"/webhook/GetSystemConfig" doc:"获取系统配置"`
	ConfigKey       string `json:"configKey" doc:"配置键" binding:"required"`
}

type SystemConfigSaveReq struct {
	mvc.RequestBody `route:"/webhook/SaveSystemConfig" doc:"保存系统配置"`
	ConfigKey       string `json:"configKey" doc:"配置键" binding:"required"`
	ConfigValue     string `json:"configValue" doc:"配置值" binding:"required"`
	Description     string `json:"description" doc:"描述"`
}

type SystemConfigResp struct {
	ConfigKey   string `json:"configKey" doc:"配置键"`
	ConfigValue string `json:"configValue" doc:"配置值"`
	Description string `json:"description" doc:"描述"`
}
