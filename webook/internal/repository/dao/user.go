package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	CreateAt int64
	UpdateAt int64
}

type UserProfile struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	UserId   int64 `gorm:"unique"`
	NickName string
	Birthday string
	AboutMe  string
	CreateAt int64
	UpdateAt int64
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreateAt = now
	user.UpdateAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error

	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) InsertProfile(ctx context.Context, userProfile UserProfile) error {
	now := time.Now().UnixMilli()
	userProfile.CreateAt = now
	userProfile.UpdateAt = now
	err := dao.db.WithContext(ctx).Create(&userProfile).Error
	return err
}

func (dao *UserDao) SetProfile(ctx context.Context, userProfile UserProfile) error {
	now := time.Now().UnixMilli()
	userProfile.UpdateAt = now
	err := dao.db.WithContext(ctx).Save(&userProfile).Error
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err

}

func (dao *UserDao) FindProfileById(ctx context.Context, userId int64) (UserProfile, error) {
	var userProfile UserProfile
	err := dao.db.WithContext(ctx).Where("user_id=?", userId).First(&userProfile).Error
	return userProfile, err

}
