package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/profile"),
			tgbotapi.NewKeyboardButton("/login"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/logout"),
			tgbotapi.NewKeyboardButton("/absensi"),
		),
	)
	menu.OneTimeKeyboard = true
	menu.Selective = true

	msg := tgbotapi.NewMessage(chatID, "Select an option:")
	msg.ReplyMarkup = menu

	bot.Send(msg)

}
