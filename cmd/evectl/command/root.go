package command

import (
	"io"

	"github.com/spf13/cobra"
)

// NewRootCommand creates and instance of the RootCommand
func NewRootCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	commands := &cobra.Command{
		Use:   "evectl",
		Short: "evectl can be used to control eve from the command line",
		Long:  `If you wish to control eve from the command line, then this program is for you.`,
	}

	commands.AddCommand(NewAuthenticateCommand(out, err))

	return commands
}
