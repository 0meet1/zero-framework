package autohttpconf

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

const (
	XSAC_DML_ADD       = "add"
	XSAC_DML_UP        = "up"
	XSAC_DML_RM        = "rm"
	XSAC_DML_TOMBSTONE = "tombstone"
	XSAC_DML_RESTORE   = "restore"

	XSAC_HTTPFETCH_READY    = "ready"
	XSAC_HTTPFETCH_ROW      = "row"
	XSAC_HTTPFETCH_COMPLETE = "complete"
)

type ZeroXsacHttpDMLTrigger interface {
	On(string, string, *structs.ZeroRequest, ...interface{}) error
}

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

	XhttpDMLTrigger() ZeroXsacHttpDMLTrigger
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

func (e *ZeroXsacXhttpStructs) XhttpDMLTrigger() ZeroXsacHttpDMLTrigger       { return nil }
func (e *ZeroXsacXhttpStructs) XhttpFetchTrigger() ZeroXsacHttpFetchTrigger   { return nil }
func (e *ZeroXsacXhttpStructs) XhttpSearchTrigger() ZeroXsacHttpSearchTrigger { return nil }

func (e *ZeroXsacXhttpStructs) makeApiWriteReq() string {
	fields := e.This().(structs.ZeroXsacFields).XsacFields()
	mapdata := make(map[string]interface{})
	for _, field := range fields {
		apiitems := strings.Split(field.Xapi(), ",")
		if len(apiitems) < 2 {
			continue
		}
		jsonitems := strings.Split(field.Xjsonopts(), ",")
		if len(jsonitems) < 1 {
			continue
		}
		if field.Writable() {
			mapdata[jsonitems[0]] = fmt.Sprintf("*%s,%s*", apiitems[0], apiitems[1])
		}
	}
	reqdata := &structs.ZeroRequest{
		Querys:  []interface{}{mapdata},
		Expands: make(map[string]interface{}),
	}
	jsonbytes, _ := json.MarshalIndent(reqdata, "", "\t")
	return string(jsonbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiUpdateReq() string {
	fields := e.This().(structs.ZeroXsacFields).XsacFields()
	mapdata := make(map[string]interface{})
	mapdata["id"] = "*唯一标识,UUID*"
	for _, field := range fields {
		apiitems := strings.Split(field.Xapi(), ",")
		if len(apiitems) < 2 {
			continue
		}
		jsonitems := strings.Split(field.Xjsonopts(), ",")
		if len(jsonitems) < 1 {
			continue
		}
		if field.Updatable() {
			mapdata[jsonitems[0]] = fmt.Sprintf("*%s,%s*", apiitems[0], apiitems[1])
		}
	}
	reqdata := &structs.ZeroRequest{
		Querys:  []interface{}{mapdata},
		Expands: make(map[string]interface{}),
	}
	jsonbytes, _ := json.MarshalIndent(reqdata, "", "\t")
	return string(jsonbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiRemoveReq() string {
	mapdata := make(map[string]interface{})
	mapdata["id"] = "*唯一标识,UUID*"
	reqdata := &structs.ZeroRequest{
		Querys:  []interface{}{mapdata},
		Expands: make(map[string]interface{}),
	}
	jsonbytes, _ := json.MarshalIndent(reqdata, "", "\t")
	return string(jsonbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiFetchReq() string {
	xQuery := &processors.ZeroQuery{
		Orderby: []*processors.ZeroOrderBy{
			{
				Column: "createTime",
				Seq:    processors.ORDER_BY_DESC,
			},
		},
		Limit: &processors.ZeroLimit{
			Start:  0,
			Length: 10,
		},
	}
	reqdata := &structs.ZeroRequest{
		Querys:  []interface{}{xQuery},
		Expands: make(map[string]interface{}),
	}
	jsonbytes, _ := json.MarshalIndent(reqdata, "", "\t")
	return string(jsonbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiSearchReq() string {
	xQuery := &database.EQuerySearch{
		Sort: []interface{}{
			map[string]interface{}{
				"createTime": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		Size: 10,
		From: 0,
	}
	reqdata := &structs.ZeroRequest{
		Querys:  []interface{}{xQuery},
		Expands: make(map[string]interface{}),
	}
	jsonbytes, _ := json.MarshalIndent(reqdata, "", "\t")
	return string(jsonbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiSuccess() string {
	respmap := make(map[string]interface{})
	respmap["code"] = 200
	respmap["message"] = "success"
	respbytes, _ := json.MarshalIndent(respmap, "", "\t")
	return string(respbytes)
}

func (e *ZeroXsacXhttpStructs) makeApiQueryOptions() [][]string {
	options := make([][]string, 0)
	fields := e.This().(structs.ZeroXsacFields).XsacFields()
	for _, field := range fields {
		if field.Inlinable() || field.Childable() {
			apiitems := strings.Split(field.Xapi(), ",")
			if len(apiitems) < 2 {
				continue
			}
			jsonitems := strings.Split(field.Xjsonopts(), ",")
			if len(jsonitems) < 1 {
				continue
			}
			options = append(options, []string{jsonitems[0], apiitems[0]})
		}
	}
	return options
}

func (e *ZeroXsacXhttpStructs) makeApiQueryExpands() [][]string {
	expands := make([][]string, 0)
	if e.This().(ZeroXsacXhttpDeclares).XsacPartition() != structs.XSAC_PARTITION_NONE {
		expands = append(expands, []string{"zone", "时间区间"})
	}
	return expands
}

func (e *ZeroXsacXhttpStructs) makeApiDatas() string {
	fields := e.This().(structs.ZeroXsacFields).XsacFields()
	mapdata := make(map[string]interface{})
	for _, field := range fields {
		apiitems := strings.Split(field.Xapi(), ",")
		if len(apiitems) < 2 {
			continue
		}
		jsonitems := strings.Split(field.Xjsonopts(), ",")
		if len(jsonitems) < 1 {
			continue
		}
		mapdata[jsonitems[0]] = fmt.Sprintf("*%s,%s*", apiitems[0], apiitems[1])
	}

	respmap := make(map[string]interface{})
	respmap["code"] = 200
	respmap["message"] = "success"
	respmap["datas"] = []interface{}{mapdata}
	respmap["expand"] = map[string]interface{}{
		"total":  1,
		"start":  0,
		"length": 10,
	}
	respbytes, _ := json.MarshalIndent(respmap, "", "\t")
	return string(respbytes)
}

func (e *ZeroXsacXhttpStructs) XsacApis(args ...string) []string {
	prefix := "/"
	if len(args) > 0 {
		prefix = args[0]
	}
	xhttp := e.This().(ZeroXsacXhttpDeclares)
	xapidec := e.This().(structs.ZeroXsacApiDeclares)
	xsacdec := e.This().(structs.ZeroXsacDeclares)

	rows := make([]string, 0)

	if xhttp.XhttpOpt()&0b1000 == 0b1000 {
		rows = append(rows, structs.NewApiContentNOE(
			fmt.Sprintf("添加%s：%s/add", xapidec.XsacApiName(), xhttp.XhttpPath()),
			path.Join(prefix, xhttp.XhttpPath(), "add"), e.makeApiWriteReq(), e.makeApiSuccess())...)
	}

	if xhttp.XhttpOpt()&0b100 == 0b100 {
		rows = append(rows, structs.NewApiContentNOE(
			fmt.Sprintf("修改%s：%s/up", xapidec.XsacApiName(), xhttp.XhttpPath()),
			path.Join(prefix, xhttp.XhttpPath(), "up"), e.makeApiUpdateReq(), e.makeApiSuccess())...)
	}

	if xhttp.XhttpOpt()&0b10 == 0b10 {
		if xhttp.XsacDeleteOpt()&0b10000000 == 0b10000000 {
			rows = append(rows, structs.NewApiContentNOE(
				fmt.Sprintf("移除%s：%s/rm (物理删除)", xapidec.XsacApiName(), xhttp.XhttpPath()),
				path.Join(prefix, xhttp.XhttpPath(), "rm"), e.makeApiRemoveReq(), e.makeApiSuccess())...)
		} else {
			rows = append(rows, structs.NewApiContentNOE(
				fmt.Sprintf("移除%s：%s/rm (逻辑删除)", xapidec.XsacApiName(), xhttp.XhttpPath()),
				path.Join(prefix, xhttp.XhttpPath(), "rm"), e.makeApiRemoveReq(), e.makeApiSuccess())...)
			if xhttp.XsacDeleteOpt()&0b00000010 == 0b00000010 {
				rows = append(rows, structs.NewApiContentNOE(
					fmt.Sprintf("强制移除%s：%s/force (物理删除)", xapidec.XsacApiName(), xhttp.XhttpPath()),
					path.Join(prefix, xhttp.XhttpPath(), "force"), e.makeApiRemoveReq(), e.makeApiSuccess())...)
			}
			if xhttp.XsacDeleteOpt()&0b00000100 == 0b00000100 {
				rows = append(rows, structs.NewApiContentNOE(
					fmt.Sprintf("恢复%s：%s/restore", xapidec.XsacApiName(), xhttp.XhttpPath()),
					path.Join(prefix, xhttp.XhttpPath(), "restore"), e.makeApiRemoveReq(), e.makeApiSuccess())...)
			}
		}
	}

	options := e.makeApiQueryOptions()
	expands := e.makeApiQueryExpands()
	if e.This().(ZeroXsacXhttpDeclares).XhttpOpt()&0b1 == 0b1 {
		rows = append(rows, structs.NewApiContent(
			fmt.Sprintf("查询%s：%s/fetch", xapidec.XsacApiName(), xhttp.XhttpPath()),
			path.Join(prefix, xhttp.XhttpPath(), "fetch"),
			e.makeApiFetchReq(), e.makeApiDatas(), options, expands)...)

		if xsacdec.XsacDeleteOpt()&0b00000001 == 0b00000001 {
			rows = append(rows, structs.NewApiContent(
				fmt.Sprintf("查询%s回收站：%s/history", xapidec.XsacApiName(), xhttp.XhttpPath()),
				path.Join(prefix, xhttp.XhttpPath(), "history"),
				e.makeApiFetchReq(), e.makeApiDatas(), options, expands)...)
		}
	}

	if xhttp.XhttpOpt()&0b10000 == 0b10000 {
		rows = append(rows, structs.NewApiContentNOE(
			fmt.Sprintf("搜索%s：%s/search", xapidec.XsacApiName(), xhttp.XhttpPath()),
			path.Join(prefix, xhttp.XhttpPath(), "search"),
			e.makeApiSearchReq(), e.makeApiDatas())...)
	}
	return rows
}

type ZeroXsacXhttp struct {
	dataSource string
	dbName     string
	coretype   reflect.Type

	fields structs.ZeroXsacFieldSet

	instance ZeroXsacXhttpDeclares

	inlinfields map[string]*structs.ZeroXsacField
}

func NewXsacXhttp(coretype reflect.Type) *ZeroXsacXhttp {
	xhttpDec := reflect.New(coretype).Interface().(ZeroXsacXhttpDeclares)
	xhttpDec.(structs.ZeroMetaDef).ThisDef(xhttpDec)
	return &ZeroXsacXhttp{
		coretype:   coretype,
		dataSource: "",
		dbName:     "",
		instance:   xhttpDec,
	}
}

func (xhttp *ZeroXsacXhttp) AddDataSource(dataSource string) *ZeroXsacXhttp {
	xhttp.dataSource = dataSource
	return xhttp
}

func (xhttp *ZeroXsacXhttp) AddDbName(dbName string) *ZeroXsacXhttp {
	xhttp.dbName = dbName
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

func (xhttp *ZeroXsacXhttp) xhttpinlines() map[string]*structs.ZeroXsacField {
	if xhttp.inlinfields == nil {
		xhttp.inlinfields = make(map[string]*structs.ZeroXsacField)
		for _, field := range xhttp.xhttpfields() {
			if field.Inlinable() {
				xhttp.inlinfields[strings.ToLower(field.FieldName())] = field
			}
		}
		for _, field := range xhttp.xhttpfields() {
			if field.Childable() {
				xhttp.inlinfields[strings.ToLower(field.FieldName())] = field
			}
		}
	}
	return xhttp.inlinfields
}

func (xhttp *ZeroXsacXhttp) xhttpProcessor(transaction *sql.Tx) processors.ZeroXsacAutoProcessor {
	processor := xhttp.instance.XhttpAutoProc()
	processor.AddFields(xhttp.xhttpfields())
	processor.Build(transaction)
	return processor
}

func (xhttp *ZeroXsacXhttp) xhttpParse(querys []interface{}, callback func(interface{}) error) error {
	for _, xQueryData := range querys {
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

		err = callback(xquery)
		if err != nil {
			return err
		}
	}
	return nil
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

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_ADD, XSAC_HTTPFETCH_READY, xRequest)
		if err != nil {
			panic(err)
		}
	}

	processor := xhttp.xhttpProcessor(transaction)

	err = xhttp.xhttpParse(xRequest.Querys, func(xquery interface{}) error {
		if xhttp.instance.XhttpDMLTrigger() != nil {
			err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_ADD, XSAC_HTTPFETCH_ROW, xRequest, xquery)
			if err != nil {
				return err
			}
		}
		err = processor.Insert(xquery)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_ADD, XSAC_HTTPFETCH_COMPLETE, xRequest)
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

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_UP, XSAC_HTTPFETCH_READY, xRequest)
		if err != nil {
			panic(err)
		}
	}

	processor := xhttp.xhttpProcessor(transaction)

	err = xhttp.xhttpParse(xRequest.Querys, func(xquery interface{}) error {
		if xhttp.instance.XhttpDMLTrigger() != nil {
			err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_UP, XSAC_HTTPFETCH_ROW, xRequest, xquery)
			if err != nil {
				return err
			}
		}

		err = processor.Update(xquery)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_UP, XSAC_HTTPFETCH_COMPLETE, xRequest)
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

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RM, XSAC_HTTPFETCH_READY, xRequest)
		if err != nil {
			panic(err)
		}
	}

	processor := xhttp.xhttpProcessor(transaction)

	err = xhttp.xhttpParse(xRequest.Querys, func(xquery interface{}) error {
		if xhttp.instance.XhttpDMLTrigger() != nil {
			err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RM, XSAC_HTTPFETCH_ROW, xRequest, xquery)
			if err != nil {
				return err
			}
		}

		err = processor.Delete(xquery)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RM, XSAC_HTTPFETCH_COMPLETE, xRequest)
		if err != nil {
			panic(err)
		}
	}

	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) tombstone(writer http.ResponseWriter, req *http.Request) {
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

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_TOMBSTONE, XSAC_HTTPFETCH_READY, xRequest)
		if err != nil {
			panic(err)
		}
	}

	processor := xhttp.xhttpProcessor(transaction)

	err = xhttp.xhttpParse(xRequest.Querys, func(xquery interface{}) error {
		if xhttp.instance.XhttpDMLTrigger() != nil {
			err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_TOMBSTONE, XSAC_HTTPFETCH_ROW, xRequest, xquery)
			if err != nil {
				return err
			}
		}

		err = processor.Tombstone(xquery)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_TOMBSTONE, XSAC_HTTPFETCH_COMPLETE, xRequest)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) restore(writer http.ResponseWriter, req *http.Request) {
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

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RESTORE, XSAC_HTTPFETCH_READY, xRequest)
		if err != nil {
			panic(err)
		}
	}

	processor := xhttp.xhttpProcessor(transaction)

	err = xhttp.xhttpParse(xRequest.Querys, func(xquery interface{}) error {
		if xhttp.instance.XhttpDMLTrigger() != nil {
			err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RESTORE, XSAC_HTTPFETCH_ROW, xRequest, xquery)
			if err != nil {
				return err
			}
		}

		err = processor.Xrestore(xquery)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if xhttp.instance.XhttpDMLTrigger() != nil {
		err = xhttp.instance.XhttpDMLTrigger().On(XSAC_DML_RESTORE, XSAC_HTTPFETCH_COMPLETE, xRequest)
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMessages(writer, 200, "success")
}

func (xhttp *ZeroXsacXhttp) checkpart(xRequest *structs.ZeroRequest, xOperation processors.ZeroQueryOperation) {
	if xhttp.instance.XsacPartition() != structs.XSAC_PARTITION_NONE {
		if xRequest.Expands == nil {
			panic("missing necessary parameter `expands.zone`")
		}

		zone, ok := xRequest.Expands["zone"]
		if !ok {
			panic("missing necessary parameter `expands.zone`")
		}

		date, err := time.Parse("2006-01-02", zone.(string))
		if err != nil {
			panic(err)
		}

		switch xhttp.instance.XsacPartition() {
		case structs.XSAC_PARTITION_YEAR:
			startTime, endTime, err := structs.YearDurationString(date, "2006-01-02 15:04:05")
			if err != nil {
				panic(err)
			}
			xOperation.AppendCondition(fmt.Sprintf("create_time BETWEEN '%s' AND '%s'", startTime, endTime))
		case structs.XSAC_PARTITION_MONTH:
			startTime, endTime, err := structs.MonthDurationString(date, "2006-01-02 15:04:05")
			if err != nil {
				panic(err)
			}
			xOperation.AppendCondition(fmt.Sprintf("create_time BETWEEN '%s' AND '%s'", startTime, endTime))
		case structs.XSAC_PARTITION_DAY:
			startTime, endTime, err := structs.DayDurationString(date, "2006-01-02 15:04:05")
			if err != nil {
				panic(err)
			}
			xOperation.AppendCondition(fmt.Sprintf("create_time BETWEEN '%s' AND '%s'", startTime, endTime))
		}
	}
}

func (xhttp *ZeroXsacXhttp) parserowdata(xoptions []string, processor processors.ZeroXsacAutoProcessor, row map[string]interface{}) interface{} {
	data := reflect.New(xhttp.coretype).Interface()
	returnValues := reflect.ValueOf(data).MethodByName("LoadRowData").Call([]reflect.Value{reflect.ValueOf(row)})
	if len(returnValues) > 0 && returnValues[0].Interface() != nil {
		panic(returnValues[0].Interface())
	}

	for _, xoption := range xoptions {
		if xoption == server.XHTTP_QUERY_OPTIONS_ALL {
			for _, field := range xhttp.xhttpinlines() {
				processor.FetchChildrens(field, data)
			}
		} else {
			field, ok := xhttp.xhttpinlines()[xoption]
			if ok {
				processor.FetchChildrens(field, data)
			}
		}
	}
	return data
}

func (xhttp *ZeroXsacXhttp) corefetch(writer http.ResponseWriter, req *http.Request, flag int) {
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

	if xRequest.Querys == nil || len(xRequest.Querys) <= 0 {
		panic("missing necessary parameter `query options -> $.querys[0]`")
	}

	xOperation, _, err := server.XhttpCompleteQueryOperation(xRequest, xhttp.instance.XhttpQueryOperation(), xhttp.XcheckTable())
	if err != nil {
		panic(err)
	}
	xOperation.Build(transaction)
	if flag >= 0 {
		xOperation.AppendCondition(fmt.Sprintf("flag = %d", flag))
	}
	xhttp.checkpart(xRequest, xOperation)

	processor := xhttp.xhttpProcessor(transaction)

	if xhttp.instance.XhttpFetchTrigger() != nil {
		err = xhttp.instance.XhttpFetchTrigger().On(XSAC_HTTPFETCH_READY, xOperation, xRequest)
		if err != nil {
			panic(err)
		}
	}

	xoptions := server.XhttpQueryOptions(xRequest)

	rows, expands := xOperation.Exec()
	datas := make([]interface{}, 0)
	for _, row := range rows {
		data := xhttp.parserowdata(xoptions, processor, row)

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

func (xhttp *ZeroXsacXhttp) fetch(writer http.ResponseWriter, req *http.Request) {
	xhttp.corefetch(writer, req, -1)
}

func (xhttp *ZeroXsacXhttp) tombfetch(writer http.ResponseWriter, req *http.Request) {
	xhttp.corefetch(writer, req, 0)
}

func (xhttp *ZeroXsacXhttp) history(writer http.ResponseWriter, req *http.Request) {
	xhttp.corefetch(writer, req, 1)
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
		if xhttp.instance.(structs.ZeroXsacDeclares).XsacDeleteOpt()&0b10000000 == 0b10000000 {
			xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.rm, fmt.Sprintf("%s/rm", xhttp.XhttpPath())))
		} else {
			xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.tombstone, fmt.Sprintf("%s/rm", xhttp.XhttpPath())))
			if xhttp.instance.(structs.ZeroXsacDeclares).XsacDeleteOpt()&0b00000010 == 0b00000010 {
				xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.rm, fmt.Sprintf("%s/force", xhttp.XhttpPath())))
			}
			if xhttp.instance.(structs.ZeroXsacDeclares).XsacDeleteOpt()&0b00000100 == 0b00000100 {
				xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.restore, fmt.Sprintf("%s/restore", xhttp.XhttpPath())))
			}
		}
	}

	if xhttp.instance.XhttpOpt()&0b1 == 0b1 {
		if xhttp.instance.(structs.ZeroXsacDeclares).XsacDeleteOpt()&0b10000000 == 0b10000000 {
			xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.fetch, fmt.Sprintf("%s/fetch", xhttp.XhttpPath())))
		} else {
			xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.tombfetch, fmt.Sprintf("%s/fetch", xhttp.XhttpPath())))
			if xhttp.instance.(structs.ZeroXsacDeclares).XsacDeleteOpt()&0b00000001 == 0b00000001 {
				xExecutors = append(xExecutors, server.XhttpFuncHandle(xhttp.history, fmt.Sprintf("%s/history", xhttp.XhttpPath())))
			}
		}
	}
	return xExecutors
}
