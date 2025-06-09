package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"
	"strings"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

var XhttpResponseMaps = func(writer http.ResponseWriter, code int, message string, datas []map[string]interface{}, expands map[string]interface{}) {
	response := make(map[string]interface{})
	response["code"] = code
	response["message"] = message
	response["datas"] = datas
	response["expands"] = expands

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bytes)
}

var XhttpResponseDatas = func(writer http.ResponseWriter, code int, message string, datas []interface{}, expands map[string]interface{}) {
	bytes, err := json.Marshal(structs.ZeroResponse{
		Code:    code,
		Message: message,
		Datas:   datas,
		Expands: expands,
	})
	if err != nil {
		panic(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bytes)
}

var XhttpResponseMessages = func(writer http.ResponseWriter, code int, message string) {
	XhttpResponseDatas(writer, code, message, nil, nil)
}

var XhttpZeroRequest = func(req *http.Request) (*structs.ZeroRequest, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var request structs.ZeroRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

var XhttpZeroQuery = func(xRequest *structs.ZeroRequest) (*processors.ZeroQuery, error) {
	if len(xRequest.Querys) <= 0 {
		return nil, errors.New("missing necessary parameter `query[0]`")
	}

	bytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return nil, err
	}

	var query processors.ZeroQuery
	err = json.Unmarshal(bytes, &query)
	if err != nil {
		return nil, err
	}

	return &query, nil
}

var XhttpMysqlQueryOperation = func(xRequest *structs.ZeroRequest, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	return processors.NewZeroMysqlQueryOperation(xQuery, tableName), xQuery, nil
}

var XhttpPostgresQueryOperation = func(xRequest *structs.ZeroRequest, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	return processors.NewZeroPostgresQueryOperation(xQuery, tableName), xQuery, nil
}

var XhttpCompleteQueryOperation = func(xRequest *structs.ZeroRequest, xProcessor processors.ZeroQueryOperation, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	xProcessor.AddQuery(xQuery)
	xProcessor.AddTableName(tableName)
	return xProcessor, xQuery, nil
}

const XHTTP_QUERY_OPTIONS_ALL = "all"

var XhttpQueryOptions = func(xRequest *structs.ZeroRequest) []string {
	xoptions := make([]string, 0)
	if xRequest.Expands == nil {
		return xoptions
	}
	if _, ok := xRequest.Expands["options"]; ok {
		xoptionItems := strings.Split(xRequest.Expands["options"].(string), "|")
		for _, xoption := range xoptionItems {
			if xoption == XHTTP_QUERY_OPTIONS_ALL {
				return []string{XHTTP_QUERY_OPTIONS_ALL}
			}
			xoptions = append(xoptions, strings.ToLower(xoption))
		}
	}
	return xoptions
}

var XhttpContainsOptions = func(xRequest *structs.ZeroRequest, option string) bool {
	if _, ok := xRequest.Expands["options"]; ok {
		return strings.Contains(xRequest.Expands["options"].(string), option) ||
			strings.Contains(xRequest.Expands["options"].(string), XHTTP_QUERY_OPTIONS_ALL)
	}
	return false
}

var XhttpEQuery = func(xRequest *structs.ZeroRequest) (*database.EQuerySearch, error) {
	if len(xRequest.Querys) <= 0 {
		return nil, errors.New("missing necessary parameter `query[0]`")
	}

	bytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return nil, err
	}

	var query database.EQuerySearch
	err = json.Unmarshal(bytes, &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

var XhttpEQueryRequest = func(xRequest *structs.ZeroRequest, indexName string) (*database.EQueryRequest, *database.EQuerySearch, error) {
	xEQuery, err := XhttpEQuery(xRequest)
	if err != nil {
		return nil, nil, err
	}
	if xEQuery.Size > 1000 {
		xEQuery.Size = 1000
	}
	eRequest := &database.EQueryRequest{Query: xEQuery}
	eRequest.InitIndex(indexName)
	return eRequest, xEQuery, nil
}

var XhttpURIParams = func(req *http.Request, xPattern string) map[string]string {
	uriparams := make(map[string]string)
	xAfter := xPattern[:strings.Index(xPattern, ":")]
	if strings.Index(req.URL.Path, xAfter) > 0 {
		xFieldItems := strings.Split(xPattern[strings.Index(xPattern, xAfter)+len(xAfter):], "/")
		xParamsURI := req.URL.Path[strings.Index(req.URL.Path, xAfter)+len(xAfter):]
		xParamsItems := strings.Split(xParamsURI, "/")

		for i, item := range xFieldItems {
			if strings.HasPrefix(item, ":") && len(xParamsItems) > i {
				uriparams[item[1:]] = xParamsItems[i]
			}
		}
	}
	return uriparams
}

type XhttpFromFile struct {
	header     *multipart.FileHeader
	filesbytes []byte
}

func (xfile *XhttpFromFile) MIMEHeader() textproto.MIMEHeader {
	return xfile.header.Header
}

func (xfile *XhttpFromFile) FileName() string {
	return xfile.header.Filename
}

func (xfile *XhttpFromFile) FileSize() int64 {
	return xfile.header.Size
}

func (xfile *XhttpFromFile) FileHeader() *multipart.FileHeader {
	return xfile.header
}

func (xfile *XhttpFromFile) FilesBytes() []byte {
	return xfile.filesbytes
}

var XhttpFromFileRequest = func(req *http.Request, maxmem int64) ([]*XhttpFromFile, error) {
	err := req.ParseMultipartForm(maxmem)
	if err != nil {
		return nil, err
	}

	formfiles := make([]*XhttpFromFile, 0)
	for formName := range req.MultipartForm.File {
		formFile, formFileHeader, err := req.FormFile(formName)
		if err != nil {
			return nil, err
		}
		defer formFile.Close()

		filebytes, err := io.ReadAll(formFile)
		if err != nil {
			return nil, err
		}

		formfiles = append(formfiles, &XhttpFromFile{
			header:     formFileHeader,
			filesbytes: filebytes,
		})
	}
	return formfiles, nil
}

var XhttpKeyValueRequest = func(req *http.Request) map[string]string {
	kv := make(map[string]string)
	if req.URL.Query() != nil {
		for k := range req.URL.Query() {
			kv[k] = req.URL.Query().Get(k)
		}
	}

	if req.PostForm != nil {
		for k := range req.PostForm {
			kv[k] = req.PostFormValue(k)
		}
	}
	return kv
}

type XhttpExecutor struct {
	funcmode     bool
	path         string
	executorfunc func(http.ResponseWriter, *http.Request)
	executor     http.Handler
}

func xhttpuri(args ...string) string {
	uri := ""
	for _, arg := range args {
		if len(strings.TrimSpace(arg)) > 0 {
			uri = path.Join(uri, arg)
		}
	}

	if strings.HasSuffix(args[len(args)-1], "/") {
		return fmt.Sprintf("/%s/", uri)
	} else {
		return fmt.Sprintf("/%s", uri)
	}
}

type XhttpInterceptor interface {
	Registry(*XhttpExecutor) http.Handler
}

var XhttpFuncHandle = func(funcx func(http.ResponseWriter, *http.Request), path ...string) *XhttpExecutor {
	return &XhttpExecutor{
		funcmode:     true,
		path:         xhttpuri(path...),
		executorfunc: funcx,
	}
}

var XhttpHandle = func(handler http.Handler, path ...string) *XhttpExecutor {
	return &XhttpExecutor{
		funcmode: false,
		path:     xhttpuri(path...),
		executor: handler,
	}
}

var XhttpPerform = func(executor *XhttpExecutor, writer http.ResponseWriter, request *http.Request) {
	if executor.funcmode {
		executor.executorfunc(writer, request)
	} else {
		executor.executor.ServeHTTP(writer, request)
	}
}

var RunHttpServer = func(handlers ...*XhttpExecutor) {
	prefix := path.Join("/", global.StringValue("zero.httpserver.prefix"))
	server := http.Server{Addr: fmt.Sprintf("%s:%d", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port"))}
	for _, handler := range handlers {
		if handler.funcmode {
			if strings.HasSuffix(handler.path, "/") {
				http.HandleFunc(fmt.Sprintf("%s/", path.Join(prefix, handler.path)), handler.executorfunc)
				global.Logger().Info(fmt.Sprintf("http server register path : %s", fmt.Sprintf("%s/", path.Join(prefix, handler.path))))
			} else {
				http.HandleFunc(path.Join(prefix, handler.path), handler.executorfunc)
				global.Logger().Info(fmt.Sprintf("http server register path : %s", path.Join(prefix, handler.path)))
			}
		} else {
			if strings.HasSuffix(handler.path, "/") {
				http.Handle(fmt.Sprintf("%s/", path.Join(prefix, handler.path)), handler.executor)
				global.Logger().Info(fmt.Sprintf("http server register path : %s", fmt.Sprintf("%s/", path.Join(prefix, handler.path))))
			} else {
				http.Handle(path.Join(prefix, handler.path), handler.executor)
				global.Logger().Info(fmt.Sprintf("http server register path : %s", path.Join(prefix, handler.path)))
			}
		}
	}
	global.Logger().Info(fmt.Sprintf("http server start on : http://%s:%d%s", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port"), prefix))
	server.ListenAndServe()
}

var RunInterceptor = func(interceptor XhttpInterceptor, handlers ...*XhttpExecutor) {
	prefix := path.Join("/", global.StringValue("zero.httpserver.prefix"))
	server := http.Server{Addr: fmt.Sprintf("%s:%d", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port"))}
	for _, handler := range handlers {
		if strings.HasSuffix(handler.path, "/") {
			http.Handle(fmt.Sprintf("%s/", path.Join(prefix, handler.path)), interceptor.Registry(handler))
			global.Logger().Info(fmt.Sprintf("http server register path : %s", fmt.Sprintf("%s/", path.Join(prefix, handler.path))))
		} else {
			http.Handle(path.Join(prefix, handler.path), interceptor.Registry(handler))
			global.Logger().Info(fmt.Sprintf("http server register path : %s", path.Join(prefix, handler.path)))
		}
	}
	global.Logger().Info(fmt.Sprintf("http server start on : http://%s:%d%s", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port"), prefix))
	server.ListenAndServe()
}
