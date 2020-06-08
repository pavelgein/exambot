package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pavelgein/exambot/models"
	"github.com/pavelgein/exambot/services"
)

// func ensureUser(db *gorm.DB, userName string) {
// 	var user models.User
// 	if err := db.Find(&user, "Name = ?", userName).Error; err != nil {
// 		log.Printf("creating user %s", userName)
// 		user.Name = userName
// 		db.NewRecord(&user)
// 		db.Create(&user)
// 	} else {
// 		log.Printf("User %s has id %d", userName, user.ID)
// 	}

// 	var telegramUser models.TelegramUser
// 	if err := db.First(&telegramUser, "Login = ?", userName).Error; err != nil {
// 		log.Printf("creating telegram user %s", userName)
// 		telegramUser.Login = userName
// 		telegramUser.User = user
// 		db.NewRecord(&telegramUser)
// 		db.Create(&telegramUser)
// 	} else {
// 		log.Printf("telegramUser %s has id %d", userName, telegramUser.ID)
// 	}

// 	courseName := "Алгебра и геометрия"
// 	var course models.Course
// 	if err := db.First(&course, "Name = ?", courseName).Error; err != nil {
// 		log.Printf("creating course")
// 		course.Name = courseName
// 		db.NewRecord(&course)
// 		db.Create(&course)
// 	} else {
// 		log.Printf("Course %s has id %d", course.Name, course.ID)
// 	}

// 	taskContent := "Условие задачи"
// 	var task models.Task
// 	if err := db.First(&task, "Content = ?", taskContent).Error; err != nil {
// 		log.Printf("creating task")
// 		task.Content = taskContent
// 		db.NewRecord(&task)
// 		db.Create(&task)
// 	} else {
// 		log.Printf("Task %s has id %d", taskContent, task.ID)
// 	}

// 	var assignment models.Assignment
// 	if err := db.First(&assignment).Error; err != nil {
// 		log.Printf("creating assigment")
// 		assignment.Task = task
// 		assignment.Course = course
// 		assignment.User = user
// 		db.NewRecord(&assignment)
// 		db.Create(&assignment)
// 	} else {
// 		log.Printf("Assignment has id %d, %+v", assignment.ID, assignment)
// 	}
// }

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

type Bot struct {
	AssignmentService services.AssignmentService
	UserProvider      TelegramUserProvider
	Api               *tgbotapi.BotAPI
}

func (bot *Bot) Serve() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.Api.GetUpdatesChan(u)
	if err != nil {
		panic("cannot get update")
	}

	for update := range updates {
		go bot.ServeUpdate(update)
	}
}

func (bot *Bot) ServeUpdate(update tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	user := bot.UserProvider.GetUser(update.Message.From.UserName)
	if user == nil {
		bot.SendUserNotFound(&update)
		return
	}

	assigment := bot.AssignmentService.GetAssignment(user)
	if assigment == nil {
		bot.SendAssignmentNotFound(&update)
		return
	}

	now := time.Now()
	bot.AssignmentService.Assign(assigment, &now)
	bot.SendAssignment(&update, assigment)
}

func (bot *Bot) SendUserNotFound(update *tgbotapi.Update) {
	bot.replyTo(update, fmt.Sprintf("Пользователь %s не найден", update.Message.From.UserName))
}

func (bot *Bot) SendAssignmentNotFound(update *tgbotapi.Update) {
	bot.replyTo(update, fmt.Sprintf("Для пользователя %s задания не найдены", update.Message.From.UserName))
}

func (bot *Bot) SendAssignment(update *tgbotapi.Update, assignment *models.Assignment) {
	bot.replyTo(update, fmt.Sprintf("Курс: %s\nЗадача №%d\n%s", assignment.Course.Name, assignment.Task.ID, assignment.Task.Content))
}

func (bot *Bot) replyTo(update *tgbotapi.Update, payload string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, payload)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Api.Send(msg)
}

func main() {
	config := MakeConfigFromEnvironemnt()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(config.DBDialect, config.DBConnectionParams)
	if err != nil {
		log.Printf("error: %s", err.Error())
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Task{}, &models.Assignment{}, &models.Course{}, &models.User{}, &models.TelegramUser{})
	taskSystem := services.FullAssignmentService{DB: db}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	myBot := Bot{
		Api:               bot,
		AssignmentService: taskSystem,
		UserProvider:      TelegramUserProvider{db},
	}

	myBot.Serve()
}
