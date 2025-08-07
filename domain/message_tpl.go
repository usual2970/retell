package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MessageTplType int8

const (
	_ MessageTplType = iota
	MessageTplTypeEmail
	MessageTplTypeSystem
	MessageTplTypeSystemPush
	MessageTplTypeFeishu
)

const messageSystemPushTopic = "new_message_system_push:%s"

const FeishuBotToken = "95e57567-7652-49cb-acfa-90f7d3b286be"

var (
	ErrHasNoEmail = errors.New("用户没有邮箱")
	ErrHasNoId    = errors.New("用户id不存在")
	ErrTplType    = errors.New("模板类型有误")
)

type MessageTpl struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID    int            `gorm:"column:task_id;default:0" json:"task_id"`
	TaskURI   string         `gorm:"column:task_uri;size:64;default:''" json:"task_uri"`
	Type      MessageTplType `gorm:"column:type;default:0" json:"type"`
	Title     string         `gorm:"column:title;size:64;default:''" json:"title"`
	Content   string         `gorm:"column:content;size:512;default:''" json:"content"`
	Link      string         `gorm:"column:link;size:128;default:''" json:"link"`
	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 指定表名
func (m *MessageTpl) TableName() string {
	return "message_tpl"
}

func (t *MessageTpl) GetTitle(param map[string]string) string {
	title := t.Title
	for k, v := range param {
		title = strings.ReplaceAll(title, "{{"+k+"}}", v)
	}
	return title
}

func (t *MessageTpl) GetContent(param *MessageSendReq) (string, error) {
	content := t.Content
	for k, v := range param.Param {
		content = strings.ReplaceAll(content, "{{"+k+"}}", v)
	}

	return content, nil
}

func (t *MessageTpl) GetLink(param *MessageSendReq) (string, error) {
	link := t.Link
	for k, v := range param.Param {
		link = strings.ReplaceAll(link, "{{"+k+"}}", v)
	}

	return link, nil
}

func (t *MessageTpl) GetToUser(userInfo *ToUserInfo) (string, error) {
	switch t.Type {
	case MessageTplTypeEmail:
		if userInfo.Email == "" {
			return "", ErrHasNoEmail
		}
		return userInfo.Email, nil

	case MessageTplTypeSystem:
		return fmt.Sprintf("%d", userInfo.UserId), nil
	case MessageTplTypeFeishu:
		if userInfo.FeishuBotToken == "" {
			return "", ErrHasNoId
		}
		return userInfo.FeishuBotToken, nil
	case MessageTplTypeSystemPush:
		return fmt.Sprintf(messageSystemPushTopic, userInfo.UserUri), nil
	}

	return "", ErrTplType

}
