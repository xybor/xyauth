package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/certificate"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/router"
)

func NewHTTPS() func() error {
	host := config.GetDefault("server.host", "0.0.0.0").MustString()
	port := config.GetDefault("server.tls_port", 8443).MustInt()
	if _, ok := config.Get("DOCKER_RUNNING"); ok {
		host = "0.0.0.0"
		port = 8443
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	// TODO: ReplaceAll commands will be removed if the PR 156 of godotenv is
	// merged.
	key := config.MustGet("SERVER_PUBLIC_KEY").MustString()
	publicKey, err := certificate.GetCertificateContent(key)
	xycond.AssertNil(err)

	key = config.MustGet("SERVER_PRIVATE_KEY").MustString()
	privateKey, err := certificate.GetCertificateContent(key)
	xycond.AssertNil(err)

	cert, err := tls.X509KeyPair(publicKey, privateKey)
	xycond.AssertNil(err)

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := http.Server{Addr: addr, Handler: router.NewHTTPS(), TLSConfig: tlsConfig}
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	xycond.AssertNil(err)

	return func() error {
		return server.Serve(listener)
	}
}
