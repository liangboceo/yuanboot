package controller

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/utils/jwt"
	"github.com/liangboceo/yuanboot/web/actionresult"
	context2 "github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/liangboceo/yuanboot/web/mvc"
	"sendex-server/internal/dto"
	"sendex-server/internal/middleware"
	"sendex-server/internal/model"
	"sendex-server/internal/response"
	"sendex-server/internal/service"
	utils2 "sendex-server/pkg/utils"
)

type SysController struct {
	mvc.ApiController `doc:"系统管理"`
	log               *middlewares.Logger
	SysService        *service.SysService
	CacheService      *service.CacheService
	Config            abstractions.IConfiguration
}

func NewSysController(sysService *service.SysService, cacheService *service.CacheService, config abstractions.IConfiguration) *SysController {
	controller := &SysController{SysService: sysService, CacheService: cacheService, Config: config}
	controller.log = middlewares.NewLogger()
	return controller
}

func (c *SysController) lockSysOperation(module, key string, id uint) (string, bool) {
	lockKey := fmt.Sprintf(service.DATALOCK, module, key, strconv.FormatUint(uint64(id), 10))
	return lockKey, c.CacheService.Lock(lockKey, 10)
}

// ==================== 登录接口 ====================

// Captcha 获取图形验证码
func (c *SysController) Captcha(req *dto.CaptchaReq) actionresult.IActionResult {
	captchaId := uuid.NewString()
	code := utils2.GenerateCaptchaCode()
	imageData, err := utils2.GenerateCaptchaImage(code)
	if err != nil {
		return response.FailureMessage(nil, "生成验证码失败")
	}
	if err := c.CacheService.PutCaptcha(captchaId, strings.ToLower(code), time.Minute*2); err != nil {
		return response.FailureMessage(nil, "保存验证码失败")
	}

	return response.SuccessMessage(dto.CaptchaResp{
		CaptchaId: captchaId,
		Image:     imageData,
	}, "获取验证码成功")
}

// Login 登录
func (c *SysController) Login(req *dto.LoginReq) actionresult.IActionResult {
	if err := validator.New().Struct(req); err != nil {
		return response.FailureMessage(nil, "用户名、密码或验证码不能为空")
	}

	cachedCode, err := c.CacheService.GetCaptcha(req.CaptchaId)
	if err != nil || cachedCode == "" {
		return response.FailureMessage(nil, "验证码已过期")
	}
	if strings.ToLower(req.VerifyCode) != cachedCode {
		return response.FailureMessage(nil, "验证码错误")
	}
	c.CacheService.DelCaptcha(req.CaptchaId)

	user, err := c.SysService.GetLoginUserByUsername(req.Username)
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	if !c.SysService.VerifyPassword(req.Password, user.Password) {
		return response.FailureMessage(nil, "密码错误")
	}

	if user.Status != 1 {
		return response.FailureMessage(nil, "账号已被禁用")
	}

	// 获取用户角色
	roleIds, _ := c.SysService.GetUserRoles(user.ID)
	var roles []string
	for _, roleId := range roleIds {
		role, _ := c.SysService.GetRoleByID(roleId)
		if role != nil {
			roles = append(roles, role.Code)
		}
	}

	// 获取用户权限
	permissions, _ := c.SysService.GetUserPermissions(user.ID)
	var hasSecretKey bool
	var secretKey string
	if c.Config != nil {
		secretKey, hasSecretKey = c.Config.Get("yuanboot.application.server.uas.auth.jwt-secret").(string)
	}

	if !hasSecretKey {
		secretKey = "5Zk2Qx8LpW7rT3eY9uB1vF4sH6dG2jK8mN3bV7cX1zA9sD4fG7hJ2kL5pR8tY3"
	}
	// 生成token (简化处理，实际应该使用JWT库)
	claim := &middleware.JwtCustomClaims{
		UserName:  user.Username,
		UserId:    user.ID,
		IsAdmin:   user.IsAdmin,
		TokenType: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(12)).Unix(),
			Issuer:    user.Username,
		},
	}

	token, _ := jwt.CreateCustomToken([]byte(secretKey), claim)
	resp := dto.LoginResp{
		AccessToken:  token,
		Expires:      claim.ExpiresAt,
		RefreshToken: token + "-refresh",
		Avatar:       user.Avatar,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Roles:        roles,
		Permissions:  permissions,
	}

	return response.SuccessMessage(resp, "登录成功")
}

// GetUserInfo 获取当前用户信息
func (c *SysController) GetUserInfo(request *struct {
	mvc.RequestBody `route:"/sys/GetUserInfo" doc:"获取当前用户信息"`
}, ctx *context2.HttpContext) actionresult.IActionResult {
	userId := utils2.GetUserId(ctx)
	user, err := c.SysService.GetUserByID(uint(userId))
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	// 获取用户角色
	roleIds, _ := c.SysService.GetUserRoles(user.ID)
	var roles []string
	for _, roleId := range roleIds {
		role, _ := c.SysService.GetRoleByID(roleId)
		if role != nil {
			roles = append(roles, role.Code)
		}
	}

	// 获取用户权限
	permissions, _ := c.SysService.GetUserPermissions(user.ID)

	// 获取部门名称
	var deptName string
	if user.DeptID > 0 {
		dept, _ := c.SysService.GetDeptByID(user.DeptID)
		if dept != nil {
			deptName = dept.Name
		}
	}

	resp := dto.UserInfoResp{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		DeptId:      user.DeptID,
		DeptName:    deptName,
		Roles:       roles,
		Permissions: permissions,
	}

	return response.SuccessMessage(resp, "获取用户信息成功")
}

// GetRoutes 获取用户路由权限
func (c *SysController) GetRoutes(request *struct {
	mvc.RequestBody `route:"/sys/GetRoutes" doc:"获取用户路由权限"`
}, ctx *context2.HttpContext) actionresult.IActionResult {
	userId := utils2.GetUserId(ctx)

	// 获取用户角色
	roleIds, _ := c.SysService.GetUserRoles(uint(userId))

	var allMenus []*model.SysMenu
	for _, roleId := range roleIds {
		menuIds, _ := c.SysService.GetRoleMenus(roleId)
		menus, _ := c.SysService.GetMenuByIds(menuIds)
		allMenus = append(allMenus, menus...)
	}

	// 去重
	menuMap := make(map[uint]*model.SysMenu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}

	var uniqueMenus []*model.SysMenu
	for _, menu := range menuMap {
		uniqueMenus = append(uniqueMenus, menu)
	}
	sort.SliceStable(uniqueMenus, func(i, j int) bool {
		if uniqueMenus[i].ParentID == uniqueMenus[j].ParentID {
			if uniqueMenus[i].OrderNum == uniqueMenus[j].OrderNum {
				return uniqueMenus[i].ID < uniqueMenus[j].ID
			}
			return uniqueMenus[i].OrderNum < uniqueMenus[j].OrderNum
		}
		return uniqueMenus[i].ParentID < uniqueMenus[j].ParentID
	})

	// 构建路由树
	routes := c.buildRoutes(uniqueMenus, 0)

	return response.SuccessMessage(routes, "获取路由成功")
}

func (c *SysController) buildRoutes(menus []*model.SysMenu, parentId uint) []dto.RouteResp {
	var routes []dto.RouteResp

	for _, menu := range menus {
		if menu.ParentID != parentId {
			continue
		}

		route := dto.RouteResp{
			Path: menu.Path,
			Name: menu.Name,
			Meta: dto.RouteMeta{
				Title:       menu.Name,
				Icon:        menu.Icon,
				IsHide:      menu.Visible == 1,
				IsKeepAlive: menu.IsCache == 1,
				MenuType:    menu.MenuType,
			},
		}

		if menu.Component != "" && (menu.MenuType == 1 || menu.MenuType == 3) {
			route.Component = menu.Component
		}

		children := c.buildRoutes(menus, menu.ID)
		if len(children) > 0 {
			if menu.MenuType == 0 {
				route.Redirect = children[0].Path
			}
			route.Children = children
		}

		routes = append(routes, route)
	}

	return routes
}

// ==================== 用户管理 ====================

// GetUserList 获取用户列表
func (c *SysController) GetUserList(req *dto.SysUserListReq) actionresult.IActionResult {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	users, total, err := c.SysService.GetUserList(req.Username, req.Nickname, req.Phone, req.Status, req.DeptId, req.PageNum, req.PageSize)
	if err != nil {
		return response.FailureMessage(nil, "获取用户列表失败")
	}

	result := &dto.SysUserListResp{
		BasePageResp: dto.BasePageResp{Page: req.PageNum, Total: total},
	}
	result.Data = make([]dto.SysUserResp, len(users))

	for i, user := range users {
		roleNames, _ := c.SysService.GetUserRoleNames(user.ID)
		roleIds, _ := c.SysService.GetUserRoles(user.ID)

		var deptName string
		if user.DeptID > 0 {
			dept, _ := c.SysService.GetDeptByID(user.DeptID)
			if dept != nil {
				deptName = dept.Name
			}
		}

		result.Data[i] = dto.SysUserResp{
			BaseResp:  dto.BaseResp{ID: user.ID, CreateTime: user.CreateTime, UpdateTime: user.UpdateTime},
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			DeptId:    user.DeptID,
			DeptName:  deptName,
			Status:    user.Status,
			IsAdmin:   user.IsAdmin,
			LoginIp:   user.LoginIp,
			LoginDate: user.LoginDate,
			RoleIds:   roleIds,
			RoleNames: roleNames,
		}
	}

	return response.SuccessMessage(result, "获取用户列表成功")
}

// GetUserDetail 获取用户详情
func (c *SysController) GetUserDetail(req *dto.SysUserDetailReq) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "用户ID不能为空")
	}

	user, err := c.SysService.GetUserByID(req.ID)
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	roleIds, _ := c.SysService.GetUserRoles(user.ID)
	roleNames, _ := c.SysService.GetUserRoleNames(user.ID)

	resp := dto.SysUserResp{
		BaseResp:  dto.BaseResp{ID: user.ID, CreateTime: user.CreateTime, UpdateTime: user.UpdateTime},
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		DeptId:    user.DeptID,
		Status:    user.Status,
		IsAdmin:   user.IsAdmin,
		RoleIds:   roleIds,
		RoleNames: roleNames,
	}

	return response.SuccessMessage(resp, "获取用户详情成功")
}

// SaveUser 保存用户
func (c *SysController) SaveUser(saveReq *dto.SysUserSaveReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(saveReq); err != nil {
		return response.FailureMessage(nil, "用户名不能为空")
	}

	lockKey, lock := c.lockSysOperation("SaveUser", saveReq.Username, saveReq.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if saveReq.ID > 0 {
		// 更新
		user, err := c.SysService.GetUserByID(saveReq.ID)
		if err != nil || user.ID == 0 {
			return response.FailureMessage(nil, "用户不存在")
		}

		user.Nickname = saveReq.Nickname
		user.Email = saveReq.Email
		user.Phone = saveReq.Phone
		user.Avatar = saveReq.Avatar
		user.DeptID = saveReq.DeptId
		user.Status = saveReq.Status
		user.UpdateBy = utils2.GetUserId(ctx)
		user.UpdateName = utils2.GetUserName(ctx)
		user.UpdateTime = time.Now()
		user.LoginDate = time.Now()

		if saveReq.Password != "" {
			user.Password = c.SysService.EncryptPassword(saveReq.Password)
		}

		if err := c.SysService.UpdateUser(user, saveReq.RoleIds); err != nil {
			return response.FailureMessage(nil, "更新用户失败")
		}
		c.SysService.DelLoginUserCache(user.Username)

		return response.SuccessMessage(nil, "更新用户成功")
	}

	// 新增
	user := &model.SysUser{
		Username:  saveReq.Username,
		Password:  c.SysService.EncryptPassword(saveReq.Password),
		Nickname:  saveReq.Nickname,
		Email:     saveReq.Email,
		Phone:     saveReq.Phone,
		Avatar:    saveReq.Avatar,
		DeptID:    saveReq.DeptId,
		LoginDate: time.Now(),
		Status:    saveReq.Status,
		CommonModel: model.CommonModel{
			CreateTime: time.Now(),
			CreateBy:   utils2.GetUserId(ctx),
			CreateName: utils2.GetUserName(ctx),
		},
	}

	if err := c.SysService.CreateUser(user, saveReq.RoleIds); err != nil {
		return response.FailureMessage(nil, "创建用户失败")
	}

	return response.SuccessMessage(nil, "创建用户成功")
}

// DelUser 删除用户
func (c *SysController) DelUser(req *dto.SysUserDelReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "用户ID不能为空")
	}

	lockKey, lock := c.lockSysOperation("DelUser", "user", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	user, err := c.SysService.GetUserByID(req.ID)
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	if user.IsAdmin == 1 {
		return response.FailureMessage(nil, "不能删除超级管理员")
	}

	user.IsDeleted = 1
	user.UpdateBy = utils2.GetUserId(ctx)
	user.UpdateName = utils2.GetUserName(ctx)
	user.UpdateTime = time.Now()

	if err := c.SysService.DeleteUser(req.ID); err != nil {
		return response.FailureMessage(nil, "删除用户失败")
	}

	return response.SuccessMessage(nil, "删除用户成功")
}

// ResetPwd 重置密码
func (c *SysController) ResetPwd(req *dto.SysUserResetPwdReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(req); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	lockKey, lock := c.lockSysOperation("ResetPwd", "user", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	user, err := c.SysService.GetUserByID(req.ID)
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	user.Password = c.SysService.EncryptPassword(req.Password)
	user.UpdateBy = utils2.GetUserId(ctx)
	user.UpdateName = utils2.GetUserName(ctx)
	user.UpdateTime = time.Now()

	if err := c.SysService.UpdateUser(user, nil); err != nil {
		return response.FailureMessage(nil, "重置密码失败")
	}

	return response.SuccessMessage(nil, "重置密码成功")
}

func (c *SysController) ChangePwd(req *dto.SysUserChangePwdReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(req); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	userId := utils2.GetUserId(ctx)
	lockKey, lock := c.lockSysOperation("ChangePwd", "user", uint(userId))
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if !isStrongPassword(req.NewPassword) {
		return response.FailureMessage(nil, "新密码需包含大小写字母和特殊字符，长度不低于8位")
	}

	user, err := c.SysService.GetUserByID(uint(userId))
	if err != nil || user.ID == 0 {
		return response.FailureMessage(nil, "用户不存在")
	}

	if !c.SysService.VerifyPassword(req.OldPassword, user.Password) {
		return response.FailureMessage(nil, "旧密码错误")
	}

	user.Password = c.SysService.EncryptPassword(req.NewPassword)
	user.UpdateBy = userId
	user.UpdateName = utils2.GetUserName(ctx)
	user.UpdateTime = time.Now()

	if err := c.SysService.UpdateUser(user, nil); err != nil {
		return response.FailureMessage(nil, "修改密码失败")
	}
	c.SysService.DelLoginUserCache(user.Username)

	return response.SuccessMessage(nil, "修改密码成功")
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
	return hasUpper && hasLower && hasSpecial
}

// ==================== 角色管理 ====================

// GetRoleList 获取角色列表
func (c *SysController) GetRoleList(req *dto.SysRoleListReq) actionresult.IActionResult {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	roles, total, err := c.SysService.GetRoleList(req.Name, req.Code, req.Status, req.PageNum, req.PageSize)
	if err != nil {
		return response.FailureMessage(nil, "获取角色列表失败")
	}

	result := &dto.SysRoleListResp{
		BasePageResp: dto.BasePageResp{Page: req.PageNum, Total: total},
	}
	result.Data = make([]dto.SysRoleResp, len(roles))

	for i, role := range roles {
		result.Data[i] = dto.SysRoleResp{
			BaseResp:          dto.BaseResp{ID: role.ID, CreateTime: role.CreateTime, UpdateTime: role.UpdateTime},
			Name:              role.Name,
			Code:              role.Code,
			Sort:              role.Sort,
			DataScope:         role.DataScope,
			MenuCheckStrictly: role.MenuCheckStrictly,
			DeptCheckStrictly: role.DeptCheckStrictly,
			Status:            role.Status,
			Remark:            role.Remark,
		}
	}

	return response.SuccessMessage(result, "获取角色列表成功")
}

// GetAllRoleList 获取所有角色列表
func (c *SysController) GetAllRoleList(request *struct {
	mvc.RequestBody `route:"/sys/GetAllRoleList" doc:"获取所有角色列表"`
}) actionresult.IActionResult {
	roles, err := c.SysService.GetAllRoles()
	if err != nil {
		return response.FailureMessage(nil, "获取角色列表失败")
	}

	var result []dto.SysRoleResp
	for _, role := range roles {
		result = append(result, dto.SysRoleResp{
			BaseResp: dto.BaseResp{ID: role.ID},
			Name:     role.Name,
			Code:     role.Code,
			Status:   role.Status,
		})
	}

	return response.SuccessMessage(result, "获取角色列表成功")
}

// GetRoleDetail 获取角色详情
func (c *SysController) GetRoleDetail(req *dto.SysRoleDetailReq) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "角色ID不能为空")
	}

	role, err := c.SysService.GetRoleByID(req.ID)
	if err != nil || role.ID == 0 {
		return response.FailureMessage(nil, "角色不存在")
	}

	menuIds, _ := c.SysService.GetRoleMenus(role.ID)

	resp := dto.SysRoleResp{
		BaseResp:          dto.BaseResp{ID: role.ID, CreateTime: role.CreateTime, UpdateTime: role.UpdateTime},
		Name:              role.Name,
		Code:              role.Code,
		Sort:              role.Sort,
		DataScope:         role.DataScope,
		MenuCheckStrictly: role.MenuCheckStrictly,
		DeptCheckStrictly: role.DeptCheckStrictly,
		Status:            role.Status,
		Remark:            role.Remark,
	}

	return response.SuccessMessage(map[string]interface{}{
		"role":    resp,
		"menuIds": menuIds,
	}, "获取角色详情成功")
}

// SaveRole 保存角色
func (c *SysController) SaveRole(saveReq *dto.SysRoleSaveReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(saveReq); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	lockKey, lock := c.lockSysOperation("SaveRole", saveReq.Code, saveReq.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if saveReq.ID > 0 {
		// 更新
		role, err := c.SysService.GetRoleByID(saveReq.ID)
		if err != nil || role.ID == 0 {
			return response.FailureMessage(nil, "角色不存在")
		}

		role.Name = saveReq.Name
		role.Code = saveReq.Code
		role.Sort = saveReq.Sort
		role.DataScope = saveReq.DataScope
		role.MenuCheckStrictly = saveReq.MenuCheckStrictly
		role.DeptCheckStrictly = saveReq.DeptCheckStrictly
		role.Status = saveReq.Status
		role.Remark = saveReq.Remark
		role.UpdateBy = utils2.GetUserId(ctx)
		role.UpdateName = utils2.GetUserName(ctx)
		role.UpdateTime = time.Now()

		if err := c.SysService.UpdateRole(role, saveReq.MenuIds, saveReq.DeptIds); err != nil {
			return response.FailureMessage(nil, "更新角色失败")
		}

		return response.SuccessMessage(nil, "更新角色成功")
	}

	// 新增
	role := &model.SysRole{
		Name:              saveReq.Name,
		Code:              saveReq.Code,
		Sort:              saveReq.Sort,
		DataScope:         saveReq.DataScope,
		MenuCheckStrictly: saveReq.MenuCheckStrictly,
		DeptCheckStrictly: saveReq.DeptCheckStrictly,
		Status:            saveReq.Status,
		Remark:            saveReq.Remark,
		CommonModel: model.CommonModel{
			CreateBy:   utils2.GetUserId(ctx),
			CreateName: utils2.GetUserName(ctx),
		},
	}

	if err := c.SysService.CreateRole(role, saveReq.MenuIds, saveReq.DeptIds); err != nil {
		return response.FailureMessage(nil, "创建角色失败")
	}

	return response.SuccessMessage(nil, "创建角色成功")
}

// DelRole 删除角色
func (c *SysController) DelRole(req *dto.SysRoleDelReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "角色ID不能为空")
	}

	lockKey, lock := c.lockSysOperation("DelRole", "role", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	role, err := c.SysService.GetRoleByID(req.ID)
	if err != nil || role.ID == 0 {
		return response.FailureMessage(nil, "角色不存在")
	}

	if role.Code == "admin" {
		return response.FailureMessage(nil, "不能删除超级管理员角色")
	}

	role.IsDeleted = 1
	role.UpdateBy = utils2.GetUserId(ctx)
	role.UpdateName = utils2.GetUserName(ctx)
	role.UpdateTime = time.Now()

	if err := c.SysService.DeleteRole(req.ID); err != nil {
		return response.FailureMessage(nil, err.Error())
	}

	return response.SuccessMessage(nil, "删除角色成功")
}

// GetRoleMenu 获取角色菜单权限
func (c *SysController) GetRoleMenu(req *dto.SysRoleMenuReq) actionresult.IActionResult {
	if req.RoleId == 0 {
		return response.FailureMessage(nil, "角色ID不能为空")
	}

	// 获取角色已有的菜单
	checkedKeys, _ := c.SysService.GetRoleMenus(req.RoleId)

	// 获取所有菜单树
	menus, _ := c.SysService.GetMenuList("", 0)
	menuTree := c.SysService.BuildMenuTree(menus)

	resp := dto.SysRoleMenuResp{
		CheckedKeys: checkedKeys,
		MenuList:    c.convertMenuTreeToResp(menuTree),
	}

	return response.SuccessMessage(resp, "获取角色菜单成功")
}

func (c *SysController) convertMenuTreeToResp(menus []*model.SysMenu) []dto.SysMenuResp {
	var result []dto.SysMenuResp
	for _, menu := range menus {
		item := dto.SysMenuResp{
			ID:         menu.ID,
			ParentID:   menu.ParentID,
			Name:       menu.Name,
			OrderNum:   menu.OrderNum,
			Path:       menu.Path,
			Component:  menu.Component,
			MenuType:   menu.MenuType,
			Visible:    menu.Visible,
			Status:     menu.Status,
			Icon:       menu.Icon,
			IsFrame:    menu.IsFrame,
			IsCache:    menu.IsCache,
			Permission: menu.Permission,
			Query:      menu.Query,
			Perms:      menu.Perms,
		}
		if len(menu.Children) > 0 {
			item.Children = c.convertMenuTreeToResp(menu.Children)
		}
		result = append(result, item)
	}
	return result
}

// ==================== 菜单管理 ====================

// GetMenuList 获取菜单列表
func (c *SysController) GetMenuList(req *dto.SysMenuListReq) actionresult.IActionResult {
	menus, err := c.SysService.GetMenuList(req.MenuName, req.Status)
	if err != nil {
		return response.FailureMessage(nil, "获取菜单列表失败")
	}

	menuTree := c.SysService.BuildMenuTree(menus)
	result := c.convertMenuTreeToResp(menuTree)

	return response.SuccessMessage(result, "获取菜单列表成功")
}

// GetMenuTree 获取菜单树
func (c *SysController) GetMenuTree(_ *struct {
	mvc.RequestBody `route:"/sys/GetMenuTree" doc:"获取菜单树"`
}) actionresult.IActionResult {
	menus, err := c.SysService.GetMenuList("", 0)
	if err != nil {
		return response.FailureMessage(nil, "获取菜单树失败")
	}

	var treeResp []dto.SysMenuTreeResp
	for _, menu := range menus {
		treeResp = append(treeResp, c.convertToTreeResp(menu))
	}

	return response.SuccessMessage(treeResp, "获取菜单树成功")
}

func (c *SysController) convertToTreeResp(menu *model.SysMenu) dto.SysMenuTreeResp {
	resp := dto.SysMenuTreeResp{
		ID:       menu.ID,
		ParentID: menu.ParentID,
		Name:     menu.Name,
		OrderNum: menu.OrderNum,
		MenuType: menu.MenuType,
		Icon:     menu.Icon,
	}

	// 获取子菜单
	var children []dto.SysMenuTreeResp
	allMenus, _ := c.SysService.GetMenuList("", 0)
	for _, m := range allMenus {
		if m.ParentID == menu.ID {
			children = append(children, c.convertToTreeResp(m))
		}
	}
	if len(children) > 0 {
		resp.Children = children
	}

	return resp
}

// GetMenuDetail 获取菜单详情
func (c *SysController) GetMenuDetail(req *dto.SysMenuDetailReq) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "菜单ID不能为空")
	}

	menu, err := c.SysService.GetMenuByID(req.ID)
	if err != nil || menu.ID == 0 {
		return response.FailureMessage(nil, "菜单不存在")
	}

	resp := dto.SysMenuResp{
		ID:         menu.ID,
		ParentID:   menu.ParentID,
		Name:       menu.Name,
		OrderNum:   menu.OrderNum,
		Path:       menu.Path,
		Component:  menu.Component,
		MenuType:   menu.MenuType,
		Visible:    menu.Visible,
		Status:     menu.Status,
		Icon:       menu.Icon,
		IsFrame:    menu.IsFrame,
		IsCache:    menu.IsCache,
		Permission: menu.Permission,
		Query:      menu.Query,
		Perms:      menu.Perms,
	}

	return response.SuccessMessage(resp, "获取菜单详情成功")
}

// SaveMenu 保存菜单
func (c *SysController) SaveMenu(saveReq *dto.SysMenuSaveReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(saveReq); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	lockKey, lock := c.lockSysOperation("SaveMenu", saveReq.Name, saveReq.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if saveReq.ID > 0 {
		// 更新
		menu, err := c.SysService.GetMenuByID(saveReq.ID)
		if err != nil || menu.ID == 0 {
			return response.FailureMessage(nil, "菜单不存在")
		}

		menu.Name = saveReq.Name
		menu.ParentID = saveReq.ParentID
		menu.OrderNum = saveReq.OrderNum
		menu.Path = saveReq.Path
		menu.Component = saveReq.Component
		menu.MenuType = saveReq.MenuType
		menu.Visible = saveReq.Visible
		menu.Status = saveReq.Status
		menu.Icon = saveReq.Icon
		menu.IsFrame = saveReq.IsFrame
		menu.IsCache = saveReq.IsCache
		menu.Permission = saveReq.Permission
		menu.Query = saveReq.Query
		menu.Perms = saveReq.Perms
		menu.UpdateBy = utils2.GetUserId(ctx)
		menu.UpdateName = utils2.GetUserName(ctx)
		menu.UpdateTime = time.Now()

		if err := c.SysService.UpdateMenu(menu); err != nil {
			return response.FailureMessage(nil, "更新菜单失败")
		}

		return response.SuccessMessage(nil, "更新菜单成功")
	}

	// 新增
	menu := &model.SysMenu{
		Name:       saveReq.Name,
		ParentID:   saveReq.ParentID,
		OrderNum:   saveReq.OrderNum,
		Path:       saveReq.Path,
		Component:  saveReq.Component,
		MenuType:   saveReq.MenuType,
		Visible:    saveReq.Visible,
		Status:     saveReq.Status,
		Icon:       saveReq.Icon,
		IsFrame:    saveReq.IsFrame,
		IsCache:    saveReq.IsCache,
		Permission: saveReq.Permission,
		Query:      saveReq.Query,
		Perms:      saveReq.Perms,
		CommonModel: model.CommonModel{
			CreateBy:   utils2.GetUserId(ctx),
			CreateName: utils2.GetUserName(ctx),
		},
	}

	if err := c.SysService.CreateMenu(menu); err != nil {
		return response.FailureMessage(nil, "创建菜单失败")
	}

	return response.SuccessMessage(nil, "创建菜单成功")
}

// DelMenu 删除菜单
func (c *SysController) DelMenu(req *dto.SysMenuDelReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "菜单ID不能为空")
	}

	lockKey, lock := c.lockSysOperation("DelMenu", "menu", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	menu, err := c.SysService.GetMenuByID(req.ID)
	if err != nil || menu.ID == 0 {
		return response.FailureMessage(nil, "菜单不存在")
	}

	menu.IsDeleted = 1
	menu.UpdateBy = utils2.GetUserId(ctx)
	menu.UpdateName = utils2.GetUserName(ctx)
	menu.UpdateTime = time.Now()

	if err := c.SysService.DeleteMenu(req.ID); err != nil {
		return response.FailureMessage(nil, err.Error())
	}

	return response.SuccessMessage(nil, "删除菜单成功")
}

// ==================== 部门管理 ====================

// GetDeptList 获取部门列表
func (c *SysController) GetDeptList(req *dto.SysDeptListReq) actionresult.IActionResult {
	depts, err := c.SysService.GetDeptList(req.DeptName, req.Status)
	if err != nil {
		return response.FailureMessage(nil, "获取部门列表失败")
	}

	deptTree := c.SysService.BuildDeptTree(depts)
	result := c.convertDeptTreeToResp(deptTree)

	return response.SuccessMessage(result, "获取部门列表成功")
}

func (c *SysController) convertDeptTreeToResp(depts []*model.SysDept) []dto.SysDeptResp {
	var result []dto.SysDeptResp
	for _, dept := range depts {
		item := dto.SysDeptResp{
			BaseResp: dto.BaseResp{ID: dept.ID, CreateTime: dept.CreateTime, UpdateTime: dept.UpdateTime},
			ParentID: dept.ParentID,
			Name:     dept.Name,
			Path:     dept.Path,
			Sort:     dept.Sort,
			Leader:   dept.Leader,
			Phone:    dept.Phone,
			Email:    dept.Email,
			Status:   dept.Status,
		}
		if len(dept.Children) > 0 {
			item.Children = c.convertDeptTreeToResp(dept.Children)
		}
		result = append(result, item)
	}
	return result
}

// GetDeptTree 获取部门树
func (c *SysController) GetDeptTree(_ *struct {
	mvc.RequestBody `route:"/sys/GetDeptTree" doc:"获取部门树"`
}) actionresult.IActionResult {
	depts, err := c.SysService.GetDeptList("", "1")
	if err != nil {
		return response.FailureMessage(nil, "获取部门树失败")
	}

	var treeResp []dto.SysDeptTreeResp
	for _, dept := range depts {
		treeResp = append(treeResp, c.convertDeptToTreeResp(dept))
	}

	return response.SuccessMessage(treeResp, "获取部门树成功")
}

func (c *SysController) convertDeptToTreeResp(dept *model.SysDept) dto.SysDeptTreeResp {
	resp := dto.SysDeptTreeResp{
		ID:       dept.ID,
		ParentID: dept.ParentID,
		Name:     dept.Name,
	}

	// 获取子部门
	var children []dto.SysDeptTreeResp
	allDepts, _ := c.SysService.GetDeptList("", "1")
	for _, d := range allDepts {
		if d.ParentID == dept.ID {
			children = append(children, c.convertDeptToTreeResp(d))
		}
	}
	if len(children) > 0 {
		resp.Children = children
	}

	return resp
}

// GetDeptDetail 获取部门详情
func (c *SysController) GetDeptDetail(req *dto.SysDeptDetailReq) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "部门ID不能为空")
	}

	dept, err := c.SysService.GetDeptByID(req.ID)
	if err != nil || dept.ID == 0 {
		return response.FailureMessage(nil, "部门不存在")
	}

	resp := dto.SysDeptResp{
		BaseResp: dto.BaseResp{ID: dept.ID, CreateTime: dept.CreateTime, UpdateTime: dept.UpdateTime},
		ParentID: dept.ParentID,
		Name:     dept.Name,
		Path:     dept.Path,
		Sort:     dept.Sort,
		Leader:   dept.Leader,
		Phone:    dept.Phone,
		Email:    dept.Email,
		Status:   dept.Status,
	}

	return response.SuccessMessage(resp, "获取部门详情成功")
}

// SaveDept 保存部门
func (c *SysController) SaveDept(saveReq *dto.SysDeptSaveReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(saveReq); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}

	lockKey, lock := c.lockSysOperation("SaveDept", saveReq.Name, saveReq.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	if saveReq.ID > 0 {
		// 更新
		dept, err := c.SysService.GetDeptByID(saveReq.ID)
		if err != nil || dept.ID == 0 {
			return response.FailureMessage(nil, "部门不存在")
		}

		dept.ParentID = saveReq.ParentID
		dept.Name = saveReq.Name
		dept.Sort = saveReq.Sort
		dept.Leader = saveReq.Leader
		dept.Phone = saveReq.Phone
		dept.Email = saveReq.Email
		dept.Status = saveReq.Status
		dept.UpdateBy = utils2.GetUserId(ctx)
		dept.UpdateName = utils2.GetUserName(ctx)
		dept.UpdateTime = time.Now()

		if err := c.SysService.UpdateDept(dept); err != nil {
			return response.FailureMessage(nil, "更新部门失败")
		}

		return response.SuccessMessage(nil, "更新部门成功")
	}

	// 新增
	dept := &model.SysDept{
		ParentID: saveReq.ParentID,
		Name:     saveReq.Name,
		Sort:     saveReq.Sort,
		Leader:   saveReq.Leader,
		Phone:    saveReq.Phone,
		Email:    saveReq.Email,
		Status:   saveReq.Status,
		CommonModel: model.CommonModel{
			CreateBy:   utils2.GetUserId(ctx),
			CreateName: utils2.GetUserName(ctx),
		},
	}

	if err := c.SysService.CreateDept(dept); err != nil {
		return response.FailureMessage(nil, "创建部门失败")
	}

	return response.SuccessMessage(nil, "创建部门成功")
}

// DelDept 删除部门
func (c *SysController) DelDept(req *dto.SysDeptDelReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if req.ID == 0 {
		return response.FailureMessage(nil, "部门ID不能为空")
	}

	lockKey, lock := c.lockSysOperation("DelDept", "dept", req.ID)
	if !lock {
		return response.FailureMessage(nil, "操作频繁")
	}
	defer c.CacheService.Unlock(lockKey)

	dept, err := c.SysService.GetDeptByID(req.ID)
	if err != nil || dept.ID == 0 {
		return response.FailureMessage(nil, "部门不存在")
	}

	dept.IsDeleted = 1
	dept.UpdateBy = utils2.GetUserId(ctx)
	dept.UpdateName = utils2.GetUserName(ctx)
	dept.UpdateTime = time.Now()

	if err := c.SysService.DeleteDept(req.ID); err != nil {
		return response.FailureMessage(nil, err.Error())
	}

	return response.SuccessMessage(nil, "删除部门成功")
}

// GetSystemConfig 获取系统配置（缓存优先读取）
func (c *SysController) GetSystemConfig(req *dto.SystemConfigGetReq) actionresult.IActionResult {
	config, err := c.SysService.GetSystemConfig(req.ConfigKey)
	if err != nil {
		return response.SuccessMessage("", "没有配置")
	}

	resp := dto.SystemConfigResp{
		ConfigKey:   req.ConfigKey,
		ConfigValue: "",
		Description: "",
	}
	if config != nil {
		resp.ConfigValue = config.ConfigValue
		resp.Description = config.Description
	}

	return response.SuccessMessage(resp, "获取配置成功")
}

// SaveSystemConfig 保存系统配置
func (c *SysController) SaveSystemConfig(req *dto.SystemConfigSaveReq, ctx *context2.HttpContext) actionresult.IActionResult {
	if err := validator.New().Struct(req); err != nil {
		return response.FailureMessage(nil, "参数错误")
	}
	config := &model.SystemConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
	}
	if err := c.SysService.SaveSystemConfig(config, ctx); err != nil {
		c.log.ALogger.Errorf("SaveSystemConfig error: %v", err)
		return response.FailureMessage(nil, "保存配置失败")
	}
	return response.SuccessMessage(nil, "保存配置成功")
}
