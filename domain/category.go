package domain

import "time"

const (
	CategoryNotShow = 1 // 不展示
	CategoryShow    = 2 // 展示
)

type Category struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                            // 主键
	URI        string     `gorm:"column:uri;type:varchar(32);not null;default:''" json:"uri"`                                              // 分类URI
	Name       string     `gorm:"column:name;type:varchar(64);not null;default:''" json:"name"`                                            // 分类名称
	ParentID   int64      `gorm:"column:parent_id;not null;default:0" json:"parentId"`                                                     // 父分类ID
	ParentName string     `gorm:"column:parent_name;type:varchar(64);not null;default:''" json:"parentName"`                               // 父分类名称
	AppShow    int8       `gorm:"column:app_show;not null;default:1" json:"appShow"`                                                       // app是否展示 1不展示 2展示
	WebShow    int8       `gorm:"column:web_show;not null;default:1" json:"webShow"`                                                       // web是否展示 1不展示 2展示
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp;default:null" json:"createdAt"`                                          // 创建时间
	UpdatedAt  *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt"` // 更新时间
	DeletedAt  *time.Time `gorm:"column:deleted_at;type:timestamp;default:null" json:"deletedAt"`                                          // 删除时间
}

func (*Category) TableName() string {
	return "category"
}

const (
	PlatformApp = "app" // app平台
	PlatformWeb = "web" // web平台
)

func (c *Category) Trans2Resp() *CategoryResp {
	return &CategoryResp{
		URI:        c.URI,
		Name:       c.Name,
		ParentID:   c.ParentID,
		ParentName: c.ParentName,
	}
}

type CategoriesReq struct {
	ParentUri string `json:"parentUri" query:"parentUri"` // 父分类URI
	Platform  string `json:"platform" query:"platform"`   // 平台类型
}

type CategoryResp struct {
	URI        string `json:"uri"`        // 分类URI
	Name       string `json:"name"`       // 分类名称
	ParentID   int64  `json:"parentId"`   // 父分类ID
	ParentName string `json:"parentName"` // 父分类名称
}

type CategoriesResp struct {
	Categories []CategoryResp `json:"categories"` // 分类列表
}
