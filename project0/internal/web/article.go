package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"project0/internal/domain"
	"project0/internal/service"
	"project0/internal/web/ijwt"
	"project0/pkg/ginx"
	"project0/pkg/loggerDefine"
	"time"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   loggerDefine.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l loggerDefine.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {

	ag := server.Group("/articles")
	ag.POST("/edit", ginx.WrapBodyAndClaims(a.Edit))
	ag.POST("/publish",ginx.WrapBody(a.Publish))
	// 自己从头到尾写一个维护帖子状态。撤回帖子
	ag.POST("/withdraw",a.Withdraw)
	// 读者接口
	ag.GET("/detail/:id",a.Detail)
	// 创作者接口
	ag.POST("/list",a.List)
}

// 新建和修改共用一个接口
func (a *ArticleHandler) Edit(ctx *gin.Context, req ArticleEdit,uc ijwt.UserClaims) (ginx.Result, error) {
	//uc := ctx.MustGet("user").(ijwt.UserClaims)
	id, err := a.svc.Save(ctx, domain.Article{
		//有没有id判断是新建还是修改？
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
		//Status: 1,
	})

	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
		}, err
	}

	return ginx.Result{
		//Msg: "帖子更新成功", 这里返回的帖子id
		Data: id,
	}, nil

}
//发布接口好像和新建更新修改没什么区别？
func (a *ArticleHandler) Publish(ctx *gin.Context,req ArticleEdit) (ginx.Result, error) {

	type ArticlePublishReq struct {
		// 没加json完蛋?
		Id int64      `json:"id"`
		Title string `json:"title"`
		Content string `json:"content"`
	}
	// ------------
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	id, err := a.svc.Publish(ctx, domain.Article{
		//有没有id判断是新建还是修改？
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
		//Status: 1,
	})

	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err

	}
	return ginx.Result{
		//Msg: "帖子更新成功", 这里返回的帖子id
		Data: id,
	}, nil

}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if  err :=ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	// 如何改掉惯性思维
	err := a.svc.Withdraw(ctx,uc.Uid,req.Id)

	if err != nil {
		ctx.JSON(http.StatusOK,Result{
			Code: 5,
			Msg: "回撤帖子，系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK,Result{
		Msg: "ok",
	})

}

func (a *ArticleHandler) List(ctx *gin.Context) {
     var  page Page
	 if err := ctx.Bind(&page);  err != nil {
		 return
	 }
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	arts, err := a.svc.GetByAuthor(ctx, uc.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK,Result{
			Code: 5,
			Msg:"系统错误",
		})
		a.l.Error("查找文章列表失败",
			loggerDefine.Int64("uid",uc.Uid),
			loggerDefine.Int("offset",page.Offset),
			loggerDefine.Int("limit",page.Limit))
		return
	}
	ctx.JSON(http.StatusOK,Result{
		Code: 2,
		Data: slice.Map[domain.Article,ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
			return ArticleVo{
				Id: src.Id,
				Title: src.Title,
				Content: src.Content,
				AuthorId: src.Author.Id,
				Status: src.Status.ToUint8(),
				Ctime: src.Ctime.Format(time.DateTime),
				Utime: src.Utime.Format(time.DateTime),
			}
		}),
	})



}

func (a *ArticleHandler) Detail(context *gin.Context) {
	
}
