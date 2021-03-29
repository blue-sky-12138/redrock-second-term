package utilities

import "gorm.io/gorm"

type Filters []Filter

type Filter struct {
	gorm.Model
	Content string
}

func (f Filter) TableName() string {
	return "filter_information"
}
