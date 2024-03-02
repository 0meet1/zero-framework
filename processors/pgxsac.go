package processors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacPostgresAutoProcessor struct {
	ZeroCoreProcessor

	fields []*structs.ZeroXsacField

	triggers []structs.ZeroXsacTrigger
}

func NewXsacPostgresProcessor(triggers ...structs.ZeroXsacTrigger) *ZeroXsacPostgresAutoProcessor {
	return &ZeroXsacPostgresAutoProcessor{
		triggers: triggers,
	}
}

func (processor *ZeroXsacPostgresAutoProcessor) AddFields(fields []*structs.ZeroXsacField) {
	processor.fields = fields
}

func (processor *ZeroXsacPostgresAutoProcessor) AddTriggers(triggers ...structs.ZeroXsacTrigger) {
	if processor.triggers == nil {
		processor.triggers = make([]structs.ZeroXsacTrigger, 0)
	}
	for _, trigger := range triggers {
		processor.triggers = append(processor.triggers, trigger)
	}
}

func (processor *ZeroXsacPostgresAutoProcessor) on(eventType string, data interface{}) error {
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

func (processor *ZeroXsacPostgresAutoProcessor) exterField(field *structs.ZeroXsacField, data reflect.Value, vdata reflect.Value) (string, []interface{}) {
	makeLinkSQL := fmt.Sprintf("INSERT INTO %s(%s, %s) VALUES ($1 ,$2)", field.Reftable(), field.Refcolumn(), field.Refbrocolumn())
	return makeLinkSQL, []interface{}{data.FieldByName("ID").Interface(), vdata.FieldByName("ID").Interface()}
}

func (processor *ZeroXsacPostgresAutoProcessor) insertWithField(fields []*structs.ZeroXsacField, data interface{}) error {
	reflect.ValueOf(data).MethodByName("InitDefault").Call([]reflect.Value{})
	elem := reflect.ValueOf(data).Elem()

	dataset := make([]interface{}, 0)
	fieldStrings := ""
	valueStrings := ""
	fieldIdx := 0

	delaydatas := make(map[interface{}][]*structs.ZeroXsacField)
	delaystmts := make([]string, 0)
	delaydataset := make(map[string][]interface{})

	addFieldString := func(field *structs.ZeroXsacField) {
		fieldIdx++
		if len(fieldStrings) <= 0 {
			fieldStrings = field.ColumnName()
			valueStrings = fmt.Sprintf("$%d", fieldIdx)
		} else {
			fieldStrings = fmt.Sprintf("%s,%s", fieldStrings, field.ColumnName())
			valueStrings = fmt.Sprintf("%s,$%d", valueStrings, fieldIdx)
		}
	}

	for _, field := range fields {
		if !field.Writable() {
			continue
		}
		vdata := elem.FieldByName(field.FieldName())

		if vdata.Kind() == reflect.Pointer && vdata.IsNil() {
			continue
		}

		if field.Inlinable() {
			if vdata.Elem().FieldByName("ID").Interface().(string) == "" {
				continue
			}
			if field.Exterable() {
				makeLinkSQL, dataLinks := processor.exterField(field, elem, vdata)
				delaystmts = append(delaystmts, makeLinkSQL)
				delaydataset[makeLinkSQL] = dataLinks
			} else {
				addFieldString(field)
				dataset = append(dataset, vdata.Elem().FieldByName("ID").Interface())
			}
		} else if field.Childable() {
			if field.Exterable() {
				if field.IsArray() {
					for i := 0; i < vdata.Len(); i++ {
						vxdatai := vdata.Index(i).Interface()
						vdatai := reflect.ValueOf(vxdatai)
						vdatai.MethodByName("InitDefault").Call([]reflect.Value{})
						delaydatas[vxdatai] = field.XLinkFields()

						makeLinkSQL, dataLinks := processor.exterField(field, elem, vdatai)
						delaystmts = append(delaystmts, makeLinkSQL)
						delaydataset[makeLinkSQL] = dataLinks
					}
				} else {
					vdata.MethodByName("InitDefault").Call([]reflect.Value{})
					delaydatas[vdata.Interface()] = field.XLinkFields()

					makeLinkSQL, dataLinks := processor.exterField(field, elem, vdata)
					delaystmts = append(delaystmts, makeLinkSQL)
					delaydataset[makeLinkSQL] = dataLinks
				}
			} else {
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
						vdata.Elem().FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
					} else {
						vdata.FieldByName(field.ChildName()).Set(reflect.ValueOf(data))
					}
					delaydatas[vdata.Interface()] = field.XLinkFields()
				}
			}
		} else {
			addFieldString(field)
			if vdata.Kind() == reflect.Map ||
				vdata.Kind() == reflect.Slice ||
				structs.FindMetaType(vdata.Type()).Kind() == reflect.Struct {
				jsonbytes, _ := json.Marshal(vdata.Interface())
				dataset = append(dataset, string(jsonbytes))
			} else {
				if field.Metatype().PkgPath() == "github.com/0meet1/zero-framework/structs" && field.Metatype().Name() == "Time" {
					dataset = append(dataset, vdata.Interface().(*structs.Time).Time().Format("2006-01-02 15:04:05"))
				} else if field.Metatype().Kind() == reflect.Map ||
					field.Metatype().Kind() == reflect.Slice ||
					field.Metatype().Kind() == reflect.Struct {
					jsonbytes, _ := json.Marshal(vdata.Interface())
					dataset = append(dataset, string(jsonbytes))
				} else {
					dataset = append(dataset, vdata.Interface())
				}
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

func (processor *ZeroXsacPostgresAutoProcessor) Insert(datas ...interface{}) error {
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

func (processor *ZeroXsacPostgresAutoProcessor) Update(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		dataset := make([]interface{}, 0)
		updatefields := ""
		fieldIdx := 0

		for _, field := range processor.fields {
			if field.Updatable() && field.FieldName() != "ID" {
				fieldIdx++
				vdata := elem.FieldByName(field.FieldName())
				if vdata.Kind() == reflect.Pointer && vdata.IsNil() {
					continue
				}

				if len(updatefields) <= 0 {
					updatefields = fmt.Sprintf("%s = $%d", field.ColumnName(), fieldIdx)
				} else {
					updatefields = fmt.Sprintf("%s,%s = $%d", updatefields, field.ColumnName(), fieldIdx)
				}

				if field.Metatype().PkgPath() == "github.com/0meet1/zero-framework/structs" && field.Metatype().Name() == "Time" {
					dataset = append(dataset, vdata.Interface().(*structs.Time).Time().Format("2006-01-02 15:04:05"))
				} else if field.Metatype().Kind() == reflect.Map ||
					field.Metatype().Kind() == reflect.Slice ||
					field.Metatype().Kind() == reflect.Struct {
					jsonbytes, _ := json.Marshal(vdata.Interface())
					dataset = append(dataset, string(jsonbytes))
				} else {
					dataset = append(dataset, vdata.Interface())
				}
			}
		}

		fieldIdx++
		dataset = append(dataset, elem.FieldByName("ID").Interface())

		err := processor.on(XSAC_BE_UPDATE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET %s WHERE ID = $%d", data.(structs.ZeroXsacDeclares).XsacTableName(), updatefields, fieldIdx)).Exec(dataset...)
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

func (processor *ZeroXsacPostgresAutoProcessor) Delete(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("DELETE FROM %s WHERE ID = $1", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
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

func (processor *ZeroXsacPostgresAutoProcessor) Tombstone(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET flag = 1 WHERE ID = $1", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
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

func (processor *ZeroXsacPostgresAutoProcessor) Xrestore(datas ...interface{}) error {
	for _, data := range datas {
		elem := reflect.ValueOf(data).Elem()
		if elem.FieldByName("ID").Interface().(string) == "" {
			continue
		}

		err := processor.on(XSAC_BE_DELETE, data)
		if err != nil {
			return err
		}
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET flag = 0 WHERE ID = $1", data.(structs.ZeroXsacDeclares).XsacTableName())).Exec(elem.FieldByName("ID").Interface())
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

func (processor *ZeroXsacPostgresAutoProcessor) FetchChildrens(field *structs.ZeroXsacField, datas interface{}) error {
	stmtChildrens := ""
	stmtdata := reflect.ValueOf(datas).Elem().FieldByName("ID").Interface()
	if field.Exterable() {
		stmtChildrens = fmt.Sprintf(
			"SELECT a.* FROM %s a, %s b WHERE WHERE a.id = b.%s AND %s = $1",
			field.SubTableName(),
			field.Reftable(),
			field.Refbrocolumn(),
			field.Refcolumn())
		if !field.IsArray() {
			stmtChildrens = fmt.Sprintf("%s LIMIT 1", stmtChildrens)
		}
	} else {
		if field.Inlinable() {
			stmtChildrens = fmt.Sprintf("SELECT * FROM %s WHERE ID = $1", field.SubTableName())
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
			stmtChildrens = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", field.SubTableName(), field.ChildColumnName())
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
				data := reflect.New(field.Metatype())
				returnValues := data.MethodByName("LoadRowData").Call([]reflect.Value{reflect.ValueOf(row)})
				if len(returnValues) > 0 && returnValues[0].Interface() != nil {
					return returnValues[0].Interface().(error)
				}
				subdatas.Index(i).Set(data)
			}
			reflect.ValueOf(datas).Elem().FieldByName(field.FieldName()).Set(subdatas)
		} else {
			data := reflect.New(field.Metatype())
			returnValues := data.MethodByName("LoadRowData").Call([]reflect.Value{reflect.ValueOf(rows[0])})
			if len(returnValues) > 0 && returnValues[0].Interface() != nil {
				return returnValues[0].Interface().(error)
			}
			reflect.ValueOf(datas).Elem().FieldByName(field.FieldName()).Set(data)
		}
	}
	return nil
}
