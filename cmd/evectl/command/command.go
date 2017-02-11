package command

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "evectl",
	Short: "evectl can be used to control eve from the command line",
	Long:  `If you wish to control eve from the command line, then this program is for you.`,
}
