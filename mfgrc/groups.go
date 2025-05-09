package mfgrc

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/autohttpconf"
	"github.com/0meet1/zero-framework/errdef"
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
	autohttpconf.ZeroXsacXhttpStructs

	worker *ZeroMfgrcGroupWorker

	UniqueCode string `json:"uniqueCode,omitempty" xhttpopt:"XX" xsacprop:"NO,VARCHAR(64),NULL" xsackey:"key" xapi:"唯一标识,String"`
	Option     string `json:"option,omitempty" xhttpopt:"XX" xsacprop:"NO,VARCHAR(64),NULL" xsackey:"key" xapi:"操作类型,String"`
	Operator   string `json:"operator,omitempty" xhttpopt:"XX" xsacprop:"NO,VARCHAR(32),NULL" xsackey:"key" xapi:"操作人,String"`

	Monos []MfgrcMono `json:"monos,omitempty"`

	status string
	reason string

	xStore    ZeroMfgrcGroupStore
	xListener ZeroMfgrcGroupEventListener
}

func (*ZeroMfgrcGroup) XhttpPath() string     { return "Xdenied" }
func (*ZeroMfgrcGroup) XsacApiName() string   { return "mfgrc组任务模型" }
func (*ZeroMfgrcGroup) XsacPartition() string { return structs.XSAC_PARTITION_MONTH }
func (*ZeroMfgrcGroup) XhttpOpt() byte {
	return structs.XahttpOpt(structs.XahttpOpt_F, structs.XahttpOpt_F, structs.XahttpOpt_F, structs.XahttpOpt_F, structs.XahttpOpt_F)
}
func (*ZeroMfgrcGroup) XsacAdjunctDeclares(args ...string) structs.ZeroXsacEntrySet {
	return make(structs.ZeroXsacEntrySet, 0)
}

func (group *ZeroMfgrcGroup) LoadRowData(rowmap map[string]interface{}) {
	group.ZeroCoreStructs.LoadRowData(rowmap)

	group.UniqueCode = structs.ParseStringField(rowmap, "unique_code")
	group.Option = structs.ParseStringField(rowmap, "option")
	group.Operator = structs.ParseStringField(rowmap, "operator")
	group.status = structs.ParseStringField(rowmap, "status")

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

func (group *ZeroMfgrcGroup) Xoperator() string {
	return group.Operator
}

func (group *ZeroMfgrcGroup) Xmonos() []MfgrcMono {
	if group.Monos == nil {
		group.Monos = make([]MfgrcMono, 0)
	}
	return group.Monos
}

func (group *ZeroMfgrcGroup) XLinkTable() string { return "" }

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
		return fmt.Errorf("could not pending mono group `%s` status `%s`", group.ID, group.status)
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
		return fmt.Errorf("could not executing group `%s` status `%s`", group.ID, group.status)
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
		return fmt.Errorf("could not complete group `%s` status `%s`", group.ID, group.status)
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

func (group *ZeroMfgrcGroup) Failed(reason error) error {
	if errdef.Is(reason) {
		group.Features["errdef"] = reason.(*errdef.ZeroExceptionDef).Export()
	}
	group.reason = reason.Error()
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
	jsonbytes, err := json.Marshal(group.This())
	if err != nil {
		return nil, err
	}

	var jsonmap map[string]interface{}
	err = json.Unmarshal(jsonbytes, &jsonmap)
	if err != nil {
		return nil, err
	}

	jsonmap["uniqueCode"] = group.UniqueCode
	jsonmap["option"] = group.Option

	if group.Monos != nil {
		monos := make([]map[string]interface{}, len(group.Monos))
		for i, mono := range group.Monos {
			if mono.This() == nil {
				mono.ThisDef(mono)
			}
			monos[i], err = mono.Export()
			if err != nil {
				return nil, err
			}
		}
		jsonmap["monos"] = monos
	}

	jsonmap["status"] = group.status
	jsonmap["reason"] = group.reason

	return jsonmap, nil
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

	global.Logger().Info(fmt.Sprintf("[%s] ready and waiting ...", worker.workName))
	for xGroup := range worker.keeper.groupChan {
		worker.statusMutex.Lock()
		xstatus := worker.status
		worker.statusMutex.Unlock()
		if xstatus != xWORKER_STATUS_RUNNING {
			break
		}

		if xGroup != nil {
			global.Logger().Info(fmt.Sprintf("[%s] exec group `%s` unique code `%s`", worker.workName, xGroup.XgroupId(), xGroup.XuniqueCode()))
			worker.executing = xGroup.XuniqueCode()

			xGroup.AddWorker(worker)
			err := xGroup.Executing()
			if err != nil {
				xGroup.Failed(err)
			} else {
				err = xGroup.Do()
				if err != nil {
					xGroup.Failed(err)
				} else {
					err = xGroup.Complete()
					if err != nil {
						xGroup.Failed(err)
					}
				}
			}

			worker.keeper.closeGroup(xGroup)

			global.Logger().Info(fmt.Sprintf("[%s] `%s` unique code `%s` work complete", worker.workName, xGroup.XgroupId(), xGroup.XuniqueCode()))
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
		group.Failed(fmt.Errorf("group unexpected termination"))
	}

	keeper.statusMutex.Lock()
	defer keeper.statusMutex.Unlock()
	keeper.status = xKEEPER_STATUS_RUNNING

	global.Logger().Info(" workergroup check and resume monos complete ")
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
		return fmt.Errorf("unique code `%s` is busy now", group.XuniqueCode())
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
		return fmt.Errorf("unique code `%s` is busy now", group.XuniqueCode())
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
