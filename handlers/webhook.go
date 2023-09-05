package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/telegram-bot/command"
	"github.com/telegram-bot/connection"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
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
	var userId int64
	if body.CallBackQuery != nil {
		userId = body.CallBackQuery.From.ID
	} else {
		userId = body.Message.Chat.ID
	}
	state, err := h.Repository.GetState(chatId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}
	if state != nil {
		switch {
		case state.State == "Request_Admin" && state.SubState == "Accept":
			command.AcceptAdminRequest(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.RequestAdmin) && state.SubState == "Confirm":
			fmt.Println("masuk req admin")
			command.ConfirmUpdateAdmin(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.InsertProfile) && state.SubState == "FirstName":
			command.SendInsertFirstname(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.InsertProfile) && state.SubState == "LastName":
			command.SendInsertLastname(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.InsertProfile) && state.SubState == "Confirm":
			if body.CallBackQuery.Data == "/confirm" {
				command.SendConfirmData(h.Repository, bot, *body, *state)
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
		case state.State == string(models.StateKindEnum.UpdateProfile) && state.SubState == "FirstName":
			command.SendInsertFirstname(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.UpdateProfile) && state.SubState == "LastName":
			command.SendInsertLastname(*body, h.Repository, bot, *state)
			return
		case state.State == string(models.StateKindEnum.UpdateProfile) && state.SubState == "Confirm":
			if body.CallBackQuery.Data == "/confirm" {
				command.UpdateProfile(h.Repository, bot, *body, *state)
				command.MainMenu(bot, chatId, userId, h.Repository)
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
		case "/accept":
			command.AcceptReqAdm(*body, h.Repository, bot)
			return
		case "/insert":
			if err := command.SendInsert(*body, h.Repository, bot, models.StateKindEnum.InsertProfile); err != nil {
				log.Panic(err)
			}
			return
		case "/back":
			command.MainMenu(bot, chatId, userId, h.Repository)
			return
		case "/update":
			if err := command.SendInsert(*body, h.Repository, bot, models.StateKindEnum.UpdateProfile); err != nil {
				log.Panic(err)
			}
			return
		default:
			bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
			command.MainMenu(bot, chatId, userId, h.Repository)
			return
		}
	}

	if body.Message.Text[0] == '/' {
		text := strings.ToLower(body.Message.Text)
		switch {
		case text == "/start" || text == "/back" || text == "/no":
			command.MainMenu(bot, chatId, userId, h.Repository)
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
		case text == "/req_admin":
			command.RequestAdm(*body, h.Repository, bot)
			command.MainMenu(bot, chatId, userId, h.Repository)

		case text == "/coba":
			command.CheckFile("data.csv")
			massage := tgbotapi.NewMessage(chatId, "messageText")
			massage.ReplyMarkup = createPaginationKeyboard(2, 4)
			bot.Send(massage)
		case text == "/absensi":
			command.Attendance(*body, bot, h.Repository)
			command.MainMenu(bot, chatId, userId, h.Repository)

		case text == "/confirm_admin":
			command.ConfirmAdm(*body, h.Repository, bot)
		default:
			bot.Send(tgbotapi.NewMessage(chatId, "I don't understand that command."))
			command.MainMenu(bot, chatId, userId, h.Repository)

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
