package mfgrc

import (
	"encoding/json"
	"fmt"

	"github.com/0meet1/zero-framework/errdef"
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

	xStore    ZeroMfgrcMonoStore
	xListener ZeroMfgrcMonoEventListener

	keeper   *ZeroMfgrcKeeper
	fromFlux *ZeroMfgrcFlux
}

func (mono *ZeroMfgrcMono) LoadRowData(rowmap map[string]interface{}) {
	mono.ZeroCoreStructs.LoadRowData(rowmap)

	mono.MonoID = structs.ParseStringField(rowmap, "mono_id")
	mono.UniqueCode = structs.ParseStringField(rowmap, "unique_code")
	mono.Option = structs.ParseStringField(rowmap, "option")
	mono.Progress = structs.ParseIntField(rowmap, "progress")
	mono.status = structs.ParseStringField(rowmap, "status")
	mono.maxExecuteTimes = structs.ParseIntField(rowmap, "max_execute_times")
	mono.executeTimes = structs.ParseIntField(rowmap, "execute_times")

	if mono.Features != nil {
		reason, ok := mono.Features["reason"]
		if ok {
			mono.reason = reason.(string)
		}
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

func (mono *ZeroMfgrcMono) Store(store ZeroMfgrcMonoStore) {
	mono.xStore = store
}

func (mono *ZeroMfgrcMono) EventListener(eventListener ZeroMfgrcMonoEventListener) {
	mono.xListener = eventListener
}

func (mono *ZeroMfgrcMono) Ready(keeper *ZeroMfgrcKeeper, store ...ZeroMfgrcMonoStore) error {
	mono.keeper = keeper
	mono.maxExecuteTimes = mono.keeper.taskRetryTimes + 1
	mono.status = WORKER_MONO_STATUS_READY
	mono.executeTimes = 0
	mono.reason = ""
	if len(store) > 0 {
		mono.xStore = store[0]
	}
	if mono.xStore != nil {
		return mono.xStore.UpdateMono(mono.This().(MfgrcMono))
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` on ready", mono.MonoID))
	return nil
}

func (mono *ZeroMfgrcMono) Pending(flux *ZeroMfgrcFlux) error {
	if mono.status != WORKER_MONO_STATUS_READY {
		return fmt.Errorf("could not pending mono `%s` status `%s`", mono.MonoID, mono.status)
	}
	mono.fromFlux = flux
	mono.status = WORKER_MONO_STATUS_PENDING
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnPending(mono.This().(MfgrcMono))
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
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRevoke(mono.This().(MfgrcMono))
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
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRevoke(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is timeout", mono.MonoID))
	return nil
}

func (mono *ZeroMfgrcMono) Executing() error {
	if mono.status != WORKER_MONO_STATUS_PENDING {
		return fmt.Errorf("could not executing mono `%s` status `%s`", mono.MonoID, mono.status)
	}
	if mono.executeTimes != 0 {
		return fmt.Errorf("mono `%s` is already executing", mono.MonoID)
	}
	mono.status = WORKER_MONO_STATUS_EXECUTING
	mono.executeTimes++
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnExecuting(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is executing in worker [%s] >> flux `%s`, option: %s", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId, mono.Option))
	return nil
}

func (mono *ZeroMfgrcMono) Retrying(reason error) error {
	if mono.status != WORKER_MONO_STATUS_EXECUTING && mono.status != WORKER_MONO_STATUS_RETRYING {
		return fmt.Errorf("could not retrying mono `%s` status `%s`", mono.MonoID, mono.status)
	}
	if mono.executeTimes >= mono.maxExecuteTimes {
		if errdef.Is(reason) {
			mono.Features["errdef"] = reason
		}
		return fmt.Errorf("exceeded maximum attempts, maxExecuteTimes:%d executeTimes:%d at lastest error: %s", mono.maxExecuteTimes, mono.executeTimes, reason)
	}
	mono.reason = reason.Error()
	mono.status = WORKER_MONO_STATUS_RETRYING
	mono.executeTimes++
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnRetrying(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is retrying in worker [%s] >> flux `%s`, option: %s maxExecuteTimes:%d executeTimes:%d", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId, mono.Option, mono.maxExecuteTimes, mono.executeTimes))
	return nil
}

func (mono *ZeroMfgrcMono) Complete() error {
	if mono.status != WORKER_MONO_STATUS_EXECUTING && mono.status != WORKER_MONO_STATUS_RETRYING {
		if mono.status == WORKER_MONO_STATUS_FAILED ||
			mono.status == WORKER_MONO_STATUS_TIMEOUT ||
			mono.status == WORKER_MONO_STATUS_REVOKE {
			return nil
		}
		return fmt.Errorf("could not complete mono `%s` status `%s`", mono.MonoID, mono.status)
	}
	mono.status = WORKER_MONO_STATUS_COMPLETE
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnComplete(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is complete in worker [%s] >> flux `%s`", mono.MonoID, mono.fromFlux.worker.workName, mono.fromFlux.UniqueId))
	return nil
}

func (mono *ZeroMfgrcMono) Failed(reason error) error {
	if errdef.Is(reason) {
		mono.Features["errdef"] = reason
	}
	mono.reason = reason.Error()
	mono.status = WORKER_MONO_STATUS_FAILED
	if mono.xStore != nil {
		err := mono.xStore.UpdateMono(mono.This().(MfgrcMono))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	if mono.xListener != nil {
		err := mono.xListener.OnFailed(mono.This().(MfgrcMono), reason)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	global.Logger().Info(fmt.Sprintf("mono `%s` is failed, reason: %s", mono.MonoID, reason))
	return nil
}

func (mono *ZeroMfgrcMono) Delete() error {
	if mono.xStore != nil {
		return mono.xStore.DeleteMono(mono.This().(MfgrcMono))
	}
	return nil
}

func (mono *ZeroMfgrcMono) Do() error {
	return fmt.Errorf("mono `%s` option `%s` not implement", mono.MonoID, mono.Option)
}

func (mono *ZeroMfgrcMono) MaxExecuteTimes() int {
	return mono.maxExecuteTimes
}

func (mono *ZeroMfgrcMono) Export() (map[string]interface{}, error) {
	mjson, err := json.Marshal(mono.This().(MfgrcMono))
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
