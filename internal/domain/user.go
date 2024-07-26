package domain

import (
	"time"
)

type User struct {
	Id       int64
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	CreateAt time.Time
}
