package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type DataSource interface {
	Connect() *sql.DB
	Transaction() *sql.Tx
}

type GormDataSource struct {
	database *gorm.DB
}

func (cp *GormDataSource) init(database *gorm.DB) {
	cp.database = database
}

func (cp *GormDataSource) Connect() *sql.DB {
	connect, err := cp.database.DB()
	if err != nil {
		panic(err)
	}
	if err := connect.Ping(); err != nil {
		panic(err)
	}
	return connect
}

func (cp *GormDataSource) Transaction() *sql.Tx {
	transaction, err := cp.Connect().Begin()
	if err != nil {
		panic(err)
	}
	return transaction
}
