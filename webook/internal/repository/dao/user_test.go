package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDao_Insert(t *testing.T) {

	tests := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		ctx     context.Context
		user    User
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "insert success",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1,1)
				mock.ExpectExec("INSERT INTO .*" ).WillReturnResult(mockRes)

				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
		},

		{
			name: "duplicate email",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				
				mock.ExpectExec("INSERT INTO .*" ).WillReturnError(&mysqlDriver.MySQLError{Number: 1062})

				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			name: "db error",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				
				mock.ExpectExec("INSERT INTO .*" ).WillReturnError(errors.New("db error"))

				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: errors.New("db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB := tt.mock(t)
			db, err := gorm.Open(mysql.New(
				mysql.Config{
					Conn:                      sqlDB,
					SkipInitializeWithVersion: true,
				}),
				&gorm.Config{
					DisableAutomaticPing:   true,
					SkipDefaultTransaction: true,
				})
			assert.NoError(t, err)
			dao := NewUserDao(db)
			err = dao.Insert(tt.ctx, tt.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
