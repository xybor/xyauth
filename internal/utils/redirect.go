package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTMLOrLogin serves the HTML pages if the access token is presented,
// otherwise, it redirects to refresh page to refresh the access token.
func HTMLOrLogin(ctx *gin.Context, code int, name string, obj any) {
	if _, ok := GetAccessToken(ctx); ok {
		ctx.HTML(code, name, obj)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "refresh?redirect_uri="+ctx.Request.RequestURI)
	}
}
