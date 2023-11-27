package processors

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	RELATION_AND = "AND"
	RELATION_OR  = "OR"

	SYMBOL_EQ   = "EQ"
	SYMBOL_NEQ  = "NEQ"
	SYMBOL_GT   = "GT"
	SYMBOL_LT   = "LT"
	SYMBOL_EQGT = "EQGT"
	SYMBOL_EQLT = "EQLT"
	SYMBOL_LIKE = "LIKE"

	ORDER_BY_ASC  = "ASC"
	ORDER_BY_DESC = "DESC"
)

func relations() map[string]string {
	return map[string]string{
		RELATION_AND: " AND ",
		RELATION_OR:  " OR ",
	}
}

func symbols() map[string]string {
	return map[string]string{
		SYMBOL_EQ:   " = ",
		SYMBOL_NEQ:  " <> ",
		SYMBOL_GT:   " > ",
		SYMBOL_LT:   " < ",
		SYMBOL_EQGT: " >= ",
		SYMBOL_EQLT: " <= ",
		SYMBOL_LIKE: " LIKE ",
	}
}

func orderby() map[string]string {
	return map[string]string{
		ORDER_BY_ASC:  " ASC ",
		ORDER_BY_DESC: " DESC ",
	}
}

func exHumpToLine(name string) string {
	_name := name
	if strings.HasPrefix(_name, "ID") {
		_name = strings.ReplaceAll(_name, "ID", "id")
	} else {
		_name = strings.ReplaceAll(_name, "ID", "_id")
	}

	namebytes := []byte(_name)
	var buf bytes.Buffer
	for i, c := range namebytes {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				buf.WriteByte('_')
			}
			buf.WriteByte(c + 32)
		} else {
			buf.WriteByte(c)
		}
	}
	return string(buf.Bytes())
}

func exLineToHump(name string) string {
	namebytes := []byte(name)
	var buf bytes.Buffer
	nextc := false
	for i, c := range namebytes {
		if c == '_' {
			if i > 0 {
				nextc = true
			}
		} else if nextc {
			buf.WriteByte(c - 32)
			nextc = false
		} else {
			buf.WriteByte(c)
		}
	}
	return string(buf.Bytes())
}

type ZeroCondition struct {
	Symbol   string           `json:"symbol,omitempty"`
	Column   string           `json:"column,omitempty"`
	Value    string           `json:"value,omitempty"`
	Relation []*ZeroCondition `json:"relation,omitempty"`
}

type ZeroOrderBy struct {
	Column string `json:"column,omitempty"`
	Seq    string `json:"seq,omitempty"`
}

type ZeroLimit struct {
	Start  int `json:"start,omitempty"`
	Length int `json:"length,omitempty"`
}

type ZeroQuery struct {
	Columns   []string       `json:"columns,omitempty"`
	Condition *ZeroCondition `json:"condition,omitempty"`
	Orderby   []*ZeroOrderBy `json:"orderby,omitempty"`
	Limit     *ZeroLimit     `json:"limit,omitempty"`
}

type ZeroQueryOperation struct {
	ZeroCoreProcessor

	Query     *ZeroQuery
	TableName string

	distinctID      string
	filterTableName string

	columns    string
	conditions string
	orderby    string
	limit      string

	Start  int
	Length int
}

func (opera *ZeroQueryOperation) Build(transaction *sql.Tx) {
	opera.ZeroCoreProcessor.Build(transaction)
	opera.makeColumns()
	opera.makeConditions()
	opera.makeOrderby()
	opera.makeLimit()
}

func (opera *ZeroQueryOperation) SetDistinctID(distinctID string) {
	opera.distinctID = distinctID
}

func (opera *ZeroQueryOperation) SetFilterTableName(filterTableName string) {
	opera.filterTableName = filterTableName
}

func (opera *ZeroQueryOperation) AppendCondition(condition string) {
	if len(opera.conditions) <= 0 {
		opera.conditions = fmt.Sprintf(" WHERE (%s)", condition)
	} else {
		opera.conditions = fmt.Sprintf(" %s AND (%s)", opera.conditions, condition)
	}
}

func parserConditions(condition *ZeroCondition) (string, error) {
	if condition.Relation == nil || len(condition.Relation) <= 0 {
		symbol, ok := symbols()[condition.Symbol]
		if !ok {
			return "", errors.New(fmt.Sprintf("symbol `%s` not found", condition.Symbol))
		}

		return fmt.Sprintf("(`%s` %s '%s')", exHumpToLine(condition.Column), symbol, condition.Value), nil
	} else {
		relatSymbol, ok := relations()[condition.Symbol]
		if !ok {
			return "", errors.New(fmt.Sprintf("relation `%s` not found", condition.Symbol))
		}

		relats := make([]string, len(condition.Relation))
		for i, relation := range condition.Relation {
			relat, err := parserConditions(relation)
			if err != nil {
				return "", nil
			}
			relats[i] = relat
		}
		return fmt.Sprintf("(%s)", strings.Join(relats, relatSymbol)), nil
	}
}

func (opera *ZeroQueryOperation) makeColumns() {

	columnsPrefix := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		columnsPrefix = "a."
	}

	if opera.Query.Columns == nil || len(opera.Query.Columns) <= 0 {
		opera.columns = fmt.Sprintf(" %s* ", columnsPrefix)
	} else {
		columns := make([]string, len(opera.Query.Columns))
		for i, column := range opera.Query.Columns {
			columns[i] = fmt.Sprintf("`%s`", exHumpToLine(column))
		}
		opera.columns = fmt.Sprintf(" %s%s ", columnsPrefix, strings.Join(columns, fmt.Sprintf(", %s", columnsPrefix)))
	}
}

func (opera *ZeroQueryOperation) makeConditions() error {
	if opera.Query.Condition == nil || len(opera.Query.Condition.Symbol) <= 0 {
		opera.conditions = ""
	} else {
		condi, err := parserConditions(opera.Query.Condition)
		if err != nil {
			return err
		}
		opera.conditions = fmt.Sprintf(" WHERE %s ", condi)
	}
	return nil
}

func (opera *ZeroQueryOperation) makeOrderby() {
	if opera.Query.Orderby == nil || len(opera.Query.Orderby) <= 0 {
		opera.orderby = ""
	} else {
		orders := make([]string, len(opera.Query.Orderby))
		for i, o := range opera.Query.Orderby {
			orders[i] = fmt.Sprintf(" `%s` %s", exHumpToLine(o.Column), o.Seq)
		}
		opera.orderby = fmt.Sprintf(" ORDER BY %s ", strings.Join(orders, ","))
	}
}

func (opera *ZeroQueryOperation) makeLimit() {

	if opera.Query.Limit.Length > 0 {
		if opera.Query.Limit.Length > 100 {
			opera.Start = opera.Query.Limit.Start
			opera.Length = 100
		} else {
			opera.Start = opera.Query.Limit.Start
			opera.Length = opera.Query.Limit.Length
		}
	} else {
		opera.Start = 0
		opera.Length = 1
	}

	opera.limit = fmt.Sprintf(" LIMIT %d ,%d ", opera.Start, opera.Length)
}

const DISTINCT_QUERY_SQL_TEMPLATE = `	
	SELECT 
		{{columns}}
	FROM 
		(SELECT 
			distinct {{distinctID}} 
			FROM 
				{{filterTableName}} 
				{{conditions}}) t,
				{{tableName}} a
	WHERE 
		t.{{distinctID}} = a.ID 
		{{orderby}} {{limit}}
`

const DISTINCT_QUERY_COUNT_SQL_TEMPLATE = `
	SELECT 
		count(distinct {{distinctID}}) AS QUERY_COUNT
	FROM 
		{{filterTableName}} 
	{{conditions}}
`

func (opera *ZeroQueryOperation) makeQuerySQL() string {
	querySQL := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		querySQL = strings.ReplaceAll(DISTINCT_QUERY_SQL_TEMPLATE, "{{columns}}", opera.columns)
		querySQL = strings.ReplaceAll(querySQL, "{{tableName}}", opera.TableName)
		querySQL = strings.ReplaceAll(querySQL, "{{conditions}}", opera.conditions)
		querySQL = strings.ReplaceAll(querySQL, "{{orderby}}", opera.orderby)
		querySQL = strings.ReplaceAll(querySQL, "{{limit}}", opera.limit)
		querySQL = strings.ReplaceAll(querySQL, "{{distinctID}}", opera.distinctID)
		querySQL = strings.ReplaceAll(querySQL, "{{filterTableName}}", opera.filterTableName)
	} else {
		querySQL = fmt.Sprintf("SELECT%sFROM %s %s %s %s",
			opera.columns,
			opera.TableName,
			opera.conditions,
			opera.orderby,
			opera.limit)
	}

	return querySQL
}

func (opera *ZeroQueryOperation) makeQueryCountSQL() string {
	queryCountSQL := ""
	if len(opera.distinctID) > 0 && len(opera.filterTableName) > 0 {
		queryCountSQL = strings.ReplaceAll(DISTINCT_QUERY_COUNT_SQL_TEMPLATE, "{{conditions}}", opera.conditions)
		queryCountSQL = strings.ReplaceAll(DISTINCT_QUERY_COUNT_SQL_TEMPLATE, "{{distinctID}}", opera.distinctID)
		queryCountSQL = strings.ReplaceAll(DISTINCT_QUERY_COUNT_SQL_TEMPLATE, "{{filterTableName}}", opera.filterTableName)
	} else {
		queryCountSQL = fmt.Sprintf("SELECT count(1) AS QUERY_COUNT FROM %s%s", opera.TableName, opera.conditions)
	}
	return queryCountSQL
}

func (opera *ZeroQueryOperation) Exec() ([]map[string]interface{}, map[string]interface{}) {
	queryCountSQL := opera.makeQueryCountSQL()

	rows, err := opera.preparedStmt(queryCountSQL).Query()
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
	rows, err = opera.preparedStmt(querySQL).Query()

	rowsmap := parser(rows)
	expands := make(map[string]interface{})
	expands["start"] = strconv.Itoa(opera.Start)
	expands["length"] = strconv.Itoa(opera.Length)
	expands["total"] = strconv.FormatInt(total, 10)
	return rowsmap, expands
}
