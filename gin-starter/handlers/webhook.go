package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-starter/command"
	"github.com/gin-starter/connection"
	"github.com/gin-starter/models"
	"github.com/gin-starter/repository"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WebhookController struct {
	Repository repository.UsersRepository
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", os.Getenv("TOKEN"))
}

func (h WebhookController) WebhookCallback(c *gin.Context) {
	bot, _ := connection.TelegramBotConnection()
	bot.Debug = true

	body := &models.WebhookReqBody{}
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}
	var chatId int64
	if body.CallBackQuery != nil {
		chatId = body.CallBackQuery.From.ID
	} else {
		chatId = body.Message.Chat.ID
	}
	state, err := h.Repository.GetState(chatId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}

	if state != nil {
		switch {
		case state.State == "Insert_Profile" && state.SubState == "FirstName":
			command.SendInsertFirstname(*body, h.Repository, bot, *state)
			return
		case state.State == "Insert_Profile" && state.SubState == "LastName":
			command.SendInsertLastname(*body, h.Repository, bot, *state)
			return
		case state.State == "Insert_Profile" && state.SubState == "Confirm":
			fmt.Println("masuk 3")
			if body.CallBackQuery.Data == "/confirm" {
				command.SendConfirmData(h.Repository, bot, *body, *state)
				command.SendMainMenu(bot, chatId)
			} else if body.CallBackQuery.Data == "/cancel" {
				command.SendCancelData(h.Repository, bot, *body, *state)
			} else {
				var data models.DataProfile
				d, _ := bson.Marshal(state.Data)
				bson.Unmarshal(d, &data)
				bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
				command.ReplyConfirmData(bot, chatId, data.FirstName, data.LasttName)
			}
			return
		case state.State == "Update_Profile" && state.SubState == "FirstName":
			command.SendInsertFirstname(*body, h.Repository, bot, *state)
			return
		case state.State == "Update_Profile" && state.SubState == "LastName":
			command.SendInsertLastname(*body, h.Repository, bot, *state)
			return
		case state.State == "Update_Profile" && state.SubState == "Confirm":
			if body.CallBackQuery.Data == "/confirm" {
				command.UpdateProfile(h.Repository, bot, *body, *state)
				command.SendMainMenu(bot, chatId)
			} else if body.CallBackQuery.Data == "/cancel" {
				command.SendCancelData(h.Repository, bot, *body, *state)
			} else {
				var data models.DataProfile
				d, _ := bson.Marshal(state.Data)
				bson.Unmarshal(d, &data)
				bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
				command.ReplyConfirmData(bot, chatId, data.FirstName, data.LasttName)
			}
			return
		default:
			fmt.Println("Test", state.State == string(models.StateKindEnum.InsertProfile) || state.State == string(models.StateKindEnum.UpdateProfile))
		}
	}

	if body.CallBackQuery != nil {
		switch body.CallBackQuery.Data {
		case "/insert":
			if err := command.SendInsert(*body, h.Repository, bot, models.StateKindEnum.InsertProfile); err != nil {
				log.Panic(err)
			}
			return
		case "/back":
			command.SendMainMenu(bot, chatId)
			return
		case "/update":
			if err := command.SendInsert(*body, h.Repository, bot, models.StateKindEnum.UpdateProfile); err != nil {
				log.Panic(err)
			}
			return
		default:
			bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
			command.SendMainMenu(bot, chatId)
		}
	}

	if body.Message.Text[0] == '/' {
		text := strings.ToLower(body.Message.Text)
		switch {
		case text == "/start" || text == "/back" || text == "/no":
			command.SendMainMenu(bot, chatId)
		case text == "/profile":
			command.Profile(h.Repository, bot, *body)
		case text == "/login":
			if command.IsWeekend() {
				bot.Send(tgbotapi.NewMessage(chatId, "cannot signIn at weekend."))
				return
			}
			command.SignIn(h.Repository, bot, *body)
		case text == "/logout":
			if command.IsWeekend() {
				bot.Send(tgbotapi.NewMessage(chatId, "cannot signOut at weekend."))
				return
			}
			command.SignOut(h.Repository, bot, *body)
		case text == "/coba":
			command.CheckFile("data.csv")
			massage := tgbotapi.NewMessage(chatId, "messageText")
			massage.ReplyMarkup = createPaginationKeyboard(2, 4)
			bot.Send(massage)
		case text == "/absensi":
			command.Attendance(*body, bot, h.Repository)
			command.SendMainMenu(bot, chatId)
		default:
			bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
			command.SendMainMenu(bot, chatId)
		}
	}
	c.String(http.StatusOK, "Working!")
	return
}

func createPaginationKeyboard(currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	btnPrev := tgbotapi.NewInlineKeyboardButtonData("👍🏻", "/yes")
	btnNext := tgbotapi.NewInlineKeyboardButtonData("👎🏻", "/no")

	row := tgbotapi.NewInlineKeyboardRow()
	if currentPage > 0 {
		row = append(row, btnPrev)
	}
	if currentPage < totalPages-1 {
		row = append(row, btnNext)
	}

	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func sendRequestInsertData(bot *tgbotapi.BotAPI, chatID int64, data []string) {
	var row []tgbotapi.KeyboardButton
	for _, i := range data {
		row = append(row, tgbotapi.NewKeyboardButton(i))
	}
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			row...,
		),
	)
	numericKeyboard.OneTimeKeyboard = true
	numericKeyboard.Selective = true
	msg := tgbotapi.NewMessage(chatID, "Select an option:")
	msg.ReplyMarkup = numericKeyboard

	bot.Send(msg)
}

func sendValidInsertData(bot *tgbotapi.BotAPI, chatID int64) {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("/yes", "/yes"),
			tgbotapi.NewInlineKeyboardButtonData("/no", "/start"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Select an option:")
	msg.ReplyMarkup = numericKeyboard

	bot.Send(msg)
}

func sendCSVData(activity []models.ActivityLog) {
	// Create a CSV file
	fileName := "data.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Write the CSV header
	headers := []string{"Name", "Sign_in_hour", "Sign_out_hour"}
	if err := csvWriter.Write(headers); err != nil {
		log.Fatal(err)
	}

	for _, i := range activity {
		row := []string{i.UserName, i.SignIn.Format("Monday, January 2, 2006 15:04:05 MST"), i.SignOut.Format("Monday, January 2, 2006 15:04:05 MST")}
		if err := csvWriter.Write(row); err != nil {
			log.Fatal(err)
		}
	}

}
