package openid

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xyauth/internal/config"
)

var Configuration = make(map[string]any)

func init() {
	domain := config.MustGet("server.domain").MustString()
	port := config.MustGet("server.tls_port").MustString()
	issuer := "https://" + domain
	if port != "443" {
		issuer += ":" + port
	}

	Configuration["issuer"] = issuer
	Configuration["authorization_endpoint"] = issuer + "/oauth2/v1/authorize"
	Configuration["token_endpoint"] = issuer + "/oauth2/v1/token"
}

func Handler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Configuration)
}
