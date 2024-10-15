package models

import "time"

type User struct {
	ID        string `gorm:"primary_key" json:"user_id"`
	Hash      string
	UpdatedAt time.Time
}

type UserCred struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
