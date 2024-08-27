package repository

import (
	"context"
	"database/sql"
	"log"

	"time"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}

}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))

}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil

}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String, 
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}

}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	//createat 和 updataat就不更新了？
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{
			String: u.Email,
			Valid: u.Email!="",

		},
		Phone:    sql.NullString{
			String: u.Phone,
			Valid: u.Phone!="",

		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.User) error {
	// 更新 DB 之后，删除
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	if err != nil {
		return err
	}
	return nil

}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {

	du, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return du, err
	}

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}

	du = repo.toDomain(u)
	// set cache
	err = repo.cache.Set(ctx, du)
	if err != nil {
		log.Println(err)
	}

	return du, nil
}

func (repo *UserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {

	du, err := repo.cache.Get(ctx, uid)


	switch err{
	case nil:
		return du, err
	case cache.ErrKeyNotExist:
		u, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}

		du = repo.toDomain(u)
		// set cache
			// set cache
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}

		return du, nil
	
	default:
		return domain.User{}, err 

	}

}


func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u,err := repo.dao.FindByPhone(ctx, phone)
	if err!=nil{
		return domain.User{}, err
	}
	return repo.toDomain(u),nil

}