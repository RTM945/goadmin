package db

import (
	"database/sql"
	"sync"

	"goadmin/modules/config"
)

// Base is a common Connection.
type Base struct {
	DbList  map[string]*sql.DB
	Once    sync.Once
	Configs config.DatabaseList
}

// Close implements the method Connection.Close.
func (db *Base) Close() []error {
	errs := make([]error, 0)
	for _, d := range db.DbList {
		errs = append(errs, d.Close())
	}
	return errs
}

// GetDB implements the method Connection.GetDB.
func (db *Base) GetDB(key string) *sql.DB {
	return db.DbList[key]
}

func (db *Base) GetConfig(name string) config.Database {
	return db.Configs[name]
}
