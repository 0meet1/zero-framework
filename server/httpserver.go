package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/structs"
)

var (
	prefix string
)

func mkuri(args ...string) string {
	uri := path.Join("/", prefix)
	for _, arg := range args {
		uri = path.Join(uri, arg)
	}
	return uri
}

func responseMaps(writer http.ResponseWriter, code int, message string, datas []map[string]interface{}, expands map[string]interface{}) {
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
	writer.Write(bytes)
}

func responseDatas(writer http.ResponseWriter, code int, message string, datas []interface{}, expands map[string]interface{}) {
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
	writer.Write(bytes)
}

func responseMessages(writer http.ResponseWriter, code int, message string) {
	responseDatas(writer, code, message, nil, nil)
}

type xHttpHandler struct {
	path  string
	funcx func(http.ResponseWriter, *http.Request)
}

func HttpHandler(funcx func(http.ResponseWriter, *http.Request), path ...string) *xHttpHandler {
	return &xHttpHandler{
		path:  mkuri(path...),
		funcx: funcx,
	}
}

func RunHttpServer(handlers ...*xHttpHandler) {
	prefix = global.StringValue("zero.server.prefix")
	if len(prefix) == 0 {
		prefix = "/zeroapi"
	}
	global.Logger().Info(fmt.Sprintf("http server start on : http://%s:%d", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port")))
	server := http.Server{Addr: fmt.Sprintf("%s:%d", global.StringValue("zero.httpserver.hostname"), global.IntValue("zero.httpserver.port"))}
	for _, handler := range handlers {
		http.HandleFunc(handler.path, handler.funcx)
	}
	server.ListenAndServe()
}
