package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	context2 "github.com/liangboceo/yuanboot/web/context"
	"gorm.io/gorm"
	"sendex-server/internal/model"
	utils2 "sendex-server/pkg/utils"
)

type SysService struct {
	DbService *DbService
	Log       xlog.ILogger
	Cache     *CacheService
	Config    abstractions.IConfiguration
}

func NewSysService(db *DbService, cache *CacheService, config abstractions.IConfiguration) *SysService {
	sys := &SysService{DbService: db, Log: xlog.GetXLogger("SysService"), Cache: cache, Config: config}
	err := sys.InitSystemData()
	if err != nil {
		sys.Log.Errorf("sys init system data err %v", err)
	}
	return sys
}

// ==================== 初始化系统数据 ====================

func (s *SysService) InitSystemData() error {
	// 初始化管理员账号
	if err := s.initAdminUser(); err != nil {
		return fmt.Errorf("初始化管理员账号失败: %v", err)
	}

	// 初始化菜单数据
	if err := s.initMenus(); err != nil {
		return fmt.Errorf("初始化菜单数据失败: %v", err)
	}

	// 初始化角色
	if err := s.initRoles(); err != nil {
		return fmt.Errorf("初始化角色失败: %v", err)
	}

	return nil
}

// 初始化管理员账号
func (s *SysService) initAdminUser() error {
	var count int64
	s.DbService.Db.Model(&model.SysUser{}).Count(&count)
	if count > 0 {
		s.Log.Info("管理员账号已存在，跳过初始化")
		return nil
	}

	// 创建管理员账号
	adminUser := &model.SysUser{
		Username:  "admin",
		Password:  s.EncryptPassword("admin123"),
		Nickname:  "管理员",
		Status:    1,
		LoginDate: time.Now(),
		IsAdmin:   1,
	}

	if err := s.DbService.Db.Create(adminUser).Error; err != nil {
		return err
	}

	s.Log.Info("管理员账号初始化成功：admin/admin123")
	return nil
}

// 初始化菜单数据
func (s *SysService) initMenus() error {
	return s.syncSystemMenus()
}

func (s *SysService) syncSystemMenus() error {
	for _, seed := range systemMenuSeeds() {
		menu := seed.toModel()
		var existing model.SysMenu
		err := s.DbService.Db.Where("id = ?", seed.ID).First(&existing).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err := s.DbService.Db.Create(&menu).Error; err != nil {
				return err
			}
			continue
		}

		updates := map[string]any{
			"name":       seed.Name,
			"parent_id":  seed.ParentID,
			"order_num":  seed.OrderNum,
			"path":       seed.Path,
			"component":  seed.Component,
			"menu_type":  seed.MenuType,
			"visible":    seed.Visible,
			"status":     seed.Status,
			"icon":       seed.Icon,
			"is_frame":   seed.IsFrame,
			"is_cache":   seed.IsCache,
			"permission": seed.Permission,
			"query":      seed.Query,
			"perms":      seed.Perms,
			"is_deleted": 0,
		}
		if err := s.DbService.Db.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
	}

	s.Log.Info("系统菜单同步完成")
	return nil
}

// 初始化角色
func (s *SysService) initRoles() error {
	var adminRole model.SysRole
	if err := s.DbService.Db.Where("code = ? and is_deleted = ?", "admin", 0).First(&adminRole).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		adminRole = model.SysRole{
			Name:      "超级管理员",
			Code:      "admin",
			Sort:      1,
			DataScope: 1,
			Status:    1,
			Remark:    "系统超级管理员，拥有所有权限",
		}
		if err := s.DbService.Db.Create(&adminRole).Error; err != nil {
			return err
		}
	}

	if err := s.syncAdminRoleMenus(adminRole.ID); err != nil {
		return err
	}

	var adminUser model.SysUser
	if err := s.DbService.Db.Where("username = ? and is_deleted = ?", "admin", 0).First(&adminUser).Error; err != nil {
		return err
	}

	var userRoleCount int64
	if err := s.DbService.Db.Model(&model.SysUserRole{}).Where("user_id = ?", adminUser.ID).Count(&userRoleCount).Error; err != nil {
		return err
	}
	if userRoleCount == 0 {
		if err := s.DbService.Db.Create(&model.SysUserRole{
			UserID: adminUser.ID,
			RoleID: adminRole.ID,
		}).Error; err != nil {
			return err
		}
	}

	s.Log.Info("角色数据初始化完成：admin 角色")
	return nil
}

func (s *SysService) syncAdminRoleMenus(adminRoleId uint) error {
	var allMenus []model.SysMenu
	if err := s.DbService.Db.Select("id").Where("is_deleted = ?", 0).Find(&allMenus).Error; err != nil {
		return err
	}
	for _, menu := range allMenus {
		var count int64
		if err := s.DbService.Db.Model(&model.SysRoleMenu{}).Where("role_id = ? and menu_id = ?", adminRoleId, menu.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := s.DbService.Db.Create(&model.SysRoleMenu{
			RoleID: adminRoleId,
			MenuID: menu.ID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// ==================== 密码加密 ====================

func (s *SysService) EncryptPassword(password string) string {
	hash := sha256.Sum256([]byte(password + "sendex-salt"))
	return hex.EncodeToString(hash[:])
}

func (s *SysService) VerifyPassword(password, encrypted string) bool {
	return s.EncryptPassword(password) == encrypted
}

// ==================== 用户管理 ====================

func (s *SysService) GetUserByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SysService) GetUserByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	if err := s.DbService.Db.Where("username=? and is_deleted=?", username, 0).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SysService) GetLoginUserByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	if err := s.DbService.Db.Where("username=? and is_deleted=?", username, 0).First(&user).Error; err != nil {
		return nil, err
	}
	_ = s.Cache.PutSysUserLoginInfo(username, &user, time.Hour*12)
	return &user, nil
}

func (s *SysService) DelLoginUserCache(username string) {
	if username != "" {
		s.Cache.DelSysUserLoginInfo(username)
	}
}

func (s *SysService) GetUserList(username, nickname, phone string, status string, deptId string, pageNum, pageSize int) ([]*model.SysUser, int64, error) {
	var users []*model.SysUser
	var total int64

	query := s.DbService.Db.Model(&model.SysUser{}).Where("is_deleted=?", 0)

	if username != "" {
		query = query.Where("username like ?", "%"+username+"%")
	}
	if nickname != "" {
		query = query.Where("nickname like ?", "%"+nickname+"%")
	}
	if phone != "" {
		query = query.Where("phone like ?", "%"+phone+"%")
	}
	if status != "" {
		query = query.Where("status=?", status)
	}
	if deptId != "" {
		query = query.Where("dept_id=?", deptId)
	}

	query.Count(&total)
	err := query.Offset((pageNum - 1) * pageSize).Order("id desc").Limit(pageSize).Find(&users).Error

	return users, total, err
}

func (s *SysService) CreateUser(user *model.SysUser, roleIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 绑定角色
		for _, roleId := range roleIds {
			if err := tx.Create(&model.SysUserRole{
				UserID: user.ID,
				RoleID: roleId,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *SysService) UpdateUser(user *model.SysUser, roleIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		// 更新用户信息
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		if roleIds == nil {
			return nil
		}

		// 删除原有角色关联
		if err := tx.Where("user_id=?", user.ID).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		for _, roleId := range roleIds {
			if err := tx.Create(&model.SysUserRole{
				UserID: user.ID,
				RoleID: roleId,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *SysService) DeleteUser(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}
	if err := s.DbService.Db.Model(&model.SysUser{}).Where("id=?", id).Update("is_deleted", 1).Error; err != nil {
		return err
	}
	s.DelLoginUserCache(user.Username)
	return nil
}

func (s *SysService) GetUserRoles(userId uint) ([]uint, error) {
	var userRoles []model.SysUserRole
	var roleIds []uint

	if err := s.DbService.Db.Where("user_id=?", userId).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleID)
	}

	return roleIds, nil
}

func (s *SysService) GetUserRoleNames(userId uint) (string, error) {
	var roleNames []string
	err := s.DbService.Db.Table("sys_role").
		Select("sys_role.name").
		Joins("LEFT JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ? AND sys_role.is_deleted = 0", userId).
		Pluck("name", &roleNames).Error

	if err != nil {
		return "", err
	}

	return strings.Join(roleNames, ","), nil
}

// ==================== 角色管理 ====================

func (s *SysService) GetRoleByID(id uint) (*model.SysRole, error) {
	var role model.SysRole
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *SysService) GetRoleList(name, code string, status string, pageNum, pageSize int) ([]*model.SysRole, int64, error) {
	var roles []*model.SysRole
	var total int64

	query := s.DbService.Db.Model(&model.SysRole{}).Where("is_deleted=?", 0)

	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code like ?", "%"+code+"%")
	}
	if status != "" {
		query = query.Where("status=?", status)
	}

	query.Count(&total)
	err := query.Offset((pageNum - 1) * pageSize).Order("sort asc, id desc").Limit(pageSize).Find(&roles).Error

	return roles, total, err
}

func (s *SysService) GetAllRoles() ([]*model.SysRole, error) {
	var roles []*model.SysRole
	err := s.DbService.Db.Where("status=1 and is_deleted=0").Order("sort asc").Find(&roles).Error
	return roles, err
}

func (s *SysService) CreateRole(role *model.SysRole, menuIds, deptIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		// 绑定菜单
		for _, menuId := range menuIds {
			if err := tx.Create(&model.SysRoleMenu{
				RoleID: role.ID,
				MenuID: menuId,
			}).Error; err != nil {
				return err
			}
		}

		// 绑定部门
		for _, deptId := range deptIds {
			if err := tx.Create(&model.SysRoleDept{
				RoleID: role.ID,
				DeptID: deptId,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *SysService) UpdateRole(role *model.SysRole, menuIds, deptIds []uint) error {
	return s.DbService.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if menuIds != nil {
			// 删除原有菜单关联
			if err := tx.Where("role_id=?", role.ID).Delete(&model.SysRoleMenu{}).Error; err != nil {
				return err
			}

			// 添加新的菜单关联
			for _, menuId := range menuIds {
				if err := tx.Create(&model.SysRoleMenu{
					RoleID: role.ID,
					MenuID: menuId,
				}).Error; err != nil {
					return err
				}
			}
		}

		if deptIds != nil {
			// 删除原有部门关联
			if err := tx.Where("role_id=?", role.ID).Delete(&model.SysRoleDept{}).Error; err != nil {
				return err
			}

			// 添加新的部门关联
			for _, deptId := range deptIds {
				if err := tx.Create(&model.SysRoleDept{
					RoleID: role.ID,
					DeptID: deptId,
				}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *SysService) DeleteRole(id uint) error {
	// 检查是否有用户使用该角色
	var count int64
	s.DbService.Db.Model(&model.SysUserRole{}).Where("role_id=?", id).Count(&count)
	if count > 0 {
		return errors.New("该角色已分配给用户，无法删除")
	}

	return s.DbService.Db.Model(&model.SysRole{}).Where("id=?", id).Update("is_deleted", 1).Error
}

func (s *SysService) GetRoleMenus(roleId uint) ([]uint, error) {
	var roleMenus []model.SysRoleMenu
	var menuIds []uint

	if err := s.DbService.Db.Where("role_id=?", roleId).Find(&roleMenus).Error; err != nil {
		return nil, err
	}

	for _, rm := range roleMenus {
		menuIds = append(menuIds, rm.MenuID)
	}

	return menuIds, nil
}

// ==================== 菜单管理 ====================

func (s *SysService) GetMenuByID(id uint) (*model.SysMenu, error) {
	var menu model.SysMenu
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (s *SysService) GetMenuList(menuName string, status int) ([]*model.SysMenu, error) {
	var menus []*model.SysMenu

	query := s.DbService.Db.Model(&model.SysMenu{}).Where("is_deleted=?", 0)

	if menuName != "" {
		query = query.Where("name like ?", "%"+menuName+"%")
	}
	if status > 0 {
		query = query.Where("status=?", status)
	}

	err := query.Order("parent_id asc, order_num asc").Find(&menus).Error
	return menus, err
}

func (s *SysService) GetMenuTree() ([]*model.SysMenu, error) {
	var menus []*model.SysMenu
	err := s.DbService.Db.Where("is_deleted=0 and status=1").Order("parent_id asc, order_num asc").Find(&menus).Error
	return menus, err
}

func (s *SysService) GetMenuByIds(menuIds []uint) ([]*model.SysMenu, error) {
	if len(menuIds) == 0 {
		return []*model.SysMenu{}, nil
	}

	idMap := make(map[uint]struct{})
	pendingIds := menuIds
	for len(pendingIds) > 0 {
		var menus []*model.SysMenu
		if err := s.DbService.Db.Where("id in ? and is_deleted=0", pendingIds).Find(&menus).Error; err != nil {
			return nil, err
		}

		pendingIds = nil
		for _, menu := range menus {
			if _, ok := idMap[menu.ID]; ok {
				continue
			}
			idMap[menu.ID] = struct{}{}
			if menu.ParentID > 0 {
				if _, ok := idMap[menu.ParentID]; !ok {
					pendingIds = append(pendingIds, menu.ParentID)
				}
			}
		}
	}

	allIds := make([]uint, 0, len(idMap))
	for id := range idMap {
		allIds = append(allIds, id)
	}

	var menus []*model.SysMenu
	err := s.DbService.Db.Where("id in ? and is_deleted=0", allIds).Order("parent_id asc, order_num asc, id asc").Find(&menus).Error
	return menus, err
}

func (s *SysService) CreateMenu(menu *model.SysMenu) error {
	return s.DbService.Db.Create(menu).Error
}

func (s *SysService) UpdateMenu(menu *model.SysMenu) error {
	return s.DbService.Db.Save(menu).Error
}

func (s *SysService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	s.DbService.Db.Model(&model.SysMenu{}).Where("parent_id=? and is_deleted=0", id).Count(&count)
	if count > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	// 删除菜单关联
	s.DbService.Db.Where("menu_id=?", id).Delete(&model.SysRoleMenu{})

	return s.DbService.Db.Model(&model.SysMenu{}).Where("id=?", id).Update("is_deleted", 1).Error
}

func (s *SysService) GetUserPermissions(userId uint) ([]string, error) {
	var permissions []string

	// 获取用户的所有菜单权限
	err := s.DbService.Db.Table("sys_menu").
		Select("sys_menu.permission").
		Joins("LEFT JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Joins("LEFT JOIN sys_user_role ON sys_user_role.role_id = sys_role_menu.role_id").
		Where("sys_user_role.user_id = ? AND sys_menu.is_deleted = 0 AND sys_menu.permission != '' AND sys_menu.permission IS NOT NULL", userId).
		Distinct("permission").
		Pluck("permission", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// ==================== 部门管理 ====================

func (s *SysService) GetDeptByID(id uint) (*model.SysDept, error) {
	var dept model.SysDept
	if err := s.DbService.Db.Where("id=? and is_deleted=?", id, 0).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (s *SysService) GetDeptList(deptName string, status string) ([]*model.SysDept, error) {
	var depts []*model.SysDept

	query := s.DbService.Db.Model(&model.SysDept{}).Where("is_deleted=?", 0)

	if deptName != "" {
		query = query.Where("name like ?", "%"+deptName+"%")
	}
	if status != "" {
		query = query.Where("status=?", status)
	}

	err := query.Order("sort asc").Find(&depts).Error
	return depts, err
}

func (s *SysService) CreateDept(dept *model.SysDept) error {
	return s.DbService.Db.Create(dept).Error
}

func (s *SysService) UpdateDept(dept *model.SysDept) error {
	return s.DbService.Db.Save(dept).Error
}

func (s *SysService) DeleteDept(id uint) error {
	// 检查是否有子部门
	var count int64
	s.DbService.Db.Model(&model.SysDept{}).Where("parent_id=? and is_deleted=0", id).Count(&count)
	if count > 0 {
		return errors.New("存在子部门，无法删除")
	}

	return s.DbService.Db.Model(&model.SysDept{}).Where("id=?", id).Update("is_deleted", 1).Error
}

// ==================== 辅助函数 ====================

func (s *SysService) BuildMenuTree(menus []*model.SysMenu) []*model.SysMenu {
	menuMap := make(map[uint]*model.SysMenu)
	var roots []*model.SysMenu

	// 先把所有菜单放入map
	for _, menu := range menus {
		menu.Children = nil
		menuMap[menu.ID] = menu
	}

	// 构建树形结构
	for _, menu := range menus {
		if menu.ParentID == 0 {
			roots = append(roots, menu)
		} else {
			if parent, ok := menuMap[menu.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []*model.SysMenu{}
				}
				parent.Children = append(parent.Children, menu)
			}
		}
	}

	return roots
}

func (s *SysService) BuildDeptTree(depts []*model.SysDept) []*model.SysDept {
	deptMap := make(map[uint]*model.SysDept)
	var roots []*model.SysDept

	for _, dept := range depts {
		dept.Children = nil
		deptMap[dept.ID] = dept
	}

	for _, dept := range depts {
		if dept.ParentID == 0 {
			roots = append(roots, dept)
		} else {
			if parent, ok := deptMap[dept.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []*model.SysDept{}
				}
				parent.Children = append(parent.Children, dept)
			}
		}
	}

	return roots
}

// ──────────── System Config ────────────

const (
	ConfigKeyWorkOrderWebhookIDs = "work_order_webhook_ids"
	ConfigKeyAgentMatchModelID   = "agent_match_model_id" // Agent 匹配使用的模型配置ID
	ConfigKeySystemName          = "system_name"          // 系统名称
	ConfigKeySystemLogo          = "system_logo"          // 系统Logo(base64)
	ConfigKeySystemTitle         = "system_title"         // 系统标题
	ConfigKeySystemDescription   = "system_description"   // 系统宣传语
)

func (s *SysService) GetSystemConfig(configKey string) (*model.SystemConfig, error) {
	// 缓存优先读取
	cacheKey := fmt.Sprintf(SYSTEM_CONFIG, configKey)
	var config model.SystemConfig
	if s.Cache != nil && s.Cache.CacheGetJSON(cacheKey, &config) {
		return &config, nil
	}

	if err := s.DbService.Db.Where("config_key=? and is_deleted=?", configKey, 0).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 回写缓存
	if s.Cache != nil {
		_ = s.Cache.CacheSetJSON(cacheKey, &config, 10*time.Minute)
	}
	return &config, nil
}

func (s *SysService) SaveSystemConfig(config *model.SystemConfig, ctx *context2.HttpContext) error {
	// Upsert: 查找是否存在
	var existing model.SystemConfig
	err := s.DbService.Db.Where("config_key=? and is_deleted=?", config.ConfigKey, 0).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 新建
		config.CreateTime = time.Now()
		config.CreateBy = utils2.GetUserId(ctx)
		config.CreateName = utils2.GetUserName(ctx)
		if err := s.DbService.Db.Create(config).Error; err != nil {
			return err
		}
	} else {
		// 更新
		config.ID = existing.ID
		config.CreateTime = existing.CreateTime
		config.UpdateTime = time.Now()
		config.UpdateBy = utils2.GetUserId(ctx)
		config.UpdateName = utils2.GetUserName(ctx)
		if err := s.DbService.Db.Save(config).Error; err != nil {
			return err
		}
	}

	// 使缓存失效
	if s.Cache != nil {
		s.Cache.CacheDelete(fmt.Sprintf(SYSTEM_CONFIG, config.ConfigKey))
	}
	return nil
}
