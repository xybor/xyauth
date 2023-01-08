package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/pkg/service"
)

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthHandler(ctx *gin.Context) {
	params := new(AuthParams)
	ctx.ShouldBind(params)

	err := service.Authenticate(params.Email, params.Password)

	if err != nil {
		switch {
		case errors.Is(err, service.NotFoundError):
			ctx.JSON(http.StatusNotFound, gin.H{"message": xyerror.Message(err)})
		case errors.Is(err, service.CredentialError):
			ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		case err != nil:
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	accessToken, err := service.CreateAccessToken(params.Email)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	refreshToken, err := service.CreateRefreshToken(params.Email)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
