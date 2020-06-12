package tgbot

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/internal/models"
)

type TelegramUserProvider struct {
	DB *gorm.DB
}

func (provider TelegramUserProvider) GetUser(userName string) *models.User {
	var telegramUser models.TelegramUser
	var user models.User
	log.Printf("Looking for telegram user %s", userName)
	res := provider.DB.First(&telegramUser, "Login = ?", userName).Related(&user)
	if res.Error != nil {
		log.Printf("error %s", res.Error.Error())
		return nil
	}

	return &user
}
