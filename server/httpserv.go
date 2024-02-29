package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/structs"
)

var (
	prefix string
)

func mkuri(args ...string) string {
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

func XhttpResponseMaps(writer http.ResponseWriter, code int, message string, datas []map[string]interface{}, expands map[string]interface{}) {
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

func XhttpResponseDatas(writer http.ResponseWriter, code int, message string, datas []interface{}, expands map[string]interface{}) {
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

func XhttpResponseMessages(writer http.ResponseWriter, code int, message string) {
	XhttpResponseDatas(writer, code, message, nil, nil)
}

func XhttpZeroRequest(req *http.Request) (*structs.ZeroRequest, error) {
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

func XhttpZeroQuery(xRequest *structs.ZeroRequest) (*processors.ZeroQuery, error) {
	if xRequest.Querys == nil || len(xRequest.Querys) <= 0 {
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

func XhttpMysqlQueryOperation(xRequest *structs.ZeroRequest, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	return processors.NewZeroMysqlQueryOperation(xQuery, tableName), xQuery, nil
}

func XhttpPostgresQueryOperation(xRequest *structs.ZeroRequest, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	return processors.NewZeroPostgresQueryOperation(xQuery, tableName), xQuery, nil
}

func XhttpCompleteQueryOperation(xRequest *structs.ZeroRequest, xProcessor processors.ZeroQueryOperation, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	xProcessor.AddQuery(xQuery)
	xProcessor.AddTableName(tableName)
	return xProcessor, xQuery, nil
}

func XhttpEQuery(xRequest *structs.ZeroRequest) (*database.EQuerySearch, error) {
	if xRequest.Querys == nil || len(xRequest.Querys) <= 0 {
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

func XhttpEQueryRequest(xRequest *structs.ZeroRequest, indexName string) (*database.EQueryRequest, *database.EQuerySearch, error) {
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

type XhttpExecutor struct {
	funcmode     bool
	path         string
	executorfunc func(http.ResponseWriter, *http.Request)
	executor     http.Handler
}

func XhttpFuncHandle(funcx func(http.ResponseWriter, *http.Request), path ...string) *XhttpExecutor {
	return &XhttpExecutor{
		funcmode:     true,
		path:         mkuri(path...),
		executorfunc: funcx,
	}
}

func XhttpHandle(handler http.Handler, path ...string) *XhttpExecutor {
	return &XhttpExecutor{
		funcmode: false,
		path:     mkuri(path...),
		executor: handler,
	}
}

func RunHttpServer(handlers ...*XhttpExecutor) {
	prefix = path.Join("/", global.StringValue("zero.httpserver.prefix"))
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
