package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, user *User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	return dao.db.WithContext(ctx).Create(user).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func (User) TableName() string {
	return "user"
}
