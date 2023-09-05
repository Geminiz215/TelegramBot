package command

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignIn(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
	chatId := body.Message.Chat.ID
	userId := body.Message.From.ID
	if _, err := p.FindUser(models.UserQuery{UserID: &userId}); err != nil {
		if err == mongo.ErrNoDocuments {
			requestInsertData(bot, chatId)
			return
		}
	}
	activity, _, err := p.FindLog(models.ActivityLogQuery{
		UserID: &userId,
	})
	if err != nil {
		log.Panic(err)
	}

	if len(activity) != 0 && (activity[0].SignOut == nil && activity[0].SignIn != nil) {
		bot.Send(tgbotapi.NewMessage(chatId, "You are still logged in."))
		SendMainMenu(bot, chatId, userId, false, nil)
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
	formattedTime := currentTime.Format("January 02, 2006 15:04:05")
	k, _ := bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("status : Sucessed.\nSign_in Time : %s\nSign_out Time :", formattedTime)))
	p.CreateLog(models.ActivityLog{
		UserID:    userId,
		UserName:  body.Message.From.UserName,
		SignIn:    &currentTime,
		MessageID: int64(k.MessageID),
	})
	SendMainMenu(bot, chatId, userId, false, nil)
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
