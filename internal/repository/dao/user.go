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

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Update(ctx context.Context, user User) error
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}
func (dao *GormUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = sql.NullInt64{Int64: now, Valid: true}
	user.UpdatedAt = sql.NullInt64{Int64: now, Valid: true}
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

func (dao *GormUserDAO) Update(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.UpdatedAt = sql.NullInt64{Int64: now, Valid: true}
	return dao.db.WithContext(ctx).Updates(&user).Error
}
func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	user := User{}
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	user := User{}
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}
func (dao *GormUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	user := User{}
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&user).Error
	return user, err
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	NickName sql.NullString
	Birthday sql.NullInt64
	AboutMe  sql.NullString `gorm:"type:varchar(1024)"`

	// wechat 字段
	WechatUnionId sql.NullString `gorm:"unique"`
	WechatOpenId  sql.NullString `gorm:"unique"`

	CreatedAt sql.NullInt64
	UpdatedAt sql.NullInt64
	DeletedAt sql.NullInt64
}
