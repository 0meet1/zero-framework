package structs

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const (
	XSAC_NAME        = "xsacname"
	XSAC_PROP        = "xsacprop"
	XSAC_REF         = "xsacref"
	XSAC_KEY         = "xsackey"
	XSAC_REF_INSPECT = "inspect"
	XSAC_CHILD       = "xsacchild"
	XSAC_FIELD       = "xsacfield"

	XHTTP_OPT = "xhttpopt"

	XSAC_NULL = "NULL"
	XSAC_YES  = "YES"
	XSAC_NO   = "NO"
)

func FindMetaType(fields reflect.StructField) reflect.Type {
	metaType := fields.Type
	for metaType.Kind() == reflect.Pointer || metaType.Kind() == reflect.Slice {
		metaType = metaType.Elem()
	}
	return metaType
}

type ZeroCoreStructs struct {
	ZeroMeta

	ID         string                 `json:"id,omitempty" xhttpopt:"OX"`
	CreateTime Date                   `json:"createTime,omitempty" xhttpopt:"XX"`
	UpdateTime Date                   `json:"updateTime,omitempty" xhttpopt:"XX"`
	Features   map[string]interface{} `json:"features,omitempty" xhttpopt:"OO"`
}

func (e *ZeroCoreStructs) XsacPrimaryType() string { return "UUID" }
func (e *ZeroCoreStructs) XsacDataSource() string  { return "" }
func (e *ZeroCoreStructs) XsacDbName() string      { return "" }
func (e *ZeroCoreStructs) XsacTableName() string   { panic("not implemented") }

func (e *ZeroCoreStructs) findXsacEntry(fields reflect.StructField) []*ZeroXsacEntry {
	entries := make([]*ZeroXsacEntry, 0)

	xrProp := fields.Tag.Get(XSAC_PROP)
	if len(xrProp) > 0 {
		xrPropItems := strings.Split(xrProp, ",")
		if len(xrPropItems) == 3 {
			columnName := fields.Tag.Get(XSAC_NAME)
			if len(columnName) <= 0 {
				columnName = exHumpToLine(fields.Name)
			}

			entries = append(entries, NewColumn(
				e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName(),
				columnName, xrPropItems[0], xrPropItems[1], xrPropItems[2]))

			xsacKey := fields.Tag.Get(XSAC_KEY)
			if len(xsacKey) > 0 {
				if strings.HasPrefix(xsacKey, "foreign") {
					xrKeyItems := strings.Split(xsacKey, ",")
					if len(xrKeyItems) == 3 {
						entries = append(entries, NewForeignKey(e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName(), columnName, xrKeyItems[1], xrKeyItems[2]))
					}
				} else {
					switch xsacKey {
					case "primary":
						entries = append(entries, NewPrimaryKey(e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName(), columnName))
					case "key":
						entries = append(entries, NewKey(e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName(), columnName))
					case "unique":
						entries = append(entries, NewUniqueKey(e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName(), columnName))
					}
				}
			}
		}
	}
	return entries
}

func (e *ZeroCoreStructs) readXsacEntries(xrType reflect.Type) []*ZeroXsacEntry {
	entries := make([]*ZeroXsacEntry, 0)
	for i := 0; i < xrType.NumField(); i++ {
		if xrType.Field(i).Anonymous {
			entries = append(entries, e.readXsacEntries(xrType.Field(i).Type)...)
		} else {
			entries = append(entries, e.findXsacEntry(xrType.Field(i))...)
		}
	}
	return entries
}

func (e *ZeroCoreStructs) XsacDeclares() ZeroXsacEntrySet {
	entries := make([]*ZeroXsacEntry, 0)
	entries = append(entries, NewTable0s(e.This().(ZeroXsacDeclares).XsacDbName(), e.This().(ZeroXsacDeclares).XsacTableName()))
	entries = append(entries, e.readXsacEntries(reflect.TypeOf(e.This()).Elem())...)
	return entries
}

func (e *ZeroCoreStructs) findXsacRefEntry(fields reflect.StructField) []*ZeroXsacEntry {
	entries := make([]*ZeroXsacEntry, 0)
	xrRefProp := fields.Tag.Get(XSAC_REF)
	metaType := FindMetaType(fields)
	if len(xrRefProp) > 0 {
		xrRefProppItems := strings.Split(xrRefProp, ",")
		if len(xrRefProppItems) == 4 && xrRefProppItems[3] == XSAC_REF_INSPECT {
			entries = append(entries, NewTable(e.This().(ZeroXsacDeclares).XsacDbName(), xrRefProppItems[0]))
			entries = append(entries, NewColumn(e.This().(ZeroXsacDeclares).XsacDbName(), xrRefProppItems[0], xrRefProppItems[1], XSAC_NO, e.XsacPrimaryType(), XSAC_NULL))
			entries = append(entries, NewColumn(e.This().(ZeroXsacDeclares).XsacDbName(), xrRefProppItems[0], xrRefProppItems[2], XSAC_NO, e.XsacPrimaryType(), XSAC_NULL))
			entries = append(entries, NewForeignKey(e.This().(ZeroXsacDeclares).XsacDbName(), xrRefProppItems[0], xrRefProppItems[1], e.This().(ZeroXsacDeclares).XsacTableName(), "id"))
			entries = append(entries, NewForeignKey(e.This().(ZeroXsacDeclares).XsacDbName(), xrRefProppItems[0], xrRefProppItems[2], reflect.New(metaType).Interface().(ZeroXsacDeclares).XsacTableName(), "id"))
		}
	}
	return entries
}

func (e *ZeroCoreStructs) readXsacRefEntries(xrType reflect.Type) []*ZeroXsacEntry {
	entries := make([]*ZeroXsacEntry, 0)
	for i := 0; i < xrType.NumField(); i++ {
		if xrType.Field(i).Anonymous {
			entries = append(entries, e.readXsacRefEntries(xrType.Field(i).Type)...)
		} else {
			entries = append(entries, e.findXsacRefEntry(xrType.Field(i))...)
		}
	}
	return entries
}

func (e *ZeroCoreStructs) XsacRefDeclares() ZeroXsacEntrySet {
	return e.readXsacRefEntries(reflect.TypeOf(e.This()).Elem())
}

func (e *ZeroCoreStructs) findXopFields(xrType reflect.Type, ignore bool) ZeroXsacFieldSet {
	fields := make([]*ZeroXsacField, 0)
	for i := 0; i < xrType.NumField(); i++ {
		if xrType.Field(i).Anonymous {
			fields = append(fields, e.findXopFields(xrType.Field(i).Type, ignore)...)
		} else if len(xrType.Field(i).Tag.Get(XHTTP_OPT)) > 0 {
			fields = append(fields, NewXsacField(xrType.Field(i), ignore))
		}
	}
	return fields
}

func (e *ZeroCoreStructs) XsacFields(xm ...int) ZeroXsacFieldSet {
	fields := make([]*ZeroXsacField, 0)
	fields = append(fields, e.findXopFields(reflect.TypeOf(e.This()).Elem(), len(xm) > 0)...)
	return fields
}

func (e *ZeroCoreStructs) InitDefault() error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	e.ID = uid.String()
	e.CreateTime = Date(time.Now())
	e.UpdateTime = Date(time.Now())
	if e.Features == nil {
		e.Features = make(map[string]interface{})
	}
	return nil
}

func (e *ZeroCoreStructs) JSONFeature() string {
	if e.Features == nil {
		e.Features = make(map[string]interface{})
	}
	mjson, _ := json.Marshal(e.Features)
	return string(mjson)
}

func (e *ZeroCoreStructs) JSONFeatureWithString(jsonString string) {
	var jsonMap map[string]interface{}
	_ = json.Unmarshal([]byte(jsonString), &jsonMap)
	e.Features = jsonMap
}

func (e *ZeroCoreStructs) LoadRowData(rowmap map[string]interface{}) {
	_, ok := rowmap["id"]
	if ok {
		e.ID = string(rowmap["id"].([]uint8))
	}

	_, ok = rowmap["create_time"]
	if ok {
		e.CreateTime = Date(rowmap["create_time"].(time.Time))
	}

	_, ok = rowmap["update_time"]
	if ok {
		e.UpdateTime = Date(rowmap["update_time"].(time.Time))
	}

	_, ok = rowmap["features"]
	if ok {
		e.JSONFeatureWithString(string(rowmap["features"].([]uint8)))
	}
}

func (e *ZeroCoreStructs) String() string {
	mjson, _ := json.Marshal(e)
	return string(mjson)
}

func (e *ZeroCoreStructs) Map() map[string]interface{} {
	mjson, _ := json.Marshal(e)
	var jsonMap map[string]interface{}
	_ = json.Unmarshal([]byte(mjson), &jsonMap)
	return jsonMap
}
