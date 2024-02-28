package processors

import "database/sql"

const (
	XSAC_BE_INSERT = "beinsert"
	XSAC_BE_UPDATE = "beupdate"
	XSAC_BE_DELETE = "bedelete"

	XSAC_AF_INSERT = "afinsert"
	XSAC_AF_UPDATE = "afupdate"
	XSAC_AF_DELETE = "afdelete"
)

type ZeroXsacTrigger interface {
	On(string, interface{}) error
}

type ZeroXsacAutoProcessor interface {
	Build(transaction *sql.Tx)

	DBName() string
	TableName() string
	AddTriggers(...ZeroXsacTrigger)

	Insert(...interface{}) error
	Update(...interface{}) error
	Delete(...interface{}) error
}
