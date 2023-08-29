package command

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func SendInsert(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI, kind models.StateKind) error {
	if err := p.StateCreate(models.State{
		UserID:   body.CallBackQuery.From.ID,
		ChatID:   body.CallBackQuery.Message.Chat.ID,
		State:    string(kind),
		SubState: "FirstName",
	}); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(body.CallBackQuery.From.ID, "Insert ur First name :"))
	return nil
}

func SendInsertFirstname(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI, state models.State) error {
	h := models.DataProfile{
		FirstName: body.Message.Text,
	}
	if err := p.UpdateState(models.State{
		State:    state.State,
		UserID:   body.Message.From.ID,
		SubState: "LastName",
		Data:     h,
	}); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(body.Message.From.ID, "Insert ur Last name :"))
	return nil
}

func SendInsertLastname(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI, state models.State) error {
	var data models.DataProfile
	d, err := bson.Marshal(state.Data)
	if err != nil {
		return err
	}
	bson.Unmarshal(d, &data)
	data.LasttName = body.Message.Text
	state.Data = data
	state.SubState = "Confirm"
	if err := p.UpdateState(state); err != nil {
		return err
	}
	bot.Send(ReplyConfirmData(bot, state.ChatID, data.FirstName, data.LasttName))
	return nil
}

func SendConfirmData(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody, state models.State) error {
	var data models.DataProfile
	d, err := bson.Marshal(state.Data)
	if err != nil {
		return err
	}
	bson.Unmarshal(d, &data)
	if _, err := p.Create(models.UserData{
		UserID:    state.UserID,
		ChatID:    state.ChatID,
		FirstName: data.FirstName,
		LastName:  data.LasttName,
		UserName:  body.CallBackQuery.From.UserName,
	}); err != nil {
		return err
	}
	p.DeleteState(state.ChatID)
	bot.Send(tgbotapi.NewMessage(state.ChatID, "Sucessed"))
	return nil
}

func UpdateProfile(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody, state models.State) error {
	var data models.DataProfile
	d, err := bson.Marshal(state.Data)
	if err != nil {
		return err
	}
	bson.Unmarshal(d, &data)
	if _, err := p.Update(models.UserData{
		UserID:    state.UserID,
		ChatID:    state.ChatID,
		FirstName: data.FirstName,
		LastName:  data.LasttName,
		UserName:  body.CallBackQuery.From.UserName,
	}); err != nil {
		return err
	}
	p.DeleteState(state.ChatID)
	bot.Send(tgbotapi.NewMessage(state.ChatID, "Sucessed"))
	return nil
}

func SendCancelData(p repository.UsersRepository, bot *tgbotapi.BotAPI, body models.WebhookReqBody, state models.State) error {
	if err := p.DeleteState(state.UserID); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(state.ChatID, "Data canceled"))
	return nil
}

func ReplyConfirmData(bot *tgbotapi.BotAPI, chatID int64, firstName string, lastName string) tgbotapi.MessageConfig {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ confirm", "/confirm"),
			tgbotapi.NewInlineKeyboardButtonData("➫ cancel", "/cancel"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("are u sure this is ur name %s %s:", firstName, lastName))
	msg.ReplyMarkup = numericKeyboard
	return msg
}
