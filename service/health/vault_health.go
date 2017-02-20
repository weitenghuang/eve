package health

import (
	"github.com/concur/rohr"
	"github.com/hashicorp/vault/api"
)

type VaultChecker struct {
	Client *api.Client
	Config *api.Config
}

func NewVaultChecker() (*VaultChecker, *rohr.Error) {
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, &rohr.Error{
			Type:        "Vault",
			Description: "Create Client error",
			Error:       err.Error(),
		}
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, &rohr.Error{
			Type:        "Vault",
			Description: "Create Client error",
			Metadata: map[string]string{
				"Address": config.Address,
			},
			Error: err.Error(),
		}
	}
	return &VaultChecker{Client: client, Config: config}, nil
}

func (v *VaultChecker) InitStatus() *rohr.Error {
	initStatus, err := v.Client.Sys().InitStatus()
	switch {
	case err != nil:
		return &rohr.Error{
			Type:        "Vault",
			Description: "Init Status error",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: err.Error(),
		}
	case !initStatus:
		return &rohr.Error{
			Type:        "Vault",
			Description: "Init Status error",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: "Vault is not initialized",
		}
	}

	return nil
}

func (v *VaultChecker) SealStatus() *rohr.Error {
	sealStatus, err := v.Client.Sys().SealStatus()
	switch {
	case err != nil:
		return &rohr.Error{
			Type:        "Vault",
			Description: "Seal Status error",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: err.Error(),
		}
	case sealStatus.Sealed:
		return &rohr.Error{
			Type:        "Vault",
			Description: "Vault is sealed",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: "Vault is sealed",
		}
	}

	return nil
}
