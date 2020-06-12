package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"

	"github.com/pavelgein/exambot/internal/db"
	"github.com/pavelgein/exambot/internal/models"
	"github.com/pavelgein/exambot/internal/services"
	"github.com/pavelgein/exambot/packages/tgbot"
)

func ensureUser(db *gorm.DB, userName string) {
	var user models.User
	if err := db.Find(&user, "Name = ?", userName).Error; err != nil {
		log.Printf("creating user %s", userName)
		user.Name = userName
		db.NewRecord(&user)
		db.Create(&user)
	} else {
		log.Printf("User %s has id %d", userName, user.ID)
	}

	var telegramUser models.TelegramUser
	if err := db.First(&telegramUser, "Login = ?", userName).Error; err != nil {
		log.Printf("creating telegram user %s", userName)
		telegramUser.Login = userName
		telegramUser.User = user
		db.NewRecord(&telegramUser)
		db.Create(&telegramUser)
	} else {
		log.Printf("telegramUser %s has id %d", userName, telegramUser.ID)
	}

	courseName := "Алгебра и геометрия"
	var course models.Course
	if err := db.First(&course, "Name = ?", courseName).Error; err != nil {
		log.Printf("creating course")
		course.Name = courseName
		db.NewRecord(&course)
		db.Create(&course)
	} else {
		log.Printf("Course %s has id %d", course.Name, course.ID)
	}

	taskContent := "Условие задачи"
	var task models.Task
	if err := db.First(&task, "Content = ?", taskContent).Error; err != nil {
		log.Printf("creating task")
		task.Content = taskContent
		db.NewRecord(&task)
		db.Create(&task)
	} else {
		log.Printf("Task %s has id %d", taskContent, task.ID)
	}

	var assignment models.Assignment
	if err := db.Model(&assignment).Related(&user).Error; err != nil {
		log.Printf("creating assigment")
		assignment.Task = task
		assignment.Course = course
		assignment.User = user
		db.NewRecord(&assignment)
		db.Create(&assignment)
	} else {
		log.Printf("Assignment has id %d, %+v", assignment.ID, assignment)
	}
}

func main() {
	config := MakeConfigFromEnvironemnt()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	db, err := db.InitWithMigrations(&config.DBConfig)
	if err != nil {
		log.Printf("error: %s", err.Error())
		panic("failed to connect to database")
	}
	defer db.Close()

	if os.Getenv("DEBUG") != "" {
		db = db.Debug()
	}

	taskSystem := services.FullAssignmentService{DB: db}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	myBot := tgbot.Bot{
		Api:               bot,
		AssignmentService: taskSystem,
		UserProvider:      tgbot.TelegramUserProvider{DB: db},
	}

	myBot.Serve()
}
