package query

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Filter interface {
	Make(*gorm.DB) *gorm.DB
}

type Limit uint

func (q Limit) Make(db *gorm.DB) *gorm.DB {
	limit := uint(q)
	if limit == 0 {
		limit = 100
	}
	return db.Limit(limit)
}

type Offset uint

func (q Offset) Make(db *gorm.DB) *gorm.DB {
	offset := uint(q)
	return db.Offset(offset)
}

type Order map[string]string

func (q Order) Make(db *gorm.DB) *gorm.DB {
	if q != nil {
		for k, v := range q {
			db = db.Order(fmt.Sprintf("%s %s", k, v))
		}
	}
	return db
}

type Where map[string]interface{}

func (q Where) Make(db *gorm.DB) *gorm.DB {
	if q != nil {
		m := make(map[string]interface{}, len(q))
		for f := range q {
			m[f] = q[f]
		}
		db = db.Where(m)
	}
	return db
}

type Search struct {
	Text   string
	Fields []string
}

func (q Search) Make(db *gorm.DB) *gorm.DB {
	if q.Text != "" && len(q.Fields) > 0 {
		search := fmt.Sprintf("%%%s%%", q.Text)
		var args []interface{}
		var builder strings.Builder
		for _, f := range q.Fields {
			args = append(args, search)
			fmt.Fprintf(&builder, "%s LIKE ? OR ", f)
		}
		query := builder.String()
		db = db.Where(query[:len(query)-3], args...) // remove last `OR `
	}
	return db
}

func Args(db *gorm.DB, args ...interface{}) *gorm.DB {
	if len(args) == 0 {
		return db
	}
	for _, arg := range args {
		if arg, ok := arg.(Filter); ok {
			db = arg.Make(db)
		}
	}
	return db
}
