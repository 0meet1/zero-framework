package processors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacMysqlAutoProcessor struct {
	ZeroCoreProcessor

	fields []*structs.ZeroXsacField

	triggers []structs.ZeroXsacTrigger
}

func NewXsacMysqlProcessor(triggers ...structs.ZeroXsacTrigger) *ZeroXsacMysqlAutoProcessor {
	return &ZeroXsacMysqlAutoProcessor{
		triggers: triggers,
	}
}

func (processor *ZeroXsacMysqlAutoProcessor) AddFields(fields []*structs.ZeroXsacField) {
	processor.fields = fields
}

func (processor *ZeroXsacMysqlAutoProcessor) AddTriggers(triggers ...structs.ZeroXsacTrigger) {
	if processor.triggers == nil {
		processor.triggers = make([]structs.ZeroXsacTrigger, 0)
	}
	processor.triggers = append(processor.triggers, triggers...)
}

func (processor *ZeroXsacMysqlAutoProcessor) on(eventType string, data interface{}) error {
	if processor.triggers != nil {
		for _, trigger := range processor.triggers {
			err := trigger.On(eventType, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) exterField(field *structs.ZeroXsacField, data reflect.Value, vdata reflect.Value) (string, []interface{}) {
	makeLinkSQL := fmt.Sprintf("INSERT INTO %s(%s, %s) VALUES (? ,?)", field.Reftable(), field.Refcolumn(), field.Refbrocolumn())
	return makeLinkSQL, []interface{}{data.FieldByName("ID").Interface(), vdata.FieldByName("ID").Interface()}
}

func (processor *ZeroXsacMysqlAutoProcessor) insertWithField(fields []*structs.ZeroXsacField, data interface{}) error {
	reflect.ValueOf(data).MethodByName("InitDefault").Call([]reflect.Value{})
	elem := reflect.ValueOf(data).Elem()

	dataset := make([]interface{}, 0)
	fieldStrings := ""
	valueStrings := ""

	delaydatas := make(map[interface{}][]*structs.ZeroXsacField)
	delaystmts := make([]string, 0)
	delaydataset := make(map[string][]interface{})

	addFieldString := func(field *structs.ZeroXsacField) {
		if len(fieldStrings) <= 0 {
			fieldStrings = field.ColumnName()
			valueStrings = "?"
		} else {
			fieldStrings = fmt.Sprintf("%s,%s", fieldStrings, field.ColumnName())
			valueStrings = fmt.Sprintf("%s,?", valueStrings)
		}
	}

	for _, field := range fields {
		if !field.Writable() {
			continue
		}
		vdata := elem.FieldByName(field.FieldName())

		if vdata.Kind() == reflect.Pointer {
			if vdata.IsNil() {
				continue
			}
			vdata = vdata.Elem()
		}

		if field.Inlinable() {
			if vdata.FieldByName("ID").Interface().(string) == "" {
				continue
			}
			if field.Exterable() {
				makeLinkSQL, dataLinks := processor.exterField(field, elem, vdata)
				delaystmts = append(delaystmts, makeLinkSQL)
				delaydataset[makeLinkSQL] = dataLinks
			} else {
				addFieldString(field)
				dataset = append(dataset, vdata.FieldByName("ID").Interface())
			}
		} else if field.Childable() {
			if field.IsArray() {
				for i := 0; i < vdata.Len(); i++ {
					vxdatai := vdata.Index(i).Interface()
					vdatai := reflect.ValueOf(vxdatai)
					if vdatai.Kind() == reflect.Pointer {
						reflect.ValueOf(vxdatai).Elem().FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
					} else {
						reflect.ValueOf(vxdatai).FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
					}
					delaydatas[vxdatai] = field.XLinkFields()
				}
			} else {
				if vdata.Kind() == reflect.Pointer {
					vdata.FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
				} else {
					vdata.FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
				}
				delaydatas[vdata.Interface()] = field.XLinkFields()
			}
		} else {
			addFieldString(field)
			if field.Metatype().PkgPath() == "github.com/0meet1/zero-framework/structs" && field.Metatype().Name() == "Time" {
				dataset = append(dataset, vdata.Interface().(*structs.Time).Time().Format("2006-01-02 15:04:05"))
			} else if vdata.Type().Kind() == reflect.Map ||
				vdata.Type().Kind() == reflect.Slice ||
				vdata.Type().Kind() == reflect.Struct {
				jsonbytes, _ := json.Marshal(vdata.Interface())
				dataset = append(dataset, string(jsonbytes))
			} else {
				dataset = append(dataset, vdata.Interface())
			}
		}
	}

	_, err := processor.PreparedStmt(fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", data.(structs.ZeroXsacDeclares).XsacTableName(), fieldStrings, valueStrings)).Exec(dataset...)
	if err != nil {
		return err
	}

	for delaydata, delayfields := range delaydatas {
		err = processor.insertWithField(delayfields, delaydata)
		if err != nil {
			return err
		}
	}

	for _, delaystmt := range delaystmts {
		_, err = processor.PreparedStmt(delaystmt).Exec(delaydataset[delaystmt]...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Insert(datas ...interface{}) error {
	for _, data := range datas {
		reflect.ValueOf(data).MethodByName("InitDefault").Call([]reflect.Value{})
		err := processor.on(XSAC_BE_INSERT, data)
		if err != nil {
			return err
		}
		err = processor.insertWithField(processor.fields, data)
		if err != nil {
			return err
		}
		err = processor.on(XSAC_AF_INSERT, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Update(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		dataset := make([]interface{}, 0)
		updatefields := ""

		for _, field := range processor.fields {
			if field.Updatable() && field.FieldName() != "ID" {
				vdata := elem.FieldByName(field.FieldName())
				if vdata.Kind() == reflect.Pointer && vdata.IsNil() {
					continue
				}

				if len(updatefields) <= 0 {
					updatefields = fmt.Sprintf("%s = ?", field.ColumnName())
				} else {
					updatefields = fmt.Sprintf("%s,%s = ?", updatefields, field.ColumnName())
				}

				if field.Metatype().PkgPath() == "github.com/0meet1/zero-framework/structs" && field.Metatype().Name() == "Time" {
					dataset = append(dataset, vdata.Interface().(*structs.Time).Time().Format("2006-01-02 15:04:05"))
				} else if vdata.Kind() == reflect.Map ||
					vdata.Kind() == reflect.Slice ||
					vdata.Kind() == reflect.Struct {
					jsonbytes, _ := json.Marshal(vdata.Interface())
					dataset = append(dataset, string(jsonbytes))
				} else {
					dataset = append(dataset, vdata.Interface())
				}
			}
		}

		dataset = append(dataset, elem.FieldByName("ID").Interface())

		err := processor.on(XSAC_BE_UPDATE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET %s WHERE ID = ?", data.(structs.ZeroXsacDeclares).XsacTableName(), updatefields)).Exec(dataset...)
		if err != nil {
			return err
		}
		err = processor.on(XSAC_AF_UPDATE, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Delete(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("DELETE FROM %s WHERE ID = ?", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
		if err != nil {
			return err
		}
		err = processor.on(XSAC_AF_DELETE, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Tombstone(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET flag = 1 WHERE ID = ?", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
		if err != nil {
			return err
		}
		err = processor.on(XSAC_AF_DELETE, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Xrestore(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET flag = 0 WHERE ID = ?", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
		if err != nil {
			return err
		}
		err = processor.on(XSAC_AF_DELETE, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *ZeroXsacMysqlAutoProcessor) Fetch(dataId string) (interface{}, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE ID = ? LIMIT 1", reflect.ValueOf(processor.fields[0].Metatype()).Elem().FieldByName("XsacTableName").Interface().(func() string)())
	rowdata, err := processor.PreparedStmt(stmt).Query(dataId)
	if err != nil {
		return nil, err
	}

	rows := processor.Parser(rowdata)
	if len(rows) <= 0 {
		return nil, nil
	}

	data, err := structs.XautoLoad(processor.fields[0].Metatype(), rows[0])
	if err != nil {
		return nil, err
	}
	return data.Interface(), nil
}

func (processor *ZeroXsacMysqlAutoProcessor) FetchChildrens(field *structs.ZeroXsacField, datas interface{}) error {
	stmtChildrens := ""
	stmtdata := reflect.ValueOf(datas).Elem().FieldByName("ID").Interface()
	if field.Exterable() {
		stmtChildrens = fmt.Sprintf(
			"SELECT a.* FROM %s a, %s b WHERE a.id = b.%s AND b.%s = ?",
			field.SubTableName(),
			field.Reftable(),
			field.Refbrocolumn(),
			field.Refcolumn())
		if !field.IsArray() {
			stmtChildrens = fmt.Sprintf("%s LIMIT 1", stmtChildrens)
		}
	} else {
		if field.Inlinable() {
			stmtChildrens = fmt.Sprintf("SELECT * FROM %s WHERE ID = ?", field.SubTableName())
			superf := reflect.ValueOf(datas).Elem().FieldByName(field.FieldName())
			if superf.Kind() == reflect.Ptr {
				if superf.IsNil() {
					return nil
				}
				superf = superf.Elem()
			}

			if superf.FieldByName("ID").Interface().(string) == "" {
				return nil
			}
			stmtdata = superf.FieldByName("ID").Interface()
		} else {
			stmtChildrens = fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", field.SubTableName(), field.ChildColumnName())
			if !field.IsArray() {
				stmtChildrens = fmt.Sprintf("%s LIMIT 1", stmtChildrens)
			}
		}
	}

	rowdatas, err := processor.PreparedStmt(stmtChildrens).Query(stmtdata)
	if err != nil {
		return err
	}

	rows := processor.Parser(rowdatas)
	if len(rows) > 0 {
		if field.IsArray() {
			subdatas := reflect.MakeSlice(reflect.ValueOf(datas).Elem().FieldByName(field.FieldName()).Type(), len(rows), len(rows))
			for i, row := range rows {
				data, err := structs.XautoLoad(processor.fields[0].Metatype(), row)
				if err != nil {
					return err
				}
				subdatas.Index(i).Set(data)
			}
			reflect.ValueOf(datas).Elem().FieldByName(field.FieldName()).Set(subdatas)
		} else {
			data, err := structs.XautoLoad(processor.fields[0].Metatype(), rows[0])
			if err != nil {
				return err
			}
			reflect.ValueOf(datas).Elem().FieldByName(field.FieldName()).Set(data)
		}
	}
	return nil
}
