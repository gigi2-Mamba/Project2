package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"project0/internal/domain"
	"project0/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Profile(ctx *gin.Context, uid int64) error
	//TBCA SQL.NULL STRING
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
	Edit(ctx *gin.Context, profile domain.UserProfile) error
	CreateUser(ctx context.Context, u domain.User) error
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

// 业务逻辑层调用 数据抽象层，直达数据访问层
type userService struct {
	repos repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repos: repo,
	}
}

// 登录业务逻辑 ,传了go的context 要干嘛呢
func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost) // 指定加密代价。
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return svc.CreateUser(ctx, u)
}

func (svc *userService) CreateUser(ctx context.Context, u domain.User) error {
	uid, err := svc.repos.CreateUser(ctx, u)
	if err != nil {
		// 解决掉这手动打日志的问题
		//log.Println("Sign up create user failed: ", err)
		return err
	}
	var up = domain.UserProfile{
		Id: uid,
	}
	return svc.repos.CreateProfile(ctx, up)
}

// 在service 手动转化string -- > sql.null string
func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	phone0 := sql.NullString{
		String: phone,
		Valid:  phone != "",
	}
	u, err := svc.repos.FindByPhone(ctx, phone0)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	// 走到这里，要注册
	//err = svc.repos.CreateUser(ctx, domain.User{
	//	Phone: phone,
	//})
	//err = svc.Create(ctx, domain.User{
	//	Phone: phone,
	//})
	err = svc.CreateUser(ctx, u)

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 能找到直接login
	// 可能会发生主从延迟
	return svc.repos.FindByPhone(ctx, phone0)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	openId := sql.NullString{
		String: wechatInfo.OpenID,
		Valid:  wechatInfo.OpenID != "",
	}

	u, err := svc.repos.FindByWechat(ctx, openId)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	// 这边意味新用户进来
	// json格式的
	zap.L().Info("新用户", zap.Any("wechatInfo", wechatInfo))

	err = svc.CreateUser(ctx, u)

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 能找到直接login
	// 可能会发生主从延迟
	return svc.repos.FindByWechat(ctx, openId)

}
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repos.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *userService) Profile(ctx *gin.Context, uid int64) error {

	profile, err := svc.repos.Profile(ctx, uid)
	// 以后再改
	ctx.JSON(http.StatusOK, profile)

	return err

}

func (svc *userService) Edit(ctx *gin.Context, profile domain.UserProfile) error {

	return svc.repos.Edit(ctx, profile)
}
