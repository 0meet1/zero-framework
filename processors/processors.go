package processors

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ZeroCoreProcessor struct {
	transaction *sql.Tx
	prepares    map[string]*sql.Stmt
}

func (processor *ZeroCoreProcessor) Build(transaction *sql.Tx) {
	processor.transaction = transaction
	processor.prepares = make(map[string]*sql.Stmt)
}

func (processor *ZeroCoreProcessor) Exec(execsql string) (sql.Result, error) {
	return processor.transaction.Exec(execsql)
}

func (processor *ZeroCoreProcessor) Parser(rows *sql.Rows) []map[string]interface{} {
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

func (processor *ZeroCoreProcessor) registeyPreparedStatement(preparedSQL string) error {
	stmt, err := processor.transaction.Prepare(preparedSQL)
	if err != nil {
		return err
	}
	processor.prepares[preparedSQL] = stmt
	return nil
}

func (processor *ZeroCoreProcessor) existsPreparedStatement(preparedSQL string) bool {
	_, ok := processor.prepares[preparedSQL]
	return ok
}

func (processor *ZeroCoreProcessor) PreparedStmt(preparedSQL string) *sql.Stmt {
	if !processor.existsPreparedStatement(preparedSQL) {
		err := processor.registeyPreparedStatement(preparedSQL)
		if err != nil {
			panic(err)
		}
	}
	stmt, ok := processor.prepares[preparedSQL]
	if !ok {
		return nil
	}
	return stmt
}

func (processor *ZeroCoreProcessor) DatabaseDatetime() (*time.Time, error) {
	const FETCH_DATE_SQL = "SELECT current_timestamp FROM DUAL"
	rows, err := processor.PreparedStmt(FETCH_DATE_SQL).Query()
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New(fmt.Sprintf("query -> %s result error", FETCH_DATE_SQL))
	}
	var datetime time.Time
	err = rows.Scan(&datetime)
	if err != nil {
		return nil, err
	}
	return &datetime, nil
}
