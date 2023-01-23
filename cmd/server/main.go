package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
	"github.com/xybor/xyauth/internal/server"
	"github.com/xybor/xyauth/pkg/service"
)

var configs []string

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  "Start the authentication and authorization server",
	Run: func(cmd *cobra.Command, args []string) {
		for i := range configs {
			xycond.AssertNil(config.Add(configs[i]))
		}

		xycond.AssertNil(database.InitPostgresDB(nil))
		xycond.AssertNil(database.InitMongoDB())

		go service.JanitorRefreshToken()

		adminEmail := config.GetDefault("XYBOR_INIT_ADMIN_EMAIL", "").MustString()
		adminPassword := config.GetDefault("XYBOR_INIT_ADMIN_PASSWORD", "").MustString()

		if adminEmail != "" && adminPassword != "" {
			err := service.Register(adminEmail, adminPassword, "admin")
			if err != nil {
				logger.Event("register-init-admin-failed").Field("error", err).Warning()
			}
		}
		host := config.GetDefault("server.host", "0.0.0.0").MustString()
		port := config.GetDefault("server.port", 8080).MustInt()
		tlsPort := config.GetDefault("server.tls_port", 8443).MustInt()

		addr := fmt.Sprintf("%s:%d", host, port)
		tlsAddr := fmt.Sprintf("%s:%d", host, tlsPort)

		go func() {
			logger.Event("server-start").Field("address", addr).Info()
			logger.Event("server-close").Field("error", server.NewHTTP()()).Info()
		}()

		logger.Event("server-start").Field("address", tlsAddr).Info()
		logger.Event("server-close").Field("error", server.NewHTTPS()()).Info()
	},
}

func main() {
	rootCmd.Flags().StringArrayVarP(&configs, "config", "c", nil, "Specify overridden config files")
	if err := rootCmd.Execute(); err != nil {
		logger.Event("server-command-error").Field("error", err).Fatal()
	}
}
