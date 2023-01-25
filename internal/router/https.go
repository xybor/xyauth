package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/config"
	apiv1 "github.com/xybor/xyauth/internal/handlers/api/v1"
	appv1 "github.com/xybor/xyauth/internal/handlers/app/v1"
	"github.com/xybor/xyauth/internal/handlers/well_known/openid"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/middlewares"
)

// NewHTTPS returns a new router for HTTPS.
func NewHTTPS() *gin.Engine {
	env := config.GetDefault("general.environment", "dev").MustString()
	switch env {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test", "staging":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		logger.Event("invalid-environment").Field("env", env).Panic()
	}

	router := gin.Default()

	router.Static("/static", "web/static")
	router.StaticFile("/favicon.ico", "web/static/favicon.ico")
	router.LoadHTMLGlob("web/template/*.html")

	router.Use(middlewares.VerifyAccessToken)

	router.GET(".well-known/openid-configuration", openid.Handler)

	router.GET("login", appv1.LoginGETHandler)
	router.POST("login", appv1.LoginPOSTHandler)

	router.GET("register", appv1.RegisterGETHandler)
	router.POST("register", appv1.RegisterPOSTHandler)

	router.GET("profile", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "profile.html", nil)
	})

	router.GET("", appv1.WelcomeHandler)
	router.GET("logout", appv1.LogoutHandler)

	router.Any("refresh", appv1.RefreshHandler)

	apiv1Group := router.Group("api/v1")
	{
		apiv1Group.POST("register", apiv1.RegisterHandler)
		apiv1Group.POST("auth", apiv1.AuthHandler)
		apiv1Group.POST("revoke", apiv1.RevokeHandler)
		apiv1Group.POST("refresh", apiv1.RefreshHandler)
	}

	return router
}
