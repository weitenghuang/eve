package rethinkdb

import (
	"github.com/scipian/eve"
	r "gopkg.in/gorethink/gorethink.v3"
)

const (
	PROVIDER_TABLE = "provider"
)

func (db *DbSession) GetProviderByName(name string) (*eve.Provider, error) {
	var provider eve.Provider
	cursor, err := r.DB(db.DbName).Table(PROVIDER_TABLE).Get(r.UUID(name)).Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return nil, err
	}
	if cursor.IsNil() {
		return nil, nil
	}
	if err = cursor.One(&provider); err != nil {
		return nil, err
	}
	return &provider, nil
}
