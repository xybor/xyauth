package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

type LoginParams struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginGETHandler(ctx *gin.Context) {
	// Redirect to the main page if the user already authenticated.
	if utils.IsAuthenticated(ctx) {
		ctx.Redirect(http.StatusTemporaryRedirect, "")
	} else {
		ctx.HTML(http.StatusOK, "login.html", nil)
	}
}

func LoginPOSTHandler(ctx *gin.Context) {
	params := new(LoginParams)
	ctx.ShouldBind(params)

	err := service.Authenticate(params.Email, params.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.NotFoundError):
			ctx.HTML(http.StatusNotFound, "login.html", gin.H{"message": xyerror.Message(err)})
		case errors.Is(err, service.CredentialError):
			ctx.HTML(http.StatusForbidden, "login.html", gin.H{"message": xyerror.Message(err)})
		case err != nil:
			logger.Event("query-user-failed").
				Field("email", params.Email).Field("error", err).Warning()
			ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"message": "Something is wrong, please contact to administrator if the issue persists",
			})
		}
		return
	}

	accessToken, err := service.CreateAccessToken(params.Email)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"message": "Something is wrong, please contact to administrator if the issue persists",
		})
		return
	}

	refreshToken, err := service.CreateRefreshToken(params.Email)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"message": "Something is wrong, please contact to administrator if the issue persists",
		})
		return
	}

	ctx.SetCookie(
		"access_token", accessToken,
		int(token.AccessTokenExpiration/time.Second), "/",
		config.GetDefault("server.domain", "localhost").MustString(),
		true, true,
	)

	ctx.SetCookie(
		"refresh_token", refreshToken,
		int(token.RefreshTokenExpiration/time.Second), "/",
		config.GetDefault("server.domain", "localhost").MustString(),
		true, true,
	)

	uri, ok := ctx.GetQuery("redirect_uri")
	if !ok {
		uri = ""
	}
	ctx.Redirect(http.StatusMovedPermanently, uri)
}
