package main

import (
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

		server, listener := server.NewServer()

		logger.Event("server-start").Field("address", server.Addr).Info()
		logger.Event("server-close").Field("error", server.Serve(listener)).Info()
	},
}

func main() {
	rootCmd.Flags().StringArrayVarP(&configs, "config", "c", nil, "Specify overridden config files")
	if err := rootCmd.Execute(); err != nil {
		logger.Event("server-command-error").Field("error", err).Fatal()
	}
}
