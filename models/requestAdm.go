package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestAdmin struct {
	DocumentBase `bson:",inline"`
	UserID       int64  `bson:"user_id"`
	UserName     string `bson:"username"`
}

type RequestAdminQuery struct {
	UserID   *int64
	ID       *primitive.ObjectID
	UserName *string
}

type ReqAdmState struct {
	UserID   int64  `bson:"user_id"`
	Accept   bool   `bson:"accept"`
	Username string `bson:"username"`
}
