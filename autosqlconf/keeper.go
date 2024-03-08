package autosqlconf

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/0meet1/zero-framework/autohttpconf"
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacKeeper struct {
	proctype   reflect.Type
	types      []reflect.Type
	dataSource string

	entries   []*structs.ZeroXsacEntry
	httpconfs []*autohttpconf.ZeroXsacXhttp

	apimks string
}

func NewKeeper(proctype reflect.Type, types ...reflect.Type) *ZeroXsacKeeper {
	keeper := &ZeroXsacKeeper{
		proctype:  structs.FindMetaType(proctype),
		types:     make([]reflect.Type, 0),
		entries:   make([]*structs.ZeroXsacEntry, 0),
		httpconfs: make([]*autohttpconf.ZeroXsacXhttp, 0),
	}
	keeper.AddTypes(types...)
	return keeper
}

func (keeper *ZeroXsacKeeper) DataSource(dataSource string) *ZeroXsacKeeper {
	keeper.dataSource = dataSource
	return keeper
}

func (keeper *ZeroXsacKeeper) AddTypes(types ...reflect.Type) *ZeroXsacKeeper {
	if types != nil {
		for _, t := range types {
			keeper.types = append(keeper.types, structs.FindMetaType(t))
		}
	}

	return keeper
}

func (keeper *ZeroXsacKeeper) DMLTables() {
	datasource := global.Value(keeper.dataSource)
	maxtrytimes := 10
	for datasource == nil && maxtrytimes > 0 {
		global.Logger().Warn(fmt.Sprintf("data source is not ready, try after 5s ..."))
		<-time.After(time.Duration(5) * time.Second)
		datasource = global.Value(keeper.dataSource)
		maxtrytimes--
	}

	if datasource == nil {
		global.Logger().Error(fmt.Sprintf("data source is not ready, give up"))
		return
	}

	transaction := global.Value(keeper.dataSource).(*database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("dml tables failed: %s", err))
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	xsacProcessor := reflect.New(keeper.proctype).Interface().(ZeroXsacProcessor)
	xsacProcessor.Build(transaction)

	for _, entry := range keeper.entries {
		switch entry.EntryType() {
		case structs.ZEOR_XSAC_ENTRY_TYPE_TABLE:
			err := xsacProcessor.DMLTable(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_TABLE0S:
			err := xsacProcessor.Create0Struct(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_TABLE0FS:
			err := xsacProcessor.Create0FlagStruct(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_COLUMN:
			err := xsacProcessor.DMLColumn(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2], entry.EntryParams()[3], entry.EntryParams()[4], entry.EntryParams()[5])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DROPCOLUMN:
			err := xsacProcessor.DropColumn(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_KEY:
			err := xsacProcessor.DMLIndex(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DROPKEY:
			err := xsacProcessor.DropIndex(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_PRIMARY_KEY:
			err := xsacProcessor.DMLPrimary(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DROP_PRIMARY_KEY:
			err := xsacProcessor.DropPrimary(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_UNIQUE_KEY:
			err := xsacProcessor.DMLUnique(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DROP_UNIQUE_KEY:
			err := xsacProcessor.DropUnique(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_FOREIGN_KEY:
			err := xsacProcessor.DMLForeign(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2], entry.EntryParams()[3], entry.EntryParams()[4])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DROP_FOREIGN_KEY:
			err := xsacProcessor.DropForeign(entry.EntryParams()[0], entry.EntryParams()[1], entry.EntryParams()[2])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_YEAR_PARTITION:
			err := xsacProcessor.DMLY0SPart(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_MONTH_PARTITION:
			err := xsacProcessor.DMLM0SPart(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		case structs.ZEOR_XSAC_ENTRY_TYPE_DAY_PARTITION:
			err := xsacProcessor.DMLD0SPart(entry.EntryParams()[0], entry.EntryParams()[1])
			if err != nil {
				panic(err)
			}
		default:
			panic(errors.New(fmt.Sprintf("unknown entry type: %s", entry.EntryType())))
		}
	}
}

func (keeper *ZeroXsacKeeper) pretreat() {
	xsacProcessor := reflect.New(keeper.proctype).Interface().(ZeroXsacProcessor)
	refDeclares := make([]*structs.ZeroXsacEntry, 0)
	for _, t := range keeper.types {
		declares := reflect.New(t).Interface().(structs.ZeroXsacDeclares)
		reflect.ValueOf(declares).MethodByName("ThisDef").Call([]reflect.Value{reflect.ValueOf(declares)})

		keeper.entries = append(keeper.entries, declares.XsacDeclares(xsacProcessor.DbName())...)
		refDeclares = append(refDeclares, declares.XsacRefDeclares(xsacProcessor.DbName())...)
		if reflect.TypeOf(declares).Implements(reflect.TypeOf((*autohttpconf.ZeroXsacXhttpDeclares)(nil)).Elem()) {
			if len(declares.(autohttpconf.ZeroXsacXhttpDeclares).XhttpPath()) > 0 {
				keeper.httpconfs = append(keeper.httpconfs, autohttpconf.NewXsacXhttp(t).AddDataSource(keeper.dataSource).AddDbName(xsacProcessor.DbName()))
			}
		}
	}
	keeper.entries = append(keeper.entries, refDeclares...)
}

func (keeper *ZeroXsacKeeper) RunKeeper() *ZeroXsacKeeper {
	keeper.pretreat()
	if global.StringValue("zero.xsac.autocheck") == "enable" {
		go keeper.DMLTables()
	}
	return keeper
}

func (keeper *ZeroXsacKeeper) api(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte(keeper.apimks))
}

func (keeper *ZeroXsacKeeper) Exports() []*server.XhttpExecutor {
	executors := make([]*server.XhttpExecutor, 0)
	if global.StringValue("zero.xsac.autoapi") == "enable" {
		appname := global.StringValue("zero.appname")
		version := global.StringValue("zero.version")
		prefix := global.StringValue("zero.httpserver.prefix")

		rows := make([]string, 0)
		rows = append(rows, structs.NewApiHeader(appname, version))

		for i, t := range keeper.types {
			declares := reflect.New(t).Interface().(structs.ZeroXsacDeclares)
			reflect.ValueOf(declares).MethodByName("ThisDef").Call([]reflect.Value{reflect.ValueOf(declares)})
			rows = append(rows, declares.(structs.ZeroXsacApiDeclares).XsacApiExports(fmt.Sprintf("%s、", structs.NumberToChinese(i+1)), prefix)...)
		}

		keeper.apimks = structs.NewMarkdown(rows...).HTML()
		executors = append(executors, server.XhttpFuncHandle(keeper.api, "api"))
	}
	for _, httpconf := range keeper.httpconfs {
		executors = append(executors, httpconf.ExportExecutors()...)
	}
	return executors
}
