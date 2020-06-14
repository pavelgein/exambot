package tgbot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/pavelgein/exambot/internal/models"
	"github.com/pavelgein/exambot/internal/services"
)

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
	formatted := make([]string, 0, len(assignments))

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
