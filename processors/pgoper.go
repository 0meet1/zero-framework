package processors

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type ZeroPostgresQueryOperation struct {
	ZeroCoreProcessor

	query     *ZeroQuery
	tableName string

	distinctID      string
	filterTableName string

	columns    string
	conditions string
	orderby    string
	limit      string

	Start  int
	Length int
}

func NewZeroPostgresQueryOperation(xQuery *ZeroQuery, xTableName string) *ZeroPostgresQueryOperation {
	return &ZeroPostgresQueryOperation{
		query:     xQuery,
		tableName: xTableName,
	}
}

func (opera *ZeroPostgresQueryOperation) Build(transaction *sql.Tx) {
	opera.ZeroCoreProcessor.Build(transaction)
	opera.makeColumns()
	opera.makeConditions()
	opera.makeOrderby()
	opera.makeLimit()
}

func (opera *ZeroPostgresQueryOperation) AddQuery(xQuery *ZeroQuery)    { opera.query = xQuery }
func (opera *ZeroPostgresQueryOperation) AddTableName(tableName string) { opera.tableName = tableName }
func (opera *ZeroPostgresQueryOperation) AddDistinctID(distinctID string) {
	opera.distinctID = distinctID
}
func (opera *ZeroPostgresQueryOperation) AddFilterTableName(filterTableName string) {
	opera.filterTableName = filterTableName
}

func (opera *ZeroPostgresQueryOperation) AppendCondition(condition string) {
	if len(opera.conditions) <= 0 {
		opera.conditions = fmt.Sprintf(" WHERE (%s)", condition)
	} else {
		opera.conditions = fmt.Sprintf(" %s AND (%s)", opera.conditions, condition)
	}
}

func (opera *ZeroPostgresQueryOperation) jsonColumnName(name string) string {
	fpidx := strings.Index(name, ".")
	if fpidx <= 0 {
		return name
	}
	return fmt.Sprintf(`"%s" #>> '{%s}'`, exHumpToLine(name[:fpidx]), strings.ReplaceAll(name[fpidx+1:], ".", ","))
}

func (opera *ZeroPostgresQueryOperation) parserConditions(condition *ZeroCondition) (string, error) {
	if len(condition.Relation) <= 0 {

		if strings.HasPrefix(condition.Column, "@!") {
			return fmt.Sprintf(`("%s" %s)`, strings.ReplaceAll(condition.Column, "@!", ""), condition.Value), nil
		} else {
			symbol, ok := symbols()[condition.Symbol]
			if !ok {
				return "", fmt.Errorf("symbol `%s` not found", condition.Symbol)
			}

			if strings.HasPrefix(condition.Column, "@") {
				return fmt.Sprintf(`("%s" %s '%s')`, strings.ReplaceAll(condition.Column, "@", ""), symbol, condition.Value), nil
			} else {
				if strings.Index(condition.Column, ".") > 1 {
					return fmt.Sprintf(`(%s %s '%s')`, opera.jsonColumnName(condition.Column), symbol, condition.Value), nil
				} else {
					return fmt.Sprintf(`("%s" %s '%s')`, exHumpToLine(condition.Column), symbol, condition.Value), nil
				}
			}
		}
	} else {
		relatSymbol, ok := relations()[condition.Symbol]
		if !ok {
			return "", fmt.Errorf("relation `%s` not found", condition.Symbol)
		}

		relats := make([]string, len(condition.Relation))
		for i, relation := range condition.Relation {
			relat, err := opera.parserConditions(relation)
			if err != nil {
				return "", nil
			}
			relats[i] = relat
		}
		return fmt.Sprintf("(%s)", strings.Join(relats, relatSymbol)), nil
	}
}

func (opera *ZeroPostgresQueryOperation) makeColumns() {

	columnsPrefix := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		columnsPrefix = "a."
	}

	if len(opera.query.Columns) <= 0 {
		opera.columns = fmt.Sprintf(" %s* ", columnsPrefix)
	} else {
		columns := make([]string, len(opera.query.Columns))
		for i, column := range opera.query.Columns {
			columns[i] = fmt.Sprintf("\"%s\"", exHumpToLine(column))
		}
		opera.columns = fmt.Sprintf(" %s%s ", columnsPrefix, strings.Join(columns, fmt.Sprintf(", %s", columnsPrefix)))
	}
}

func (opera *ZeroPostgresQueryOperation) makeConditions() error {
	if opera.query.Condition == nil || len(opera.query.Condition.Symbol) <= 0 {
		opera.conditions = ""
	} else {
		condi, err := opera.parserConditions(opera.query.Condition)
		if err != nil {
			return err
		}
		opera.conditions = fmt.Sprintf(" WHERE %s ", condi)
	}
	return nil
}

func (opera *ZeroPostgresQueryOperation) makeOrderby() {
	if len(opera.query.Orderby) <= 0 {
		opera.orderby = ""
	} else {
		orders := make([]string, len(opera.query.Orderby))
		for i, o := range opera.query.Orderby {
			orders[i] = fmt.Sprintf(" \"%s\" %s", exHumpToLine(o.Column), o.Seq)
		}
		opera.orderby = fmt.Sprintf(" ORDER BY %s ", strings.Join(orders, ","))
	}
}

func (opera *ZeroPostgresQueryOperation) makeLimit() {

	if opera.query.Limit.Length > 0 {
		if opera.query.Limit.Length > 5000 {
			opera.Start = opera.query.Limit.Start
			opera.Length = 5000
		} else {
			opera.Start = opera.query.Limit.Start
			opera.Length = opera.query.Limit.Length
		}
	} else {
		opera.Start = 0
		opera.Length = 1
	}

	opera.limit = fmt.Sprintf(" OFFSET %d LIMIT %d ", opera.Start, opera.Length)
}

const DISTINCT_POSTGRES_QUERY_SQL_TEMPLATE = `	
	SELECT 
		{{columns}}
	FROM 
		(SELECT 
			distinct {{distinctID}} AS row095c_id
			FROM 
				{{filterTableName}} 
				{{conditions}}) t,
				{{tableName}} a
	WHERE 
		t.row0c_id = a.{{distinctID}} 
		{{orderby}} {{limit}}
`

const DISTINCT_POSTGRES_QUERY_COUNT_SQL_TEMPLATE = `
	SELECT 
		count(distinct {{distinctID}}) AS QUERY_COUNT
	FROM 
		{{filterTableName}} 
	{{conditions}}
`

func (opera *ZeroPostgresQueryOperation) makeQuerySQL() string {
	querySQL := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		querySQL = strings.ReplaceAll(DISTINCT_POSTGRES_QUERY_SQL_TEMPLATE, "{{columns}}", opera.columns)
		querySQL = strings.ReplaceAll(querySQL, "{{tableName}}", opera.tableName)
		querySQL = strings.ReplaceAll(querySQL, "{{conditions}}", opera.conditions)
		querySQL = strings.ReplaceAll(querySQL, "{{orderby}}", opera.orderby)
		querySQL = strings.ReplaceAll(querySQL, "{{limit}}", opera.limit)
		querySQL = strings.ReplaceAll(querySQL, "{{distinctID}}", opera.distinctID)
		querySQL = strings.ReplaceAll(querySQL, "{{filterTableName}}", opera.filterTableName)
	} else {
		querySQL = fmt.Sprintf("SELECT%sFROM %s %s %s %s",
			opera.columns,
			opera.tableName,
			opera.conditions,
			opera.orderby,
			opera.limit)
	}

	return querySQL
}

func (opera *ZeroPostgresQueryOperation) makeQueryCountSQL() string {
	queryCountSQL := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		queryCountSQL = strings.ReplaceAll(DISTINCT_POSTGRES_QUERY_COUNT_SQL_TEMPLATE, "{{conditions}}", opera.conditions)
		queryCountSQL = strings.ReplaceAll(queryCountSQL, "{{distinctID}}", opera.distinctID)
		queryCountSQL = strings.ReplaceAll(queryCountSQL, "{{filterTableName}}", opera.filterTableName)
	} else {
		queryCountSQL = fmt.Sprintf("SELECT count(1) AS QUERY_COUNT FROM %s%s", opera.tableName, opera.conditions)
	}
	return queryCountSQL
}

func (opera *ZeroPostgresQueryOperation) Exec() ([]map[string]interface{}, map[string]interface{}) {
	queryCountSQL := opera.makeQueryCountSQL()

	rows, err := opera.PreparedStmt(queryCountSQL).Query()
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		panic(fmt.Sprintf("query -> %s result error", queryCountSQL))
	}
	var total int64
	err = rows.Scan(&total)
	if err != nil {
		panic(err)
	}
	rows.Close()

	querySQL := opera.makeQuerySQL()
	rows, err = opera.PreparedStmt(querySQL).Query()
	if err != nil {
		panic(err)
	}

	rowsmap := opera.Parser(rows)
	expands := make(map[string]interface{})
	expands["start"] = strconv.Itoa(opera.Start)
	expands["length"] = strconv.Itoa(opera.Length)
	expands["total"] = strconv.FormatInt(total, 10)
	return rowsmap, expands
}
