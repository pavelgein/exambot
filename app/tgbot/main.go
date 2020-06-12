package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"

	"github.com/pavelgein/exambot/internal/db"
	"github.com/pavelgein/exambot/internal/models"
	"github.com/pavelgein/exambot/services"
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

func (bot *Bot) ServeTask(update *tgbotapi.Update) {
	user := bot.UserProvider.GetUser(update.Message.From.UserName)
	if user == nil {
		bot.SendUserNotFound(update)
		return
	}

	assignments := bot.AssignmentService.GetAllAssignments(user)
	if len(assignments) == 0 {
		bot.SendAssignmentNotFound(update)
		return
	}

	now := time.Now()
	log.Printf("Need to send %d assignments", len(assignments))
	bot.AssignmentService.AssignMany(assignments, &now)
	bot.SendAssignment(update, assignments)
}

func (bot *Bot) ServeHelp(update *tgbotapi.Update) {
	bot.replyTo(update, "/task для получения задачи")
}

func (bot *Bot) ServeUnknownCommand(update *tgbotapi.Update) {
	bot.replyTo(update, "/help для справки")
}

func (bot *Bot) ServeCommand(update *tgbotapi.Update, command string) {
	if command == "task" {
		bot.ServeTask(update)
		return
	}

	if command == "help" {
		bot.ServeHelp(update)
		return
	}

	bot.ServeUnknownCommand(update)
}

func (bot *Bot) ServeUpdate(update tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	if command := update.Message.Command(); command != "" {
		log.Printf("recieve command %s", command)
		bot.ServeCommand(&update, command)
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	bot.replyTo(&update, "К сожалению, я понимаю только комманды. /help для справки")
}

func (bot *Bot) SendUserNotFound(update *tgbotapi.Update) {
	bot.replyTo(update, fmt.Sprintf("Пользователь %s не найден", update.Message.From.UserName))
}

func (bot *Bot) SendAssignmentNotFound(update *tgbotapi.Update) {
	bot.replyTo(update, fmt.Sprintf("Для пользователя %s задания не найдены", update.Message.From.UserName))
}

func (bot *Bot) SendAssignment(update *tgbotapi.Update, assignments []models.Assignment) {
	formatted := make([]string, 0)

	for _, assignment := range assignments {
		formatted = append(formatted, fmt.Sprintf("Курс: %s\nЗадача №%d\n%s", assignment.Course.Name, assignment.Task.Number, assignment.Task.Content))
	}

	log.Printf("Ready to send %d items", len(formatted))

	bot.replyTo(update, strings.Join(formatted, "\n\n"))
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

	myBot := Bot{
		Api:               bot,
		AssignmentService: taskSystem,
		UserProvider:      TelegramUserProvider{db},
	}

	myBot.Serve()
}
