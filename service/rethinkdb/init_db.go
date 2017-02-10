package rethinkdb

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"log"
)

func (db *DbSession) Initialization() error {
	if err := db.dbInit(); err != nil {
		return err
	}
	if err := db.tableInit(); err != nil {
		return err
	}
	return nil
}

func (db *DbSession) dbInit() error {
	cursor, err := r.DBList().Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return err
	}
	var row interface{}
	for cursor.Next(&row) {
		if row.(string) == db.DbName {
			return nil
		}
	}
	log.Println("Creating DB:", db.DbName)
	if _, err := r.DBCreate(db.DbName).RunWrite(db.Session); err != nil {
		return err
	}
	log.Println("DB", db.DbName, "is created!")
	return nil
}

func (db *DbSession) tableInit() error {
	cursor, err := r.DB(db.DbName).TableList().Run(db.Session)
	defer cursor.Close()
	if err != nil {
		return err
	}
	tables := map[string]byte{
		QUOIN_TABLE:         0,
		QUOIN_ARCHIVE_TABLE: 0,
		INFRA_TABLE:         0,
	}
	var row interface{}
	for cursor.Next(&row) {
		name := row.(string)
		if _, ok := tables[name]; ok {
			delete(tables, name)
		}
	}
	for key, _ := range tables {
		log.Println("Creating table:", key)
		if _, err := r.DB(db.DbName).TableCreate(key, r.TableCreateOpts{PrimaryKey: "Id"}).RunWrite(db.Session); err != nil {
			return err
		}
		log.Println("Table", key, "is created!")
	}
	return nil
}
