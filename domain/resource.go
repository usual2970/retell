package domain

import (
	"time"

	"github.com/usual2970/retell/pkg/encode"
	"github.com/usual2970/retell/pkg/utils"
	"gorm.io/gorm"
)

const (
	BannerPreRelease    = 2 // 预发布
	BannerNotPreRelease = 1 // 非预发布
)

type Resource struct {
	ID          int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`               // 主键
	URI         string         `gorm:"column:uri;type:varchar(64);not null;default:''" json:"uri"` // uri
	PositionURI string         `gorm:"column:position_uri;type:varchar(191);not null;default:''" json:"positionUri"`
	PositionID  int64          `gorm:"column:position_id;not null;default:0" json:"positionId"`
	Img         string         `gorm:"column:img;type:varchar(191);not null;default:''" json:"img"`
	Title       string         `gorm:"column:title;type:varchar(191);not null;default:''" json:"title"`
	Description string         `gorm:"column:description;type:text;not null;default:''" json:"description"`
	SortOrder   int64          `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
	URLType     int64          `gorm:"column:url_type;not null;default:0" json:"urlType"`
	URL         string         `gorm:"column:url;type:text;not null;default:''" json:"url"`
	Extend      string         `gorm:"column:extend;type:varchar(191);not null;d	efault:''" json:"extend"`
	IsShow      int8           `gorm:"column:is_show;not null;default:0" json:"isShow"`
	PreRelease  int8           `gorm:"column:pre_release;not null;default:0" json:"preRelease"` // 是否预发布 2 是 1 否
	ShowStartAt *time.Time     `gorm:"column:show_start_at" json:"showStartAt"`                 // 显示开始时间
	ShowEndAt   *time.Time     `gorm:"column:show_end_at" json:"showEndAt"`                     // 显示结束时间
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`       // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`       // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`                      // 删除时间
}

// TableName 指定表名
func (Resource) TableName() string {
	return "resource"
}

func InitResources(items []ResourceInfoItem, position *ResourcePosition) []Resource {
	var resources []Resource
	for _, item := range items {
		resources = append(resources, InitResource(item, position))
	}
	return resources
}

func InitResource(item ResourceInfoItem, position *ResourcePosition) Resource {
	return Resource{
		URI:         encode.Uri(),
		PositionURI: position.URI,
		PositionID:  position.ID,
		Img:         utils.ParseUrl(item.Img),
		Title:       item.Title,
		Description: item.Description,
		URL:         item.URL,
	}
}

func (r *Resource) Update(item ResourceInfoItem) {
	r.Img = utils.ParseUrl(item.Img)
	r.Title = item.Title
	r.Description = item.Description
	r.URL = item.URL
}

func (r *Resource) Trans2Resp() *ResourceInfoItem {
	return &ResourceInfoItem{
		URI:         r.URI,
		PositionURI: r.PositionURI,
		Img:         utils.FullUrl(r.Img),
		Title:       r.Title,
		Description: r.Description,
		URL:         r.URL,
	}
}

type ResourceInfoItem struct {
	URI         string `json:"uri"`
	PositionURI string `json:"positionUri"`
	Img         string `json:"img"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
