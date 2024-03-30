package web

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project0/internal/domain"
	"project0/internal/errs"
	"project0/internal/service"
	"project0/internal/web/ijwt"
	"project0/pkg/ginx"
	"time"
)

// 构建一个实例
const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "bizLogin"
)

// 往handler里注入服务
// 验证码服务，用户服务，邮箱，密码格式验证
type UserHandler struct { // 怎么选择把正则表达式放在 UserHandler里
	jwtHandler     ijwt.Handler
	svc            service.UserService // 注入了 对应的service对象
	emailRegexp    *regexp2.Regexp
	passwordRegexp *regexp2.Regexp
	code           service.CodeService
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, handler ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRegexp:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegexp: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		svc:            svc,
		code:           codeSvc,
		jwtHandler:     handler,
	}

}
func (u *UserHandler) RegisterRoute(server *gin.Engine) {

	ug := server.Group("/users")
	// 利用函数装饰器的魅力
	ug.POST("/signup", ginx.WrapBody(u.Signup))
	ug.POST("/login", ginx.WrapBody(u.LoginJWT))
	//ug.POST("/login", u.Login)
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", ginx.WrapBodyAndClaims(u.Edit))
	ug.GET("/refresh_token", u.RefreshToken)
	//ug.POST("/jwtLogin", u.LoginJWT)

	// 手机验证码登录
	ug.POST("/login_sms/send/code0", ginx.WrapBody(u.SendLoginCode))
	ug.POST("/login_sms0", ginx.WrapBody(u.LoginSms)) // 验证验证码
	ug.POST("/logout", u.LogoutJWT)
}

// 注册 注册路由
func (u *UserHandler) Signup(ctx *gin.Context, req SignUpReq) (ginx.Result, error) { //  因为 type HandleFunc  func(*context){}

	isEmail, err := u.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "非法邮箱格式",
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入密码不一致",
		}, nil
	}

	isPwd, err := u.passwordRegexp.MatchString(req.Password)
	//err = errors.New("create an error")
	if err != nil {
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	//log.Println("ISpwd ,", isPwd, req.Password)
	if !isPwd {
		return ginx.Result{
			// 那这里还status ok？
			Code: errs.UserInvalidInput,
			Msg:  "请至少包含数字，特殊字符，字母，整体长度8到16",
		}, nil
	}
	// 保证万无一失，使用这个context
	err = u.svc.Signup(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	//log.Println("sign req", req)
	switch err {
	case nil:
		//sessions.Default(ctx) tbc  omit
		//ctx.String(http.StatusOK, "注册成功")
		return ginx.Result{
			Msg: "OK",
		}, nil
	case service.ErrDuplicateEmail:
		return ginx.Result{
			//Code: errs.UserDuplicateEmail,
			Code: errs.UserDuplicateEmail,
			Msg:  "邮箱冲突",
		}, nil
	default:
		return ginx.Result{
			Msg: "系统错误",
		}, err
	}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	// 接受前端数据要son格式  使用bind（）
	if err := ctx.Bind(&req); err != nil {
		return //  不用给信息，gin会给
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", user.Id)
		sess.Options(sessions.Options{
			MaxAge: 60, // 控制了 redis对的过期时间
			//HttpOnly:
		})
		err = sess.Save() //  gin session设置要求主动保存
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (u *UserHandler) Edit(ctx *gin.Context, req UserProfileReq, uc ijwt.UserClaims) (ginx.Result, error) {

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		return ginx.Result{
			Code: 4,
			Msg:  "日期格式错误",
		}, err
	}

	if req.Nickname == "" || len(req.Nickname) > 24 {
		return ginx.Result{
			Code: 4,
			Msg:  "昵称为空，或者超出范围",
		}, err
	}

	//id,err := strconv.Atoi(ctx.Query("id"))
	// 通过session--
	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId")
	//v, ok := uid.(int64)
	//if !ok {
	//	ctx.String(http.StatusOK, "请求参数有误")
	//	return
	//}
	// 通过session--

	//uc := ctx.MustGet("user")
	//log.Println("不可能来到这里")
	//uc0 := uc.(ijwt.UserClaims)

	err = u.svc.Edit(ctx, domain.UserProfile{
		Id:           uc.Uid,
		Gender:       req.Gender,
		NickName:     req.Nickname,
		Introduction: req.Introduction,
		BirthDate:    birthday,
	})

	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {

	err := u.jwtHandler.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		//log.Println("Logout JWT faield ",err)
		zap.L().Error("Logout JWT failed", zap.Error(err))
	}
	ctx.String(http.StatusOK, "退出登录成功")
}

// 前端访问这个路由刷新token
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 先从前端去token，并解析
	tokenStr := u.jwtHandler.ExtraToken(ctx)
	var uc ijwt.RefreshCliams
	// 研究 解析claim方法它做了什么？
	token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RefreshKey, nil
	})

	if err != nil {
		log.Println("refresh token 不正确")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	if token != nil || !token.Valid {
		log.Println("refresh token invalid")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = u.jwtHandler.CheckSession(ctx, uc.Ssid)
	if err != nil {
		//token 无效 或者redis出错
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//长token自动延期
	// 如果设置7天登录一次，就要做变化
	err = u.jwtHandler.SetJWT(ctx, uc.Uid, uc.Ssid)
	if err != nil {
		log.Println("setJwt failed: ", err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})

}
func (u *UserHandler) Profile(ctx *gin.Context) {
	//birthday, err := time.Parse(time.DateOnly, req.Birthday)

	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId")
	//v, ok := uid.(int64)

	//us := ctx.MustGet("user").(ijwt.UserClaims)
	us := ctx.MustGet("user")

	uc := us.(ijwt.UserClaims)
	//v, ok := us.Uid.(int64)
	//if !ok || v == 0 {
	//	ctx.String(http.StatusOK, "请求参数有误")
	//	return
	//}
	//err := u.svc.Profile(ctx, v)

	err := u.svc.Profile(ctx, uc.Uid)
	if err != nil {
		log.Println("系统错误，个人简介获取", err)
		ctx.String(http.StatusOK, "系统错误，个人简介")
		return
	}
	ctx.String(http.StatusOK, "个人简介展示成功")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context, req LoginReq) (ginx.Result, error) {

	user, err := u.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		//sess := sessions.Default(ctx)
		//sess.Set("userId", user.Id)
		//sess.Options(sessions.Options{
		//	MaxAge: 10, // 控制了 redis对的过期时间
		//	//HttpOnly:
		//})
		//err = sess.Save() //  gin session设置要求主动保存
		//if err != nil {
		//	ctx.String(http.StatusOK, "系统错误")
		//}
		err := u.jwtHandler.SetLoginJWTToken(ctx, user.Id)
		if err != nil {
			return ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			}, err
		}

		return ginx.Result{
			Msg: "登录成功",
		}, err
	case service.ErrInvalidUserOrPassword:
		return ginx.Result{
			Code: 4,
			Msg:  "用户名或密码不对",
		}, nil
	default:
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
}

func (u *UserHandler) SendLoginCode(ctx *gin.Context, req SendLoginCodeReq) (ginx.Result, error) {

	if req.Phone == "" {
		return ginx.Result{
			Code: 4,
			Msg:  "请输入手机号",
		}, nil
	}
	err, _ := u.code.SendFaker(ctx, bizLogin, req.Phone)

	switch err {
	case nil:
		return ginx.Result{
			Msg: "发送成功",
		}, nil
	case service.ErrCodeSendTooMany:
		//ctx.JSON(http.StatusOK, Result{
		//	Msg: "短信发送太频繁，请稍后再试",
		//})
		//频繁但是正常找不到原因
		//一直有warning的话要排查
		//zap.L().Warn("频繁发送验证码")
		return ginx.Result{
			Code: 5,
			Msg:  "发送验证码服务端繁忙",
		}, err
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		zap.L().Error("短信发送系统错误", zap.Error(err))
	}

	if err != nil {

		fmt.Println("zap 原生")
		zap.L().Error("发送验证码有误", zap.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "发送验证码服务端繁忙",
		}, err
	}

	//log.Println(fakerCode)
	return ginx.Result{
		Msg: "发送成功",
	}, nil
}

func (u *UserHandler) LoginSms(ctx *gin.Context, req LoginSmsReq) (ginx.Result, error) {

	log.Println("a.code : ", req.Code)
	if req.Phone == "" {
		ctx.String(http.StatusOK, "请输入手机号")
		return ginx.Result{
			Code: 4,
			Msg:  "请输入手机号",
		}, nil
	}

	ok, err := u.code.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "验证发送繁忙",
		}, err

	}
	if !ok {
		return ginx.Result{
			Code: 5,
			Msg:  "验证码错误，请重新输入",
		}, err
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	err = u.jwtHandler.SetLoginJWTToken(ctx, user.Id)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}

	return ginx.Result{
		Msg: "登录成功",
	}, nil

}
