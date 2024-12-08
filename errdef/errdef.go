package errdef

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0meet1/zero-framework/autohttpconf"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroExceptionDef struct {
	autohttpconf.ZeroXsacXhttpStructs

	Code        string                 `json:"code,omitempty" xhttpopt:"OX" xsacprop:"NO,VARCHAR(64),NULL" xsackey:"unique" xapi:"异常编号,String"`
	Description string                 `json:"description,omitempty" xhttpopt:"OX" xsacprop:"NO,VARCHAR(1024),NULL" xapi:"异常描述,String"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

func (*ZeroExceptionDef) XhttpPath() string     { return "errdef" }
func (*ZeroExceptionDef) XsacTableName() string { return "zero_exception_def" }
func (*ZeroExceptionDef) XhttpOpt() byte        { return 0x01 }
func (*ZeroExceptionDef) XsacApiName() string   { return "异常定义" }
func (*ZeroExceptionDef) XhttpAutoProc() processors.ZeroXsacAutoProcessor {
	return global.Value(EXCEPTION_AUTO_PROC).(processors.ZeroXsacAutoProcessor)
}

func (*ZeroExceptionDef) XhttpQueryOperation() processors.ZeroQueryOperation {
	return global.Value(EXCEPTION_OPERATION).(processors.ZeroQueryOperation)
}
func (errdef *ZeroExceptionDef) LoadRowData(rowmap map[string]interface{}) {
	errdef.ZeroXsacXhttpStructs.LoadRowData(rowmap)

	errdef.Code = structs.ParseStringField(rowmap, "code")
	errdef.Description = structs.ParseStringField(rowmap, "description")
}

func (errdef *ZeroExceptionDef) String() string {
	mjson, _ := json.Marshal(errdef)
	return string(mjson)
}

func (errdef *ZeroExceptionDef) transfer(parameterName string) string {
	return fmt.Sprintf("#{{%s}}", parameterName)
}

func (errdef *ZeroExceptionDef) Error() string {
	errdes := errdef.Code
	if errdef.Description == "" && errdef.Parameters != nil {
		exceptionKeeper := global.Value(EXCEPTION_KEEPER)
		if exceptionKeeper != nil {
			errdef.Description = exceptionKeeper.(ZeroExceptionKeeper).DescriptionByCode(errdef.Code)
		}
	}
	if errdef.Description != "" {
		if errdef.Parameters != nil {
			description := errdef.Description
			for parameterName, parameterValue := range errdef.Parameters {
				description = strings.ReplaceAll(description, errdef.transfer(parameterName), parameterValue.(string))
			}
			errdes = fmt.Sprintf("%s: %s", errdes, description)
		} else {
			errdes = fmt.Sprintf("%s: %s", errdes, errdef.Description)
		}
	}
	return errdes
}
