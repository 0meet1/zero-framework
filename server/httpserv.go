package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/0meet1/zero-framework/global"
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

	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
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

	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
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
