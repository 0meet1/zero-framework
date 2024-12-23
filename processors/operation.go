package processors

import (
	"bytes"
	"database/sql"
	"fmt"
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

func Orderby() map[string]string {
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
	return buf.String()
}

func ExLineToHump(name string) string {
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
	return buf.String()
}

func parseJSONColumnName(name string) string {
	fpidx := strings.Index(name, ".")
	if fpidx <= 0 {
		return name
	}
	return fmt.Sprintf(`%s ->> "$.%s"`, exHumpToLine(name[:fpidx]), name[fpidx+1:])
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

type ZeroQueryOperation interface {
	Build(transaction *sql.Tx)

	AddQuery(xQuery *ZeroQuery)
	AddTableName(tableName string)
	AddDistinctID(distinctID string)
	AddFilterTableName(filterTableName string)
	AppendCondition(condition string)

	Exec() ([]map[string]interface{}, map[string]interface{})
}
