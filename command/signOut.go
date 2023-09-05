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

func SignOut(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
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
	if len(activity) == 0 || activity[0].SignOut != nil {
		bot.Send(tgbotapi.NewMessage(chatId, "You are not logged in."))
		SendMainMenu(bot, chatId, userId, true, nil)
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
		UserID:  userId,
		SignOut: &currentTime,
	})
	formattedTimeIn := activity[0].SignIn.Format("January 02, 2006 15:04:05")
	formattedTimeOut := currentTime.Format("January 02, 2006 15:04:05")
	editMessage(bot, chatId, int(activity[0].MessageID), fmt.Sprintf("status : Sucessed.\nSign_in Time : %s\n Sign_out Time : %s", formattedTimeIn, formattedTimeOut))
	SendMainMenu(bot, chatId, userId, true, nil)
}

func editMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, newMessageText string) error {
	editMessage := tgbotapi.NewEditMessageText(chatID, messageID, newMessageText)
	editMessage.ParseMode = "Markdown"

	_, err := bot.Send(editMessage)
	return err
}

func IsWeekend() bool {
	currentTime := time.Now()
	dayOfWeek := currentTime.Weekday()

	if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
		return true
	}

	return false
}
