package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/token"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
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
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "500 Internal Server Error - " + err.Error(),
			})
		}
		return
	}

	accessToken, err := service.CreateAccessToken(params.Email)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "500 Internal Server Error - " + err.Error(),
		})
		return
	}

	refreshToken, err := service.CreateRefreshToken(params.Email)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "500 Internal Server Error - " + err.Error(),
		})
		return
	}

	utils.SetCookie(ctx, "access_token", accessToken, token.AccessTokenExpiration)
	utils.SetCookie(ctx, "refresh_token", refreshToken, token.RefreshTokenExpiration)

	uri, ok := ctx.GetQuery("redirect_uri")
	if !ok {
		uri = ""
	}
	ctx.Redirect(http.StatusMovedPermanently, uri)
}
