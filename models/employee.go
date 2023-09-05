package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserData struct {
	DocumentBase `bson:",inline"`
	UserID       int64   `bson:"user_id"`
	ChatID       int64   `bson:"chat_id"`
	FirstName    string  `bson:"first_name"`
	LastName     string  `bson:"last_name"`
	UserName     string  `bson:"username"`
	Status       *string `bson:"status"`
}

type UserQuery struct {
	ID        *primitive.ObjectID
	UserID    *int64  `bson:"user_id"`
	ChatID    *int64  `bson:"chat_id"`
	FirstName *string `bson:"first_name"`
	LastName  *string `bson:"last_name"`
}
