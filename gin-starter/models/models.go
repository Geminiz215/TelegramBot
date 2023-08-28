package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	UpdateID int64  `json:"update_id"`
	Entities *[]ent `json:"entities"`
	Message  struct {
		Text string `json:"text"`
		Date int64  `json:"date"`
		From struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:""`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Lang      string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
	} `json:"message"`
	CallBackQuery *CallBackQuery `json:"callback_query"`
}

type CallBackQuery struct {
	From struct {
		ID        int64  `json:"id"`
		IsBot     bool   `json:""`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"username"`
		Lang      string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserName  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Message struct {
		Text string `json:"text"`
		Date int64  `json:"date"`
		From struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:""`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Lang      string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
	} `json:"message"`
	Text string `json:"text"`
	Date int64  `json:"date"`
	Data string `json:"data"`
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

type ActivityLogKind string

type ActivityLogTypeKinds struct {
	Login  ActivityLogKind
	Logout ActivityLogKind
}

var ActivityLogTypeEnum = ActivityLogTypeKinds{
	Login:  "SIGN_IN",
	Logout: "SIGN_OUT",
}

type StateKind string

type StateKinds struct {
	InsertProfile StateKind
	UpdateProfile StateKind
	Idle          StateKind
}

var StateKindEnum = StateKinds{
	InsertProfile: "Insert_Profile",
	Idle:          "Idle",
	UpdateProfile: "Update_Profile",
}

type ActivityLog struct {
	DocumentBase `bson:",inline"`
	UserID       int64      `bson:"user_id"`
	UserName     string     `bson:"username"`
	SignIn       *time.Time `bson:"sign_in_hour"`
	SignOut      *time.Time `bson:"sign_out_hour"`
}
type ActivityLogQuery struct {
	DocumentBase `bson:",inline"`
	UserID       *int64          `bson:"user_id"`
	UserName     *string         `bson:"username"`
	Type         ActivityLogKind `bson:"type"`
}

type State struct {
	DocumentBase `bson:",inline"`
	UserID       int64       `bson:"user_id"`
	ChatID       int64       `bson:"chat_id"`
	State        string      `bson:"state"`
	SubState     string      `bson:"sub_state"`
	Data         interface{} `bson:"data"`
}

type DataProfile struct {
	FirstName string `bson:"first_name"`
	LasttName string `bson:"last_name"`
}

func (component *State) GetDataByState() (*State, error) {
	d, err := bson.Marshal(component.Data)
	if err != nil {
		return nil, err
	}
	switch component.State {
	case string(StateKindEnum.InsertProfile):
		var data DataProfile
		bson.Unmarshal(d, &data)
		component.Data = data
	default:
		component.Data = map[string]interface{}{}
	}

	return component, nil
}
