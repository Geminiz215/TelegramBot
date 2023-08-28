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

func SignOut(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
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
	if activity[0].SignOut != nil {
		bot.Send(tgbotapi.NewMessage(chatId, "You are not logged in."))
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
	p.UpdateLog(models.ActivityLog{
		DocumentBase: models.DocumentBase{
			ID: activity[0].ID,
		},
		UserID:  body.Message.From.ID,
		SignOut: &currentTime,
	})
	bot.Send(tgbotapi.NewMessage(chatId, "Sucessed"))
	SendMainMenu(bot, chatId)
}

func IsWeekend() bool {
	currentTime := time.Now()
	dayOfWeek := currentTime.Weekday()

	if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
		return true
	}

	return false
}
