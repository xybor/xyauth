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

type RevokeParams struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RevokeHandler(ctx *gin.Context) {
	params := new(RevokeParams)
	if err := ctx.ShouldBind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameters"})
		return
	}

	refreshToken := token.RefreshToken{}
	if err := token.Verify(params.RefreshToken, &refreshToken); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"message": xyerror.Message(err)})
		return
	}

	if utils.Check(ctx).Perform("delete").On(refreshToken) != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}

	if err := service.RevokeRefreshToken(refreshToken); err != nil {
		if errors.Is(err, service.NotFoundError) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "not found the refresh token"})
		} else {
			logger.Event("revoke-refresh-token-failed", ctx).Field("error", err).Warning()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unknown error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "refresh token is revoked successfully"})
}
