package mfgrc

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/structs"
)

const (
	WORKER_MONOGROUP_STATUS_READY     = "monogroup.status.ready"
	WORKER_MONOGROUP_STATUS_PENDING   = "monogroup.status.pending"
	WORKER_MONOGROUP_STATUS_EXECUTING = "monogroup.status.executing"
	WORKER_MONOGROUP_STATUS_COMPLETE  = "monogroup.status.complete"
	WORKER_MONOGROUP_STATUS_FAILED    = "monogroup.status.failed"
)

type ZeroMfgrcGroup struct {
	structs.ZeroCoreStructs

	worker *ZeroMfgrcGroupWorker

	UniqueCode string `json:"uniqueCode,omitempty"`
	Option     string `json:"option,omitempty"`

	Monos []MfgrcMono `json:"monos,omitempty"`

	status string
	reason string

	xStore    ZeroMfgrcGroupStore
	xListener ZeroMfgrcGroupEventListener
}

func (group *ZeroMfgrcGroup) LoadRowData(rowmap map[string]interface{}) {
	group.ZeroCoreStructs.LoadRowData(rowmap)

	_, ok := rowmap["unique_code"]
	if ok {
		group.UniqueCode = group.UInt8ToString(rowmap["unique_code"].([]uint8))
	}

	_, ok = rowmap["option"]
	if ok {
		group.Option = group.UInt8ToString(rowmap["option"].([]uint8))
	}

	_, ok = rowmap["status"]
	if ok {
		group.status = group.UInt8ToString(rowmap["status"].([]uint8))
	}

	reason, ok := group.Features["reason"]
	if ok {
		group.reason = reason.(string)
	}
}

func (group *ZeroMfgrcGroup) XgroupId() string {
	return group.ID
}

func (group *ZeroMfgrcGroup) XuniqueCode() string {
	return group.UniqueCode
}

func (group *ZeroMfgrcGroup) Xoption() string {
	return group.Option
}

func (group *ZeroMfgrcGroup) Xmonos() []MfgrcMono {
	return group.Monos
}

func (group *ZeroMfgrcGroup) State() string {
	return group.status
}

func (group *ZeroMfgrcGroup) Store(store ZeroMfgrcGroupStore) {
	group.xStore = store
}

func (group *ZeroMfgrcGroup) UseStore() ZeroMfgrcGroupStore {
	return group.xStore
}

func (group *ZeroMfgrcGroup) EventListener(xListener ZeroMfgrcGroupEventListener) {
	group.xListener = xListener
}

func (group *ZeroMfgrcGroup) AddWorker(worker *ZeroMfgrcGroupWorker) {
	group.worker = worker
}

func (group *ZeroMfgrcGroup) Do() error {
	return errors.New("not implemented")
}

func (group *ZeroMfgrcGroup) Ready(store ZeroMfgrcGroupStore) error {
	group.status = WORKER_MONOGROUP_STATUS_READY
	group.reason = ""
	group.xStore = store
	if group.xStore != nil {
		return group.xStore.UpdateGroup(group.This().(MfgrcGroup))
	}
	global.Logger().Info(fmt.Sprintf("mono group `%s` on ready", group.ID))
	return nil
}

func (group *ZeroMfgrcGroup) Pending() error {
	if group.status != WORKER_MONOGROUP_STATUS_READY {
		return errors.New(fmt.Sprintf("could not pending mono group `%s` status `%s`", group.ID, group.status))
	}
	group.status = WORKER_MONOGROUP_STATUS_PENDING
	if group.xStore != nil {
		group.xStore.UpdateGroup(group.This().(MfgrcGroup))
	}
	global.Logger().Info(fmt.Sprintf("group `%s` is pending`", group.ID))
	if group.xListener != nil {
		err := group.xListener.OnPending(group.This().(MfgrcGroup))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	return nil
}

func (group *ZeroMfgrcGroup) Executing() error {
	if group.status != WORKER_MONOGROUP_STATUS_PENDING {
		return errors.New(fmt.Sprintf("could not executing group `%s` status `%s`", group.ID, group.status))
	}

	group.status = WORKER_MONOGROUP_STATUS_EXECUTING

	if group.xStore != nil {
		group.xStore.UpdateGroup(group.This().(MfgrcGroup))
	}

	global.Logger().Info(fmt.Sprintf("group `%s` is executing in worker [%s] , option: %s", group.ID, group.worker.workName, group.Option))

	if group.xListener != nil {
		err := group.xListener.OnExecuting(group.This().(MfgrcGroup))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	return nil
}

func (group *ZeroMfgrcGroup) Complete() error {
	if group.status != WORKER_MONOGROUP_STATUS_EXECUTING {
		return errors.New(fmt.Sprintf("could not complete group `%s` status `%s`", group.ID, group.status))
	}
	group.status = WORKER_MONOGROUP_STATUS_COMPLETE
	if group.xStore != nil {
		group.xStore.UpdateGroup(group.This().(MfgrcGroup))
	}
	global.Logger().Info(fmt.Sprintf("group `%s` is complete in worker [%s]", group.ID, group.worker.workName))
	if group.xListener != nil {
		err := group.xListener.OnComplete(group.This().(MfgrcGroup))
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	return nil
}

func (group *ZeroMfgrcGroup) Failed(reason string) error {
	group.reason = reason
	group.status = WORKER_MONOGROUP_STATUS_FAILED
	if group.xStore != nil {
		group.xStore.UpdateGroup(group.This().(MfgrcGroup))
	}
	if group.worker != nil {
		global.Logger().Info(fmt.Sprintf("group `%s` is failed in worker [%s] , reason: %s", group.ID, group.worker.workName, reason))
	} else {
		global.Logger().Info(fmt.Sprintf("group `%s` is failed in worker [checker] , reason: %s", group.ID, reason))
	}

	if group.xListener != nil {
		err := group.xListener.OnFailed(group.This().(MfgrcGroup), reason)
		if err != nil {
			global.Logger().Error(err.Error())
		}
	}
	return nil
}

func (group *ZeroMfgrcGroup) Delete() error {
	if group.xStore != nil {
		return group.xStore.DeleteGroup(group)
	}
	return nil
}

func (group *ZeroMfgrcGroup) Export() (map[string]interface{}, error) {

	jsonbytes, err := json.Marshal(group.This().(MfgrcGroup))
	if err != nil {
		return nil, err
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonbytes, &jsonMap)
	if err != nil {
		return nil, err
	}

	jsonMap["uniqueCode"] = group.UniqueCode
	jsonMap["option"] = group.Option

	if group.Monos != nil {
		monos := make([]map[string]interface{}, len(group.Monos))
		for i, mono := range group.Monos {
			monos[i], err = mono.Export()
			if err != nil {
				return nil, err
			}
		}
		jsonMap["monos"] = monos
	}

	jsonMap["status"] = group.status
	jsonMap["reason"] = group.reason

	return jsonMap, nil
}

type ZeroMfgrcGroupWorker struct {
	workName string

	status      string
	statusMutex sync.RWMutex

	executing string

	keeper *ZeroMfgrcGroupKeeper
}

func newMfgrcGroupWorker(workName string, keeper *ZeroMfgrcGroupKeeper) *ZeroMfgrcGroupWorker {
	worker := ZeroMfgrcGroupWorker{}
	worker.workName = workName
	worker.keeper = keeper
	return &worker
}

func (worker *ZeroMfgrcGroupWorker) Start() {
	worker.statusMutex.Lock()
	worker.status = xWORKER_STATUS_RUNNING
	worker.statusMutex.Unlock()

	global.Logger().Info(fmt.Sprintf("[%s] workergroup is ready and waiting", worker.workName))
	for xGroup := range worker.keeper.groupChan {
		worker.statusMutex.Lock()
		xstatus := worker.status
		worker.statusMutex.Unlock()
		if xstatus != xWORKER_STATUS_RUNNING {
			break
		}

		if xGroup != nil {
			global.Logger().Info(fmt.Sprintf("[%s] workergroup with group `%s` device `%s`", worker.workName, xGroup.XgroupId(), xGroup.XuniqueCode()))
			worker.executing = xGroup.XuniqueCode()

			xGroup.AddWorker(worker)
			err := xGroup.Do()
			if err != nil {
				xGroup.Failed(err.Error())
			}

			err = xGroup.Complete()
			if err != nil {
				xGroup.Failed(err.Error())
			}
			worker.keeper.closeGroup(xGroup)

			global.Logger().Info(fmt.Sprintf("[%s] workergroup `%s` device `%s` work complete", worker.workName, xGroup.XgroupId(), xGroup.XuniqueCode()))
			worker.executing = ""
		}
	}

	worker.keeper.closeWorker(worker)
	global.Logger().Info(fmt.Sprintf("[%s] warning! workergroup is shutdown now", worker.workName))
}

func (worker *ZeroMfgrcGroupWorker) Stop() {
	worker.statusMutex.Lock()
	worker.status = xWORKER_STATUS_STOPPED
	worker.statusMutex.Unlock()
}

func (worker *ZeroMfgrcGroupWorker) Export() map[string]interface{} {
	worker.statusMutex.RLock()
	defer worker.statusMutex.RUnlock()

	workerMap := make(map[string]interface{})
	workerMap["workName"] = worker.workName
	workerMap["status"] = worker.status
	workerMap["executing"] = worker.executing
	return workerMap
}

type ZeroMfgrcGroupKeeper struct {
	keeperName string

	workerMap   map[string]*ZeroMfgrcGroupWorker
	workerMutex sync.RWMutex

	groupMap   map[string]MfgrcGroup
	groupChan  chan MfgrcGroup
	groupMutex sync.RWMutex

	maxGroupQueues int

	status      string
	statusMutex sync.RWMutex

	keeperOpts ZeroMfgrcGroupKeeperOpts
}

func NewGroupKeeper(keeperName string, keeperOpts ZeroMfgrcGroupKeeperOpts, maxGroupQueues int) *ZeroMfgrcGroupKeeper {
	return &ZeroMfgrcGroupKeeper{
		keeperName:     keeperName,
		workerMap:      make(map[string]*ZeroMfgrcGroupWorker),
		groupMap:       make(map[string]MfgrcGroup),
		groupChan:      make(chan MfgrcGroup, maxGroupQueues),
		maxGroupQueues: maxGroupQueues,
		status:         xKEEPER_STATUS_STOPPED,
		keeperOpts:     keeperOpts,
	}
}

func (keeper *ZeroMfgrcGroupKeeper) RunGroupWorker() {
	global.Logger().Info(fmt.Sprintf("workergroup start with maxGroupQueues: %d", keeper.maxGroupQueues))

	if len(keeper.keeperName) == 0 {
		keeper.keeperName = "default"
	}

	for i := 0; i < keeper.maxGroupQueues; i++ {
		worker := newMfgrcGroupWorker(fmt.Sprintf("%s-group-%03d::", keeper.keeperName, i), keeper)

		keeper.workerMutex.Lock()
		keeper.workerMap[worker.workName] = worker
		keeper.workerMutex.Unlock()
		go worker.Start()
	}
	go keeper.revokeMonoGroups()
}

func (keeper *ZeroMfgrcGroupKeeper) revokeMonoGroups() {
	<-time.After(time.Second * time.Duration(3))

	groups, err := keeper.keeperOpts.FetchUncompleteGroups()
	if err != nil {
		global.Logger().Error(fmt.Sprintf(" fetch uncomplete groups err : %s", err.Error()))
	}

	for _, group := range groups {
		group.Store(keeper.keeperOpts.GroupStore())
		for _, mono := range group.Xmonos() {
			mono.Store(keeper.keeperOpts.MonoStore())
			mono.Revoke()
		}
		group.Failed("group unexpected termination")
	}

	keeper.statusMutex.Lock()
	defer keeper.statusMutex.Unlock()
	keeper.status = xKEEPER_STATUS_RUNNING

	global.Logger().Info(fmt.Sprintf(" workergroup check and resume monos complete "))
}

func (keeper *ZeroMfgrcGroupKeeper) closeWorker(worker *ZeroMfgrcGroupWorker) {
	keeper.workerMutex.Lock()
	delete(keeper.workerMap, worker.workName)
	workerLen := len(keeper.workerMap)
	keeper.workerMutex.Unlock()

	if workerLen == 0 {
		keeper.statusMutex.Lock()
		defer keeper.statusMutex.Unlock()
		keeper.status = xKEEPER_STATUS_STOPPED
	} else {
		global.Logger().Warn(fmt.Sprintf("[%s] warning! groupworkers limit %d plans, actually %d", keeper.keeperName, keeper.maxGroupQueues, workerLen))
	}
}

func (keeper *ZeroMfgrcGroupKeeper) closeGroup(group MfgrcGroup) {
	keeper.groupMutex.Lock()
	delete(keeper.groupMap, group.XuniqueCode())
	keeper.groupMutex.Unlock()
}

func (keeper *ZeroMfgrcGroupKeeper) AddGroup(group MfgrcGroup) error {
	keeper.statusMutex.Lock()
	xStatus := keeper.status
	keeper.statusMutex.Unlock()
	if xStatus == xKEEPER_STATUS_STOPPED {
		return errors.New("keeper not yet ready")
	} else if xStatus == xKEEPER_STATUS_STOPPING {
		return errors.New("keeper is stopping now")
	}

	keeper.groupMutex.Lock()
	_, ok := keeper.groupMap[group.XuniqueCode()]
	keeper.groupMutex.Unlock()
	if !ok {
		keeper.groupMutex.Lock()
		keeper.groupMap[group.XuniqueCode()] = group
		keeper.groupMutex.Unlock()

		group.Pending()
		keeper.groupChan <- group
	} else {
		return errors.New(fmt.Sprintf("device `%s` is busy now", group.XuniqueCode()))
	}

	return nil
}

func (keeper *ZeroMfgrcGroupKeeper) Check(group MfgrcGroup) error {
	keeper.statusMutex.Lock()
	xStatus := keeper.status
	keeper.statusMutex.Unlock()
	if xStatus == xKEEPER_STATUS_STOPPED {
		return errors.New("keeper not yet ready")
	} else if xStatus == xKEEPER_STATUS_STOPPING {
		return errors.New("keeper is stopping now")
	}

	keeper.groupMutex.Lock()
	_, ok := keeper.groupMap[group.XuniqueCode()]
	groupMapLen := len(keeper.groupMap)
	keeper.groupMutex.Unlock()

	if ok {
		return errors.New(fmt.Sprintf("device `%s` is busy now", group.XuniqueCode()))
	}

	if groupMapLen >= keeper.maxGroupQueues {
		return errors.New("exceeding maximum limit")
	}

	return nil
}

func (keeper *ZeroMfgrcGroupKeeper) Export() (map[string]interface{}, error) {
	configs := make(map[string]interface{})

	configs["maxGroupQueues"] = keeper.maxGroupQueues
	keeper.statusMutex.RLock()
	configs["status"] = keeper.status
	keeper.statusMutex.RUnlock()

	workers := make(map[string]interface{})
	keeper.workerMutex.RLock()
	for workerGroupName, worker := range keeper.workerMap {
		workers[workerGroupName] = worker.Export()
	}
	keeper.workerMutex.RUnlock()

	groups := make(map[string]interface{})
	keeper.groupMutex.RLock()
	defer keeper.groupMutex.RUnlock()
	for groupID, group := range keeper.groupMap {
		groupMap, err := group.Export()
		if err != nil {
			return nil, err
		}
		groups[groupID] = groupMap
	}

	exports := make(map[string]interface{})
	exports["configs"] = configs
	exports["workers"] = workers
	exports["groups"] = groups

	return exports, nil
}
