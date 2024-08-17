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

func (repo *UserRepository) toDomainProfile(u dao.UserProfile) domain.UserProfile{
	return domain.UserProfile{
		UserId: u.UserId,
		NickName: u.NickName,
		Birthday: u.Birthday,
		AboutMe: u.AboutMe,
		RestParam: domain.RestParam{
			Id: u.Id,
			CreateAt: u.CreateAt,
			UpdateAt: u.UpdateAt,
		},
	}

}

func (repo *UserRepository) FindProfileById(ctx context.Context, userId int64) (domain.UserProfile, error){
	userProfile, err  := repo.dao.FindProfileById(ctx, userId)
	if err!=nil{
		  return domain.UserProfile{}, err
	}
	return  repo.toDomainProfile(userProfile), nil

}

func (repo *UserRepository) CreateProfile(ctx context.Context, userProfile domain.UserProfile) error{
	err := repo.dao.InsertProfile(ctx, dao.UserProfile{
		 UserId: userProfile.UserId,
		 NickName: userProfile.NickName,
		 Birthday: userProfile.Birthday,
		 AboutMe: userProfile.AboutMe,
	})
	return err
}

func (repo *UserRepository) OverwriteProfle(ctx context.Context, userProfile domain.UserProfile) error{

	err := repo.dao.SetProfile(ctx, dao.UserProfile{
		Id: userProfile.RestParam.Id,
		UserId: userProfile.UserId,
		NickName: userProfile.NickName,
		Birthday: userProfile.Birthday,
		AboutMe: userProfile.AboutMe,
		CreateAt: userProfile.RestParam.CreateAt,
   })
   return err
	
}