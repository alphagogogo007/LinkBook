package dao

import (
	"context"
	"time"

	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

type AsyncSms struct {
	// 不标注是因为gorm规则会自动替换
	Id       int64
	Config   sqlx.JsonColumn[SmsConfig]
	RetryCnt int
	RetryMax int
	Status   uint8
	Ctime    int64
	Utime    int64 `gorm:"index"` //increase search speed
}

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}

//go:generate mockgen -source=./async_sms.go -package=daomocks -destination=mocks/async_sms.mock.go AsyncSmsDAO
type AsyncSmsDAO interface {
	Insert(ctx context.Context, s AsyncSms) error
	GetWaitingSMS(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

const (
	asyncStatusWaiting = iota
	asyncStatusFailed
	asyncStatusSuccess
)

type GORMAsyncSmsDAO struct {
	db *gorm.DB
}

func NewGORMAsyncSmsDAO(db *gorm.DB) AsyncSmsDAO {
	return &GORMAsyncSmsDAO{
		db: db,
	}
}

func (g *GORMAsyncSmsDAO) Insert(ctx context.Context, s AsyncSms) error {
	return g.db.Create(&s).Error
}

func (g *GORMAsyncSmsDAO) GetWaitingSMS(ctx context.Context) (AsyncSms, error) {
	
	var s AsyncSms
	err := g.db.WithContext(ctx).Transaction(func (tx *gorm.DB) error{
		now := time.Now().UnixMilli()
		endTime := now-time.Minute.Milliseconds()
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("utime < ? and status = ?",
			endTime, asyncStatusWaiting).First(&s).Error
		if err!=nil{
			return err
		}

		//there is risk that many machines read the same record
		// Then they update the record several times.
		err = tx.Model(&AsyncSms{}).
			Where("id=?", s.Id).
			Updates(map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt+1"),
				"utime": now,

			}).Error
		return err
	})
	return s,err
}

func (g *GORMAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&AsyncSms{}).
		Where("id=?", id).Updates(
		map[string]any{
			"utime":  now,
			"status": asyncStatusSuccess,
		}).Error
}

func (g *GORMAsyncSmsDAO) MarkFailed(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&AsyncSms{}).
		Where("id =? and `retry_cnt`>=`retry_max`", id).Updates(
		map[string]any{
			"utime":  now,
			"status": asyncStatusFailed,
		}).Error
}
