package domain

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        int64          `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Uri       string         `gorm:"column:uri;type:varchar(64);not null;uniqueIndex" json:"uri"`
	Url       string         `gorm:"column:url;type:varchar(256);not null" json:"url"`
	Name      string         `gorm:"column:name;type:varchar(64);not null" json:"name"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
}

// TableName 指定表名
func (File) TableName() string {
	return "file"
}

type FileUploadResp struct {
	URL string `json:"url"`
}

// FileStructureResp 文件结构响应
type FileStructureResp struct {
	Name  string      `json:"name"`  // 文件名
	Type  string      `json:"type"`  // 文件类型 (archive, image, video, etc.)
	Size  int64       `json:"size"`  // 文件大小（字节）
	Files []*FileNode `json:"files"` // 子文件列表（如果是压缩文件）
}

// FileNode 文件节点
type FileNode struct {
	Name     string      `json:"name"`     // 节点名称
	Path     string      `json:"path"`     // 节点路径
	Size     int64       `json:"size"`     // 节点大小
	Type     string      `json:"type"`     // 节点类型 (file, directory)
	Children []*FileNode `json:"children"` // 子节点（如果是目录）
}

type FileValidateSQLReq struct {
	URL string `json:"url"` // SQL 文件的 URL
}

type FileValidateSQLResp struct {
	Error string `json:"error"` // 错误信息，如果有的话
}

type FilePolicyTokenResp struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Dir              string `json:"dir"`
}
