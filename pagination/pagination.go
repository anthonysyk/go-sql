package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

type Param struct {
	DB      *gorm.DB
	Page    int
	Size    int
	OrderBy []string
	DBIsRaw bool
}

type Meta struct {
	TotalRecord int `json:"total"`
	TotalPage   int `json:"totalPage"`
	Offset      int `json:"offset"`
	Size        int `json:"size"`
	Page        int `json:"page"`
	PrevPage    int `json:"prevPage"`
	NextPage    int `json:"nextPage"`
}

func Paging(p *Param, result interface{}) *Meta {
	db := p.DB

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var count int
	var offset int

	go countRecords(db, result, done, &count, p.DBIsRaw)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Size
	}

	db.Limit(p.Size).Offset(offset).Find(result)
	<-done

	paginator := &Meta{
		TotalRecord: count,
		TotalPage:   int(math.Ceil(float64(count) / float64(p.Size))),
		Offset:      offset,
		Size:        p.Size,
		Page:        p.Page,
		PrevPage:    p.Page,
		NextPage:    p.Page + 1,
	}

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	}

	if p.Page == paginator.TotalPage || paginator.TotalPage == 0 {
		paginator.NextPage = p.Page
	}

	return paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int, dbIsRaw bool) {
	if dbIsRaw {
		db.New().Raw("SELECT count(*) FROM ? rawquery", db.SubQuery()).Count(count)
	} else {
		db.Model(anyType).Count(count)
	}
	done <- true
}
