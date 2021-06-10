package model

import "gorm.io/gorm"

type XsBrandSell struct {
	id int
	UserId int
	Logo string
	Title string
	Company string
	Country string
	Types string
	No string
	Price float32
	BrandSmall string
	ContactName string
	ContactPhone string
	Status int8
	ReasonDesc string
	Draft int8
	CreateTime int64
	UpdateTime int64
	HandleTime int64
	CompleteTime int64
	DeleteTime gorm.DeletedAt
	Mark string
	Cost string
}