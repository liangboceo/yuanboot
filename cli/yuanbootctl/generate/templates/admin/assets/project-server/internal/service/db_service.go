package service

import (
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/datasources/mysql"
	"gorm.io/gorm"
	"sendex-server/internal/model"
)

type DbService struct {
	Db    *gorm.DB
	Log   xlog.ILogger
	Cache *CacheService
}

func NewDbService(source *mysql.MySqlDataSource, cache *CacheService) *DbService {
	db := mysql.NewGormDb(source)
	_ = db.Set("gorm:table_options", "COMMENT='系统用户表'").AutoMigrate(&model.SysUser{})
	_ = db.Set("gorm:table_options", "COMMENT='系统角色表'").AutoMigrate(&model.SysRole{})
	_ = db.Set("gorm:table_options", "COMMENT='系统菜单表'").AutoMigrate(&model.SysMenu{})
	_ = db.Set("gorm:table_options", "COMMENT='系统部门表'").AutoMigrate(&model.SysDept{})
	_ = db.Set("gorm:table_options", "COMMENT='用户角色关联表'").AutoMigrate(&model.SysUserRole{})
	_ = db.Set("gorm:table_options", "COMMENT='角色菜单关联表'").AutoMigrate(&model.SysRoleMenu{})
	_ = db.Set("gorm:table_options", "COMMENT='角色部门关联表'").AutoMigrate(&model.SysRoleDept{})
	_ = db.Set("gorm:table_options", "COMMENT='系统配置表'").AutoMigrate(&model.SystemConfig{})
	return &DbService{Db: db, Log: xlog.GetXLogger("DbService"), Cache: cache}
}
