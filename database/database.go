package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type xConnectPool struct {
	database *gorm.DB
}

func (cp *xConnectPool) init(database *gorm.DB) {
	cp.database = database
}

func (cp *xConnectPool) Connect() *sql.DB {
	connect, err := cp.database.DB()
	if err != nil {
		panic(err)
	}
	if err := connect.Ping(); err != nil {
		panic(err)
	}
	return connect
}

func (cp *xConnectPool) Transaction() *sql.Tx {
	transaction, err := cp.Connect().Begin()
	if err != nil {
		panic(err)
	}
	return transaction
}
