package domain

import (
	"time"

	"github.com/usual2970/retell/domain/constant"
	"gorm.io/gorm"
)

type MessageLevel int8

const (
	MessageLevelSuccess MessageLevel = iota + 1
	MessageLevelNotice
	MessageLevelWarning
	MessageLevelError
)

type MessageHistory struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID    int            `gorm:"column:task_id;default:0" json:"task_id"`
	TplID     int            `gorm:"column:tpl_id;default:0" json:"tpl_id"`
	UserID    int64          `gorm:"column:user_id;default:0" json:"user_id"`
	Receiver  string         `gorm:"column:reciever;size:64;default:''" json:"receiver"`
	Title     string         `gorm:"column:title;size:128;default:''" json:"title"`
	Content   string         `gorm:"column:content;size:512;default:''" json:"content"`
	Link      string         `gorm:"column:link;size:128;default:''" json:"link"`
	ReadedAt  *time.Time     `gorm:"column:readed_at" json:"readed_at"`
	Type      MessageTplType `gorm:"column:type;default:0" json:"type"`
	Level     MessageLevel   `gorm:"column:level;default:0" json:"level"`
	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 指定表名
func (m *MessageHistory) TableName() string {
	return "message_history"
}

type InitMessageHistoryParams struct {
	Task     *MessageTask `json:"task"`
	Tpl      *MessageTpl  `json:"tpl"`
	UserID   int64        `json:"user_id"`
	UserUri  string       `json:"user_uri"`
	Receiver string       `json:"receiver"`
	Title    string       `json:"title"`
	Content  string       `json:"content"`
	Link     string       `json:"link"`
}

func NewMessageHistory(params *InitMessageHistoryParams) *MessageHistory {
	now := time.Now().UTC()
	return &MessageHistory{
		TaskID: params.Task.ID,
		TplID:  params.Tpl.ID,
		Type:   params.Tpl.Type,

		UserID: params.UserID,

		Receiver: params.Receiver,
		Title:    params.Title,
		Content:  params.Content,
		Link:     params.Link,
		Level:    params.Task.Level,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

type MessageListReq struct {
	constant.Pagination
	KWD            string       `json:"kwd,omitempty" query:"kwd"`                       // 关键字
	Level          MessageLevel `json:"level,omitempty" query:"level"`                   // 消息等级
	Readed         *bool        `json:"readed,omitempty" query:"readed"`                 // 是否已读
	CreatedAtStart int64        `json:"createdAtStart,omitempty" query:"createdAtStart"` // 创建时间开始
	CreatedAtEnd   int64        `json:"createdAtEnd,omitempty" query:"createdAtEnd"`     // 创建时间结束
}

type MessageInfoResp struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Level     MessageLevel `json:"level"`
	Content   string       `json:"content"`
	Readed    bool         `json:"readed"`
	CreatedAt int64        `json:"createdAt"` // 创建时间
}

type MessageSetRetentionDaysReq struct {
	Days int `json:"days"` // 保留天数
}

func (m *MessageHistory) Trans2InfoResp() *MessageInfoResp {
	return &MessageInfoResp{
		ID:        m.ID,
		Title:     m.Title,
		Level:     m.Level,
		Content:   m.Content,
		Readed:    m.ReadedAt != nil,
		CreatedAt: m.CreatedAt.Unix(),
	}
}

type MessageReadReq struct {
	ID int `json:"id"` // 消息ID
}
