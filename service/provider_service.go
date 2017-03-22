package service

import (
	"github.com/concur/eve"
	"github.com/concur/eve/service/rethinkdb"
)

type ProviderService struct {
	*eve.User
}

func NewProviderService(user *eve.User) *ProviderService {
	return &ProviderService{
		User: user,
	}
}

// GetProvider returns Provider information from database
func (p ProviderService) GetProvider(name string) (*eve.Provider, error) {
	db := rethinkdb.DefaultSession()
	provider, err := db.GetProviderByName(name)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
