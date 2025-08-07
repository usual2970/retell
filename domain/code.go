package domain

import (
	"time"

	"github.com/usual2970/retell/pkg/encode"
	"gorm.io/gorm"
)

const expireDuration = time.Minute * 5

type CodePurpose uint8

const (
	_                            CodePurpose = iota
	CodePurposeLogin                         // 登录
	CodePurposeRegister                      // 注册
	CodePurposeForget                        // 忘记密码
	CodePurposeViewAppSecret                 // 查看应用密钥
	CodePurposeGenerateAppSecret             // 生成应用密钥

)

const (
	_ = iota
	CodeStateSent
	CodeStateUsed
)

type Code struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Receiver  string         `gorm:"type:varchar(32);not null;default:''" json:"receiver"`
	Code      string         `gorm:"type:varchar(8);not null;default:''" json:"code"`
	State     uint8          `gorm:"type:tinyint;not null;default:1" json:"state"`
	Purpose   CodePurpose    `gorm:"type:tinyint;not null;default:1" json:"purpose"`
	ExpiredAt *time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"expiredAt"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (Code) TableName() string {
	return "code"
}

func NewEmailCode(email string, purpose CodePurpose) *Code {
	expiredAt := time.Now().Add(expireDuration).UTC()
	code := encode.Code(4)

	return &Code{
		Receiver: email,

		Purpose:   purpose,
		Code:      code,
		State:     CodeStateSent,
		ExpiredAt: &expiredAt,
	}
}

func (c *Code) IsUsed() bool {
	return c.State == CodeStateUsed
}

type CodeGenerateReq struct {
	Receiver string `json:"receiver"`

	Purpose CodePurpose `json:"purpose"`
}

type CodeCheckReq struct {
	Receiver string      `json:"receiver"`
	Code     string      `json:"code"`
	Purpose  CodePurpose `json:"purpose"`
}

func (c *Code) SetState(state uint8) {
	c.State = state
}
