package command

import (
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "To handle message from queue",
	Long:  `Listener will handle message from queue asynchronousely`,
}
