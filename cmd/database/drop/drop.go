package drop

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/logger"
)

var sql bool
var nosql bool

var Command = &cobra.Command{
	Use:   "drop",
	Short: "Drop all tables",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if config.GetDefault("general.environment", "dev").MustString() != "dev" {
			fmt.Println("Only allow dropping the table if the environment is dev")
			return
		}

		if sql {
			if err := database.DropAllSQL(); err != nil {
				logger.Event("drop-sql-error").Field("error", err).Fatal()
			}
			logger.Event("drop-sql-success").Info()
		}

		if nosql {
			if err := database.DropAllNoSQL(); err != nil {
				logger.Event("drop-nosql-error").Field("error", err).Fatal()
			}
			logger.Event("drop-nosql-success").Info()
		}
	},
}

func init() {
	Command.Flags().BoolVar(&sql, "sql", false, "Drop all SQL tables")
	Command.Flags().BoolVar(&nosql, "nosql", false, "Drop all NoSQL collections")
}
