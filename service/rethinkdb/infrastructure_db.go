package rethinkdb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/concur/rohr"
	r "gopkg.in/gorethink/gorethink.v3"
	"time"
)

const (
	INFRA_TABLE = "infrastructure"
)

func (db *DbSession) InsertInfrastructure(infra *rohr.Infrastructure) error {
	res, err := r.DB(db.DbName).Table(INFRA_TABLE).Insert(
		map[string]interface{}{
			"Id":   r.UUID(infra.Name),
			"Name": infra.Name,
			"Quoin": map[string]interface{}{
				"Name":       infra.Quoin.Name,
				"ArchiveUri": infra.Quoin.ArchiveUri,
				"Variables":  infra.Quoin.Variables,
			},
			"Status":    infra.Status,
			"Variables": infra.Variables,
			"Authorization": map[string]interface{}{
				"Owner":       infra.Authorization.Owner,
				"GroupAccess": infra.Authorization.GroupAccess,
			},
			"Timestamp": r.EpochTime(time.Now().Unix()),
		},
	).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Printf("%d row inserted. \n", res.Inserted)
	return nil
}

func (db *DbSession) UpdateInfrastructureState(name string, state map[string]interface{}) error {
	res, err := r.DB(db.DbName).Table(INFRA_TABLE).Get(r.UUID(name)).Update(
		map[string]interface{}{
			"State": state,
		}).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Printf("%d row replaced. \n", res.Replaced)
	return nil
}

func (db *DbSession) UpdateInfrastructureStatus(name string, status rohr.Status) error {
	res, err := r.DB(db.DbName).Table(INFRA_TABLE).Get(r.UUID(name)).Update(map[string]interface{}{
		"Status": status,
	}).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Printf("%d row replaced. \n", res.Replaced)
	return nil
}

func (db *DbSession) GetInfrastructureByName(name string) (*rohr.Infrastructure, error) {
	var infrastructure rohr.Infrastructure
	cursor, err := r.DB(db.DbName).Table(INFRA_TABLE).Get(r.UUID(name)).Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return nil, err
	}
	if cursor.IsNil() {
		return nil, nil
	}
	if err = cursor.One(&infrastructure); err != nil {
		return nil, err
	}
	return &infrastructure, nil
}
