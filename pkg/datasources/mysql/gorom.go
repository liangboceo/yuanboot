package mysql

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// GormPoolConfig GORM 连接池配置
type GormPoolConfig struct {
	MaxOpenConns    int // 最大打开连接数
	MaxIdleConns    int // 最大空闲连接数
	ConnMaxLifetime int // 连接最大生命周期（秒）
	ConnMaxIdleTime int // 连接最大空闲时间（秒）
}

// NewGormDb 创建 GORM 数据库实例
func NewGormDb(source *MySqlDataSource) *gorm.DB {
	timestr := time.Now().Format("2006/01/02 - 15:04:05.00")
	logPrefix := fmt.Sprintf("%s - [yuanboot] - [DEBUG] ", timestr)
	dbLogger := logger.New(
		log.New(os.Stdout, logPrefix, log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	conn, _, _ := source.Open()
	sqlDB := conn.(*sql.DB)

	// 应用 GORM 连接池配置
	poolConfig := source.GetGormPoolConfig()
	applyGormPoolConfig(sqlDB, poolConfig)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:   dbLogger,
		ConnPool: sqlDB,
	})

	if err != nil {
		panic(err)
	}

	if source.isDebug {
		return gormDB.Debug()
	}

	return gormDB
}

// applyGormPoolConfig 应用连接池配置到 sql.DB
func applyGormPoolConfig(sqlDB *sql.DB, config GormPoolConfig) {
	// 设置最大打开连接数（0 表示无限制）
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}

	// 设置最大空闲连接数
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}

	// 设置连接最大生命周期
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}

	// 设置连接最大空闲时间
	if config.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	}
}
