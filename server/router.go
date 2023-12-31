package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/telegram-bot/connection"
	"github.com/telegram-bot/handlers"
	"github.com/telegram-bot/repository"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	logger := log.New(os.Stdout, "go-api:", log.LstdFlags)
	res, err := connection.TelegramConnection()
	if err != nil {
		logger.Println(err)
	}
	logger.Println("Response from the initial webhook status:", res.Status)

	//SET MONGODB CONNECTION
	db, err := connection.ConnectMongoDB()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
	}

	TeleRepository := repository.UsersRepositoryMongo{
		ConnectionDB: db,
	}

	webhookHand := handlers.WebhookController{
		Repository: &TeleRepository,
	}

	health := new(handlers.HealthController)
	router.GET("/health", health.Status)
	router.POST("/ping", health.Status)
	router.POST("/webhook", webhookHand.WebhookCallback)

	return router

}
