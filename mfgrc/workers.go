package mfgrc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
)

type ZeroMfgrcFlux struct {
	UniqueId  string
	monoMap   map[string]MfgrcMono
	monos     chan MfgrcMono
	monoMutex sync.RWMutex

	keeper *ZeroMfgrcKeeper
	worker *ZeroMfgrcWorker

	active bool
}

func createMfgrcf(mono MfgrcMono, keeper *ZeroMfgrcKeeper) *ZeroMfgrcFlux {
	flux := &ZeroMfgrcFlux{
		keeper:  keeper,
		monoMap: make(map[string]MfgrcMono),
		monos:   make(chan MfgrcMono, keeper.maxQueueLimit),
		active:  true,
	}
	flux.UniqueId = mono.XuniqueCode()
	return flux
}

func (flux *ZeroMfgrcFlux) Push(mono MfgrcMono, keeper *ZeroMfgrcKeeper) (bool, error) {

	flux.monoMutex.Lock()
	if !flux.active {
		flux.monoMutex.Unlock()
		return false, nil
	} else {
		_, ok := flux.monoMap[mono.XmonoId()]
		if ok {
			return true, fmt.Errorf("flux `%s` mono `%s` is already exists", flux.UniqueId, mono.XmonoId())
		}

		if len(flux.monoMap) >= flux.keeper.maxQueueLimit {
			return true, fmt.Errorf("flux `%s` has been exceeded maximum number of mono = %d", flux.UniqueId, flux.keeper.maxQueueLimit)
		}
		flux.monoMap[mono.XmonoId()] = mono
		err := mono.Pending(flux)
		if err != nil {
			delete(flux.monoMap, mono.XmonoId())
			mono.Failed(err)
		}
		flux.monoMutex.Unlock()
		if err != nil {
			return true, err
		}
		go func() { flux.monos <- mono }()
		return true, nil
	}
}

func (flux *ZeroMfgrcFlux) Revoke(mono MfgrcMono) error {
	flux.monoMutex.Lock()
	var err error = nil
	mono, ok := flux.monoMap[mono.XmonoId()]
	if !ok {
		err = fmt.Errorf("mono `%s` not found", mono.XmonoId())
	} else {
		err = mono.Revoke()
	}
	flux.monoMutex.Unlock()
	return err
}

func (flux *ZeroMfgrcFlux) Check(mono MfgrcMono) bool {
	flux.monoMutex.Lock()
	monoLen := len(flux.monoMap)
	flux.monoMutex.Unlock()
	return !(monoLen >= flux.keeper.maxQueueLimit)
}

func (flux *ZeroMfgrcFlux) close() {
	flux.keeper.mfgrcMutex.Lock()
	delete(flux.keeper.mfgrcMap, flux.UniqueId)
	flux.keeper.mfgrcMutex.Unlock()
}

func (flux *ZeroMfgrcFlux) cleanMono(mono MfgrcMono) {
	flux.monoMutex.Lock()
	delete(flux.monoMap, mono.XmonoId())
	flux.monoMutex.Unlock()
}

func (flux *ZeroMfgrcFlux) completeMono(mono MfgrcMono, err error) {
	if err == nil {
		err = mono.Complete()
		if err != nil {
			mono.Failed(err)
		}
	} else {
		global.Logger().Error(fmt.Sprintf("flux `%s` mono `%s` error : %s", flux.UniqueId, mono.XmonoId(), err.Error()))
		if mono.MaxExecuteTimes() > 1 {
			err := mono.Retrying(err)
			if err != nil {
				mono.Failed(err)
			} else {
				if flux.keeper.taskIntervalSeconds > 0 {
					<-time.After(time.Second * time.Duration(flux.keeper.taskIntervalSeconds))
				}
				flux.completeMono(mono, mono.Do())
			}
		} else {
			mono.Failed(err)
		}
	}
}

func (flux *ZeroMfgrcFlux) runLoop() bool {
	select {
	case mono := <-flux.monos:
		if mono.State() != WORKER_MONO_STATUS_PENDING && mono.State() != WORKER_MONO_STATUS_EXECUTING && mono.State() != WORKER_MONO_STATUS_RETRYING {
			flux.cleanMono(mono)
		} else {
			if mono.State() == WORKER_MONO_STATUS_PENDING {
				err := mono.Executing()
				if err != nil {
					mono.Failed(err)
					flux.cleanMono(mono)
				}
			}
			if mono.State() == WORKER_MONO_STATUS_EXECUTING {
				flux.completeMono(mono, mono.Do())
				flux.cleanMono(mono)
			}
		}
		if flux.keeper.taskIntervalSeconds > 0 {
			<-time.After(time.Second * time.Duration(flux.keeper.taskIntervalSeconds))
		}
		return true
	case <-time.After(time.Millisecond * time.Duration(500)):
		flux.monoMutex.Lock()
		if len(flux.monoMap) > 0 {
			flux.monoMutex.Unlock()
			return true
		} else {
			flux.active = false
			flux.monoMutex.Unlock()
			flux.close()
			return false
		}

	}
}

func (flux *ZeroMfgrcFlux) Start(worker *ZeroMfgrcWorker) {
	flux.worker = worker
	for {
		if !flux.runLoop() {
			break
		}
	}
}

func (flux *ZeroMfgrcFlux) Export() (map[string]interface{}, error) {
	exportMap := make(map[string]interface{})

	exportMap["uniqueId"] = flux.UniqueId
	exportMap["workName"] = flux.worker.workName

	monosMap := make(map[string]interface{})
	flux.monoMutex.RLock()
	defer flux.monoMutex.RUnlock()

	for key, value := range flux.monoMap {
		monoMap, err := value.Export()
		if err != nil {
			return nil, err
		}
		monosMap[key] = monoMap
	}
	exportMap["monos"] = monosMap
	return exportMap, nil
}

type ZeroMfgrcWorker struct {
	workName string

	status      string
	statusMutex sync.RWMutex

	executing string

	keeper *ZeroMfgrcKeeper
}

func newMfgrcWorker(workName string, keeper *ZeroMfgrcKeeper) *ZeroMfgrcWorker {
	worker := ZeroMfgrcWorker{}
	worker.workName = workName
	worker.keeper = keeper
	return &worker
}

func (worker *ZeroMfgrcWorker) Start() {
	worker.statusMutex.Lock()
	worker.status = xWORKER_STATUS_RUNNING
	worker.statusMutex.Unlock()

	global.Logger().Info(fmt.Sprintf("[%s] ready and waiting ...", worker.workName))
	for xQueue := range worker.keeper.mfgrcChan {
		worker.statusMutex.Lock()
		xstatus := worker.status
		worker.statusMutex.Unlock()

		if xQueue == nil || xstatus != xWORKER_STATUS_RUNNING {
			break
		}

		global.Logger().Info(fmt.Sprintf("[%s] exec flux `%s`", worker.workName, xQueue.UniqueId))
		worker.executing = xQueue.UniqueId

		xQueue.Start(worker)

		global.Logger().Info(fmt.Sprintf("[%s] flux `%s` complete", worker.workName, xQueue.UniqueId))
		worker.executing = ""
	}
	worker.keeper.closeWorker(worker)
	global.Logger().Info(fmt.Sprintf("[%s] warning! worker is shutdown now", worker.workName))
}

func (worker *ZeroMfgrcWorker) Stop() {
	worker.statusMutex.Lock()
	worker.status = xWORKER_STATUS_STOPPED
	worker.statusMutex.Unlock()
}

func (worker *ZeroMfgrcWorker) Export() map[string]interface{} {
	worker.statusMutex.RLock()
	defer worker.statusMutex.RUnlock()

	workerMap := make(map[string]interface{})
	workerMap["workName"] = worker.workName
	workerMap["status"] = worker.status
	workerMap["executing"] = worker.executing
	return workerMap
}

type ZeroMfgrcKeeper struct {
	keeperName  string
	workerMap   map[string]*ZeroMfgrcWorker
	workerMutex sync.RWMutex

	mfgrcMap   map[string]*ZeroMfgrcFlux
	mfgrcChan  chan *ZeroMfgrcFlux
	mfgrcMutex sync.RWMutex

	maxQueues           int
	maxQueueLimit       int
	taskWaitSeconds     int
	taskIntervalSeconds int
	taskRetryTimes      int
	taskRetryInterval   int

	status      string
	statusMutex sync.RWMutex

	keeperOpts ZeroMfgrcKeeperOpts
}

func NewWorker(
	keeperName string,
	keeperOpts ZeroMfgrcKeeperOpts,
	maxQueues int,
	maxQueueLimit int,
	taskWaitSeconds int,
	taskIntervalSeconds int,
	taskRetryTimes int,
	taskRetryInterval int) *ZeroMfgrcKeeper {

	return &ZeroMfgrcKeeper{
		keeperName:          keeperName,
		workerMap:           make(map[string]*ZeroMfgrcWorker),
		mfgrcMap:            make(map[string]*ZeroMfgrcFlux),
		mfgrcChan:           make(chan *ZeroMfgrcFlux, maxQueues),
		maxQueues:           maxQueues,
		maxQueueLimit:       maxQueueLimit,
		taskWaitSeconds:     taskWaitSeconds,
		taskIntervalSeconds: taskIntervalSeconds,
		taskRetryTimes:      taskRetryTimes,
		taskRetryInterval:   taskRetryInterval,
		status:              xKEEPER_STATUS_STOPPED,
		keeperOpts:          keeperOpts,
	}
}

func (keeper *ZeroMfgrcKeeper) RunWorker() {
	global.Logger().Info(fmt.Sprintf("worker start with maxQueues: %d, maxGroupLimit: %d, taskRetryTimes: %d, taskWaitSeconds: %ds",
		keeper.maxQueues,
		keeper.maxQueueLimit,
		keeper.taskRetryTimes,
		keeper.taskWaitSeconds))

	if len(keeper.keeperName) == 0 {
		keeper.keeperName = "default"
	}

	for i := 0; i < keeper.maxQueues; i++ {
		worker := newMfgrcWorker(fmt.Sprintf("%s-worker-%03d::", keeper.keeperName, i), keeper)

		keeper.workerMutex.Lock()
		keeper.workerMap[worker.workName] = worker
		keeper.workerMutex.Unlock()
		go worker.Start()
	}
	go keeper.resumeMonos()
}

func (keeper *ZeroMfgrcKeeper) ShutdownWorker() {
	for _, worker := range keeper.workerMap {
		worker.Stop()
	}

	keeper.mfgrcMutex.Lock()
	for i := 0; i < keeper.maxQueues-len(keeper.mfgrcMap); i++ {
		keeper.mfgrcChan <- nil
	}
	keeper.mfgrcMutex.Unlock()

	keeper.statusMutex.Lock()
	defer keeper.statusMutex.Unlock()
	keeper.status = xKEEPER_STATUS_STOPPING
}

func (keeper *ZeroMfgrcKeeper) checkStare() error {
	keeper.statusMutex.RLock()
	defer keeper.statusMutex.RUnlock()
	switch keeper.status {
	case xKEEPER_STATUS_STOPPED:
		return errors.New("keeper not yet ready")
	case xKEEPER_STATUS_STOPPING:
		return errors.New("keeper is stopping now")
	default:
		return nil
	}
}

func (keeper *ZeroMfgrcKeeper) resumeMonos() {
	if keeper.keeperOpts != nil {
		<-time.After(time.Second * time.Duration(3))
		monos, err := keeper.keeperOpts.FetchUncompleteMonos()
		if err != nil {
			global.Logger().Error(fmt.Sprintf(" resume monos err : %s", err.Error()))
		} else {
			for _, mono := range monos {
				mono.Store(keeper.keeperOpts.MonoStore())
				mono.Revoke()
			}
		}
	}
	keeper.statusMutex.Lock()
	keeper.status = xKEEPER_STATUS_RUNNING
	keeper.statusMutex.Unlock()

	global.Logger().Info(" worker check and resume monos complete ")
}

func (keeper *ZeroMfgrcKeeper) closeWorker(worker *ZeroMfgrcWorker) {
	keeper.workerMutex.Lock()
	delete(keeper.workerMap, worker.workName)
	workerLen := len(keeper.workerMap)
	keeper.workerMutex.Unlock()

	if workerLen == 0 {
		keeper.statusMutex.Lock()
		keeper.status = xKEEPER_STATUS_STOPPED
		keeper.statusMutex.Unlock()
	} else {
		global.Logger().Warn(fmt.Sprintf("[%s] warning! workers limit %d plans, actually %d", keeper.keeperName, keeper.maxQueues, workerLen))
	}
}

func (keeper *ZeroMfgrcKeeper) TaskWaitSeconds() int {
	return keeper.taskWaitSeconds
}

func (keeper *ZeroMfgrcKeeper) AddMono(mono MfgrcMono) error {
	err := keeper.checkStare()
	if err != nil {
		return err
	}
	keeper.mfgrcMutex.Lock()
	flux, ok := keeper.mfgrcMap[mono.XuniqueCode()]
	if !ok {
		keeper.mfgrcMap[mono.XuniqueCode()] = createMfgrcf(mono, keeper)
		flux = keeper.mfgrcMap[mono.XuniqueCode()]
		go func() { keeper.mfgrcChan <- flux }()
	}
	onPending, err := flux.Push(mono, keeper)
	keeper.mfgrcMutex.Unlock()
	if !onPending {
		return keeper.AddMono(mono)
	}
	return err
}

func (keeper *ZeroMfgrcKeeper) RevokeMono(mono MfgrcMono) error {
	err := keeper.checkStare()
	if err != nil {
		return err
	}
	keeper.mfgrcMutex.Lock()
	flux, ok := keeper.mfgrcMap[mono.XuniqueCode()]
	if !ok {
		err = fmt.Errorf("unique code `%s` flux not found", mono.XuniqueCode())
	} else {
		err = flux.Revoke(mono)
	}
	keeper.mfgrcMutex.Unlock()
	return err
}

func (keeper *ZeroMfgrcKeeper) AddMonosQueue(monos ...MfgrcMono) error {
	err := keeper.checkStare()
	if err != nil {
		return err
	}
	if len(monos) > keeper.maxQueueLimit {
		return fmt.Errorf("exceeding maximum number of monos = %d", keeper.maxQueueLimit)
	}
	errors := ""
	for _, mono := range monos {
		err = keeper.AddMono(mono)
		if err != nil {
			if errors == "" {
				errors = fmt.Sprintf("%s:%s", mono.XuniqueCode(), err.Error())
			} else {
				errors = fmt.Sprintf("%s,%s:%s", errors, mono.XuniqueCode(), err.Error())
			}
		}
	}
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func (keeper *ZeroMfgrcKeeper) Check(mono MfgrcMono) error {
	err := keeper.checkStare()
	if err != nil {
		return err
	}
	keeper.mfgrcMutex.Lock()
	flux, ok := keeper.mfgrcMap[mono.XuniqueCode()]
	if ok && !flux.Check(mono) {
		err = errors.New("exceeding maximum limit")
	}
	keeper.mfgrcMutex.Unlock()
	return err
}

func (keeper *ZeroMfgrcKeeper) Export() (map[string]interface{}, error) {
	configs := make(map[string]interface{})

	configs["keeperName"] = keeper.keeperName
	configs["maxQueues"] = keeper.maxQueues
	configs["maxQueueLimit"] = keeper.maxQueueLimit
	configs["taskWaitSeconds"] = keeper.taskWaitSeconds
	configs["taskIntervalSeconds"] = keeper.taskIntervalSeconds
	configs["taskRetryTimes"] = keeper.taskRetryTimes
	configs["taskRetryInterval"] = keeper.taskRetryInterval
	configs["status"] = keeper.status

	workers := make(map[string]interface{})
	keeper.workerMutex.RLock()
	for workerName, worker := range keeper.workerMap {
		workers[workerName] = worker.Export()
	}
	keeper.workerMutex.RUnlock()

	fluxs := make(map[string]interface{})
	keeper.mfgrcMutex.RLock()
	defer keeper.mfgrcMutex.RUnlock()
	for fluxID, flux := range keeper.mfgrcMap {
		fluxMap, err := flux.Export()
		if err != nil {
			return nil, err
		}
		fluxs[fluxID] = fluxMap
	}

	exports := make(map[string]interface{})
	exports["configs"] = configs
	exports["workers"] = workers
	exports["fluxs"] = fluxs

	return exports, nil
}
