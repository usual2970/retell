package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	UserAccountPlatformEmail int8 = 1
)

type UserAccount struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                         // 主键
	UserID         int64          `gorm:"column:user_id;not null;default:0" json:"userId"`                      // 用户id
	Openid         string         `gorm:"column:openid;type:varchar(64);not null;default:''" json:"openid"`     // openid
	PlatformID     int8           `gorm:"column:platform_id;type:tinyint;not null;default:1" json:"platformId"` // 注册来源1手机号
	DeregisteredAt *time.Time     `gorm:"column:deregistered_at;default:NULL" json:"deregisteredAt"`            // 注销时间
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                    // 创建时间
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                    // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`                             // 删除时间

	UserProfile     *UserProfile     `gorm:"foreignKey:UserID;references:UserID" json:"profile"`
	UserPrivateInfo *UserPrivateInfo `gorm:"foreignKey:UserID;references:UserID" json:"privateInfo"`
}

// TableName 指定表名
func (UserAccount) TableName() string {
	return "user_account"
}

func NewAccount(openid string, platformId int8) *UserAccount {

	return &UserAccount{
		Openid:     openid,
		PlatformID: platformId,
	}
}
