package main

import (
	"github.com/spf13/cobra"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/cmd/database/drop"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:   "database",
	Short: "Migrate, delete, modify database",
	Long:  `A helper tool for managing database`,
}

func main() {
	rootCmd.AddCommand(drop.Command)

	xycond.AssertNil(database.InitPostgresDB(nil))
	xycond.AssertNil(database.InitMongoDB())

	if err := rootCmd.Execute(); err != nil {
		utils.GetLogger().Event("database-command-error").Field("error", err).Fatal()
	}
}
