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
	const COLUMN_EXISTS_SQL = "SELECT COLUMN_DIFF(? ,? ,? ,? ,? ,?)"
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
	const INDEX_EXISTS_SQL = "SELECT INDEX_EXISTS(? ,? ,?)"
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
	const DML_CONSTRAINT_SQL = "CALL DML_INDEX(? ,? ,?, ?)"
	_, err := processor.PreparedStmt(DML_CONSTRAINT_SQL).Exec(tableSchema, tableName, indexName, defineIndexSQL)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DropConstraint(tableSchema string, tableName string, indexName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, indexName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DMLIndex(tableSchema string, tableName string, colnumName string) error {
	const DML_INDEX_SQL = "CALL DML_INDEX(? ,? ,?, ?)"
	_, err := processor.PreparedStmt(DML_INDEX_SQL).Exec(tableSchema, tableName, fmt.Sprintf("idx_%s_%s", tableName, colnumName), fmt.Sprintf("ADD KEY (`%s`)", colnumName))
	return err
}

func (processor *ZeroXsacMysqlProcessor) DropIndex(tableSchema string, tableName string, colnumName string) error {
	const DROP_INDEX_SQL = "CALL DROP_INDEX(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_INDEX_SQL).Exec(tableSchema, tableName, colnumName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) TriggerExists(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) (int, error) {
	const TRIGGER_EXISTS_SQL = "SELECT TRIGGER_EXISTS(? ,? ,? ,? ,? ,?)"
	rows, err := processor.PreparedStmt(TRIGGER_EXISTS_SQL).Query(tableSchema, tableName, triggerTiming, triggerEvent, triggerName, triggerAction)
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

func (processor *ZeroXsacMysqlProcessor) DMLTrigger(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) error {
	_, err := processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s`", tableName))
	if err != nil {
		return err
	}
	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s` %s %s ON `%s` FOR EACH ROW BEGIN %s END", triggerName, triggerTiming, triggerEvent, tableName, triggerAction))
	return err
}

func (processor *ZeroXsacMysqlProcessor) DropTrigger(tableSchema string, tableName string, triggerName string) error {
	_, err := processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s`", tableName))
	return err
}

func (processor *ZeroXsacMysqlProcessor) DMLPrimary(tableSchema string, tableName string, columnName string) error {
	const DML_PRIMARY_SQL = "CALL DML_PRIMARY(? ,? ,?)"
	_, err := processor.PreparedStmt(DML_PRIMARY_SQL).Exec(tableSchema, tableName, columnName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DropPrimary(tableSchema string, tableName string, columnName string) error {
	const DROP_PRIMARY_SQL = "CALL DROP_PRIMARY(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_PRIMARY_SQL).Exec(tableSchema, tableName, columnName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DMLUnique(tableSchema string, tableName string, columnName string) error {
	const DML_UNIQUE_SQL = "CALL DML_UNIQUE(? ,? ,?)"
	_, err := processor.PreparedStmt(DML_UNIQUE_SQL).Exec(tableSchema, tableName, columnName)
	return err
}

func (processor *ZeroXsacMysqlProcessor) DropUnique(tableSchema string, tableName string, columnName string) error {
	const DROP_UNIQUE_SQL = "CALL DROP_UNIQUE(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_UNIQUE_SQL).Exec(tableSchema, tableName, columnName)
	return err
}
func (processor *ZeroXsacMysqlProcessor) DMLForeign(tableSchema string, tableName string, columnName string, relTableName string, relColumnName string) error {
	const DML_FOREIGN_SQL = "CALL DML_FOREIGN(? ,? ,? ,? ,?)"
	_, err := processor.PreparedStmt(DML_FOREIGN_SQL).Exec(tableSchema, tableName, columnName, relTableName, relColumnName)
	return err
}
func (processor *ZeroXsacMysqlProcessor) DropForeign(tableSchema string, tableName string, columnName string) error {
	const DROP_FOREIGN_SQL = "CALL DROP_FOREIGN(? ,? ,?)"
	_, err := processor.PreparedStmt(DROP_FOREIGN_SQL).Exec(tableSchema, tableName, columnName)
	return err
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
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s_uuid`", tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s_uuid` BEFORE INSERT ON `%s` FOR EACH ROW BEGIN IF new.id = '-' THEN SET new.id = (SELECT uuid()); END IF; END", tableName, tableName))
	return err
}

func (processor *ZeroXsacMysqlProcessor) Create0Struct(tableSchema string, tableName string) error {
	const CREATE_0STRUCT_SQL = "CALL create_0struct(? ,?)"
	_, err := processor.PreparedStmt(CREATE_0STRUCT_SQL).Exec(tableSchema, tableName)
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s_uuid`", tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s_uuid` BEFORE INSERT ON `%s` FOR EACH ROW BEGIN IF new.id = '-' THEN SET new.id = (SELECT uuid()); END IF; END", tableName, tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s_update`", tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s_update` BEFORE UPDATE ON `%s` FOR EACH ROW BEGIN SET new.update_time = (SELECT now()); END", tableName, tableName))
	return err
}

func (processor *ZeroXsacMysqlProcessor) Create0FlagStruct(tableSchema string, tableName string) error {
	const CREATE_0FLAGSTRUCT_SQL = "CALL create_0flagstruct(? ,?)"
	_, err := processor.PreparedStmt(CREATE_0FLAGSTRUCT_SQL).Exec(tableSchema, tableName)
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s_uuid`", tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s_uuid` BEFORE INSERT ON `%s` FOR EACH ROW BEGIN IF new.id = '-' THEN SET new.id = (SELECT uuid()); END IF; END", tableName, tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s_update`", tableName))
	if err != nil {
		return err
	}

	_, err = processor.Exec(fmt.Sprintf("CREATE TRIGGER `%s_update` BEFORE UPDATE ON `%s` FOR EACH ROW BEGIN SET new.update_time = (SELECT now()); END", tableName, tableName))
	return err
}

func (processor *ZeroXsacMysqlProcessor) DMLY0SPart(tableSchema string, tableName string) error {
	return errors.New("not support")
}

func (processor *ZeroXsacMysqlProcessor) DMLM0SPart(tableSchema string, tableName string) error {
	return errors.New("not support")
}

func (processor *ZeroXsacMysqlProcessor) DMLD0SPart(tableSchema string, tableName string) error {
	return errors.New("not support")
}

func (processor *ZeroXsacMysqlProcessor) DropPartitionTable(tableSchema string, tableName string) error {
	return errors.New("not support")
}
