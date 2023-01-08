package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/pkg/token"
)

func WelcomeHandler(ctx *gin.Context) {
	name := "Stranger"
	if val, ok := ctx.Get("accessToken"); ok {
		name = strings.Split(val.(token.AccessToken).Email, "@")[0]
	}

	ctx.HTML(http.StatusOK, "welcome.html", gin.H{"name": name})
}
