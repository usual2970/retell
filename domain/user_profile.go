package domain

import (
	"fmt"
	"time"

	"github.com/usual2970/retell/pkg/encode"
	"github.com/usual2970/retell/pkg/utils"
	"gorm.io/gorm"
)

type UserProfile struct {
	UserID     int64          `gorm:"column:user_id;primaryKey" json:"userId"`                                   // 用户id
	URI        string         `gorm:"column:uri;type:varchar(32);not null;default:''" json:"uri"`                // 用户编码
	Nickname   string         `gorm:"column:nickname;type:varchar(32);not null;default:''" json:"nickname"`      // 昵称
	Headimgurl string         `gorm:"column:headimgurl;type:varchar(256);not null;default:''" json:"headimgurl"` // 头像
	Sex        int8           `gorm:"column:sex;type:tinyint;not null;default:1" json:"sex"`                     // 性别
	Region     string         `gorm:"column:region;type:varchar(32);not null;default:''" json:"region"`          // 区域
	Country    string         `gorm:"column:country;type:varchar(32);not null;default:''" json:"country"`        // 国家
	Province   string         `gorm:"column:province;type:varchar(32);not null;default:''" json:"province"`      // 省份
	City       string         `gorm:"column:city;type:varchar(32);not null;default:''" json:"city"`              // 城市
	Role       AuthRole       `gorm:"column:role;type:tinyint;not null;default:1" json:"role"`                   // 角色 1 接入商 2 开发者 3 运营者
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                         // 创建时间
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                         // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`                                  // 删除时间
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profile"
}

func InitProfile(email string, role AuthRole) *UserProfile {

	return &UserProfile{
		URI:      encode.Uri(),
		Nickname: getDefaultNickname(email),
		Role:     role,
	}
}

func getDefaultNickname(email string) string {

	pre := email[0:5]
	rs := fmt.Sprintf("%s%s", pre, encode.GenerateRandomString(4))

	return rs
}

type UserProfileResp struct {
	UserUri    string `json:"userUri"`    // 用户编码
	Headimgurl string `json:"headimgurl"` // 头像
}

func (p *UserProfile) Trans2Resp() *UserProfileResp {
	return &UserProfileResp{
		UserUri:    p.URI,
		Headimgurl: utils.FullUrl(p.Headimgurl),
	}
}
