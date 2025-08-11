package util

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int32 // 页码，从1开始
	PageSize int32 // 每页条数
}

// NormalizePagination 标准化分页参数
// 默认值：page=1, pageSize=20
// 限制：pageSize最大100
func NormalizePagination(page, pageSize int32) PaginationParams {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// CalculateOffset 计算数据库查询的offset
func (p PaginationParams) CalculateOffset() int {
	return int((p.Page - 1) * p.PageSize)
}

// CalculateLimit 计算数据库查询的limit
func (p PaginationParams) CalculateLimit() int {
	return int(p.PageSize)
}
