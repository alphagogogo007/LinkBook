package domain

type User struct{
	Id int64
	Email string
	Password string
}

type UserProfile struct{
	UserId int64
	NickName string
	Birthday string
	AboutMe string
	RestParam RestParam
}

type RestParam struct{
	Id int64
	CreateAt int64
	UpdateAt int64
}

type FrontProfile struct{

	NickName string
	Birthday string
	AboutMe string
}

