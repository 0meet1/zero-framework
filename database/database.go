package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type DataSource struct {
	database *gorm.DB
}

func (cp *DataSource) init(database *gorm.DB) {
	cp.database = database
}

func (cp *DataSource) Connect() *sql.DB {
	connect, err := cp.database.DB()
	if err != nil {
		panic(err)
	}
	if err := connect.Ping(); err != nil {
		panic(err)
	}
	return connect
}

func (cp *DataSource) Transaction() *sql.Tx {
	transaction, err := cp.Connect().Begin()
	if err != nil {
		panic(err)
	}
	return transaction
}
