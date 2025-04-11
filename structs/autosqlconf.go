package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/0meet1/zero-framework/global"
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

	ZEOR_XSAC_ENTRY_TYPE_YEAR_PARTITION   = "year"
	ZEOR_XSAC_ENTRY_TYPE_MONTH_PARTITION  = "month"
	ZEOR_XSAC_ENTRY_TYPE_DAY_PARTITION    = "day"
	ZEOR_XSAC_ENTRY_TYPE_CUSTOM_PARTITION = "custom"
)

const (
	XSAC_PARTITION_NONE   = "none"
	XSAC_PARTITION_YEAR   = "year"
	XSAC_PARTITION_MONTH  = "month"
	XSAC_PARTITION_DAY    = "day"
	XSAC_PARTITION_CUSTOM = "custom"
)

type ZeroXsacTrigger interface {
	On(string, interface{}) error
}

type ZeroXsacDeclares interface {
	This() interface{}
	ThisDef(interface{})

	XsacPrimaryType() string
	XsacDataSource() string
	XsacDbName() string
	XsacTableName() string
	XsacDeleteOpt() byte
	XsacDeclares(...string) ZeroXsacEntrySet
	XsacRefDeclares(...string) ZeroXsacEntrySet
	XsacPartition() string
	XsacCustomPartTrigger() string
	XsacTriggers() []ZeroXsacTrigger

	XsacAutoParser() []ZeroXsacAutoParser
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

func NewYearPartition(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_YEAR_PARTITION,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewMonthPartition(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_MONTH_PARTITION,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewDayPartition(tableSchema string, tableName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_DAY_PARTITION,
		entryParams: []string{tableSchema, tableName},
	}
}

func NewCustomPartition(tableSchema string, tableName string, partTriggerName string) *ZeroXsacEntry {
	return &ZeroXsacEntry{
		entryType:   ZEOR_XSAC_ENTRY_TYPE_CUSTOM_PARTITION,
		entryParams: []string{tableSchema, tableName, partTriggerName},
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
	return buf.String()
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
	xapi     string
	jsonopts string

	fieldName  string
	columnName string
	inlineName string

	childName       string
	childColumnName string

	subTableName string

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
		xapi:       field.Tag.Get(XHTTP_API),
		jsonopts:   field.Tag.Get("json"),
		fieldName:  field.Name,
		columnName: columnName,
		isArray:    field.Type.Kind() == reflect.Slice,
		writable:   len(xhttpopt) > 0 && xhttpopt[0] == 'O',
		updatable:  len(xhttpopt) > 1 && xhttpopt[1] == 'O',
	}

	xsacref := field.Tag.Get(XSAC_REF)
	if len(xsacref) > 0 {
		xsacrefItems := strings.Split(xsacref, ",")
		if len(xsacrefItems) >= 3 {
			xfield.reftable = xsacrefItems[0]
			xfield.refcolumn = xsacrefItems[1]
			xfield.refbrocolumn = xsacrefItems[2]
		}
	}

	if field.Tag.Get(XSAC_CHILD) != "" {
		xfield.childName = field.Tag.Get(XSAC_CHILD)
		xfield.inlineName = ""

		if !ignore {
			xLinkFields := xfield.XLinkFields()
			if !xfield.Exterable() {
				for _, xLField := range xLinkFields {
					if xLField.FieldName() == xfield.childName {
						xfield.childColumnName = xLField.ColumnName()
						break
					}
				}
			}
		}
	} else if field.Tag.Get(XSAC_FIELD) != "" {
		xfield.inlineName = field.Tag.Get(XSAC_FIELD)
		xfield.childName = ""
		if !ignore {
			xfield.XLinkFields()
		}
	}
	if xfield.Inlinable() || xfield.Childable() {
		xfield.subTableName = reflect.New(FindStructFieldMetaType(field)).Interface().(ZeroXsacDeclares).XsacTableName()
	}
	return xfield
}

func (xf *ZeroXsacField) Metatype() reflect.Type {
	return xf.metatype
}

func (xf *ZeroXsacField) Xapi() string {
	return xf.xapi
}

func (xf *ZeroXsacField) Xjsonopts() string {
	return xf.jsonopts
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

func (xf *ZeroXsacField) ChildColumnName() string {
	return xf.childColumnName
}

func (xf *ZeroXsacField) SubTableName() string {
	return xf.subTableName
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

type ZeroXsacApiDeclares interface {
	XsacApiName() string
	XsacApiFields() [][]string
	XsacApiEnums() []string
	XsacApis(...string) []string

	XsacApiExports(...string) []string
}

type ZeroXsacAutoParser interface {
	Parse(map[string]any, any) error
}

type xZeroXsacAutoParser struct {
	ColumnName string
	FieldName  string
}

func NewAutoParser(columnName string, fieldName string) ZeroXsacAutoParser {
	return &xZeroXsacAutoParser{
		ColumnName: columnName,
		FieldName:  fieldName,
	}
}

func (autoParser *xZeroXsacAutoParser) intValue(row map[string]any, data reflect.Value) error {
	_, ok := row[autoParser.ColumnName]
	if ok {
		data.Elem().FieldByName(autoParser.FieldName).SetInt(int64(ParseIntField(row, autoParser.ColumnName)))
	}
	return nil
}

func (autoParser *xZeroXsacAutoParser) uintValue(row map[string]any, data reflect.Value) error {
	_, ok := row[autoParser.ColumnName]
	if ok {
		data.Elem().FieldByName(autoParser.FieldName).SetUint(uint64(ParseIntField(row, autoParser.ColumnName)))
	}
	return nil
}

func (autoParser *xZeroXsacAutoParser) floatValue(row map[string]any, data reflect.Value) error {
	_, ok := row[autoParser.ColumnName]
	if ok {
		data.Elem().FieldByName(autoParser.FieldName).SetFloat(float64(ParseFloatField(row, autoParser.ColumnName)))
	}
	return nil
}

func (autoParser *xZeroXsacAutoParser) stringValue(row map[string]any, data reflect.Value) error {
	_, ok := row[autoParser.ColumnName]
	if ok {
		data.Elem().FieldByName(autoParser.FieldName).SetString(ParseStringField(row, autoParser.ColumnName))
	}
	return nil
}

func (autoParser *xZeroXsacAutoParser) ptrValue(row map[string]any, data reflect.Value, field reflect.StructField) error {
	_, ok := row[autoParser.ColumnName]
	if ok {
		fieldtype := field.Type
		if fieldtype.Kind() == reflect.Pointer {
			fieldtype = fieldtype.Elem()
		}

		if fieldtype.String() == reflect.TypeOf(Time{}).String() {
			tm := Time(row[autoParser.ColumnName].(time.Time))
			if field.Type.Kind() == reflect.Pointer {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(&tm))
			} else {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(tm))
			}
		} else if fieldtype.String() == reflect.TypeOf(time.Time{}).String() {
			tm := row[autoParser.ColumnName].(time.Time)
			if field.Type.Kind() == reflect.Pointer {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(&tm))
			} else {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(tm))
			}
		} else {
			contents := ParseStringField(row, autoParser.ColumnName)

			newstruct := reflect.New(fieldtype).Interface()
			err := json.Unmarshal([]byte(contents), newstruct)
			if err != nil {
				return err
			}
			if field.Type.Kind() == reflect.Pointer {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(newstruct))
			} else {
				data.Elem().FieldByName(autoParser.FieldName).Set(reflect.ValueOf(newstruct).Elem())
			}
		}
	}
	return nil
}

func (autoParser *xZeroXsacAutoParser) Parse(row map[string]any, data any) error {
	datarf := reflect.ValueOf(data)
	if datarf.Kind() != reflect.Pointer {
		return fmt.Errorf(" data need ptr type ")
	}

	field, ok := reflect.TypeOf(data).Elem().FieldByName(autoParser.FieldName)
	if !ok {
		return fmt.Errorf(" field `%s` not found ", autoParser.FieldName)
	}

	switch field.Type.String() {
	case reflect.Int.String():
		fallthrough
	case reflect.Int8.String():
		fallthrough
	case reflect.Int16.String():
		fallthrough
	case reflect.Int32.String():
		fallthrough
	case reflect.Int64.String():
		return autoParser.intValue(row, datarf)

	case reflect.Uint.String():
		fallthrough
	case reflect.Uint8.String():
		fallthrough
	case reflect.Uint16.String():
		fallthrough
	case reflect.Uint32.String():
		fallthrough
	case reflect.Uint64.String():
		return autoParser.uintValue(row, datarf)

	case reflect.Float32.String():
		fallthrough
	case reflect.Float64.String():
		return autoParser.floatValue(row, datarf)

	case reflect.String.String():
		return autoParser.stringValue(row, datarf)
	default:
		return autoParser.ptrValue(row, datarf, field)
	}
}

const XSAC_AUTO_PARSER_KEEPER = "XsacAutoParserKeeper"

type ZeroXsacAutoParserKeeper interface {
	FindAutoParser(string) ([]ZeroXsacAutoParser, bool)
}

func XautoLoad(meta reflect.Type, row map[string]any) (any, error) {
	data := reflect.New(meta)
	atParserKeeper := global.Value(XSAC_AUTO_PARSER_KEEPER)
	if atParserKeeper != nil {
		parsers, ok := atParserKeeper.(ZeroXsacAutoParserKeeper).FindAutoParser(data.Interface().(ZeroXsacDeclares).XsacTableName())
		if ok {
			for _, parser := range parsers {
				err := parser.Parse(row, data.Interface())
				if err != nil {
					return nil, err
				}
			}
			return data.Interface(), nil
		}
	}

	returnValues := data.MethodByName("LoadRowData").Call([]reflect.Value{reflect.ValueOf(row)})
	if len(returnValues) > 0 && returnValues[0].Interface() != nil {
		return nil, returnValues[0].Interface().(error)
	}
	return data.Interface(), nil
}
