package autohttpconf

import (
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
	On(string, *processors.ZeroQueryOperation, *structs.ZeroRequest, ...interface{}) error
}

type ZeroXsacHttpSearchTrigger interface {
	On(string, *database.EQueryRequest, *structs.ZeroRequest, ...interface{}) error
}

type ZeroXsacXhttpDeclares interface {
	structs.ZeroXsacDeclares

	XhttpPath() string
	XhttpAutoProc() reflect.Type
	XhttpOpt() byte
	XhttpDataSource() string

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
func (e *ZeroXsacXhttpStructs) XhttpDataSource() string  { return "" }
func (e *ZeroXsacXhttpStructs) XhttpCheckTable() string  { return "" }
func (e *ZeroXsacXhttpStructs) XhttpSearchIndex() string { return "" }

func (e *ZeroXsacXhttpStructs) XhttpAutoProc() reflect.Type {
	return reflect.TypeOf(&processors.ZeroXsacPostgresAutoProcessor{})
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

	instance ZeroXsacXhttpDeclares
}

func NewXsacXhttp(coretype reflect.Type) *ZeroXsacXhttp {
	return &ZeroXsacXhttp{
		coretype:   coretype,
		dataSource: "",
		instance:   reflect.New(coretype.Elem()).Interface().(ZeroXsacXhttpDeclares),
	}
}

func (xhttp *ZeroXsacXhttp) AddDataSource(dataSource string) *ZeroXsacXhttp {
	xhttp.dataSource = dataSource
	return xhttp
}

func (xhttp *ZeroXsacXhttp) XDataSource() string {
	if len(xhttp.instance.XhttpDataSource()) > 0 {
		return xhttp.instance.XhttpDataSource()
	}
	return xhttp.dataSource
}

func (xhttp *ZeroXsacXhttp) XhttpPath() string {
	return xhttp.instance.XhttpPath()
}

func (xhttp *ZeroXsacXhttp) XdbName() string {
	return xhttp.instance.XsacDbName()
}

func (xhttp *ZeroXsacXhttp) XtableName() string {
	return xhttp.instance.XsacTableName()
}

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

	processor := reflect.New(xhttp.instance.XhttpAutoProc().Elem()).Interface().(processors.ZeroXsacAutoProcessor)
	processor.Build(transaction)
	err = processor.Insert(xRequest.Querys...)
	if err != nil {
		panic(err)
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

	processor := reflect.New(xhttp.instance.XhttpAutoProc().Elem()).Interface().(processors.ZeroXsacAutoProcessor)
	processor.Build(transaction)
	err = processor.Update(xRequest.Querys...)
	if err != nil {
		panic(err)
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

	processor := reflect.New(xhttp.instance.XhttpAutoProc().Elem()).Interface().(processors.ZeroXsacAutoProcessor)
	processor.Build(transaction)
	err = processor.Delete(xRequest.Querys...)
	if err != nil {
		panic(err)
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

	xOperation, _, err := server.XhttpZeroQueryOperation(xRequest, xhttp.XcheckTable())
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
		data := reflect.New(xhttp.coretype.Elem()).Interface()
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
		server.XhttpFuncHandle(xhttp.search, fmt.Sprintf("%s/search", xhttp.XhttpPath()))
	}

	if xhttp.instance.XhttpOpt()&0b1000 == 0b1000 {
		server.XhttpFuncHandle(xhttp.add, fmt.Sprintf("%s/add", xhttp.XhttpPath()))
	}

	if xhttp.instance.XhttpOpt()&0b100 == 0b100 {
		server.XhttpFuncHandle(xhttp.up, fmt.Sprintf("%s/up", xhttp.XhttpPath()))
	}

	if xhttp.instance.XhttpOpt()&0b10 == 0b10 {
		server.XhttpFuncHandle(xhttp.rm, fmt.Sprintf("%s/rm", xhttp.XhttpPath()))
	}

	if xhttp.instance.XhttpOpt()&0b1 == 0b1 {
		server.XhttpFuncHandle(xhttp.fetch, fmt.Sprintf("%s/fetch", xhttp.XhttpPath()))
	}
	return xExecutors
}
