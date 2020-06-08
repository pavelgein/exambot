package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name    string
	Surname string
}

type TelegramUser struct {
	gorm.Model
	Login  string `gorm:"primary_key"`
	UserID int
	User   User
}
