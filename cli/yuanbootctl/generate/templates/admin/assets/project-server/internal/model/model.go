package model

import "time"

type BaseModel struct {
	ID            uint   `gorm:"type:int(11);primaryKey;column:id;comment:id;AUTO_INCREMENT" `
	CreateOrgCode string `gorm:"type:varchar(11);column:create_org_code;comment:创建组织;default:'';not null" doc:"创建组织"`
}
type CommonModel struct {
	IsDeleted  int       `gorm:"type:tinyint(2);column:is_deleted;comment:是否删除【0：未删除 1：删除】;default:0;not null" doc:"是否删除【0：未删除 1：删除】"`
	CreateBy   int       `gorm:"type:int(11);column:create_by;comment:创建人;default:0;not null" doc:"创建人"`
	CreateTime time.Time `gorm:"type:datetime;column:create_time;comment:创建时间;autoCreateTime;not null;default:CURRENT_TIMESTAMP" doc:"创建时间"`
	UpdateBy   int       `gorm:"type:int(11);column:update_by;comment:更新人;default:0;not null" doc:"更新人"`
	UpdateTime time.Time `gorm:"type:datetime;column:update_time;comment:更新时间;autoUpdateTime;not null;default:CURRENT_TIMESTAMP" doc:"更新时间"`
	CreateName string    `gorm:"type:varchar(50);column:create_name;comment:创建人名称;default:'';not null" doc:"创建人名称"`
	UpdateName string    `gorm:"type:varchar(50);column:update_name;comment:更新人名称;default:'';not null" doc:"更新人名称"`
}
