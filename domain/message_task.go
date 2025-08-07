package domain

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	jsoniter "github.com/json-iterator/go"
)

var ErrTaskParam = errors.New("param wrong")

type MessageTask struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	Uri       string         `gorm:"column:uri;type:varchar(32);not null;default:''"`
	Param     string         `gorm:"column:param;type:varchar(128);not null;default:''"`
	Title     string         `gorm:"column:title;type:varchar(128);not null;default:''"`
	Level     MessageLevel   `gorm:"column:level;default:0"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	MessaageTpls []MessageTpl `gorm:"foreignKey:TaskID;"`
}

// TableName 指定表名
func (m *MessageTask) TableName() string {
	return "message_task"
}

type MessageSendReq struct {
	UserID     string            `json:"userId"`
	TaskURI    string            `json:"taskUri"`
	Param      map[string]string `json:"param"`
	ToUserInfo *ToUserInfo       `json:"toUserInfo"`
}

type MessageSendResp struct{}

type taskParam struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}

func (t *MessageTask) CheckParam(param map[string]string) error {
	tps := make([]taskParam, 0)
	if err := jsoniter.UnmarshalFromString(t.Param, &tps); err != nil {
		return err
	}

	for _, tp := range tps {
		if _, ok := param[tp.Key]; !ok {
			return fmt.Errorf("%w,key:%s", ErrTaskParam, tp.Key)
		}
	}

	return nil
}

type ToUserInfo struct {
	Tel             string `json:"tel,omitempty"`
	Email           string `json:"email,omitempty"`
	UserId          int64  `json:"userId,omitempty"`
	UserUri         string `json:"userUri,omitempty"`
	FeishuBotToken  string `json:"feishuBotToken,omitempty"`
	SystemPushTopic string `json:"systemPushTopic,omitempty"`
}
