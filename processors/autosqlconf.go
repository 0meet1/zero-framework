package processors

import (
	"database/sql"

	"github.com/0meet1/zero-framework/structs"
)

const (
	XSAC_BE_INSERT = "beinsert"
	XSAC_BE_UPDATE = "beupdate"
	XSAC_BE_DELETE = "bedelete"

	XSAC_AF_INSERT = "afinsert"
	XSAC_AF_UPDATE = "afupdate"
	XSAC_AF_DELETE = "afdelete"
)

type ZeroXsacAutoProcessor interface {
	Build(transaction *sql.Tx)

	AddFields(fields []*structs.ZeroXsacField)
	AddTriggers(...structs.ZeroXsacTrigger)

	Insert(...interface{}) error
	Update(...interface{}) error
	Delete(...interface{}) error
	Tombstone(...interface{}) error
	Xrestore(...interface{}) error
	FetchChildrens(*structs.ZeroXsacField, interface{}) error
}
