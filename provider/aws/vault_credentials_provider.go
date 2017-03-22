package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/concur/eve/pkg/vault"
)

// VaultCredentialsProvider retrieves credentials from vault at the specified key
// Implements credentials.Provider
type VaultCredentialsProvider struct {
	Key       string
	retrieved bool
}

// VaultCredentialsProviderName provides a name of Vault provider
const VaultCredentialsProviderName = "VaultCredentialsProvider"

// Retrieve reads and extracts the credentials from vault
func (v *VaultCredentialsProvider) Retrieve() (credentials.Value, error) {
	v.retrieved = false

	data, err := vault.GetLogicalData(v.Key)
	if err != nil {
		return credentials.Value{ProviderName: VaultCredentialsProviderName}, err
	}

	// TODO(dwr): Error if any of the below keys are not found
	creds := credentials.Value{
		AccessKeyID:     data["accessKeyID"].(string),
		SecretAccessKey: data["secretAccessKey"].(string),
		SessionToken:    data["sessionToken"].(string),
		ProviderName:    VaultCredentialsProviderName,
	}

	v.retrieved = true
	return creds, nil
}

// IsExpired returns if the Vault credentials have expired
func (v *VaultCredentialsProvider) IsExpired() bool {
	return !v.retrieved
}
