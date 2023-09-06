package command

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func MainMenu(bot *tgbotapi.BotAPI, chatID int64, userID int64, p repository.UsersRepository) error {
	_, err := p.FindUser(models.UserQuery{
		UserID: &userID,
	})
	status := true
	if err != nil && err == mongo.ErrNoDocuments {
		SendMainMenu(bot, chatID, userID, false, &status)
	}
	activity, _, _ := p.FindLog(models.ActivityLogQuery{
		UserID: &userID,
	})

	if len(activity) != 0 && activity[0].SignIn != nil && activity[0].SignOut == nil {
		SendMainMenu(bot, chatID, userID, false, nil)
	}
	SendMainMenu(bot, chatID, userID, true, nil)
	return nil

}

func SendMainMenu(bot *tgbotapi.BotAPI, chatID int64, userID int64, status bool, profile *bool) {
	//true login
	//false logout
	msg := tgbotapi.NewMessage(chatID, "Select an menu option:")
	if profile != nil && *profile {
		msg = tgbotapi.NewMessage(chatID, "please insert profile: /profile")
	}

	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/login"),
		),
	)

	if !status {
		menu = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/logout"),
			),
		)
	}
	menu.OneTimeKeyboard = true
	menu.Selective = true

	msg.ReplyMarkup = menu

	bot.Send(msg)

}

func DeleteMarkReply(bot *tgbotapi.BotAPI, body models.WebhookReqBody) {
	editConfig := tgbotapi.NewEditMessageText(body.CallBackQuery.Message.Chat.ID, int(body.CallBackQuery.Message.MessageId), body.CallBackQuery.Message.Text)
	_, err := bot.Send(editConfig)
	if err != nil {
		log.Panicln(err)
	}
}
