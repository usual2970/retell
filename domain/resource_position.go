package domain

import (
	"time"

	"github.com/usual2970/retell/pkg/encode"
	"gorm.io/gorm"
)

const (
	_ int8 = iota
	ResourcePositionPurposeGameBanner
)

type ResourcePosition struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`               // 主键
	UserID      int64          `gorm:"column:user_id;not null;default:0" json:"userId"`            // 用户id
	URI         string         `gorm:"column:uri;type:varchar(64);not null;default:''" json:"uri"` // uri
	Title       string         `gorm:"column:title;type:varchar(255);not null;default:''" json:"title"`
	Description string         `gorm:"column:description;type:varchar(1024);not null;default:''" json:"description"`
	ResourceNum int            `gorm:"column:resource_num;not null;default:0" json:"resourceNum"`                   // 资源位数量
	Purpose     int8           `gorm:"column:purpose;type:tinyint;not null;default:1" json:"purpose"`               // 用途 1 游戏 banner
	RelationURI string         `gorm:"column:relation_uri;type:varchar(64);not null;default:''" json:"relationUri"` // 关联 uri
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                           // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                           // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`                                          // 删除时间

	Resources []Resource `gorm:"foreignKey:PositionID;references:ID" json:"resources"`
}

// TableName 指定表名
func (ResourcePosition) TableName() string {
	return "resource_position"
}

func InitResourcePosition(param *ResourcePositionCreatePositionReq, userID int64) *ResourcePosition {
	return &ResourcePosition{
		UserID:      userID,
		URI:         encode.Uri(),
		Title:       param.Title,
		Description: param.Description,
		Purpose:     param.Purpose,
		RelationURI: param.RelationURI,
	}
}

func (p *ResourcePosition) Trans2Resp() *ResourcePositionResp {

	resourceItems := make([]ResourceInfoItem, 0, len(p.Resources))
	now := time.Now().UTC()
	for _, resource := range p.Resources {

		// 如果为预发布，则要判断是当前时间是否在时间范围内
		if resource.PreRelease == BannerPreRelease && (resource.ShowStartAt == nil || resource.ShowEndAt == nil ||
			(resource.ShowStartAt.After(now) || resource.ShowEndAt.Before(now))) {
			continue
		}

		resourceItems = append(resourceItems, *resource.Trans2Resp())
	}

	return &ResourcePositionResp{
		URI:         p.URI,
		Title:       p.Title,
		Description: p.Description,
		Items:       resourceItems,
	}
}

type ResourcePositionGetReq struct {
	URI string `json:"uri"`
}

type ResourcePositionResp struct {
	URI         string             `json:"uri"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Items       []ResourceInfoItem `json:"items"`
}

type ResourcePositionCreatePositionReq struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Purpose     int8               `json:"purpose"`
	RelationURI string             `json:"relationUri"`
	Items       []ResourceInfoItem `json:"items"`
}

type ResourcePositionUpdatePositionReq struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Purpose     int8               `json:"purpose"`
	RelationURI string             `json:"relationUri"`
	Items       []ResourceInfoItem `json:"items"`
}

type ResourcePositionCreateResourceReq struct {
	PositionURI string           `json:"positionUri"`
	Item        ResourceInfoItem `json:"item"`
}

type ResourcePositionUpdateResourceReq struct {
	PositionURI string           `json:"positionUri"`
	Item        ResourceInfoItem `json:"item"`
}

type ResourcePositionDeleteResourceReq struct {
	PositionURI string `json:"positionUri"`
	Uri         string `json:"uri"`
}

type ResourcePositiionGetGameBannerReq struct {
	GameURI string `json:"gameUri" query:"gameUri"`
}
