package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"project0/internal/domain"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"time"
)

// 预定义错误
var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) (int64, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.UserProfile, error)
	CreateProfile(ctx context.Context, profile domain.UserProfile) error
	Edit(ctx *gin.Context, profile domain.UserProfile) error
	FindByPhone(ctx *gin.Context, phone sql.NullString) (domain.User, error)
	FindByWechat(ctx context.Context, openId sql.NullString) (domain.User, error)
}
type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func (repos *CacheUserRepository) FindByWechat(ctx context.Context, openId sql.NullString) (domain.User, error) {
	u, err := repos.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}

	repos.dao.Insert(ctx, u)

	return toDomain(u), err
}

func NewCacheUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {

	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

// 用户注册  tbcsql.null
func (repos *CacheUserRepository) CreateUser(ctx context.Context, u domain.User) (int64, error) {
	// 在这里分离插入user profile的操作
	return repos.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			u.Email,
			u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			u.Phone,
			u.Phone != "",
		},
		WechatOpenId: sql.NullString{
			u.WechatInfo.OpenID,
			u.WechatInfo.OpenID != "",
		},
		WechatUnionId: sql.NullString{
			u.WechatInfo.UnionID,
			u.WechatInfo.UnionID != "",
		},
	})

	//if err != nil {
	//
	//}
	//return nil
}

/*func (repos *CacheUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repos.dao.FindById(ctx, uid)
}*/

// aka find  by email
func (repos *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//u, err := repos.cache.Get(ctx, email)

	u, err := repos.dao.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}
	//直接早期就引入zap吗，不用log.Println
	//log.Println("find by email u whether is nil ", u)
	return toDomainLogin(u), nil
}

func toDomainLogin(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
	}
}

func toDomain(u dao.User) domain.User {
	return domain.User{
		Id:    u.Id,
		Email: u.Email.String,
		Phone: u.Phone.String,
		WechatInfo: domain.WechatInfo{
			OpenID:  u.WechatOpenId.String,
			UnionID: u.WechatUnionId.String,
		},
	}
}

// 优化profile  先走缓存， 有种特别的情况就使用 switch 确保高可用，避免数据库被打崩。
// 缓存击穿/穿透的解决方法 当set 确认redis不正常就不查数据库
func (repos *CacheUserRepository) Profile(ctx context.Context, id int64) (domain.UserProfile, error) {
	du, err := repos.cache.Get(ctx, id)

	// 有缓存就返回
	if err == nil {
		return du, nil
	}

	profile, err := repos.dao.Profile(ctx, id)
	if err != nil {
		log.Println("repository profile, :", err)
		return domain.UserProfile{}, err
	}
	du = toDomainLoginProfile(profile)
	err = repos.cache.Set(ctx, du)
	if err != nil {
		// 网络可能蹦了
		log.Println(err)
	}
	return du, nil
}

// 这里没有使用
func (repos *CacheUserRepository) CreateProfile(ctx context.Context, profile domain.UserProfile) error {
	id := profile.Id
	var upInstance dao.UserProfile
	return repos.dao.CreateProfile(ctx, upInstance, id)
}

func toDomainLoginProfile(u dao.UserProfile) domain.UserProfile {

	birthUnix := u.BirthDate * int64((time.Millisecond))
	t := time.Unix(0, birthUnix)

	return domain.UserProfile{
		Id:           u.Id,
		Gender:       u.Gender,
		NickName:     u.NickName,
		Introduction: u.Introduction,
		BirthDate:    t,
	}
}

func (repos *CacheUserRepository) Edit(ctx *gin.Context, profile domain.UserProfile) error {

	err := repos.dao.Edit(ctx, profile)

	return err
}

// 在repository 层改造sql null string 应该可以 加一个toEntity
func (repos *CacheUserRepository) FindByPhone(ctx *gin.Context, phone sql.NullString) (domain.User, error) {

	u, err := repos.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	//var u dao.User
	u.Phone = phone
	// 这样可以吗
	repos.dao.Insert(ctx, dao.User{
		Phone: phone,
	})

	return toDomain(u), err

}
