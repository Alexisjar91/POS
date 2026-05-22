// Package pagination provee tipos genéricos para paginación de resultados.
package pagination

// Pagination define los parámetros de paginación de una consulta.
// Page y PageSize con valores < 1 se normalizan a 1 y 10 respectivamente.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Offset calcula el desplazamiento para la consulta SQL.
func (p Pagination) Offset() int {
	page := p.Page
	if page < 1 {
		page = 1
	}
	pageSize := p.Limit()
	return (page - 1) * pageSize
}

// Limit retorna el tamaño de página normalizado.
func (p Pagination) Limit() int {
	if p.PageSize < 1 {
		return 10
	}
	return p.PageSize
}

// PaginatedResult contiene los resultados paginados de una consulta.
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResult construye un PaginatedResult a partir de los datos, total y paginación.
func NewPaginatedResult[T any](data []T, total int64, page Pagination) PaginatedResult[T] {
	pageSize := page.Limit()
	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginatedResult[T]{
		Data:       data,
		Page:       page.Page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
