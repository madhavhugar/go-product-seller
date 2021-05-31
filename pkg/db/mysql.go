package database

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

// Adapter provides the DB functionality
type Adapter interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

type SqlInstance struct {
	*sql.DB
}

var (
	productDB  *SqlInstance
	initDBOnce sync.Once
)

func initConnection(driverName string, connectionURL string) {
	db, err := sql.Open(driverName, connectionURL)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create server")
		os.Exit(1)
	}
	productDB = &SqlInstance{db}
}

func InitConnection(driverName string, connectionURL string) Adapter {
	initDBOnce.Do(func() {
		initConnection(driverName, connectionURL)
	})
	return productDB
}

func (s *SqlInstance) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return productDB.DB.Query(query, args...)
}

func (s *SqlInstance) Close() error {
	if productDB == nil {
		return nil
	}
	return productDB.DB.Close()
}
