package models

import "time"

type Audit struct{
	Action string `json:"action"`
	UserId int `json:"user_id"`
	Time time.Time `json:"time"`
	UserName string `json:"user_name"`
	UserAge int `json:"user_age"`
}

