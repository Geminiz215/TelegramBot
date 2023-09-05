package command

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RequestAdm(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI) error {
	userId := body.Message.From.ID
	chatId := body.Message.Chat.ID
	_, count, err := p.FindReqAdmin(models.RequestAdminQuery{
		UserID: &userId,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		bot.Send(tgbotapi.NewMessage(chatId, "request still exist"))
		return err
	}

	if count >= 1 {
		bot.Send(tgbotapi.NewMessage(chatId, "request still exist"))
		return nil
	}

	if err := p.CreateReqAdmin(models.RequestAdmin{
		UserID:   userId,
		UserName: body.Message.From.UserName,
	}); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(chatId, "Sucessed"))
	return nil
}

func ConfirmAdm(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI) error {
	adm, err := p.FindUser(models.UserQuery{
		UserID: &body.Message.From.ID,
	})
	if err != nil {
		return err
	}
	if adm.Status == nil || *adm.Status != "ADMIN" {
		bot.Send(tgbotapi.NewMessage(body.Message.Chat.ID, "cannot use command"))
	}
	_, count, err := p.FindReqAdmin(models.RequestAdminQuery{})
	if err != nil {
		return err
	}
	bot.Send(confirmAdmin(bot, body.Message.Chat.ID, count))
	return nil

}

func confirmAdmin(bot *tgbotapi.BotAPI, chatID int64, count int64) tgbotapi.MessageConfig {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ confirm", "/accept"),
			tgbotapi.NewInlineKeyboardButtonData("➫ back", "/back"),
		),
	)
	if count < 1 {
		numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➫ back", "/cancel"),
			),
		)
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("you have %d request admin\n do you wanna start it:", count))
	msg.ReplyMarkup = numericKeyboard
	return msg
}

func AcceptReqAdm(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI) error {
	index := int(1)
	if err := p.StateCreate(models.State{
		UserID:   body.CallBackQuery.From.ID,
		ChatID:   body.CallBackQuery.Message.Chat.ID,
		State:    string(models.StateKindEnum.RequestAdmin),
		SubState: "Accept",
		Index:    &index}); err != nil {
		return err
	}
	data, _, err := p.FindReqAdmin(models.RequestAdminQuery{})
	if err != nil {
		return err
	}
	_, err = bot.Send(confirmReqAdmin(bot, body.CallBackQuery.Message.Chat.ID, data[index-1]))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func confirmReqAdmin(bot *tgbotapi.BotAPI, chatID int64, data models.RequestAdmin) tgbotapi.MessageConfig {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ confirm", "/confirm"),
			tgbotapi.NewInlineKeyboardButtonData("➫ reject", "/reject"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("username request: %s\nuserId request: %d", data.UserName, data.UserID))
	msg.ReplyMarkup = numericKeyboard
	return msg
}

func AcceptAdminRequest(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI, state models.State) error {
	userId := body.CallBackQuery.From.ID
	chatId := body.CallBackQuery.Message.Chat.ID
	//get data state
	data := make([]models.ReqAdmState, 0)
	if state.Data != nil {
		for _, dt := range state.Data.(bson.A) {
			d, _ := bson.Marshal(dt)
			var item models.ReqAdmState
			bson.Unmarshal(d, &item)
			data = append(data, item)
		}
	}

	req, count, err := p.FindReqAdmin(models.RequestAdminQuery{})
	if err != nil {
		return nil
	}

	//get value
	accept := false
	if body.CallBackQuery.Data == "/confirm" {
		accept = true
	}

	//insert data by index
	data = append(data, models.ReqAdmState{
		UserID:   req[*state.Index-1].UserID,
		Username: req[*state.Index-1].UserName,
		Accept:   accept,
	})

	if count == int64(*state.Index) {
		bot.Send(replyConfirmAdminRequest(bot, chatId, data))
	} else {
		fmt.Println(*state.Index-1, "this is index dex")
		bot.Send(confirmReqAdmin(bot, chatId, req[*state.Index]))
	}
	//update state
	if count == int64(*state.Index) {
		state.SubState = "Confirm"
	} else {
		*state.Index += 1
	}
	if err := p.UpdateState(models.State{
		State:    state.State,
		UserID:   userId,
		Data:     data,
		SubState: state.SubState,
		Index:    state.Index,
	}); err != nil {
		return err
	}

	return nil
}

func replyConfirmAdminRequest(bot *tgbotapi.BotAPI, chatID int64, data []models.ReqAdmState) tgbotapi.MessageConfig {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➫ confirm", "/confirm"),
			tgbotapi.NewInlineKeyboardButtonData("➫ cancel", "/cancel"),
		),
	)
	r := []string{}
	for _, p := range data {
		r = append(r, fmt.Sprintf("confirm %s as admin %v", p.Username, p.Accept))
	}
	msg := tgbotapi.NewMessage(chatID, strings.Join(r, "\n"))
	msg.ReplyMarkup = numericKeyboard
	return msg
}

func ConfirmUpdateAdmin(body models.WebhookReqBody, p repository.UsersRepository, bot *tgbotapi.BotAPI, state models.State) error {
	var data []models.ReqAdmState
	for _, dt := range state.Data.(bson.A) {
		d, _ := bson.Marshal(dt)
		var item models.ReqAdmState
		bson.Unmarshal(d, &item)
		data = append(data, item)
	}
	_, err := p.UpdateManyUser(data)
	if err != nil {
		return nil
	}
	err = p.DeleteState(state.UserID)
	if err != nil {
		return nil
	}

	bot.Send(tgbotapi.NewMessage(body.CallBackQuery.Message.Chat.ID, "sucessed"))
	return nil
}
