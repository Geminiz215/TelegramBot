package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaginationState struct {
	Page      int
	TotalPage int
	Data      []string // The items to be paginated
}

type ent struct {
	Type string `json:"bot_command"`
}

type WebhookReqBody struct {
	UpdateID      int64          `json:"update_id"`
	Entities      *[]ent         `json:"entities"`
	Message       Message        `json:"message"`
	CallBackQuery *CallBackQuery `json:"callback_query"`
}

type From struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:""`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Lang      string `json:"language_code"`
}

type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Type      string `json:"type"`
}

type Message struct {
	MessageId int64  `json:"message_id"`
	Text      string `json:"text"`
	Date      int64  `json:"date"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type CallBackQuery struct {
	From    From    `json:"from"`
	Chat    Chat    `json:"chat"`
	Message Message `json:"message"`
	Text    string  `json:"text"`
	Date    int64   `json:"date"`
	Data    string  `json:"data"`
}

type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type DocumentBase struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Created  time.Time          `bson:"_created" json:"created,omitempty"`
	Modified time.Time          `bson:"_modified" json:"modified,omitempty"`
}

type StateKind string

type StateKinds struct {
	InsertProfile StateKind
	UpdateProfile StateKind
	Idle          StateKind
	RequestAdmin  StateKind
}

var StateKindEnum = StateKinds{
	InsertProfile: "Insert_Profile",
	RequestAdmin:  "Request_Admin",
	Idle:          "Idle",
	UpdateProfile: "Update_Profile",
}

type State struct {
	DocumentBase `bson:",inline"`
	UserID       int64       `bson:"user_id"`
	ChatID       int64       `bson:"chat_id"`
	State        string      `bson:"state"`
	SubState     string      `bson:"sub_state"`
	Data         interface{} `bson:"data"`
	Index        *int        `bson:"index"`
}

type MyData struct {
	Data []ReqAdmState `bson:"data"`
}

type DataProfile struct {
	FirstName string `bson:"first_name"`
	LasttName string `bson:"last_name"`
}
