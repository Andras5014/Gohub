package domain

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	CreateAt time.Time
}
