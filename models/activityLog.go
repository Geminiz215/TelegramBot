package models

import "time"

type ActivityLogKind string

type ActivityLogTypeKinds struct {
	Login  ActivityLogKind
	Logout ActivityLogKind
}

var ActivityLogTypeEnum = ActivityLogTypeKinds{
	Login:  "SIGN_IN",
	Logout: "SIGN_OUT",
}

type ActivityLog struct {
	DocumentBase `bson:",inline"`
	UserID       int64      `bson:"user_id"`
	UserName     string     `bson:"username"`
	SignIn       *time.Time `bson:"sign_in_hour"`
	SignOut      *time.Time `bson:"sign_out_hour"`
	MessageID    int64      `bson:"message_id"`
}
type ActivityLogQuery struct {
	DocumentBase `bson:",inline"`
	UserID       *int64          `bson:"user_id"`
	UserName     *string         `bson:"username"`
	Type         ActivityLogKind `bson:"type"`
	MessageID    int64           `bson:"message_id"`
}
