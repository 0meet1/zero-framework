package mfgrc

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

var XmonoComplete = func(xRequest *structs.ZeroRequest, keeper *ZeroMfgrcKeeper, mono MfgrcMono) error {
	if len(xRequest.Querys) != 1 {
		return errors.New(" no support multiple tasks or task is empty ")
	}

	bytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, mono)
	if err != nil {
		return err
	}
	err = keeper.Check(mono)
	if err != nil {
		return err
	}
	mono.ThisDef(mono)
	return keeper.Check(mono)
}

var XmonoPerformed = func(
	writer http.ResponseWriter,
	xRequest *structs.ZeroRequest,
	keeper *ZeroMfgrcKeeper,
	mono MfgrcMono,
	onReady func(map[string]any),
	onSuccess func() any,
	onFailed func() string) {

	expands := make(map[string]interface{})
	if onReady != nil {
		onReady(expands)
	} else {
		expands["monoId"] = mono.XmonoId()
	}

	waittime, ok := xRequest.Expands["waittime"]
	if ok && waittime.(float64) > 0 {
		act := &ZeroMfgrcMonoActuator{Keeper: keeper}
		select {
		case err := <-act.Exec(mono):
			if err != nil {
				expands["state"] = "error"
				if onFailed != nil {
					_err := onFailed()
					if _err != "" {
						err = errors.New(_err)
					}
				}
				server.XhttpResponseMessages(writer, 500, err.Error())
				return
			} else {
				expands["state"] = "success"
				if onSuccess != nil {
					result := onSuccess()
					if result != nil {
						expands["result"] = result
					}
				}

			}
		case <-time.After(time.Second * time.Duration(waittime.(float64))):
			expands["state"] = "timeout"
		}
	} else {
		go keeper.AddMono(mono)
		expands["state"] = "created"
	}
	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}

var XgroupComplete = func(xRequest *structs.ZeroRequest, keeper *ZeroMfgrcGroupKeeper, group MfgrcGroup) error {
	if len(xRequest.Querys) != 1 {
		return errors.New(" no support multiple tasks or task is empty ")
	}
	bytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, group)
	if err != nil {
		return err
	}

	err = keeper.Check(group)
	if err != nil {
		return err
	}
	group.ThisDef(&group)
	group.Xmonos()
	return nil
}

var XgroupPerformed = func(
	writer http.ResponseWriter,
	xRequest *structs.ZeroRequest,
	keeper *ZeroMfgrcGroupKeeper,
	group MfgrcGroup,
	onReady func(map[string]any),
	onSuccess func() any,
	onFailed func(map[string]any) string) {

	expands := make(map[string]interface{})
	if onReady != nil {
		onReady(expands)
	} else {
		expands["groupId"] = group.XgroupId()
	}

	waittime, ok := xRequest.Expands["waittime"]
	if ok && waittime.(float64) > 0 {
		act := &ZeroMfgrcGroupActuator{Keeper: keeper}
		select {
		case err := <-act.Exec(group):
			if err != nil {
				expands["state"] = "error"
				if onFailed != nil {
					_err := onFailed(expands)
					if _err != "" {
						err = errors.New(_err)
					}
				}
				server.XhttpResponseMessages(writer, 500, err.Error())
				return
			} else {
				expands["state"] = "success"
				if onSuccess != nil {
					result := onSuccess()
					if result != nil {
						expands["result"] = result
					}
				}

			}
		case <-time.After(time.Second * time.Duration(waittime.(float64))):
			expands["state"] = "timeout"
		}
	} else {
		go keeper.AddGroup(group)
		expands["state"] = "created"
	}
	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}
