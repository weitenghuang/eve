package aws

import (
	"reflect"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	"github.com/mitchellh/mapstructure"
)

// accountList Implements the sort interface to sort an array of
// accounts by name alphabetically
type accountList []*Account

func (a accountList) Len() int           { return len(a) }
func (a accountList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a accountList) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Provider represents an AWS Provider
type Provider struct {
	provider *eve.Provider
	accounts accountList
}

// NewProvider returns an instance of an AWS provider
func NewProvider(provider *eve.Provider) *Provider {
	p := &Provider{
		provider: provider,
	}
	p.loadAccounts()
	return p
}

// GetAccountNames returns an array of account names
func (p *Provider) GetAccountNames() []string {
	sort.Sort(p.accounts)
	names := make([]string, 0, p.accounts.Len())
	for _, k := range p.accounts {
		names = append(names, k.Name)
	}
	return names
}

// GetAccount returns an Account with given name
func (p *Provider) GetAccount(name string) *Account {
	var account *Account
	for _, v := range p.accounts {
		if v.Name != name {
			continue
		}
		account = v
	}

	return account
}

// loadAccounts parses the provider schema to hydrate an array of accounts
func (p *Provider) loadAccounts() {
	var accounts []*Account
	if p.provider.Schema.Type == reflect.TypeOf(accounts).String() {
		err := mapstructure.Decode(p.provider.Schema.Data, &accounts)
		if err != nil {
			log.Fatalf("loadAccounts: %s", err)
		}
	} else {
		accounts = make([]*Account, 0, 0)
	}

	p.accounts = accounts
}
