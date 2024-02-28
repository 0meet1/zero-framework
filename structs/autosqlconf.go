package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	ZEOR_XSAC_ENTRY_TYPE_TABLE            = "table"
	ZEOR_XSAC_ENTRY_TYPE_TABLE0S          = "table0s"
	ZEOR_XSAC_ENTRY_TYPE_TABLE0FS         = "table0fs"
	ZEOR_XSAC_ENTRY_TYPE_COLUMN           = "column"
	ZEOR_XSAC_ENTRY_TYPE_DROPCOLUMN       = "dropcolumn"
	ZEOR_XSAC_ENTRY_TYPE_KEY              = "key"
	ZEOR_XSAC_ENTRY_TYPE_DROPKEY          = "dropkey"
	ZEOR_XSAC_ENTRY_TYPE_PRIMARY_KEY      = "primary"
	ZEOR_XSAC_ENTRY_TYPE_DROP_PRIMARY_KEY = "dropprimary"
	ZEOR_XSAC_ENTRY_TYPE_UNIQUE_KEY       = "unique"
	ZEOR_XSAC_ENTRY_TYPE_DROP_UNIQUE_KEY  = "dropunique"
	ZEOR_XSAC_ENTRY_TYPE_FOREIGN_KEY      = "foreign"
	ZEOR_XSAC_ENTRY_TYPE_DROP_FOREIGN_KEY = "dropforeign"
)

type ZeroXsacDeclares interface {
	XsacDataSource() string
	XsacDbName() string
	XsacTableName() string
	XsacDeclares() ZeroXsacEntrySet
	XsacRefDeclares() ZeroXsacEntrySet
}

type ZeroXsacEntrySet []*ZeroXsacEntry

func (entrySet ZeroXsacEntrySet) String() string {
	output := make([]string, 0)
	for _, entry := range entrySet {
		output = append(output, entry.String())
	}
	if len(output) <= 0 {
		return "{}"
	}
	return fmt.Sprintf("{\n\t%s\n}", strings.Join(output, ",\n\t"))
}

type ZeroXsacEntry struct {
	entryType   string
	entryParams []string
}

func (xe *ZeroXsacEntry) EntryType() string {
	return xe.entryType
}

func (xe *ZeroXsacEntry) EntryParams() []string {
	return xe.entryParams
}

func (xe *ZeroXsacEntry) String() string {
	return fmt.Sprintf("%s(`%s`)", xe.entryType, strings.Join(xe.entryParams, "`, `"))
}

func NewTable(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_TABLE,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewTable0s(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_TABLE0S,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewTable0fs(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_TABLE0FS,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewColumn(tableSchema string, tableName string, columnName string, isNullable string, columnType string, columnDefault string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_COLUMN,
		entryParams: []string{tableSchema, tableName, columnName, isNullable, columnType, columnDefault},
	}
}

func NewDropColumn(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DROPCOLUMN,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func NewKey(tableSchema string, tableName string, indexName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_KEY,
		entryParams: []string{tableSchema, tableName, indexName},
	}
}

func NewDropKey(tableSchema string, tableName string, indexName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DROPKEY,
		entryParams: []string{tableSchema, tableName, indexName},
	}
}

func NewPrimaryKey(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_PRIMARY_KEY,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func NewDropPrimaryKey(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DROP_PRIMARY_KEY,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func NewUniqueKey(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_UNIQUE_KEY,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func NewDropUniqueKey(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DROP_UNIQUE_KEY,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func NewForeignKey(tableSchema string, tableName string, columnName string, relTableName string, relColumnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_FOREIGN_KEY,
		entryParams: []string{tableSchema, tableName, columnName, relTableName, relColumnName},
	}
}

func NewDropForeignKey(tableSchema string, tableName string, columnName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DROP_FOREIGN_KEY,
		entryParams: []string{tableSchema, tableName, columnName},
	}
}

func exHumpToLine(name string) string {
	_name := name
	if strings.HasPrefix(_name, "ID") {
		_name = strings.ReplaceAll(_name, "ID", "id")
	} else {
		_name = strings.ReplaceAll(_name, "ID", "_id")
	}

	namebytes := []byte(_name)
	var buf bytes.Buffer
	for i, c := range namebytes {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				buf.WriteByte('_')
			}
			buf.WriteByte(c + 32)
		} else {
			buf.WriteByte(c)
		}
	}
	return string(buf.Bytes())
}

type ZeroXsacFieldSet []*ZeroXsacField

func (xfs ZeroXsacFieldSet) String() string {
	output := make([]string, 0)
	for _, xf := range xfs {
		output = append(output, xf.String())
	}
	if len(output) <= 0 {
		return "[]"
	}
	return fmt.Sprintf("[%s]", strings.Join(output, ",\n"))
}

type ZeroXsacField struct {
	metatype reflect.Type

	fieldName  string
	columnName string
	inlineName string
	childName  string

	reftable     string
	refcolumn    string
	refbrocolumn string

	isArray   bool
	writable  bool
	updatable bool

	xLinkFields ZeroXsacFieldSet
}

func NewXsacField(field reflect.StructField, ignore bool) *ZeroXsacField {
	columnName := field.Tag.Get(XSAC_NAME)
	if len(columnName) <= 0 {
		columnName = exHumpToLine(field.Name)
	}

	xhttpopt := field.Tag.Get(XHTTP_OPT)
	xfield := &ZeroXsacField{
		metatype:   FindStructFieldMetaType(field),
		fieldName:  field.Name,
		columnName: columnName,
		isArray:    field.Type.Kind() == reflect.Slice,
		writable:   len(xhttpopt) > 0 && xhttpopt[0] == 'O',
		updatable:  len(xhttpopt) > 1 && xhttpopt[1] == 'O',
	}

	xsacref := field.Tag.Get(XSAC_REF)
	if len(xsacref) > 0 {
		xsacrefItems := strings.Split(xsacref, ",")
		if len(xsacrefItems) == 3 {
			xfield.reftable = xsacrefItems[0]
			xfield.refcolumn = xsacrefItems[1]
			xfield.refbrocolumn = xsacrefItems[2]
		}
	}

	if field.Tag.Get(XSAC_CHILD) != "" {
		xfield.childName = field.Tag.Get(XSAC_CHILD)
		xfield.inlineName = ""
		if !ignore {
			xfield.XLinkFields()
		}
	} else if field.Tag.Get(XSAC_FIELD) != "" {
		xfield.inlineName = field.Tag.Get(XSAC_FIELD)
		xfield.childName = ""
		if !ignore {
			xfield.XLinkFields()
		}
	}
	return xfield
}

func (xf *ZeroXsacField) Metatype() reflect.Type {
	return xf.metatype
}

func (xf *ZeroXsacField) FieldName() string {
	return xf.fieldName
}

func (xf *ZeroXsacField) ColumnName() string {
	return xf.columnName
}

func (xf *ZeroXsacField) InlineName() string {
	return xf.inlineName
}

func (xf *ZeroXsacField) ChildName() string {
	return xf.childName
}

func (xf *ZeroXsacField) Reftable() string {
	return xf.reftable
}

func (xf *ZeroXsacField) Refcolumn() string {
	return xf.refcolumn
}

func (xf *ZeroXsacField) Refbrocolumn() string {
	return xf.refbrocolumn
}

func (xf *ZeroXsacField) IsArray() bool {
	return xf.isArray
}

func (xf *ZeroXsacField) Writable() bool {
	return xf.writable
}

func (xf *ZeroXsacField) Updatable() bool {
	return xf.updatable
}

func (xf *ZeroXsacField) Exterable() bool {
	return len(xf.reftable) > 0
}

func (xf *ZeroXsacField) Inlinable() bool {
	return len(xf.inlineName) > 0
}

func (xf *ZeroXsacField) Childable() bool {
	return len(xf.childName) > 0
}

func (xf *ZeroXsacField) XLinkFields() ZeroXsacFieldSet {
	if xf.xLinkFields == nil {
		newField := reflect.New(xf.metatype).Interface()
		reflect.ValueOf(newField).MethodByName("ThisDef").Call([]reflect.Value{reflect.ValueOf(newField)})
		xf.xLinkFields = newField.(ZeroXsacFields).XsacFields(0)
	}
	return xf.xLinkFields
}

func (xf *ZeroXsacField) Map() map[string]interface{} {
	xmap := make(map[string]interface{})
	xmap["fieldName"] = xf.fieldName
	xmap["columnName"] = xf.columnName
	xmap["inlineName"] = xf.inlineName
	xmap["childName"] = xf.childName
	xmap["isArray"] = xf.isArray
	xmap["writable"] = xf.writable
	xmap["updatable"] = xf.updatable
	xmap["reftable"] = xf.reftable
	xmap["refcolumn"] = xf.refcolumn
	xmap["refbrocolumn"] = xf.refbrocolumn
	return xmap
}

func (xf *ZeroXsacField) String() string {
	xmap := xf.Map()
	if xf.xLinkFields != nil {
		xLinkFields := make([]map[string]interface{}, 0)
		for _, xlf := range xf.xLinkFields {
			xLinkFields = append(xLinkFields, xlf.Map())
		}
		xmap["xLinkFields"] = xLinkFields
	}
	jsonbytes, _ := json.Marshal(xmap)
	return string(jsonbytes)
}

type ZeroXsacFields interface {
	XsacFields(...int) ZeroXsacFieldSet
}