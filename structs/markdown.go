package structs

import (
	"fmt"
	"strings"
)

type ZeroMarkdown struct {
	rows []string
}

func NewMarkdown(rows ...string) *ZeroMarkdown {
	if rows == nil {
		return &ZeroMarkdown{
			rows: []string{},
		}
	}
	return &ZeroMarkdown{
		rows: rows,
	}
}

func (md *ZeroMarkdown) AddRows(row ...string) {
	md.rows = append(md.rows, row...)
}

func (md *ZeroMarkdown) Rows() []string {
	return md.rows
}

func (md *ZeroMarkdown) Clear() {
	md.rows = make([]string, 0)
}

func (md *ZeroMarkdown) String() string {
	return strings.Join(md.rows, "\n")
}

func (md *ZeroMarkdown) HTML() string {
	return XmarkdownTemplete(MACDOWN_CSS, md.String())
}

func NewApiHeader(xHeaderName string, xVersion string) string {
	return fmt.Sprintf("# __%s `%s`__", xHeaderName, xVersion)
}

func NewApiContentHeader(xContentHeader string) string {
	return fmt.Sprintf("## __%s__", xContentHeader)
}

func NewApiLittleHeader(xLittleHeader string) string {
	return fmt.Sprintf("### %s", xLittleHeader)
}

func NewApiLable(xLableText string) string {
	return fmt.Sprintf("##### %s", xLableText)
}

func NewApiCode(xCodeText string) string {
	return fmt.Sprintf("`%s`\n", xCodeText)
}

func NewApiMultipleCode(xCodeText string) string {
	return fmt.Sprintf("```\n%s\n```", xCodeText)
}

func NewApiTable(heads []string, rows [][]string) []string {
	table := make([]string, 0)
	table = append(table, fmt.Sprintf("| %s |", strings.Join(heads, " | ")))
	table = append(table, fmt.Sprintf("|%s", strings.Repeat(" :-: |", len(heads))))
	for _, row := range rows {
		table = append(table, fmt.Sprintf("| %s |", strings.Join(row, " | ")))
	}
	return table
}

func NewApiDataMod(name string, rows [][]string) []string {
	datamod := make([]string, 0)
	datamod = append(datamod, NewApiLittleHeader(fmt.Sprintf("%s：", name)))
	datamod = append(datamod, NewApiTable([]string{"参数名称", "参数类型", "参数描述", "可写入", "可更新", "备注"}, rows)...)
	return datamod
}

func ApiDataMods(fields ...string) [][]string {
	if len(fields) == 0 {
		return make([][]string, 0)
	}
	rows := make([][]string, 0)
	for i := 0; i < len(fields); i += 6 {
		if len(fields[i+5]) > 0 {
			rows = append(rows, []string{fields[i], fmt.Sprintf("`%s`", fields[i+1]), fields[i+2], fields[i+3], fields[i+4], fmt.Sprintf("`%s`", fields[i+5])})
		} else {
			rows = append(rows, []string{fields[i], fmt.Sprintf("`%s`", fields[i+1]), fields[i+2], fields[i+3], fields[i+4], fields[i+5]})
		}
	}
	return rows
}

func NewApiEnums(name string, rows [][]string) []string {
	enums := make([]string, 0)
	enums = append(enums, NewApiLittleHeader(fmt.Sprintf("%s：", name)))
	datas := make([][]string, 0)
	for _, row := range rows {
		datas = append(datas, []string{fmt.Sprintf("`%s`", row[0]), fmt.Sprintf("`%s`", row[1])})
	}
	enums = append(enums, NewApiTable([]string{"值", "描述"}, datas)...)
	return enums
}

func ApiEnums(fields ...string) [][]string {
	if len(fields) == 0 {
		return make([][]string, 0)
	}
	rows := make([][]string, 0)
	for i := 0; i < len(fields); i += 2 {
		rows = append(rows, []string{fields[i], fields[i+1]})
	}
	return rows
}

func NewApiOptions(rows [][]string) []string {
	options := make([]string, 0)
	options = append(options, NewApiLable("Options参数说明:"))
	options = append(options, NewApiCode("* options可多选用 | 分隔 例如: \"baisc|option1\""))
	if len(rows) <= 0 {
		options = append(options, NewApiTable([]string{"options值", "options值说明"}, [][]string{{"-", "-"}})...)
	} else {
		options = append(options, NewApiTable([]string{"options值", "options值说明"}, rows)...)
	}
	return options
}

func NewApiExpands(rows [][]string) []string {
	expands := make([]string, 0)
	expands = append(expands, NewApiLable("扩展参数说明："))
	if len(rows) <= 0 {
		expands = append(expands, NewApiTable([]string{"扩展参数", "扩展参数说明"}, [][]string{{"-", "-"}})...)
	} else {
		expands = append(expands, NewApiTable([]string{"扩展参数", "扩展参数说明"}, rows)...)
	}
	return expands
}

func NewApiRequest(body string) []string {
	req := make([]string, 0)
	req = append(req, NewApiLable("请求示例："))
	req = append(req, NewApiMultipleCode(body))
	return req
}

func NewApiResponse(body string) []string {
	resp := make([]string, 0)
	resp = append(resp, NewApiLable("响应体："))
	resp = append(resp, NewApiMultipleCode(body))
	return resp
}

func NewApiContent(header, uri, reqbody, respbody string, options, expands [][]string) []string {
	rows := make([]string, 0)
	rows = append(rows, NewApiLittleHeader(header))
	rows = append(rows, NewApiMultipleCode(uri))
	rows = append(rows, NewApiOptions(options)...)
	rows = append(rows, NewApiExpands(expands)...)
	rows = append(rows, NewApiRequest(reqbody)...)
	rows = append(rows, NewApiResponse(respbody)...)
	return rows
}

func NewApiContentNOE(header, uri, reqbody, respbody string) []string {
	rows := make([]string, 0)
	rows = append(rows, NewApiLittleHeader(header))
	rows = append(rows, NewApiMultipleCode(uri))
	rows = append(rows, NewApiOptions(make([][]string, 0))...)
	rows = append(rows, NewApiExpands(make([][]string, 0))...)
	rows = append(rows, NewApiRequest(reqbody)...)
	rows = append(rows, NewApiResponse(respbody)...)
	return rows
}

func NewApiContentNE(header, uri, reqbody, respbody string, options [][]string) []string {
	rows := make([]string, 0)
	rows = append(rows, NewApiLittleHeader(header))
	rows = append(rows, NewApiMultipleCode(uri))
	rows = append(rows, NewApiOptions(options)...)
	rows = append(rows, NewApiExpands(make([][]string, 0))...)
	rows = append(rows, NewApiRequest(reqbody)...)
	rows = append(rows, NewApiResponse(respbody)...)
	return rows
}

func NewApiContentNO(header, uri, reqbody, respbody string, expands [][]string) []string {
	rows := make([]string, 0)
	rows = append(rows, NewApiLittleHeader(header))
	rows = append(rows, NewApiMultipleCode(uri))
	rows = append(rows, NewApiOptions(make([][]string, 0))...)
	rows = append(rows, NewApiExpands(expands)...)
	rows = append(rows, NewApiRequest(reqbody)...)
	rows = append(rows, NewApiResponse(respbody)...)
	return rows
}
