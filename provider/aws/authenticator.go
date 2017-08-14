package aws

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Authenticator represents the necessary data
type Authenticator struct {
	Account         *Account
	IamRole         string
	RoleSessionName string
}

// NewAuthenticator creates an instance of the Authenticator
func NewAuthenticator(account *Account, iamRole, roleSessionName string) *Authenticator {
	return &Authenticator{
		Account:         account,
		IamRole:         iamRole,
		RoleSessionName: roleSessionName,
	}
}

// Authenticate authenticates against AWS using the account instance
func (a *Authenticator) Authenticate() (*Credentials, error) {
	switch a.Account.AuthType {
	case IAM:
		return a.iamDirect()
	case AssumeRole:
		return a.assumeRole()
	default:
		return nil, errors.New("Authenticate: Unsupported AuthType")
	}
}

// iamDirect uses AWS IAM key and secret to operate in AWS directly without assuming role
func (a *Authenticator) iamDirect() (*Credentials, error) {

	// set credentials provider
	v := &VaultCredentialsProvider{
		Key: "secret/quoin/providers/aws/credentials",
	}

	// retrieve creds from vault
	creds, err := v.Retrieve()

	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyId:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
	}, nil
}

// assumeRole assumes a specified role in an AWS account based on credentials
func (a *Authenticator) assumeRole() (*Credentials, error) {
	// Setup the AWS session
	s := a.getSession()

	// Create the Role ARN
	arn := fmt.Sprintf("arn:aws:iam::%d:role/%s", a.Account.Id, a.IamRole)

	// Use the AssumeRoleProvider
	p := &stscreds.AssumeRoleProvider{
		Client:          sts.New(s),
		RoleARN:         arn,
		RoleSessionName: a.RoleSessionName,
		Duration:        time.Duration(60) * time.Minute,
	}

	// Get the credentials
	roleCreds, err := credentials.NewCredentials(p).Get()
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyId:     roleCreds.AccessKeyID,
		SecretAccessKey: roleCreds.SecretAccessKey,
		SessionToken:    roleCreds.SessionToken,
	}, nil
}

func (a *Authenticator) getSession() *session.Session {
	// Setup a chain of credential providers
	providers := []credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
		&VaultCredentialsProvider{
			Key: "secret/quoin/providers/aws/credentials",
		},
	}

	chain := credentials.NewChainCredentials(providers)
	c := &aws.Config{Credentials: chain}
	return session.New(c)
}
