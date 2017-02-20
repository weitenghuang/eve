package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

// GetLogicalData returns data stored at the specified path
func GetLogicalData(path string) (map[string]interface{}, error) {
	// TODO: Bubble up errors in a more informative way
	config := api.DefaultConfig()
	err := config.ReadEnvironment()
	if err != nil {
		// Error reading environment
		return nil, err
	}
	client, err := api.NewClient(config)
	if err != nil {
		// Error creating vault client
		return nil, err
	}

	secret, err := client.Logical().Read(path)
	if err != nil {
		// Error retrieving secret
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("No value found at %s", path)
	}

	return secret.Data, nil
}

// WriteLogicalData writes data at a given path in Vault
func WriteLogicalData(path string, data map[string]interface{}) (map[string]interface{}, error) {
	config := api.DefaultConfig()
	err := config.ReadEnvironment()
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	return secret.Data, nil
}

// DeleteLogicalData remove data at a given path in Vault
func DeleteLogicalData(path string) error {
	config := api.DefaultConfig()
	err := config.ReadEnvironment()
	if err != nil {
		return err
	}

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	if _, err := client.Logical().Delete(path); err != nil {
		return err
	}

	return nil
}
