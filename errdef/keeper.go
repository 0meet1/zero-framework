package errdef

import (
	"fmt"
	"reflect"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
)

const (
	EXCEPTION_KEEPER    = "zero.exception.keeper"
	EXCEPTION_AUTO_PROC = "zero.exception.auto.proc"
	EXCEPTION_OPERATION = "zero.exception.operation"

	ES00500 = "ES00500"
)

type ZeroExceptionKeeper interface {
	DescriptionByCode(string) string
}

type xZeroExceptionProcessor struct {
	processors.ZeroCoreProcessor
}

func (processor *xZeroExceptionProcessor) AddMysqlException(errdef *ZeroExceptionDef) error {
	errdef.InitDefault()
	const ADD_EXCEPTION_SQL = "INSERT INTO zero_exception_def(id, features, code, description) VALUES (?, ?, ?, ?)"
	_, err := processor.PreparedStmt(ADD_EXCEPTION_SQL).Exec(errdef.ID, errdef.JSONFeature(), errdef.Code, errdef.Description)
	return err
}

func (processor *xZeroExceptionProcessor) FetchMysqlException(code string) ([]*ZeroExceptionDef, error) {
	const FETCH_EXCEPTION_SQL = `SELECT * FROM zero_exception_def WHERE code = ?`
	rows, err := processor.PreparedStmt(FETCH_EXCEPTION_SQL).Query(code)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	rowsmap := processor.Parser(rows)
	datas := make([]*ZeroExceptionDef, len(rowsmap))
	for i, row := range rowsmap {
		data := &ZeroExceptionDef{}
		data.LoadRowData(row)
		datas[i] = data
	}
	return datas, err
}

func (processor *xZeroExceptionProcessor) AddPostgresException(errdef *ZeroExceptionDef) error {
	errdef.InitDefault()
	const ADD_EXCEPTION_SQL = "INSERT INTO zero_exception_def(id, features, code, description) VALUES ($1, $2, $3, $4)"
	_, err := processor.PreparedStmt(ADD_EXCEPTION_SQL).Exec(errdef.ID, errdef.JSONFeature(), errdef.Code, errdef.Description)
	return err
}

func (processor *xZeroExceptionProcessor) FetchPostgresException(code string) ([]*ZeroExceptionDef, error) {
	const FETCH_EXCEPTION_SQL = `SELECT * FROM zero_exception_def WHERE code = $1`
	rows, err := processor.PreparedStmt(FETCH_EXCEPTION_SQL).Query(code)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	rowsmap := processor.Parser(rows)
	datas := make([]*ZeroExceptionDef, len(rowsmap))
	for i, row := range rowsmap {
		data := &ZeroExceptionDef{}
		data.LoadRowData(row)
		datas[i] = data
	}
	return datas, err
}

func (processor *xZeroExceptionProcessor) FetchExceptions() ([]*ZeroExceptionDef, error) {
	const FETCH_EXCEPTION_SQL = `SELECT * FROM zero_exception_def`
	rows, err := processor.PreparedStmt(FETCH_EXCEPTION_SQL).Query()
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	rowsmap := processor.Parser(rows)
	datas := make([]*ZeroExceptionDef, len(rowsmap))
	for i, row := range rowsmap {
		data := &ZeroExceptionDef{}
		data.LoadRowData(row)
		datas[i] = data
	}
	return datas, err
}

type xZeroExceptionKeeper struct {
	dataSource   string
	descriptions map[string]string
}

func (keeper *xZeroExceptionKeeper) DescriptionByCode(code string) string {
	description, ok := keeper.descriptions[code]
	if !ok {
		description = ""
	}
	return description
}

func (keeper *xZeroExceptionKeeper) runKeeper(errdefs ...*ZeroExceptionDef) {
	transaction := global.Value(keeper.dataSource).(database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("%s", err))
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	exceptionProcessor := &xZeroExceptionProcessor{}
	exceptionProcessor.Build(transaction)

	defs := []*ZeroExceptionDef{{Code: ES00500, Description: "系统异常"}}
	defs = append(defs, errdefs...)

	for _, errdef := range defs {
		if keeper.dataSource == database.DATABASE_MYSQL {
			historys, err := exceptionProcessor.FetchMysqlException(errdef.Code)
			if err != nil {
				panic(err)
			}
			if len(historys) <= 0 {
				err = exceptionProcessor.AddMysqlException(errdef)
				if err != nil {
					panic(err)
				}
			}
		} else {
			historys, err := exceptionProcessor.FetchPostgresException(errdef.Code)
			if err != nil {
				panic(err)
			}
			if len(historys) <= 0 {
				err = exceptionProcessor.AddPostgresException(errdef)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	exceptions, err := exceptionProcessor.FetchExceptions()
	if err != nil {
		panic(err)
	}

	for _, errdef := range exceptions {
		keeper.descriptions[errdef.Code] = errdef.Description
	}
}

func Exports(autoProcessor processors.ZeroXsacAutoProcessor, operation processors.ZeroQueryOperation) []reflect.Type {
	global.Key(EXCEPTION_AUTO_PROC, autoProcessor)
	global.Key(EXCEPTION_OPERATION, operation)
	return []reflect.Type{reflect.TypeOf(&ZeroExceptionDef{})}
}

func RunExceptionKeeper(dataSource string, errdefs ...*ZeroExceptionDef) {
	keeper := &xZeroExceptionKeeper{
		dataSource:   dataSource,
		descriptions: make(map[string]string),
	}
	keeper.runKeeper(errdefs...)
	global.Key(EXCEPTION_KEEPER, keeper)
}
