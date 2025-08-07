package constant

type Pagination struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"pageSize" query:"pageSize"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// PaginatedResponse 是一个通用的分页响应结构体，可用于任何类型的数据集合
type PaginatedResponse[T any] struct {
	Page       int   `json:"page"`       // 当前页码
	PageSize   int   `json:"pageSize"`   // 每页条目数
	TotalPages int   `json:"totalPages"` // 总页数
	TotalCount int64 `json:"totalCount"` // 总条目数
	Data       []T   `json:"data"`       // 实际数据列表，使用泛型 T
}

// NewPaginatedResponse 创建一个新的分页响应
// pagination: 分页参数
// data: 当前页的数据列表
// totalCount: 数据总条数
func NewPaginatedResponse[T any](pagination *Pagination, data []T, totalCount int64) *PaginatedResponse[T] {
	pageSize := pagination.GetLimit()

	// 计算总页数
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	}

	return &PaginatedResponse[T]{
		Page:       pagination.GetPage(),
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
		Data:       data,
	}
}
