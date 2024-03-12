package autohttpconf

import (
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacXhttpApi struct {
	structs.ZeroMeta
}

func (_ *ZeroXsacXhttpApi) XsacPrimaryType() string                            { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacDataSource() string                             { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacDbName() string                                 { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacTableName() string                              { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacDeleteOpt() byte                                { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacPartition() string                              { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacCustomPartTrigger() string                      { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacTriggers() []structs.ZeroXsacTrigger            { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacApiName() string                                { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacApiEnums() []string                             { panic("not support") }
func (_ *ZeroXsacXhttpApi) XsacApi(...string) []string                         { return make([]string, 0) }
func (_ *ZeroXsacXhttpApi) XhttpPath() string                                  { return "" }
func (_ *ZeroXsacXhttpApi) XhttpAutoProc() processors.ZeroXsacAutoProcessor    { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpQueryOperation() processors.ZeroQueryOperation { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpOpt() byte                                     { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpCheckTable() string                            { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpSearchIndex() string                           { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpCustomPartChecker() ZeroXsacCustomPartChecker  { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpDMLTrigger() ZeroXsacHttpDMLTrigger            { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpFetchTrigger() ZeroXsacHttpFetchTrigger        { panic("not support") }
func (_ *ZeroXsacXhttpApi) XhttpSearchTrigger() ZeroXsacHttpSearchTrigger      { panic("not support") }

func (_ *ZeroXsacXhttpApi) XsacDeclares(...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}
func (_ *ZeroXsacXhttpApi) XsacRefDeclares(...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}
func (_ *ZeroXsacXhttpApi) XsacApiExports(...string) []string { return make([]string, 0) }
