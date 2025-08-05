package database

import (
	"database/sql"

	"github.com/0meet1/zero-framework/global"
	"gorm.io/gorm"
)

type SecureDataSource interface {
	SecureTransaction(func(*sql.Tx) any, ...func(error)) any
}

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

func (cp *GormDataSource) SecureTransaction(performer func(*sql.Tx) any, onevents ...func(error)) any {
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			if len(onevents) > 0 {
				onevents[0](err.(error))
			}
		}
	}()
	connect, err := cp.database.DB()
	if err != nil {
		global.Logger().ErrorS(err)
		if len(onevents) > 0 {
			onevents[0](err)
		}
		return nil
	}
	if err := connect.Ping(); err != nil {
		global.Logger().ErrorS(err)
		if len(onevents) > 0 {
			onevents[0](err)
		}
		return nil
	}
	transaction, err := connect.Begin()
	if err != nil {
		global.Logger().ErrorS(err)
		if len(onevents) > 0 {
			onevents[0](err)
		}
		return nil
	}
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			transaction.Rollback()
			if len(onevents) > 0 {
				onevents[0](err.(error))
			}
		} else {
			transaction.Commit()
			if len(onevents) > 0 {
				onevents[0](nil)
			}
		}
	}()
	return performer(transaction)
}
