package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface{
	Insert(ctx context.Context, user User) error
	UpdateById(ctx context.Context, entity User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, uid int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	Phone sql.NullString `gorm:"unique"`

	WechatOpenId sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString

	CreateAt int64
	UpdateAt int64
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (dao *GORMUserDao) Insert(ctx context.Context, user User) error {
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

func (dao *GORMUserDao) UpdateById(ctx context.Context, entity User) error {

	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"create_at": time.Now().UnixMilli(),
			"nickname":  entity.Nickname,
			"birthday":  entity.Birthday,
			"about_me":  entity.AboutMe,
		}).Error
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err

}

func (dao *GORMUserDao) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}


func (dao *GORMUserDao) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}