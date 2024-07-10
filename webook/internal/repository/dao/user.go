package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱已被注册")
	ErrUserNotFound       = gorm.ErrRecordNotFound
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

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//同一UTC毫秒数，消除时区问题
	u.Utime = time.Now().UnixMilli()
	u.Ctime = time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueConflictsErrNo uint16 = 1062
		if me.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}
