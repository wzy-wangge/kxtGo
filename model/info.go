package model

import "gorm.io/gorm"

type InfoModel struct {
	Id int
	Cid int
	Title string
	DeleteTime gorm.DeletedAt
}
