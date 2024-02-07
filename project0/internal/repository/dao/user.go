package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
	"project0/internal/domain"
	"time"
)

// 预定义错误
var (
	ErrDuplicateUser  = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound //  利用中间件提供的错误码，对数据库报错信息进行反馈
)

type User struct {
	Id    int64          `gorm:"primaryKey,autoIncrement"`
	Email sql.NullString `gorm:"unique"`
	//Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	//  时区
	//
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
	Ctime         int64
	Utime         int64
}

type UserProfile struct {
	Id       int64  `gorm:"primary"`
	NickName string `gorm:"unique"`
	Gender   string

	Introduction string
	BirthDate    int64
}
type UserDao interface {
	Insert(ctx context.Context, u User) (int64, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone sql.NullString) (User, error)
	FindByWechat(ctx context.Context, openId sql.NullString) (User, error)
	CreateProfile(ctx context.Context, up UserProfile, id int64) error
	Edit(ctx *gin.Context, profile domain.UserProfile) error
	Profile(ctx context.Context, id int64) (UserProfile, error)
}
type GORMUserDao struct {
	db *gorm.DB
}

func (dao *GORMUserDao) Insert(ctx context.Context, u User) (int64, error) {
	//log.Println("dao is nil ", dao == nil)
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	//log.Println("can here ")
	log.Println("user whether is nil ", u)
	err := dao.db.WithContext(ctx).Create(&u).Error
	log.Println("can here 2")
	if me, ok := err.(*mysql.MySQLError); ok { // 驱动包的error
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			//用户冲突吧
			return 0, ErrDuplicateUser
		}
	}

	return u.Id, nil
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone sql.NullString) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).Find(&u).Error

	return u, err
}

func (dao *GORMUserDao) FindByWechat(ctx context.Context, openId sql.NullString) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id", openId).Find(&u).Error
	return u, err
}
func NewUserDAO(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (u *GORMUserDao) CreateProfile(ctx context.Context, up UserProfile, id int64) error {
	up.Id = id
	return u.db.WithContext(ctx).Create(&up).Error
}

func (u *GORMUserDao) Profile(ctx context.Context, id int64) (UserProfile, error) {

	var uprofile UserProfile
	err := u.db.WithContext(ctx).Where("id=?", id).First(&uprofile).Error
	return uprofile, err

}

func (dao *GORMUserDao) Edit(ctx *gin.Context, profile domain.UserProfile) error {
	uprofile := &UserProfile{
		Id:           profile.Id,
		Gender:       profile.Gender,
		BirthDate:    profile.BirthDate.Unix(),
		Introduction: profile.Introduction,
		NickName:     profile.NickName,
	}

	return dao.db.WithContext(ctx).Updates(&uprofile).Error

}

//func (dao *GORMUserDao) ProfileFindById(ctx context.Context, uid int64) (interface{}, interface{}) {
//
//}

//func (u *GORMUserDao) Edit(ctx context.Context, up)

//func (u *GORMUserDao) Edit(ctx context.Context, up *UserProfile) (domain.UserProfile, error) {
//
//}
