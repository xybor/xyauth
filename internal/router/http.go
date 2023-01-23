package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/xybor/xyauth/internal/config"
)

// NewHTTP returns a new router for HTTP. It redirects all requests to HTTPS
// server.
func NewHTTP() *gin.Engine {
	domain := config.GetDefault("server.domain", "localhost").MustString()
	tlsPort := config.GetDefault("server.tls_port", 8443).MustInt()
	tlsAddr := fmt.Sprintf("%s:%d", domain, tlsPort)

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			secureMiddleware := secure.New(secure.Options{
				SSLRedirect: true,
				SSLHost:     tlsAddr,
			})

			secureMiddleware.Process(c.Writer, c.Request)
		}
	}()

	router := gin.Default()
	router.Use(secureFunc)

	return router
}
