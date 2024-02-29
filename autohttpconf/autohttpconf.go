package autohttpconf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacHttpFetchTrigger interface {
	On(string, processors.ZeroQueryOperation, *structs.ZeroRequest, ...interface{}) error
}

type ZeroXsacHttpSearchTrigger interface {
	On(string, *database.EQueryRequest, *structs.ZeroRequest, ...interface{}) error
}

type ZeroXsacXhttpDeclares interface {
	structs.ZeroXsacDeclares

	XhttpPath() string
	XhttpAutoProc() processors.ZeroXsacAutoProcessor
	XhttpQueryOperation() processors.ZeroQueryOperation
	XhttpOpt() byte

	XhttpCheckTable() string
	XhttpSearchIndex() string

	XhttpFetchTrigger() ZeroXsacHttpFetchTrigger
	XhttpSearchTrigger() ZeroXsacHttpSearchTrigger
}

type ZeroXsacXhttpStructs struct {
	structs.ZeroCoreStructs
}

func (e *ZeroXsacXhttpStructs) XhttpPath() string        { return "" }
func (e *ZeroXsacXhttpStructs) XhttpOpt() byte           { return 0b00001111 }
func (e *ZeroXsacXhttpStructs) XhttpCheckTable() string  { return "" }
func (e *ZeroXsacXhttpStructs) XhttpSearchIndex() string { return "" }

func (e *ZeroXsacXhttpStructs) XhttpAutoProc() processors.ZeroXsacAutoProcessor {
	return processors.NewXsacPostgresProcessor()
}

func (e *ZeroXsacXhttpStructs) XhttpQueryOperation() processors.ZeroQueryOperation {
	return &processors.ZeroPostgresQueryOperation{}
}

func (e *ZeroXsacXhttpStructs) XhttpFetchTrigger() ZeroXsacHttpFetchTrigger   { return nil }
func (e *ZeroXsacXhttpStructs) XhttpSearchTrigger() ZeroXsacHttpSearchTrigger { return nil }

const (
	XSAC_HTTPFETCH_READY    = "ready"
	XSAC_HTTPFETCH_ROW      = "row"
	XSAC_HTTPFETCH_COMPLETE = "complete"
)

type ZeroXsacXhttp struct {
	dataSource string
	coretype   reflect.Type

	fields structs.ZeroXsacFieldSet

	instance ZeroXsacXhttpDeclares
}

func NewXsacXhttp(coretype reflect.Type) *ZeroXsacXhttp {
	xhttpDec := reflect.New(coretype).Interface().(ZeroXsacXhttpDeclares)
	xhttpDec.(structs.ZeroMetaDef).ThisDef(xhttpDec)
	return &ZeroXsacXhttp{
		coretype:   coretype,
		dataSource: "",
		instance:   xhttpDec,
	}
}

func (xhttp *ZeroXsacXhttp) AddDataSource(dataSource string) *ZeroXsacXhttp {
	xhttp.dataSource = dataSource
	return xhttp
}

func (xhttp *ZeroXsacXhttp) XDataSource() string {
	if len(xhttp.instance.XsacDataSource()) > 0 {
		return xhttp.instance.XsacDataSource()
	}
	return xhttp.dataSource
}

func (xhttp *ZeroXsacXhttp) XhttpPath() string  { return xhttp.instance.XhttpPath() }
func (xhttp *ZeroXsacXhttp) XdbName() string    { return xhttp.instance.XsacDbName() }
func (xhttp *ZeroXsacXhttp) XtableName() string { return xhttp.instance.XsacTableName() }

func (xhttp *ZeroXsacXhttp) XcheckTable() string {
	if len(xhttp.instance.XhttpCheckTable()) == 0 {
		return xhttp.XtableName()
	}
	return xhttp.instance.XhttpCheckTable()
}

func (xhttp *ZeroXsacXhttp) XsearchIndex() string {
	if len(xhttp.instance.XhttpSearchIndex()) == 0 {
		return xhttp.XtableName()
	}
	return xhttp.instance.XhttpSearchIndex()
}

func (xhttp *ZeroXsacXhttp) xhttpfields() structs.ZeroXsacFieldSet {
	if xhttp.fields == nil {
		xhttp.fields = xhttp.instance.(structs.ZeroXsacFields).XsacFields()
	}
	return xhttp.fields
}

func (xhttp *ZeroXsacXhttp) add(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttp.XDataSource()).(*database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, fmt.Sprintf("%s", err))
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	processor := xhttp.instance.XhttpAutoProc()
	processor.AddFields(xhttp.xhttpfields())
	processor.Build(transaction)

	for _, xQueryData := range xRequest.Querys {
		jsonbytes, err := json.Marshal(xQueryData)
		if err != nil {
			panic(err)
		}

		xquery := reflect.New(xhttp.coretype).Interface()
		err = json.Unmarshal(jsonbytes, xquery)
		if err != nil {
			panic(err)
		}
		xquery.(structs.ZeroMetaDef).ThisDef(xquery)
		err = processor.Insert(xquery)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) up(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttp.XDataSource()).(*database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, fmt.Sprintf("%s", err))
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	processor := xhttp.instance.XhttpAutoProc()
	processor.AddFields(xhttp.xhttpfields())
	processor.Build(transaction)

	for _, xQueryData := range xRequest.Querys {
		jsonbytes, err := json.Marshal(xQueryData)
		if err != nil {
			panic(err)
		}

		xquery := reflect.New(xhttp.coretype).Interface()
		err = json.Unmarshal(jsonbytes, xquery)
		if err != nil {
			panic(err)
		}
		xquery.(structs.ZeroMetaDef).ThisDef(xquery)
		err = processor.Update(xquery)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) rm(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttp.XDataSource()).(*database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, fmt.Sprintf("%s", err))
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	processor := xhttp.instance.XhttpAutoProc()
	processor.AddFields(xhttp.xhttpfields())
	processor.Build(transaction)

	for _, xQueryData := range xRequest.Querys {
		jsonbytes, err := json.Marshal(xQueryData)
		if err != nil {
			panic(err)
		}

		xquery := reflect.New(xhttp.coretype).Interface()
		err = json.Unmarshal(jsonbytes, xquery)
		if err != nil {
			panic(err)
		}
		xquery.(structs.ZeroMetaDef).ThisDef(xquery)
		err = processor.Delete(xquery)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) fetch(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttp.XDataSource()).(*database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, fmt.Sprintf("%s", err))
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	xOperation, _, err := server.XhttpCompleteQueryOperation(xRequest, xhttp.instance.XhttpQueryOperation(), xhttp.XcheckTable())
	if err != nil {
		panic(err)
	}
	xOperation.Build(transaction)

	if xhttp.instance.XhttpFetchTrigger() != nil {
		err = xhttp.instance.XhttpFetchTrigger().On(XSAC_HTTPFETCH_READY, xOperation, xRequest)
		if err != nil {
			panic(err)
		}
	}

	rows, expands := xOperation.Exec()
	datas := make([]interface{}, 0)
	for _, row := range rows {
		data := reflect.New(xhttp.coretype).Interface()
		returnValues := reflect.ValueOf(data).MethodByName("LoadRowData").Call([]reflect.Value{reflect.ValueOf(row)})
		if len(returnValues) > 0 && returnValues[0].Interface() != nil {
			panic(returnValues[0].Interface())
		}

		if xhttp.instance.XhttpFetchTrigger() != nil {
			err = xhttp.instance.XhttpFetchTrigger().On(XSAC_HTTPFETCH_ROW, xOperation, xRequest, data)
			if err != nil {
				panic(err)
			}
		}
		datas = append(datas, data)
	}

	if xhttp.instance.XhttpFetchTrigger() != nil {
		err = xhttp.instance.XhttpFetchTrigger().On(XSAC_HTTPFETCH_COMPLETE, xOperation, xRequest, datas...)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseDatas(writer, 200, "success", datas, expands)
}

func (xhttp *ZeroXsacXhttp) search(writer http.ResponseWriter, req *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			writer.WriteHeader(500)
			server.XhttpResponseMessages(writer, 500, fmt.Sprintf("%s", err))
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	eRequest, xEQuery, err := server.XhttpEQueryRequest(xRequest, xhttp.XsearchIndex())
	if err != nil {
		panic(err)
	}
	xEQuery.TrackTotalHits = 100000000

	if xhttp.instance.XhttpSearchTrigger() != nil {
		xhttp.instance.XhttpSearchTrigger().On(XSAC_HTTPFETCH_READY, eRequest, xRequest)
	}

	resp, err := eRequest.Search()
	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpSearchTrigger() != nil {
		xhttp.instance.XhttpSearchTrigger().On(XSAC_HTTPFETCH_COMPLETE, eRequest, xRequest, resp.Datas...)
	}

	expands := make(map[string]interface{})
	expands["from"] = xEQuery.From
	expands["size"] = xEQuery.Size
	expands["total"] = resp.Total

	server.XhttpResponseDatas(writer, 200, "success", resp.Datas, expands)
}

func (xhttp *ZeroXsacXhttp) ExportExecutors() []*server.XhttpExecutor {
	xExecutors := make([]*server.XhttpExecutor, 0)
	if xhttp.instance.XhttpOpt()&0b10000 == 0b10000 {
		xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.search, fmt.Sprintf("%s/search", xhttp.XhttpPath())))
	}

	if xhttp.instance.XhttpOpt()&0b1000 == 0b1000 {
		xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.add, fmt.Sprintf("%s/add", xhttp.XhttpPath())))
	}

	if xhttp.instance.XhttpOpt()&0b100 == 0b100 {
		xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.up, fmt.Sprintf("%s/up", xhttp.XhttpPath())))
	}

	if xhttp.instance.XhttpOpt()&0b10 == 0b10 {
		xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.rm, fmt.Sprintf("%s/rm", xhttp.XhttpPath())))
	}

	if xhttp.instance.XhttpOpt()&0b1 == 0b1 {
		xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.fetch, fmt.Sprintf("%s/fetch", xhttp.XhttpPath())))
	}
	return xExecutors
}
