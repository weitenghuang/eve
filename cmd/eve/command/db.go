package command

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "To control eve db",
	Long:  `db will help to initialize eve's db (rethinkdb)`,
}
