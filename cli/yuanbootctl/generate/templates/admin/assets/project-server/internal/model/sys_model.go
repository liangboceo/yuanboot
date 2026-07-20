package model

import "time"

// ==================== 系统用户表 ====================
type SysUser struct {
	ID        uint      `gorm:"type:int(11);primaryKey;column:id;comment:用户id;AUTO_INCREMENT" json:"id"`
	Username  string    `gorm:"type:varchar(50);column:username;comment:用户名;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"type:varchar(100);column:password;comment:密码;not null" json:"password"`
	Nickname  string    `gorm:"type:varchar(50);column:nickname;comment:昵称;default:''" json:"nickname"`
	Email     string    `gorm:"type:varchar(100);column:email;comment:邮箱" json:"email"`
	Phone     string    `gorm:"type:varchar(20);column:phone;comment:手机号" json:"phone"`
	Avatar    string    `gorm:"type:mediumtext;column:avatar;comment:头像（系统头像key或自定义上传base64）" json:"avatar"`
	DeptID    uint      `gorm:"type:int(11);column:dept_id;comment:部门id;default:0" json:"deptId"`
	Status    int       `gorm:"type:tinyint(2);column:status;comment:状态【0：禁用 1：正常】;default:1" json:"status"`
	IsAdmin   int       `gorm:"type:tinyint(2);column:is_admin;comment:是否管理员【0：否 1：是】;default:0" json:"isAdmin"`
	LoginIp   string    `gorm:"type:varchar(50);column:login_ip;comment:最后登录IP;default:''" json:"loginIp"`
	LoginDate time.Time `gorm:"type:datetime;column:login_date;comment:最后登录时间" json:"loginDate"`
	CommonModel
}

func (SysUser) TableName() string {
	return "sys_user"
}

// ==================== 系统角色表 ====================
type SysRole struct {
	BaseModel
	ID                uint   `gorm:"type:int(11);primaryKey;column:id;comment:角色id;AUTO_INCREMENT" json:"id"`
	Name              string `gorm:"type:varchar(50);column:name;comment:角色名称;not null;uniqueIndex" json:"name"`
	Code              string `gorm:"type:varchar(50);column:code;comment:角色编码;not null;uniqueIndex" json:"code"`
	Sort              int    `gorm:"type:int(11);column:sort;comment:排序;default:1" json:"sort"`
	DataScope         int    `gorm:"type:int(2);column:data_scope;comment:数据权限【1：全部数据 2：本部门及以下数据 3：本部门数据 4：仅本人数据 5：自定义】;default:1" json:"dataScope"`
	MenuCheckStrictly bool   `gorm:"column:menu_check_strictly;comment:菜单树选择范围是否关联显示;default:true" json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `gorm:"column:dept_check_strictly;comment:部门树选择范围是否关联显示;default:true" json:"deptCheckStrictly"`
	Status            int    `gorm:"type:tinyint(2);column:status;comment:状态【0：禁用 1：正常】;default:1" json:"status"`
	Remark            string `gorm:"type:varchar(500);column:remark;comment:备注;default:''" json:"remark"`
	CommonModel
}

func (SysRole) TableName() string {
	return "sys_role"
}

// ==================== 系统菜单表 ====================
type SysMenu struct {
	BaseModel
	ID         uint       `gorm:"type:int(11);primaryKey;column:id;comment:菜单id;AUTO_INCREMENT" json:"id"`
	Name       string     `gorm:"type:varchar(50);column:name;comment:菜单名称;not null" json:"name"`
	ParentID   uint       `gorm:"type:int(11);column:parent_id;comment:父菜单id;default:0" json:"parentId"`
	OrderNum   int        `gorm:"type:int(11);column:order_num;comment:显示顺序;default:1" json:"orderNum"`
	Path       string     `gorm:"type:varchar(200);column:path;comment:路由地址;default:''" json:"path"`
	Component  string     `gorm:"type:varchar(255);column:component;comment:组件路径;default:''" json:"component"`
	MenuType   int        `gorm:"type:tinyint(2);column:menu_type;comment:菜单类型【0：目录 1：菜单 2：按钮 3：内页】;default:0" json:"menuType"`
	Visible    int        `gorm:"type:tinyint(2);column:visible;comment:显示状态【0：显示 1：隐藏】;default:0" json:"visible"`
	Status     int        `gorm:"type:tinyint(2);column:status;comment:状态【0：禁用 1：正常】;default:1" json:"status"`
	Icon       string     `gorm:"type:varchar(100);column:icon;comment:菜单图标;default:''" json:"icon"`
	IsFrame    int        `gorm:"type:tinyint(2);column:is_frame;comment:是否外链【0：否 1：是】;default:0" json:"isFrame"`
	IsCache    int        `gorm:"type:tinyint(2);column:is_cache;comment:是否缓存【0：否 1：是】;default:0" json:"isCache"`
	Permission string     `gorm:"type:varchar(100);column:permission;comment:权限标识;default:''" json:"permission"`
	Query      string     `gorm:"type:varchar(255);column:query;comment:路由参数;default:''" json:"query"`
	Perms      string     `gorm:"type:varchar(100);column:perms;comment:权限字符;default:''" json:"perms"`
	Children   []*SysMenu `gorm:"-"`
	CommonModel
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

// ==================== 系统部门表 ====================
type SysDept struct {
	BaseModel
	ID       uint       `gorm:"type:int(11);primaryKey;column:id;comment:部门id;AUTO_INCREMENT" json:"id"`
	ParentID uint       `gorm:"type:int(11);column:parent_id;comment:父部门id;default:0" json:"parentId"`
	Path     string     `gorm:"type:varchar(200);column:path;comment:部门路径;default:''" json:"path"`
	Name     string     `gorm:"type:varchar(50);column:name;comment:部门名称;not null" json:"name"`
	Sort     int        `gorm:"type:int(11);column:sort;comment:显示顺序;default:1" json:"sort"`
	Leader   string     `gorm:"type:varchar(50);column:leader;comment:负责人;default:''" json:"leader"`
	Phone    string     `gorm:"type:varchar(20);column:phone;comment:联系电话;default:''" json:"phone"`
	Email    string     `gorm:"type:varchar(100);column:email;comment:邮箱;default:''" json:"email"`
	Status   int        `gorm:"type:tinyint(2);column:status;comment:状态【0：禁用 1：正常】;default:1" json:"status"`
	Children []*SysDept `gorm:"-"`
	CommonModel
}

func (SysDept) TableName() string {
	return "sys_dept"
}

// ==================== 用户角色关联表 ====================
type SysUserRole struct {
	BaseModel
	ID     uint `gorm:"type:int(11);primaryKey;column:id;comment:id;AUTO_INCREMENT" json:"id"`
	UserID uint `gorm:"type:int(11);column:user_id;comment:用户id;not null;index" json:"userId"`
	RoleID uint `gorm:"type:int(11);column:role_id;comment:角色id;not null;index" json:"roleId"`
	CommonModel
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}

// ==================== 角色菜单关联表 ====================
type SysRoleMenu struct {
	BaseModel
	ID     uint `gorm:"type:int(11);primaryKey;column:id;comment:id;AUTO_INCREMENT" json:"id"`
	RoleID uint `gorm:"type:int(11);column:role_id;comment:角色id;not null;index" json:"roleId"`
	MenuID uint `gorm:"type:int(11);column:menu_id;comment:菜单id;not null;index" json:"menuId"`
	CommonModel
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

// ==================== 角色部门关联表 ====================
type SysRoleDept struct {
	BaseModel
	ID     uint `gorm:"type:int(11);primaryKey;column:id;comment:id;AUTO_INCREMENT" json:"id"`
	RoleID uint `gorm:"type:int(11);column:role_id;comment:角色id;not null;index" json:"roleId"`
	DeptID uint `gorm:"type:int(11);column:dept_id;comment:部门id;not null;index" json:"deptId"`
	CommonModel
}

func (SysRoleDept) TableName() string {
	return "sys_role_dept"
}

type SystemConfig struct {
	BaseModel
	ID          uint   `gorm:"type:int(11);primaryKey;column:id;comment:配置ID;AUTO_INCREMENT" json:"id"`
	ConfigKey   string `gorm:"type:varchar(100);column:config_key;comment:配置键;not null;uniqueIndex" json:"configKey"`
	ConfigValue string `gorm:"type:longtext;column:config_value;comment:配置值" json:"configValue"`
	Description string `gorm:"type:varchar(500);column:description;comment:描述;default:''" json:"description"`
	CommonModel
}

func (SystemConfig) TableName() string {
	return "system_config"
}
