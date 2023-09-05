package command

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telegram-bot/models"
	"github.com/telegram-bot/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func Attendance(body models.WebhookReqBody, bot *tgbotapi.BotAPI, p repository.UsersRepository) error {
	chatId := body.Message.Chat.ID
	user, err := p.FindUser(models.UserQuery{UserID: &body.Message.From.ID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			requestInsertData(bot, chatId)
			return nil
		}
	}

	logQuery := models.ActivityLogQuery{
		UserID: &chatId,
	}
	if user.Status != nil && *user.Status == "ADMIN" {
		logQuery = models.ActivityLogQuery{}
	}
	data, count, err := p.FindLog(logQuery)
	if err != nil {
		return err
	}

	if count < 1 {
		msg := tgbotapi.NewMessage(chatId, "sign_in data not found")
		if _, err = bot.Send(msg); err != nil {
			return err
		}
		return nil
	}

	filePath := "./data/data.csv"
	err = createCSVData(data, chatId, bot, filePath)
	if err != nil {
		return err
	}
	csv := tgbotapi.NewDocument(chatId, tgbotapi.FilePath(filePath))
	if _, err = bot.Send(csv); err != nil {
		return err
	}

	// // Remove the CSV file after it's sent
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func createCSVData(activity []models.ActivityLog, chatId int64, bot *tgbotapi.BotAPI, filePath string) error {
	// Create a CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Write the CSV header
	headers := []string{"Name", "Sign_in_hour", "Sign_out_hour"}
	if err := csvWriter.Write(headers); err != nil {
		return err

	}

	for _, i := range activity {
		if i.SignIn == nil || i.SignOut == nil {
			continue
		}
		row := []string{i.UserName, i.SignIn.Format("Monday, January 2, 2006 15:04:05 MST"), i.SignOut.Format("Monday, January 2, 2006 15:04:05 MST")}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func CheckFile(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist.")
		} else {
			log.Fatal(err)
		}
		return
	}

	// Check if the file is empty
	if fileInfo.Size() == 0 {
		fmt.Println("File is empty.")
	} else {
		fmt.Printf("File size: %d bytes\n", fileInfo.Size())
	}
}
