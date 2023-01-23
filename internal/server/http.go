package server

import (
	"fmt"

	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/router"
)

func NewHTTP() func() error {
	host := config.GetDefault("server.host", "0.0.0.0").MustString()
	port := config.GetDefault("server.port", 8080).MustInt()
	if _, ok := config.Get("DOCKER_RUNNING"); ok {
		host = "0.0.0.0"
		port = 8080
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	return func() error {
		return router.NewHTTP().Run(addr)
	}
}
