package autosqlconf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacMysqlProcessor struct {
	processors.ZeroCoreProcessor
}

func (processor *ZeroXsacMysqlProcessor) ColumnExists(tableSchema string, tableName string, columName string) (int, error) {
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

func (processor *ZeroXsacMysqlProcessor) ColumnDiff(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) (int, error) {
	const COLUMN_EXISTS_SQL = "CALL COLUMN_DIFF(? ,? ,? ,? ,? ,?)"
	if strings.ToUpper(columnDefault) == structs.XSAC_NULL {
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

func (processor *ZeroXsacMysqlProcessor) DMLColumn(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) error {
	const DML_COLUMN_SQL = "CALL DML_COLUMN(? ,? ,?, ?, ?, ?)"
	if strings.ToUpper(columnDefault) == "NULL" {
		_, err := processor.PreparedStmt(DML_COLUMN_SQL).Exec(tableSchema, tableName, columName, isNullable, columnType, nil)
		return err
	} else {
		_, err := processor.PreparedStmt(DML_COLUMN_SQL).Exec(tableSchema, tableName, columName, isNullable, columnType, columnDefault)
		return err
	}
}
func (processor *ZeroXsacMysqlProcessor) DropColumn(tableSchema string, tableName string, columName string) error {
	const DROP_COLUMN_SQL = "CALL DROP_COLUMN(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_COLUMN_SQL).Exec(tableSchema, tableName, columName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) IndexExists(tableSchema string, tableName string, indexName string) (int, error) {
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

func (processor *ZeroXsacMysqlProcessor) DMLConstraint(tableSchema string, tableName string, indexName string, defineIndexSQL string) error {
	// const DML_CONSTRAINT_SQL = "CALL DML_INDEX($1 ,$2 ,$3, $4)"
	// _, err := processor.PreparedStmt(DML_CONSTRAINT_SQL).Exec(tableSchema, tableName, indexName, defineIndexSQL)
	// return err
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DropConstraint(tableSchema string, tableName string, indexName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, indexName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DMLIndex(tableSchema string, tableName string, colnumName string) error {
	// const DML_INDEX_SQL = "CALL DML_INDEX(? ,? ,?)"
	// _, err := processor.PreparedStmt(DML_INDEX_SQL).Exec(tableSchema, tableName, colnumName)
	// return err
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DropIndex(tableSchema string, tableName string, colnumName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, colnumName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) TriggerExists(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) (int, error) {
	return 0, nil
}
func (processor *ZeroXsacMysqlProcessor) DMLTrigger(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DropTrigger(tableSchema string, tableName string, triggerName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DMLPrimary(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroXsacMysqlProcessor) DropPrimary(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroXsacMysqlProcessor) DMLUnique(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroXsacMysqlProcessor) DropUnique(tableSchema string, tableName string, columnName string) error {
	return nil
}
func (processor *ZeroXsacMysqlProcessor) DMLForeign(tableSchema string, tableName string, columnName string, relTableName string, relColumnName string) error {
	return nil
}
func (processor *ZeroXsacMysqlProcessor) DropForeign(tableSchema string, tableName string, columnName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) TableExists(tableSchema string, tableName string) (int, error) {
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

func (processor *ZeroXsacMysqlProcessor) DMLTable(tableSchema string, tableName string) error {
	const DML_TABLE_SQL = "CALL DML_TABLE(? ,?)"
	_, err := processor.PreparedStmt(DML_TABLE_SQL).Exec(tableSchema, tableName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) Create0Struct(tableSchema string, tableName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) Create0FlagStruct(tableSchema string, tableName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DMLY0SPart(tableSchema string, tableName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DMLM0SPart(tableSchema string, tableName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DMLD0SPart(tableSchema string, tableName string) error {
	return nil
}

func (processor *ZeroXsacMysqlProcessor) DropPartitionTable(tableSchema string, tableName string) error {
	return nil
}
