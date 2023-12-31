package route

import (
	"time"

	"gin-app/api/v1/middleware"
	_ "gin-app/docs"
	"gin-app/internal/bootstrap"

	"github.com/gin-gonic/gin"
	swagfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func Setup(app *bootstrap.Application, timeout time.Duration, en *gin.Engine) {
	en.Use(middleware.LimitRequestRate(app.Limiter))
	publicRouter := en.Group("/api/v1")
	publicRouter.GET("/docs/*any", ginswagger.WrapHandler(swagfiles.Handler))
	publicRouter.Use(middleware.HandlerHeadersCtx(), middleware.HandlerError(app.Log))
	NewSignupCtl(app, timeout, publicRouter)
	if app.Conf.JWTEnabled {
		publicRouter.Use(middleware.HandlerAuth(app.Rdb))
	}
	NewUserHimSelfCtrl(app, timeout, publicRouter)
	NewAdminCtrl(app, timeout, publicRouter)
}
