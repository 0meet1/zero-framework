package database

import (
	"database/sql"
	"os"

	"github.com/0meet1/zero-framework/global"
	ora "github.com/sijms/go-ora/v2"
)

const (
	DATABASE_ORACLE = "zero.database.oracle"
)

type xOracleConnectsKeeper struct {
	dsn string
}

func newOracleConnectsKeeper(dsn string) *xOracleConnectsKeeper {
	return &xOracleConnectsKeeper{
		dsn: dsn,
	}
}

func (cp *xOracleConnectsKeeper) Connect() *sql.DB {
	err := os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	if err != nil {
		panic(err)
	}

	connect, err := sql.Open("oracle", cp.dsn)
	if err != nil {
		panic(err)
	}

	if err := connect.Ping(); err != nil {
		panic(err)
	}
	return connect
}

func (cp *xOracleConnectsKeeper) Transaction() *sql.Tx {
	transaction, err := cp.Connect().Begin()
	if err != nil {
		panic(err)
	}
	return transaction
}

var OracleDatabase = func() {

	err := os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	if err != nil {
		panic(err)
	}

	dataSource := newOracleConnectsKeeper(ora.BuildUrl(
		global.StringValue("zero.oracle.hostname"),
		global.IntValue("zero.oracle.hostport"),
		global.StringValue("zero.oracle.servOsid"),
		global.StringValue("zero.oracle.dbname"),
		global.StringValue("zero.oracle.password"), nil))
	global.Key(DATABASE_ORACLE, dataSource)
}

var CustomOracleDatabase = func(registerName, hostname string, hostport int, servOsid, dbname, password string) {

	err := os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	if err != nil {
		panic(err)
	}

	dataSource := newOracleConnectsKeeper(ora.BuildUrl(
		hostname,
		hostport,
		servOsid,
		dbname,
		password, nil))
	global.Key(registerName, dataSource)
}
