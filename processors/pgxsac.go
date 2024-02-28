package processors

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacPostgresAutoProcessor struct {
	ZeroCoreProcessor

	dbName    string
	tableName string
	fields    []*structs.ZeroXsacField

	triggers []ZeroXsacTrigger
}

func NewXsacPostgresProcessor(dbName string, tableName string, triggers ...ZeroXsacTrigger) *ZeroXsacPostgresAutoProcessor {
	return &ZeroXsacPostgresAutoProcessor{
		dbName:    dbName,
		tableName: tableName,
		triggers:  triggers,
	}
}

func (processor *ZeroXsacPostgresAutoProcessor) DBName() string {
	return processor.dbName
}

func (processor *ZeroXsacPostgresAutoProcessor) TableName() string {
	return processor.tableName
}

func (processor *ZeroXsacPostgresAutoProcessor) AddFields(fields []*structs.ZeroXsacField) {
	processor.fields = fields
}

func (processor *ZeroXsacPostgresAutoProcessor) AddTriggers(triggers ...ZeroXsacTrigger) {
	if processor.triggers == nil {
		processor.triggers = make([]ZeroXsacTrigger, 0)
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
	elem := reflect.ValueOf(data).Elem()

	dataset := make([]interface{}, 0)
	fieldStrings := ""
	valueStrings := ""
	fieldIdx := 0

	delaydatas := make(map[interface{}][]*structs.ZeroXsacField)
	delaystmts := make([]string, 0)
	delaydataset := make(map[string][]interface{})

	addFieldString := func(field *structs.ZeroXsacField, vdata reflect.Value) {
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
				addFieldString(field, vdata)
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
			addFieldString(field, vdata)
			if vdata.Kind() == reflect.Map ||
				vdata.Kind() == reflect.Slice ||
				structs.FindMetaType(vdata.Type()).Kind() == reflect.Struct {
				jsonbytes, _ := json.Marshal(vdata.Interface())
				dataset = append(dataset, string(jsonbytes))
			} else {
				dataset = append(dataset, vdata.Interface())
			}
		}
	}

	_, err := processor.PreparedStmt(fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", processor.tableName, fieldStrings, valueStrings)).Exec(dataset...)
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

				if vdata.Kind() == reflect.Map ||
					vdata.Kind() == reflect.Slice ||
					structs.FindMetaType(vdata.Type()).Kind() == reflect.Struct {
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
		_, err = processor.PreparedStmt(fmt.Sprintf("UPDATE %s SET %s WHERE ID = $%d", processor.tableName, updatefields, fieldIdx)).Exec(dataset...)
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
		_, err = processor.PreparedStmt(fmt.Sprintf("DELETE FROM %s WHERE ID = $1", processor.tableName)).Exec(elem.FieldByName("ID").Interface())
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
