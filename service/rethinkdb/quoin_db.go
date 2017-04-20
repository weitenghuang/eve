package rethinkdb

import (
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve"
	r "gopkg.in/gorethink/gorethink.v3"
	"strings"
	"time"
)

const (
	QUOIN_TABLE         = "quoin"
	QUOIN_ARCHIVE_TABLE = "quoinArchive"
)

func (db *DbSession) InsertQuoin(quoin *eve.Quoin) error {
	res, err := r.DB(db.DbName).Table(QUOIN_TABLE).Insert(
		map[string]interface{}{
			"Id":         r.UUID(quoin.Name),
			"Name":       quoin.Name,
			"ArchiveUri": quoin.ArchiveUri,
			"Variables":  quoin.Variables,
			"Authorization": map[string]interface{}{
				"Owner":       quoin.Authorization.Owner,
				"GroupAccess": quoin.Authorization.GroupAccess,
			},
			"Timestamp": r.EpochTime(time.Now().Unix()),
		}).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Printf("%d row inserted. \n", res.Inserted)
	quoinData, err := db.GetQuoinByName(quoin.Name)
	if err != nil {
		return err
	}
	*quoin = *quoinData
	return nil
}

func (db *DbSession) UpdateQuoin(quoinName string, value interface{}) error {
	res, err := r.DB(db.DbName).Table(QUOIN_TABLE).Get(r.UUID(quoinName)).Update(value).RunWrite(db.Session)
	if err != nil {
		return err
	}
	log.Printf("%d row replaced. \n", res.Replaced)
	return nil
}

func (db *DbSession) GetQuoinByName(name string) (*eve.Quoin, error) {
	var quoin eve.Quoin
	cursor, err := r.DB(db.DbName).Table(QUOIN_TABLE).Get(r.UUID(name)).Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return nil, err
	}
	if cursor.IsNil() {
		return nil, nil
	}
	if err = cursor.One(&quoin); err != nil {
		return nil, err
	}
	return &quoin, nil
}

func (db *DbSession) InsertQuoinArchive(quoinArchive *eve.QuoinArchive) error {
	res, err := r.DB(db.DbName).Table(QUOIN_ARCHIVE_TABLE).Insert(
		map[string]interface{}{
			"QuoinName": quoinArchive.QuoinName,
			"Modules":   quoinArchive.Modules,
			"Authorization": map[string]interface{}{
				"Owner":       quoinArchive.Authorization.Owner,
				"GroupAccess": quoinArchive.Authorization.GroupAccess,
			},
			"Timestamp": r.EpochTime(time.Now().Unix()),
		}).RunWrite(db.Session)
	if err != nil {
		return err
	}
	if res.Inserted == 1 {
		quoinArchive.Id = res.GeneratedKeys[0]
	}
	quoinData, err := db.GetQuoinByName(quoinArchive.QuoinName)
	if err != nil {
		return err
	}
	quoinData.ArchiveUri = strings.Join(
		[]string{
			strings.SplitAfter(quoinData.ArchiveUri, "/upload")[0],
			"/",
			quoinArchive.Id,
		}, "")
	quoinData.Status = eve.VALIDATED
	// Update Quoin with Archive's id value
	if err := db.UpdateQuoin(quoinArchive.QuoinName, quoinData); err != nil {
		return err
	}
	log.Printf("%d row inserted. \n", res.Inserted)
	return nil
}

func (db *DbSession) GetQuoinArchiveById(id string) (*eve.QuoinArchive, error) {
	var quoinArchive eve.QuoinArchive
	cursor, err := r.DB(db.DbName).Table(QUOIN_ARCHIVE_TABLE).Get(id).Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return &quoinArchive, err
	}
	if cursor.IsNil() {
		return nil, nil
	}
	if err = cursor.One(&quoinArchive); err != nil {
		return &quoinArchive, err
	}
	return &quoinArchive, nil
}
