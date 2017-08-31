package health

import (
	"github.com/hashicorp/vault/api"
	"github.com/scipian/eve"
)

type VaultChecker struct {
	Client *api.Client
	Config *api.Config
}

func NewVaultChecker() (*VaultChecker, *eve.Error) {
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, &eve.Error{
			Type:        "Vault",
			Description: "Create Client error",
			Error:       err.Error(),
		}
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, &eve.Error{
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

func (v *VaultChecker) InitStatus() *eve.Error {
	initStatus, err := v.Client.Sys().InitStatus()
	switch {
	case err != nil:
		return &eve.Error{
			Type:        "Vault",
			Description: "Init Status error",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: err.Error(),
		}
	case !initStatus:
		return &eve.Error{
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

func (v *VaultChecker) SealStatus() *eve.Error {
	sealStatus, err := v.Client.Sys().SealStatus()
	switch {
	case err != nil:
		return &eve.Error{
			Type:        "Vault",
			Description: "Seal Status error",
			Metadata: map[string]string{
				"Address": v.Config.Address,
			},
			Error: err.Error(),
		}
	case sealStatus.Sealed:
		return &eve.Error{
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
