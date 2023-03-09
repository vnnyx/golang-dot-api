package web

import (
	"math"

	"gorm.io/gorm"
)

const (
	DEFAULT_LIMIT = 10
	DEFAULT_PAGE  = 1
	DEFAULT_SORT  = "created_at desc"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty"`
	Page       int         `json:"page,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	switch p.Limit {
	case 0:
		return DEFAULT_LIMIT
	default:
		return p.Limit
	}
}

func (p *Pagination) GetPage() int {
	switch p.Page {
	case 0:
		return DEFAULT_PAGE
	default:
		return p.Page
	}
}

func (p *Pagination) GetSort() string {
	switch p.Sort {
	case "":
		return DEFAULT_SORT
	default:
		return p.Sort
	}
}

func (p *Pagination) GetTotalRows() int64 {
	return p.TotalRows
}

func (p *Pagination) GetTotalPages() int {
	return int(p.TotalPages)
}

func (p *Pagination) SetPagination(value interface{}, db *gorm.DB) {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	p.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.GetLimit())))
	p.TotalPages = totalPages
}
