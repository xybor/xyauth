package router

import (
	"github.com/gin-gonic/gin"
	apiv1 "github.com/xybor/xyauth/internal/handlers/api/v1"
	appv1 "github.com/xybor/xyauth/internal/handlers/app/v1"
	"github.com/xybor/xyauth/internal/handlers/well_known/openid"
	"github.com/xybor/xyauth/internal/middlewares"
	"github.com/xybor/xyauth/internal/utils"
)

func New() *gin.Engine {
	env := utils.GetConfig().GetDefault("general.environment", "dev").MustString()
	switch env {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test", "staging":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		utils.GetLogger().Event("invalid-environment").Field("env", env).Panic()
	}

	router := gin.Default()

	router.Static("/static", "web/static")
	router.StaticFile("/favicon.ico", "web/static/favicon.ico")
	router.LoadHTMLGlob("web/template/*.html")

	router.GET(".well-known/openid-configuration", openid.Handler)
	router.GET("login", appv1.LoginGETHandler)
	router.GET("register", appv1.RegisterGETHandler)
	router.GET("refresh", appv1.RefreshHandler)
	router.GET("logout", appv1.LogoutHandler)

	router.POST("login", appv1.LoginPOSTHandler)
	router.POST("register", appv1.RegisterPOSTHandler)

	mustAuthGroup := router.Group("")
	mustAuthGroup.Use(middlewares.VerifyAccessToken)
	{
		mustAuthGroup.GET("", appv1.WelcomeHandler)
	}

	apiv1Group := router.Group("api/v1")
	{
		apiv1Group.POST("register", apiv1.RegisterHandler)
		apiv1Group.POST("auth", apiv1.AuthHandler)
		apiv1Group.POST("revoke", apiv1.RevokeHandler)
	}

	return router
}
