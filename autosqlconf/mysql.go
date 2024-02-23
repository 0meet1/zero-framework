package autosqlconf

import (
	"errors"
	"fmt"
	"strings"

	x0meet1 "github.com/0meet1/zero-framework"
)

type ZeroDbAutoMysqlProcessor struct {
	x0meet1.ZeroCoreProcessor
}

func (processor *ZeroDbAutoMysqlProcessor) ColumnExists(tableSchema string, tableName string, columName string) (int, error) {
	const COLUMN_EXISTS_SQL = "SELECT COLUMN_EXISTS(? ,? ,?)"
	rows, err := processor.PreparedStmt(COLUMN_EXISTS_SQL).Query(tableSchema, tableName, columName)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, errors.New(fmt.Sprintf("query `COLUMN_EXISTS_SQL` failed"))
	}
	var _state int64
	err = rows.Scan(&_state)
	if err != nil {
		return 0, err
	}
	return int(_state), nil
}

func (processor *ZeroDbAutoMysqlProcessor) ColumnDiff(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) (int, error) {
	const COLUMN_EXISTS_SQL = "CALL COLUMN_DIFF(? ,? ,? ,? ,? ,?)"
	if strings.ToUpper(columnDefault) == ZDA_NULL {
		rows, err := processor.PreparedStmt(COLUMN_EXISTS_SQL).Query(tableSchema, tableName, columName, isNullable, columnType, nil)
		defer func() {
			if rows != nil {
				rows.Close()
			}
		}()
		if err != nil {
			return 0, err
		}
		if !rows.Next() {
			return 0, errors.New(fmt.Sprintf("query `COLUMN_EXISTS_SQL` failed"))
		}
		var _state int64
		err = rows.Scan(&_state)
		if err != nil {
			return 0, err
		}
		return int(_state), nil
	} else {
		rows, err := processor.PreparedStmt(COLUMN_EXISTS_SQL).Query(tableSchema, tableName, columName, isNullable, columnType, columnDefault)
		defer func() {
			if rows != nil {
				rows.Close()
			}
		}()
		if err != nil {
			return 0, err
		}
		if !rows.Next() {
			return 0, errors.New(fmt.Sprintf("query `COLUMN_EXISTS_SQL` failed"))
		}
		var _state int64
		err = rows.Scan(&_state)
		if err != nil {
			return 0, err
		}
		return int(_state), nil
	}
}

func (processor *ZeroDbAutoMysqlProcessor) DMLColumn(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) error {
	const DML_COLUMN_SQL = "CALL DML_COLUMN(? ,? ,?, ?, ?, ?)"
	if strings.ToUpper(columnDefault) == "NULL" {
		_, err := processor.PreparedStmt(DML_COLUMN_SQL).Exec(tableSchema, tableName, columName, isNullable, columnType, nil)
		return err
	} else {
		_, err := processor.PreparedStmt(DML_COLUMN_SQL).Exec(tableSchema, tableName, columName, isNullable, columnType, columnDefault)
		return err
	}
}
func (processor *ZeroDbAutoMysqlProcessor) DropColumn(tableSchema string, tableName string, columName string) error {
	const DROP_COLUMN_SQL = "CALL DROP_COLUMN(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_COLUMN_SQL).Exec(tableSchema, tableName, columName)
	return err
}

func (processor *ZeroDbAutoMysqlProcessor) IndexExists(tableSchema string, tableName string, indexName string) (int, error) {
	const INDEX_EXISTS_SQL = "CALL INDEX_EXISTS(? ,? ,?)"
	rows, err := processor.PreparedStmt(INDEX_EXISTS_SQL).Query(tableSchema, tableName, indexName)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, errors.New(fmt.Sprintf("query `COLUMN_EXISTS_SQL` failed"))
	}
	var _state int64
	err = rows.Scan(&_state)
	if err != nil {
		return 0, err
	}
	return int(_state), nil
}

func (processor *ZeroDbAutoMysqlProcessor) DMLConstraint(tableSchema string, tableName string, indexName string, defineIndexSQL string) error {
	// const DML_CONSTRAINT_SQL = "CALL DML_INDEX($1 ,$2 ,$3, $4)"
	// _, err := processor.PreparedStmt(DML_CONSTRAINT_SQL).Exec(tableSchema, tableName, indexName, defineIndexSQL)
	// return err
	return nil
}

func (processor *ZeroDbAutoMysqlProcessor) DropConstraint(tableSchema string, tableName string, indexName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, indexName)
	return err
}

func (processor *ZeroDbAutoMysqlProcessor) DMLIndex(tableSchema string, tableName string, indexName string) error {
	// const DML_INDEX_SQL = "CALL DML_INDEX(? ,? ,?)"
	// _, err := processor.PreparedStmt(DML_INDEX_SQL).Exec(tableSchema, tableName, indexName)
	// return err
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropIndex(tableSchema string, tableName string, indexName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, indexName)
	return err
}

func (processor *ZeroDbAutoMysqlProcessor) TriggerExists(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) (int, error) {
	return 0, nil
}
func (processor *ZeroDbAutoMysqlProcessor) DMLTrigger(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropTrigger(tableSchema string, tableName string, triggerName string) error

func (processor *ZeroDbAutoMysqlProcessor) DMLPrimary(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropPrimary(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DMLUnique(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropUnique(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DMLForeign(tableSchema string, tableName string, columnName string, relTableName string, relColumnName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropForeign(tableSchema string, tableName string, columnName string) error {
	return nil
}

func (processor *ZeroDbAutoMysqlProcessor) TableExists(tableSchema string, tableName string) (int, error) {
	const TABLE_EXISTS_SQL = "CALL TABLE_EXISTS(? ,?)"
	rows, err := processor.PreparedStmt(TABLE_EXISTS_SQL).Query(tableSchema, tableName)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, errors.New(fmt.Sprintf("query `COLUMN_EXISTS_SQL` failed"))
	}
	var _state int64
	err = rows.Scan(&_state)
	if err != nil {
		return 0, err
	}
	return int(_state), nil
}

func (processor *ZeroDbAutoMysqlProcessor) DMLTable(tableSchema string, tableName string) error {
	const DML_TABLE_SQL = "CALL DML_TABLE(? ,?)"
	_, err := processor.PreparedStmt(DML_TABLE_SQL).Exec(tableSchema, tableName)
	return err
}

func (processor *ZeroDbAutoMysqlProcessor) Create0Struct(tableSchema string, tableName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) Create0FlagStruct(tableSchema string, tableName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DML0SPart(tableSchema string, tableName string) error {
	return nil
}
func (processor *ZeroDbAutoMysqlProcessor) DropPartitionTable(tableSchema string, tableName string) error {
	return nil
}
