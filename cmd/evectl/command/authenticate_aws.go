package command

import (
	"errors"
	"io"
	"os"

	"github.com/scipian/eve/client"
	"github.com/scipian/eve/provider/aws"
	"github.com/spf13/cobra"
	"github.com/tj/go-prompt"
)

// Account for AWS.
var account string

// IAM Role for AWS.
var iamRole string

// Region for AWS.
var region string

// Name given to profile
var profileName string

// NewAuthenticateAwsCommand creates an instance of the AuthenticateAwsCommand
func NewAuthenticateAwsCommand(out, err io.Writer) *cobra.Command {
	command := &cobra.Command{
		Use:   "aws [--account ACCOUNT] [--iam-role ROLE] [--region REGION]",
		Short: "Authencticate with AWS",
		Long:  `Used for authenticating against AWS`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authenticateAws(cmd, args, out, err)
		},
	}

	command.Flags().StringVarP(&account, "account", "a", "", "AWS account to authenticate against.")
	command.Flags().StringVarP(&iamRole, "iam-role", "i", "", "AWS IAM role to assume.")
	command.Flags().StringVarP(&region, "region", "r", "", "AWS region to use.")
	command.Flags().StringVarP(&profileName, "profile-name", "p", "", "Name to give profile.")

	return command
}

func authenticateAws(cmd *cobra.Command, args []string, out, err io.Writer) error {

	// Retrieve flags
	a := getAccount()
	i := getIamRole(a)
	r := getRegion(a)
	p := getProfileName()

	// Get logged in user
	sessionName, e := getLoggedInUser()
	if e != nil {
		return e
	}

	// Authenticate
	auth := aws.NewAuthenticator(a, i, sessionName)
	c, e := auth.Authenticate()
	if e != nil {
		return e
	}

	// Create the profile
	profile := &aws.Profile{
		Name:        p,
		Credentials: c,
		Region:      r,
	}

	// Write the profile to disk
	e = profile.Save()
	if e != nil {
		return e
	}

	return nil
}

func getLoggedInUser() (string, error) {
	// Try the USER environment variable for *nix
	user := os.Getenv("USER")
	if user == "" {
		// Try the USERNAME envrionment variable for windows
		user = os.Getenv("USERNAME")
	}

	// Still empty? Return with error
	if user == "" {
		return "", errors.New("UserNotFound: logged in user not found")
	}

	return user, nil
}

func printStatement(s string) {
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(s)
	os.Stdout.WriteString("\n")
}

func getFlagValue(flag string, environmentKey string, ask func() string) string {
	if flag == "" {
		flag = os.Getenv(environmentKey)
		if flag == "" {
			flag = ask()
		}
	}
	return flag
}

func getProfileName() string {
	// region from flag, env, prompt
	return getFlagValue(profileName, "AWS_PROFILE", func() string {
		printStatement("Please provide a name for this profile:")
		profileName = prompt.StringRequired("  Profile Name: ")
		return profileName
	})
}

func getRegion(account *aws.Account) string {
	// region from flag, env, prompt
	return getFlagValue(region, "AWS_REGION", func() string {
		printStatement("Please choose an region:")
		regions := account.Regions
		i := prompt.Choose("  Region: ", regions)
		return regions[i]
	})
}

func getIamRole(account *aws.Account) string {
	// iamRole from flag, env, prompt
	return getFlagValue(iamRole, "AWS_IAM_ROLE", func() string {
		printStatement("Please choose an IAM role:")
		roles := account.Roles
		i := prompt.Choose("  Role: ", roles)
		return roles[i]
	})
}

func getAccount() *aws.Account {

	c := client.NewDefaultClient()
	p := c.GetProvider("aws")
	a := aws.NewProvider(p)

	// account from flag, env, prompt
	account = getFlagValue(account, "AWS_ACCOUNT", func() string {
		accountNames := a.GetAccountNames()
		printStatement("Please choose an account:")
		i := prompt.Choose("  Account: ", accountNames)
		return accountNames[i]
	})

	return a.GetAccount(account)
}
