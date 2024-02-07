package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// mock模拟的实际依赖项
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				mockRes := sqlmock.NewResult(123, 1)
				assert.NoError(t, err)
				// 有分隔符就可以换行
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Ctime: 267766,
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				//mockRes := sqlmock.NewResult(123,1)
				assert.NoError(t, err)
				// 有分隔符就可以换行
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysql2.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Ctime: 26776,
			},
			wantErr: ErrDuplicateUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDb := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDb,
				// 跳过查询mysql的版本
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)

			dao := NewUserDAO(db)
			_, err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}
