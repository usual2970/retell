package domain

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const SettingMessageRetentionDaysKey = "message_retention_days"

type Setting struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                 // 主键
	Key         string         `gorm:"column:key;type:varchar(128);not null;default:''" json:"key"`                  // 键
	Value       string         `gorm:"column:value;type:varchar(1024);not null;default:''" json:"value"`             // 值
	UserID      int64          `gorm:"column:user_id;default:0" json:"userId"`                                       // 用户ID
	Description string         `gorm:"column:description;type:varchar(1024);not null;default:''" json:"description"` // 描述
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                            // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                            // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`                                           // 删除时间
}

func (Setting) TableName() string {
	return "setting"
}

func InitSetting(key string, userId int64) *Setting {

	setting := &Setting{
		Key:    key,
		UserID: userId,
	}

	return setting
}

func (s *Setting) SetValue(value any) error {
	valueBts, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.Value = string(valueBts)
	return nil
}

func (s *Setting) Trans2Resp() *SettingResp {
	var value interface{}
	if err := json.Unmarshal([]byte(s.Value), &value); err != nil {
		// 如果解析失败，则返回原始字符串
		value = s.Value
	}

	return &SettingResp{
		Key:   s.Key,
		Value: value,
	}
}

type SettingAppType struct {
	AppType string `json:"appType"` // 应用类型
	Label   string `json:"label"`   // 标签
}

type SettingMessageRetentionDays struct {
	Days int `json:"days"` // 消息保留天数
}

func (s *Setting) GetRetentionDays() int {
	var retention SettingMessageRetentionDays
	if err := json.Unmarshal([]byte(s.Value), &retention); err != nil {
		return 30 // 默认值
	}
	return retention.Days
}

type SettingResp struct {
	Key   string `json:"key"`   // 键
	Value any    `json:"value"` // 值
}

//{"balance":500,"interval":1,"scope":1}

type LowBalanceThresholdSetting struct {
	Balance  int `json:"balance"`  // 余额阈值
	Interval int `json:"interval"` // 通知间隔，单位天
	Scope    int `json:"scope"`    // 通知范围，1: 全部
}
