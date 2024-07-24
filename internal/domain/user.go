package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Id       int64
	Email    string `json:"email"`
	Password string `json:"password"`
	CreateAt time.Time
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
