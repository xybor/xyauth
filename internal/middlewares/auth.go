package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/token"
)

func redirectToRefresh(ctx *gin.Context) {
	uri, ok := ctx.GetQuery("redirect_uri")
	if !ok {
		uri = ctx.Request.RequestURI
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/refresh?redirect_uri="+uri)
	ctx.Abort()
}

func VerifyAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		redirectToRefresh(ctx)
		return
	}

	accessToken := token.AccessToken{}
	if err := token.Verify(cookie, &accessToken); err != nil {
		utils.GetLogger().Event("access-token-invalid").
			Field("cookie", cookie).Field("error", err).Debug()
		redirectToRefresh(ctx)
		return
	}

	ctx.Set("accessToken", accessToken)
	ctx.Next()
}
