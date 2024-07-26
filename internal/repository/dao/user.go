package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("email duplicated")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	return user, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (*User, error) {
	user := &User{}
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	return user, err
}

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (*User, error) {
	user := &User{}
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(user).Error
	return user, err
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}
