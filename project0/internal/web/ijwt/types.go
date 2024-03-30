package ijwt

import "github.com/gin-gonic/gin"

type Handler interface {
	ExtraToken(ctx *gin.Context) string

	CheckSession(ctx *gin.Context, ssid string) error
	SetLoginJWTToken(ctx *gin.Context, uid int64) error
	SetJWT(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
}
