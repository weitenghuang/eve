package command

import (
	"github.com/concur/eve/cmd/eve/command/agent"
	"github.com/concur/eve/cmd/eve/command/db"
	"github.com/spf13/cobra"
	// flag "github.com/spf13/pflag"
)

func Execute() error {
	addCommands()
	return eveCmd.Execute()
}

var eveCmd = &cobra.Command{
	Use:   "eve",
	Short: "API server services REST operations",
	Long:  `API server services REST operations and perform user's request to create, delete, and operate user's infrastructures`,
}

func addCommands() {
	eveCmd.AddCommand(upCmd)
	eveCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agent.CreateCmd(apiServer))
	agentCmd.AddCommand(agent.DeleteCmd(apiServer))
	eveCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(db.InitCmd)
}
