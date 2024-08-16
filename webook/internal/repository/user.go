package repository

import (
	"context"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound = dao.ErrRecordNotFound
)

type UserRepository struct{
	dao *dao.UserDao
}

 
func NewUserRepository(dao *dao.UserDao) *UserRepository{
	return &UserRepository{
		dao: dao,
	}

}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error{
	return  repo.dao.Insert(ctx, dao.User{
		Email: user.Email,
		Password: user.Password,
	})

}


func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error){
	user, err  := repo.dao.FindByEmail(ctx, email)
	if err!=nil{
		  return domain.User{}, err
	}
	return  repo.toDomain(user), nil

}

func (repo *UserRepository) toDomain(u dao.User) domain.User{
	return domain.User{
		Id: u.Id,
		Email: u.Email,
		Password: u.Password,
	}

}