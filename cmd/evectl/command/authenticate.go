package command

import (
	"io"

	"github.com/spf13/cobra"
)

// NewAuthenticateCommand creates and instance of the AuthenticateCommand
func NewAuthenticateCommand(out, err io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:   "authenticate <provider>",
		Short: "Authenticate against approved providers",
		Long:  `Used for authenticating against approved providers`,
	}

	command.AddCommand(NewAuthenticateAwsCommand(out, err))

	return command
}
