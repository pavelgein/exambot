package models

import "github.com/jinzhu/gorm"

type ApiUser struct {
	gorm.Model
	Name  string
	Token string
}

type Page struct {
	gorm.Model
	Name string
}

type Role struct {
	gorm.Model
	Comment string
	UserID  uint
	User    ApiUser `gorm:"foreignkey:UserID"`
	PageID  uint
	Page    Page `gorm:"foreignkey:PageID"`
}
