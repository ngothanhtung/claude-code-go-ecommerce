package paging

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// Query holds parsed pagination/filter params.
type Query struct {
	Page     int
	PageSize int
	Sort     string // e.g. "-created_at" or "name"
	Search   string
}

// Parse reads pagination params from the gin context query.
func Parse(c *gin.Context) Query {
	q := Query{Page: defaultPage, PageSize: defaultPageSize}
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			q.Page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= maxPageSize {
			q.PageSize = n
		}
	}
	q.Sort = c.Query("sort")
	q.Search = c.Query("q")
	return q
}

// Result is the standard paginated payload.
type Result struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func NewResult(items interface{}, total int64, q Query) Result {
	return Result{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}
}

// OrderBy converts a Sort string into a safe SQL order clause fragment.
func (q Query) OrderBy() string {
	if q.Sort == "" {
		return "created_at DESC"
	}
	field := q.Sort
	dir := "ASC"
	if len(field) > 0 && field[0] == '-' {
		dir = "DESC"
		field = field[1:]
	}
	switch field {
	case "created_at", "name", "email", "updated_at", "id":
		return field + " " + dir
	default:
		return "created_at DESC"
	}
}
