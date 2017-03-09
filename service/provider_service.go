package service

import (
	"github.com/concur/rohr"
	"github.com/concur/rohr/service/rethinkdb"
)

type ProviderService struct {
	*rohr.User
}

func NewProviderService(user *rohr.User) *ProviderService {
	return &ProviderService{
		User: user,
	}
}

// GetProvider returns Provider information from database
func (p ProviderService) GetProvider(name string) (*rohr.Provider, error) {
	db := rethinkdb.DefaultSession()
	provider, err := db.GetProviderByName(name)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
