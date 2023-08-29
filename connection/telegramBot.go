package connection

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TelegramBotConnection() (*tgbotapi.BotAPI, error) {
	TOKEN := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	return bot, err
}

func TelegramConnection() (*http.Response, error) {
	TOKEN := os.Getenv("TOKEN")
	SERVER_URL := os.Getenv("SERVER_URL")
	TELEGRAM_API := "https://api.telegram.org/bot" + TOKEN
	URI := "/webhook"
	WEBHOOK_URL := SERVER_URL + URI
	return http.Get(TELEGRAM_API + "/setWebhook?url=" + WEBHOOK_URL)
}
