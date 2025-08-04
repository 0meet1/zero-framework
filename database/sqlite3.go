package database

import (
	"database/sql"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/0meet1/zero-framework/global"
	_ "github.com/mattn/go-sqlite3"
)

type table struct {
	Name string
	Core string
	DML  []string
}

var NewSQLiteTable = func(name, core string, dmls ...string) *table {
	return &table{
		Name: name,
		Core: core,
		DML:  dmls,
	}
}

type SqliteDataSource struct {
	database *sql.DB
	tables   []*table
	mutex    sync.Mutex

	dbaddr string
}

func (st *SqliteDataSource) open(dbaddr string, tables ...*table) {
	st.dbaddr = dbaddr
	st.tables = tables

	_db, err := sql.Open("sqlite3", st.dbaddr)
	if err != nil {
		panic(err)
	}
	st.database = _db

	if len(st.tables) > 0 {
		for _, _table := range st.tables {
			rows, err := _db.Query(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type = 'table' AND name = '%s'", _table.Name))
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			rowmap := st.parse(rows)
			if len(rowmap) <= 0 {
				_, err := _db.Exec(_table.Core)
				if err != nil {
					panic(err)
				}
				for _, _dm := range _table.DML {
					_, err := _db.Exec(_dm)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

func (st *SqliteDataSource) parse(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	values := make([]interface{}, len(columns))
	for index := range values {
		var value interface{}
		values[index] = &value
	}
	var rowsArray []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			panic(err)
		}

		rowmap := make(map[string]interface{})
		for i, data := range values {
			rowmap[columns[i]] = *data.(*interface{})
		}
		rowsArray = append(rowsArray, rowmap)
	}
	if err != nil {
		panic(err)
	}
	return rowsArray
}

func (st *SqliteDataSource) SecureTransaction(performer func(*sql.Tx) any, onevents ...func(error)) any {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	transaction, err := st.database.Begin()
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

const (
	DATABASE_SQLITE = "zero.database.sqlite3"
)

var SQLiteDatabase = func(tables ...*table) {
	dbaddr := global.StringValue("zero.sqlite3.dbaddr")
	if !strings.HasPrefix(dbaddr, "/") {
		dbaddr = path.Join(global.ServerAbsPath(), dbaddr)
	}
	s := &SqliteDataSource{}
	s.open(dbaddr, tables...)
	global.Logger().Infof("sqlite open success with %s", dbaddr)
	global.Key(DATABASE_SQLITE, s)
}

var CustomSQLiteDatabase = func(registerName, prefix string, tables ...*table) {
	dbaddr := global.StringValue(fmt.Sprintf("zero.%s.dbaddr", prefix))
	if !strings.HasPrefix(dbaddr, "/") {
		dbaddr = path.Join(global.ServerAbsPath(), dbaddr)
	}
	s := &SqliteDataSource{}
	s.open(dbaddr, tables...)
	global.Logger().Infof("sqlite open success with %s", dbaddr)
	global.Key(registerName, s)
}
