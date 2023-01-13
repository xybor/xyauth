package drop

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xybor/xyauth/internal/database"
	"github.com/xybor/xyauth/internal/utils"
)

var sql bool
var nosql bool

var Command = &cobra.Command{
	Use:   "drop",
	Short: "Drop all tables",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.GetConfig().GetDefault("general.environment", "dev").MustString() != "dev" {
			fmt.Println("Only allow dropping the table if the environment is dev")
			return
		}

		if sql {
			if err := database.DropAllSQL(); err != nil {
				utils.GetLogger().Event("drop-sql-error").Field("error", err).Fatal()
			}
			fmt.Println("Drop all SQL tables successfully")
		}

		if nosql {
			if err := database.DropAllNoSQL(); err != nil {
				utils.GetLogger().Event("drop-nosql-error").Field("error", err).Fatal()
			}
			fmt.Println("Drop all NoSQL collections successfully")
		}
	},
}

func init() {
	Command.Flags().BoolVar(&sql, "sql", false, "Drop all SQL tables")
	Command.Flags().BoolVar(&nosql, "nosql", false, "Drop all NoSQL collections")
}
