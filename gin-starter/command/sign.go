package command

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-starter/models"
	"github.com/gin-starter/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignIn(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
	chatId := body.Message.Chat.ID
	if _, err := p.FindUser(models.UserQuery{UserID: &body.Message.From.ID}); err != nil {
		if err == mongo.ErrNoDocuments {
			requestInsertData(bot, chatId)
			return
		}
	}
	activity, _, err := p.FindLog(models.ActivityLogQuery{
		UserID: &body.Message.From.ID,
	})
	if err != nil {
		log.Panic(err)
	}

	if activity[0].SignOut == nil && activity[0].SignIn != nil {
		bot.Send(tgbotapi.NewMessage(chatId, "You are still logged in."))
		SendMainMenu(bot, chatId)
		return
	}

	// Load the location
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	// Get the current time in the specified location
	currentTime := time.Now().In(loc)
	p.CreateLog(models.ActivityLog{
		UserID:   body.Message.From.ID,
		UserName: body.Message.From.UserName,
		SignIn:   &currentTime,
	})
	bot.Send(tgbotapi.NewMessage(chatId, "Sucessed. "))
	SendMainMenu(bot, chatId)
}

func requestInsertData(bot *tgbotapi.BotAPI, chatID int64) {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ insert", "/insert"),
			tgbotapi.NewInlineKeyboardButtonData("➫ back", "/back"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "ur data doesn't exist please insert profile:")
	msg.ReplyMarkup = numericKeyboard

	bot.Send(msg)
}
