package command

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/telegram-bot/connection"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
)

func TestReqAdmin(t *testing.T) {
	godotenv.Load("../.env")

	db, err := connection.ConnectMongoDB()
	if err != nil {
		log.Panic(err)
	}
	bot, err := connection.TelegramBotConnection()
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}
	conn := repository.UsersRepositoryMongo{
		ConnectionDB: db,
	}

	var webhook models.WebhookReqBody
	webhook.CallBackQuery = &models.CallBackQuery{}
	webhook.CallBackQuery.Message.Chat.ID = int64(1286701115)
	log.Println(webhook.CallBackQuery.Message.Chat.ID, "tetstt")
	userID := webhook.CallBackQuery.Message.Chat.ID
	state, err := repository.UsersRepository.GetState(&conn, userID)
	if err != nil {
		log.Panic(err)
	}
	err = ConfirmUpdateAdmin(webhook, &conn, bot, *state)
	if err != nil {
		log.Panic(err)
	}

}
