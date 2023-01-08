package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/router"
	"github.com/xybor/xyauth/internal/utils"
	"github.com/xybor/xyauth/pkg/service"
)

var configs []string

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  "Start the authentication and authorization server",
	Run: func(cmd *cobra.Command, args []string) {
		for i := range configs {
			xycond.AssertNil(utils.AddConfig(configs[i]))
		}

		xycond.AssertNil(database.InitPostgresDB(nil))
		xycond.AssertNil(database.InitMongoDB())

		host := utils.GetConfig().GetDefault("server.host", "0.0.0.0").MustString()
		port := utils.GetConfig().GetDefault("server.port", 8443).MustInt()
		if _, ok := utils.GetConfig().Get("DOCKER_RUNNING"); ok {
			host = "0.0.0.0"
			port = 8443
		}
		addr := fmt.Sprintf("%s:%d", host, port)

		go service.JanitorRefreshToken()

		// TODO: ReplaceAll commands will be removed if the PR 156 of godotenv
		// is merged.
		key := utils.GetConfig().MustGet("SERVER_PUBLIC_KEY").MustString()
		public_key := strings.ReplaceAll(key, `\n`, "\n")

		key = utils.GetConfig().MustGet("SERVER_PRIVATE_KEY").MustString()
		private_key := strings.ReplaceAll(key, `\n`, "\n")

		cert, err := tls.X509KeyPair([]byte(public_key), []byte(private_key))
		xycond.AssertNil(err)
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

		server := http.Server{Addr: addr, Handler: router.New(), TLSConfig: tlsConfig}
		tlsListener, err := tls.Listen("tcp", addr, tlsConfig)
		xycond.AssertNil(err)

		utils.GetLogger().Event("server-start").Field("address", addr).Info()
		utils.GetLogger().Event("server-close").Field("error", server.Serve(tlsListener)).Info()
	},
}

func main() {
	rootCmd.Flags().StringArrayVarP(&configs, "config", "c", nil, "Specify overridden config files")
	if err := rootCmd.Execute(); err != nil {
		utils.GetLogger().Event("server-command-error").Field("error", err).Fatal()
	}
}
