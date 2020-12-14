package utils

import (
	"math"
)

const (
	Limit_Default = 10
	Skip_Default  = 1
)

type Page struct {
	Total int
	Pages int
	Page  int
	Limit int
	Skip  int
	Next  bool
}

//app 不用count console count
func (this *Page) Build(total int, page int, limit int) (o *Page) {
	this.Limit = Limit_Default
	this.Page = Skip_Default
	if limit > 0 {
		this.Limit = limit
	}
	if page > 0 {
		this.Page = page
	}

	this.Skip = this.Limit * (this.Page - 1)
	this.Total = total
	this.Pages = int(math.Ceil(float64(this.Total) / float64(this.Limit)))
	if this.Page < this.Pages {
		this.Next = true
	}
	return this
}

func (this *Page) Info() map[string]interface{} {
	pageInfo := make(map[string]interface{})
	pageInfo["totalItems"] = this.Total
	pageInfo["pageNumber"] = this.Page
	pageInfo["resultsPerPage"] = this.Limit
	pageInfo["nextPageAvailable"] = this.Next
	return pageInfo
}

func (this *Page) Default(_page int, _limit int) (o *Page) {
	this.Limit = Limit_Default
	this.Page = Skip_Default
	if _limit > 0 {
		this.Limit = _limit
	}
	if _page > 0 {
		this.Page = _page
	}

	this.Skip = this.Limit * (this.Page - 1)

	return this
}
