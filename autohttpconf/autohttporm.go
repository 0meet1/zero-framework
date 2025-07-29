package autohttpconf

import (
	"database/sql"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
)

func (e *ZeroXsacXhttpStructs) Vnew(ds ...string) error {
	_ds := e.XsacDataSource()
	if len(ds) > 0 {
		_ds = ds[0]
	}
	if _ds == "" {
		_ds = database.DATABASE_POSTGRES
	}
	transaction := global.Value(_ds).(database.DataSource).Transaction()
	defer transaction.Commit()
	err := e.VTnew(transaction)
	if err != nil {
		_err := transaction.Rollback()
		if _err != nil {
			global.Logger().ErrorS(_err)
		}
		return err
	}
	return nil
}

func (e *ZeroXsacXhttpStructs) Vupdate(ds ...string) error {
	_ds := e.XsacDataSource()
	if len(ds) > 0 {
		_ds = ds[0]
	}
	if _ds == "" {
		_ds = database.DATABASE_POSTGRES
	}
	transaction := global.Value(_ds).(database.DataSource).Transaction()
	defer transaction.Commit()
	err := e.VTupdate(transaction)
	if err != nil {
		_err := transaction.Rollback()
		if _err != nil {
			global.Logger().ErrorS(_err)
		}
		return err
	}
	return nil
}

func (e *ZeroXsacXhttpStructs) Vremove(ds ...string) error {
	_ds := e.XsacDataSource()
	if len(ds) > 0 {
		_ds = ds[0]
	}
	if _ds == "" {
		_ds = database.DATABASE_POSTGRES
	}
	transaction := global.Value(_ds).(database.DataSource).Transaction()
	defer transaction.Commit()
	err := e.VTremove(transaction)
	if err != nil {
		_err := transaction.Rollback()
		if _err != nil {
			global.Logger().ErrorS(_err)
		}
		return err
	}
	return nil
}

func (e *ZeroXsacXhttpStructs) VTnew(transaction *sql.Tx) error {
	processor := e.XhttpAutoProc()
	processor.AddFields(e.XsacFields())
	processor.Build(transaction)
	return processor.Insert(e.This())
}

func (e *ZeroXsacXhttpStructs) VTupdate(transaction *sql.Tx) error {
	processor := e.XhttpAutoProc()
	processor.AddFields(e.XsacFields())
	processor.Build(transaction)
	return processor.Update(e.This())
}

func (e *ZeroXsacXhttpStructs) VTremove(transaction *sql.Tx) error {
	processor := e.XhttpAutoProc()
	processor.AddFields(e.XsacFields())
	processor.Build(transaction)
	return processor.Delete(e.This())
}
