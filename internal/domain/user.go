package domain

import (
	"time"
)

type User struct {
	Id        int64
	Email     string
	NickName  string
	Phone     string
	Password  string
	AboutMe   string
	Birthday  time.Time
	CreatedAt time.Time
}
