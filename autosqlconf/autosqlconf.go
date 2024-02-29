package autosqlconf

import (
	"database/sql"
)

const (
	XSAC_TRIGGER_TIMING_BEFORE = "BEFORE"
	XSAC_TRIGGER_TIMING_AFTER  = "AFTER"

	XSAC_TRIGGER_EVENT_INSERT = "INSERT"
	XSAC_TRIGGER_EVENT_UPDATE = "UPDATE"
	XSAC_TRIGGER_EVENT_DELETE = "DELETE"
	XSAC_TRIGGER_EVENT_SELECT = "SELECT"
)

type ZeroXsacProcessor interface {
	Build(transaction *sql.Tx)

	ColumnExists(tableSchema string, tableName string, columName string) (int, error)
	ColumnDiff(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) (int, error)
	DMLColumn(tableSchema string, tableName string, columName string, isNullable string, columnType string, columnDefault string) error
	DropColumn(tableSchema string, tableName string, columName string) error

	IndexExists(tableSchema string, tableName string, indexName string) (int, error)
	DMLConstraint(tableSchema string, tableName string, indexName string, defineIndexSQL string) error
	DropConstraint(tableSchema string, tableName string, indexName string) error
	DMLIndex(tableSchema string, tableName string, colnumName string) error
	DropIndex(tableSchema string, tableName string, colnumName string) error

	TriggerExists(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) (int, error)
	DMLTrigger(tableSchema string, tableName string, triggerTiming string, triggerEvent string, triggerName string, triggerAction string) error
	DropTrigger(tableSchema string, tableName string, triggerName string) error

	DMLPrimary(tableSchema string, tableName string, columnName string) error
	DropPrimary(tableSchema string, tableName string, columnName string) error
	DMLUnique(tableSchema string, tableName string, columnName string) error
	DropUnique(tableSchema string, tableName string, columnName string) error
	DMLForeign(tableSchema string, tableName string, columnName string, relTableName string, relColumnName string) error
	DropForeign(tableSchema string, tableName string, columnName string) error

	TableExists(tableSchema string, tableName string) (int, error)
	DMLTable(tableSchema string, tableName string) error

	Create0Struct(tableSchema string, tableName string) error
	Create0FlagStruct(tableSchema string, tableName string) error

	DMLY0SPart(tableSchema string, tableName string) error
	DMLM0SPart(tableSchema string, tableName string) error
	DMLD0SPart(tableSchema string, tableName string) error
	DropPartitionTable(tableSchema string, tableName string) error
}
