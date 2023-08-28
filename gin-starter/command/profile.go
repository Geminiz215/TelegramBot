package command

import (
	"fmt"

	"github.com/gin-starter/models"
	"github.com/gin-starter/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Profile(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
	user, err := p.FindUser(models.UserQuery{UserID: &body.Message.From.ID})
	chatId := body.Message.From.ID
	if err != nil {
		if err == mongo.ErrNoDocuments {
			bot.Send(tgbotapi.NewMessage(chatId, "Insert profile data first"))
			replyInsertData(bot, chatId)
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("FirstName : %s \nLastName : %s\nUserName : %s", user.FirstName, user.LastName, user.UserName)))
	replyUpdateData(bot, chatId)
}

func replyInsertData(bot *tgbotapi.BotAPI, chatID int64) {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ insert", "/insert"),
			tgbotapi.NewInlineKeyboardButtonData("➫ back", "/back"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Select an option:")
	msg.ReplyMarkup = numericKeyboard

	bot.Send(msg)
}

func replyUpdateData(bot *tgbotapi.BotAPI, chatID int64) {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ update", "/update"),
			tgbotapi.NewInlineKeyboardButtonData("➫ back", "/back"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Select an option:")
	msg.ReplyMarkup = numericKeyboard

	bot.Send(msg)
}
