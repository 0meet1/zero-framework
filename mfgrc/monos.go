package mfgrc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/structs"
)

const (
	WORKER_MONO_STATUS_READY     = "mono.status.ready"
	WORKER_MONO_STATUS_PENDING   = "mono.status.pending"
	WORKER_MONO_STATUS_EXECUTING = "mono.status.executing"
	WORKER_MONO_STATUS_RETRYING  = "mono.status.retrying"
	WORKER_MONO_STATUS_COMPLETE  = "mono.status.complete"
	WORKER_MONO_STATUS_FAILED    = "mono.status.failed"
	WORKER_MONO_STATUS_REVOKE    = "mono.status.revoke"
	WORKER_MONO_STATUS_TIMEOUT   = "mono.status.timeout"
)

type ZeroMfgrcMono struct {
	structs.ZeroCoreStructs

	MonoID     string `json:"monoID,omitempty"`
	UniqueCode string `json:"uniqueCode,omitempty"`
	Option     string `json:"option,omitempty"`

	Progress int `json:"progress,omitempty"`

	status          string
	reason          string
	maxExecuteTimes int
	executeTimes    int
	response        interface{}

	xStore    ZeroMfgrcMonoStore
	xListener ZeroMfgrcMonoEventListener
	xProgress ZeroMfgrcMonoProgressListener
	fromFlux  *ZeroMfgrcFlux
}

func (mono *ZeroMfgrcMono) LoadRowData(rowmap map[string]interface{}) {
	mono.ZeroCoreStructs.LoadRowData(rowmap)

	_, ok := rowmap["mono_id"]
	if ok {
		mono.MonoID = mono.UInt8ToString(rowmap["mono_id"].([]uint8))
	}

	_, ok = rowmap["unique_code"]
	if ok {
		mono.UniqueCode = mono.UInt8ToString(rowmap["unique_code"].([]uint8))
	}

	_, ok = rowmap["option"]
	if ok {
		mono.Option = mono.UInt8ToString(rowmap["option"].([]uint8))
	}

	_, ok = rowmap["progress"]
	if ok {
		mono.Progress = int(rowmap["progress"].(int64))
	}

	_, ok = rowmap["status"]
	if ok {
		mono.status = mono.UInt8ToString(rowmap["status"].([]uint8))
	}

	_, ok = rowmap["execute_times"]
	if ok {
		mono.executeTimes = int(rowmap["execute_times"].(int64))
	}

	reason, ok := mono.Features["reason"]
	if ok {
		mono.reason = reason.(string)
	}
}

func (mono *ZeroMfgrcMono) XmonoId() string {
	return mono.MonoID
}

func (mono *ZeroMfgrcMono) XuniqueCode() string {
	return mono.UniqueCode
}

func (mono *ZeroMfgrcMono) Xoption() string {
	return mono.Option
}

func (mono *ZeroMfgrcMono) Xprogress() int {
	return mono.Progress
}

func (mono *ZeroMfgrcMono) State() string {
	return mono.status
}

func (mono *ZeroMfgrcMono) FromFlux() *ZeroMfgrcFlux {
	return mono.fromFlux
}

func (mono *ZeroMfgrcMono) Response() interface{} {
	return mono.response
}

func (mono *ZeroMfgrcMono) Store(store ZeroMfgrcMonoStore) {
	mono.xStore = store
}

func (mono *ZeroMfgrcMono) EventListener(eventListener ZeroMfgrcMonoEventListener) {
	mono.xListener = eventListener
}

func (mono *ZeroMfgrcMono) ProgressListener(progressListener ZeroMfgrcMonoProgressListener) {
	mono.xProgress = progressListener
}

func (mono *ZeroMfgrcMono) Ready(store ...ZeroMfgrcMonoStore) error {
	mono.status = WORKER_MONO_STATUS_READY

	mono.maxExecuteTimes = mono.fromFlux.worker.keeper.taskRetryTimes + 1
	mono.executeTimes = 0
	mono.reason = ""
	if store != nil && len(store) > 0 {
		mono.xStore = store[0]
	}
	if mono.xStore != nil {
		return mono.xStore.UpdateMono(mono)
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` on ready", mono.MonoID))
	return nil
}

func (mono *ZeroMfgrcMono) Pending(flux *ZeroMfgrcFlux) error {
	if mono.status != WORKER_MONO_STATUS_READY {
		return errors.New(fmt.Sprintf("could not pending mono `%s` status `%s`", mono.MonoID, mono.status))
	}
	mono.fromFlux = flux
	mono.status = WORKER_MONO_STATUS_PENDING
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnPending(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is pending in flux `%s`", mono.MonoID, mono.fromFlux.UniqueId))
	return nil
}

func (mono *ZeroMfgrcMono) Revoke() error {
	mono.status = WORKER_MONO_STATUS_REVOKE

	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRevoke(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is revoke ", mono.MonoID))
	return nil
}

func (mono *ZeroMfgrcMono) Timeout() error {
	mono.status = WORKER_MONO_STATUS_TIMEOUT

	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRevoke(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is timeout", mono.MonoID))
	return nil
}

func (mono *ZeroMfgrcMono) Executing() error {
	if mono.status != WORKER_MONO_STATUS_PENDING {
		return errors.New(fmt.Sprintf("could not executing mono `%s` status `%s`", mono.MonoID, mono.status))
	}
	if mono.executeTimes != 0 {
		return errors.New(fmt.Sprintf("mono `%s` is already executing", mono.MonoID))
	}
	mono.status = WORKER_MONO_STATUS_EXECUTING
	mono.executeTimes++
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnExecuting(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is executing in worker [%s] >> flux `%s`, option: %s", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId, mono.Option))
	return nil
}

func (mono *ZeroMfgrcMono) Retrying(reason string) error {
	if mono.status != WORKER_MONO_STATUS_EXECUTING && mono.status != WORKER_MONO_STATUS_RETRYING {
		return errors.New(fmt.Sprintf("could not retrying mono `%s` status `%s`", mono.MonoID, mono.status))
	}
	if mono.executeTimes >= mono.maxExecuteTimes {
		return errors.New(fmt.Sprintf("exceeded maximum attempts, maxExecuteTimes:%d executeTimes:%d", mono.maxExecuteTimes, mono.executeTimes))
	}
	mono.reason = reason
	mono.status = WORKER_MONO_STATUS_RETRYING
	mono.executeTimes++
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRetrying(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is retrying in worker [%s] >> flux `%s`, option: %s maxExecuteTimes:%d executeTimes:%d", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId, mono.Option, mono.maxExecuteTimes, mono.executeTimes))
	return nil
}

func (mono *ZeroMfgrcMono) Complete() error {
	if mono.status != WORKER_MONO_STATUS_EXECUTING && mono.status != WORKER_MONO_STATUS_RETRYING {
		return errors.New(fmt.Sprintf("could not complete mono `%s` status `%s`", mono.MonoID, mono.status))
	}
	mono.status = WORKER_MONO_STATUS_COMPLETE
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnComplete(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is complete in worker [%s] >> flux `%s`", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId))
	return nil
}

func (mono *ZeroMfgrcMono) Failed(reason string) error {
	mono.reason = reason
	mono.status = WORKER_MONO_STATUS_FAILED
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnFailed(mono, reason)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is failed, reason: %s", mono.MonoID, reason))
	return nil
}

func (mono *ZeroMfgrcMono) Delete() error {
	if mono.xStore != nil {
		return mono.xStore.DeleteMono(mono)
	}
	return nil
}

func (mono *ZeroMfgrcMono) Do() error {
	return errors.New(fmt.Sprintf("mono `%s` unknow option `%s`", mono.MonoID, mono.Option))
}

func (mono *ZeroMfgrcMono) Export() (map[string]interface{}, error) {
	mjson, err := json.Marshal(mono)
	if err != nil {
		return nil, err
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(mjson), &jsonMap)
	if err != nil {
		return nil, err
	}

	jsonMap["status"] = mono.status
	jsonMap["reason"] = mono.reason
	jsonMap["maxExecuteTimes"] = mono.maxExecuteTimes
	jsonMap["executeTimes"] = mono.executeTimes
	if mono.response != nil {
		jsonMap["response"] = mono.response
	}

	return jsonMap, nil
}

func (mono *ZeroMfgrcMono) String() (string, error) {
	jsonMap, err := mono.Export()
	if err != nil {
		return "", err
	}
	mjson, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}
	return string(mjson), nil
}
