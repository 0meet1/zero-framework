package autohttpconf

import (
	"errors"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacXhttpApi struct {
	structs.ZeroMeta
}

func (*ZeroXsacXhttpApi) XsacPrimaryType() string                 { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacDataSource() string                  { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacDbName() string                      { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacTableName() string                   { return "" }
func (*ZeroXsacXhttpApi) XsacDeleteOpt() byte                     { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacPartition() string                   { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacCustomPartTrigger() string           { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacTriggers() []structs.ZeroXsacTrigger { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacApiName() string                     { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacApiFields() [][]string               { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacApiEnums() []string                  { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XsacApis(...string) []string             { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XhttpPath() string                       { return "" }
func (*ZeroXsacXhttpApi) XhttpAutoProc() processors.ZeroXsacAutoProcessor {
	panic(errors.New("not support"))
}
func (*ZeroXsacXhttpApi) XhttpQueryOperation() processors.ZeroQueryOperation {
	panic(errors.New("not support"))
}
func (*ZeroXsacXhttpApi) XhttpOpt() byte           { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XhttpCheckTable() string  { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XhttpSearchIndex() string { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XhttpCustomPartChecker() ZeroXsacCustomPartChecker {
	panic(errors.New("not support"))
}
func (*ZeroXsacXhttpApi) XhttpDMLTrigger() ZeroXsacHttpDMLTrigger { panic(errors.New("not support")) }
func (*ZeroXsacXhttpApi) XhttpFetchTrigger() ZeroXsacHttpFetchTrigger {
	panic(errors.New("not support"))
}
func (*ZeroXsacXhttpApi) XhttpSearchTrigger() ZeroXsacHttpSearchTrigger {
	panic(errors.New("not support"))
}
func (*ZeroXsacXhttpApi) XsacAutoParser() []structs.ZeroXsacAutoParser { return nil }
func (*ZeroXsacXhttpApi) XsacDeclares(...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}
func (*ZeroXsacXhttpApi) XsacRefDeclares(...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}
func (*ZeroXsacXhttpApi) XsacAdjunctDeclares(...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}
func (*ZeroXsacXhttpApi) XsacApiExports(...string) []string { return make([]string, 0) }

var XautoProcessor = func(declare ZeroXsacXhttpDeclares) processors.ZeroXsacAutoProcessor {
	declare.ThisDef(declare)
	_ds := declare.XsacDataSource()
	if _ds == "" {
		_ds = database.DATABASE_POSTGRES
	}
	processor := declare.XhttpAutoProc()
	processor.AddFields(declare.(structs.ZeroXsacFields).XsacFields())
	transaction := global.Value(_ds).(database.DataSource).Transaction()
	processor.Build(transaction)
	return processor
}
