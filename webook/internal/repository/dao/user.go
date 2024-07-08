package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int    `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	//创建时间、修改时间
	Ctime int64
	Utime int64
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (ud *UserDAO) Insert(ctx context.Context, u User) error {
	//同一UTC毫秒数，消除时区问题
	u.Utime = time.Now().UnixMilli()
	u.Ctime = time.Now().UnixMilli()
	return ud.db.WithContext(ctx).Create(&u).Error
}
