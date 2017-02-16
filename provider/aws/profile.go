package aws

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var directions = `
	Your new profile [%s] is stored in the following files:

		Config: %s
		Credentials: %s

	Note that the credentials will expire in 1 hour.
	After this time, you may safely re-run the "authenticate" command to refresh your access credentials.
	To use the credentials, call the AWS CLI with the --profile option.

		$ aws --profile %s ec2 describe-instances
`

// Credentials represents static AWS credentials
type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

// Profile represents an AWS profile
type Profile struct {
	Name        string
	Credentials *Credentials
	Region      string
}

// Save writes the profile to disk
func (p *Profile) Save() error {
	log.Info(fmt.Sprintf("Writing profile '%s' to config...", p.Name))
	configFilePath, err := p.saveConfig()
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Writing profile '%s' to credentials...", p.Name))
	credentialsFilePath, err := p.saveCredentials()
	if err != nil {
		return err
	}

	log.Infof(directions, p.Name, configFilePath, credentialsFilePath, p.Name)

	return nil
}

func (p *Profile) saveConfig() (string, error) {
	configFilePath, err := p.getConfigFilePath()
	if err != nil {
		return "", err
	}

	config, err := ini.Load(configFilePath)
	if err != nil {
		return "", err
	}

	var profileName string
	if p.Name != "default" {
		profileName = fmt.Sprintf("profile %s", p.Name)
	} else {
		profileName = p.Name
	}

	s := config.Section(profileName)
	s.Key("region").SetValue(p.Region)
	config.SaveTo(configFilePath)

	return configFilePath, nil
}

func (p *Profile) saveCredentials() (string, error) {
	credentialsFilePath, err := p.getCredentialsFilePath()
	if err != nil {
		return "", err
	}

	credentials, err := ini.Load(credentialsFilePath)
	if err != nil {
		return "", err
	}

	s := credentials.Section(p.Name)
	s.Key("aws_access_key_id").SetValue(p.Credentials.AccessKeyId)
	s.Key("aws_secret_access_key").SetValue(p.Credentials.SecretAccessKey)
	s.Key("aws_session_token").SetValue(p.Credentials.SessionToken)
	credentials.SaveTo(credentialsFilePath)

	return credentialsFilePath, nil
}

func (p *Profile) getConfigFilePath() (string, error) {
	homeDir, err := p.getHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".aws", "config"), nil
}

func (p *Profile) getCredentialsFilePath() (string, error) {

	homeDir, err := p.getHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".aws", "credentials"), nil
}

func (p *Profile) getHomeDir() (string, error) {
	// Try the HOME environment variable for *nix
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		// Try the USERPROFILE envrionment variable for windows
		homeDir = os.Getenv("USERPROFILE")
	}

	// Still empty? Return with error
	if homeDir == "" {
		return "", errors.New("UserHomeNotFound: user home directory not found")
	}

	return homeDir, nil
}
