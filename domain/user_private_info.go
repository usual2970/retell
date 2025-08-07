package domain

import (
	"time"

	"github.com/usual2970/retell/pkg/encode"
	"gorm.io/gorm"
)

type UserPrivateInfo struct {
	UserID    int64          `gorm:"column:user_id;primaryKey" json:"userId"`                    // 用户id
	Tel       string         `gorm:"column:tel;type:varchar(20);not null;default:''" json:"tel"` // 手机号
	Email     string         `gorm:"column:email;type:varchar(64);not null;default:''" json:"email"`
	Password  string         `gorm:"column:password;type:varchar(64);not null;default:''" json:"password"` // 密码
	Salt      string         `gorm:"column:salt;type:varchar(64);not null;default:''" json:"salt"`         // 密码盐
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                    // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                    // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`                             // 删除时间
}

// TableName 指定表名
func (UserPrivateInfo) TableName() string {
	return "user_private_info"
}

func NewPrivateInfo(email string) *UserPrivateInfo {
	return &UserPrivateInfo{
		Email: email,
	}
}

func NewPrivateInfoWithPassword(email, password string) *UserPrivateInfo {
	salt := encode.GenerateRandomString(6)
	password = encode.Md5(encode.Md5(password) + salt)
	return &UserPrivateInfo{
		Email:    email,
		Password: password,
		Salt:     salt,
	}
}

func (u *UserPrivateInfo) SetPassword(password string) {
	salt := encode.GenerateRandomString(6)
	u.Salt = salt
	u.Password = encode.Md5(encode.Md5(password) + salt)
}

func (u *UserPrivateInfo) CheckPassword(password string) bool {
	return u.Password == encode.Md5(encode.Md5(password)+u.Salt)
}
