package domain

import "time"

type User struct {
	Id       int64
	Email    string `json:"email"`
	Password string `json:"password"`
	CreateAt time.Time
}
