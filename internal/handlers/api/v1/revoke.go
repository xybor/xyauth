package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xyerror"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
	"github.com/xybor/xyauth/pkg/token"
)

type RevokeParams struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RevokeHandler(ctx *gin.Context) {
	params := RevokeParams{}

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "expected an access_token"})
		return
	}

	accessToken := token.AccessToken{}
	if err := token.Verify(params.AccessToken, &accessToken); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": xyerror.Message(err)})
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(params.RefreshToken, &refreshToken); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"message": xyerror.Message(err)})
		return
	}

	if accessToken.Email != refreshToken.Email {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "you do not have the permission to revoke this token"})
		return
	}

	if err := service.RevokeRefreshToken(params.RefreshToken); err != nil {
		if errors.Is(err, service.NotFoundError) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "not found the refresh token"})
		} else {
			utils.GetLogger().Event("revoke-refresh-token-failed").Field("error", err).Warning()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unknown error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "refresh token is revoked successfully"})
}
