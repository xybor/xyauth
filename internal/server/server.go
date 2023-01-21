package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/router"
)

type AuthServer struct {
	host string
	port int

	privateKey []byte
	publicKey  []byte

	router *gin.Engine
}

func NewServer() (http.Server, net.Listener) {
	host := config.GetDefault("server.host", "0.0.0.0").MustString()
	port := config.GetDefault("server.port", 8443).MustInt()
	if _, ok := config.Get("DOCKER_RUNNING"); ok {
		host = "0.0.0.0"
		port = 8443
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	// TODO: ReplaceAll commands will be removed if the PR 156 of godotenv is
	// merged.
	key := config.MustGet("SERVER_PUBLIC_KEY").MustString()
	publicKey := []byte(strings.ReplaceAll(key, `\n`, "\n"))

	key = config.MustGet("SERVER_PRIVATE_KEY").MustString()
	privateKey := []byte(strings.ReplaceAll(key, `\n`, "\n"))

	cert, err := tls.X509KeyPair(publicKey, privateKey)
	xycond.AssertNil(err)

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	httpServer := http.Server{Addr: addr, Handler: router.New(), TLSConfig: tlsConfig}
	tlsListener, err := tls.Listen("tcp", addr, tlsConfig)
	xycond.AssertNil(err)

	return httpServer, tlsListener
}
