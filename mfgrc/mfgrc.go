package mfgrc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/structs"
)

const (
	xKEEPER_STATUS_RUNNING  = "running"
	xKEEPER_STATUS_STOPPING = "stopping"
	xKEEPER_STATUS_STOPPED  = "stopped"
)

const (
	xWORKER_STATUS_RUNNING = "running"
	xWORKER_STATUS_STOPPED = "stopped"
)

type MfgrcMono interface {
	structs.ZeroMetaDef

	XmonoId() string
	XuniqueCode() string
	Xoption() string
	Xprogress() int

	State() string

	Pending(*ZeroMfgrcFlux) error
	Revoke() error
	Timeout() error
	Executing() error
	Retrying(string) error
	Complete() error
	Failed(string) error

	Do() error
	Export() (map[string]interface{}, error)

	Store(ZeroMfgrcMonoStore)
	EventListener(ZeroMfgrcMonoEventListener)
}

type MfgrcGroup interface {
	structs.ZeroMetaDef

	XgroupId() string
	XuniqueCode() string
	Xoption() string
	Xmonos() []MfgrcMono

	State() string
	AddWorker(*ZeroMfgrcGroupWorker)

	Ready(ZeroMfgrcGroupStore) error
	Pending() error
	Executing() error
	Complete() error
	Failed(string) error

	Do() error
	Export() (map[string]interface{}, error)

	Store(ZeroMfgrcGroupStore)
	UseStore() ZeroMfgrcGroupStore
	EventListener(ZeroMfgrcGroupEventListener)
}

type ZeroMfgrcMonoStore interface {
	UpdateMono(MfgrcMono) error
	DeleteMono(MfgrcMono) error
}

type ZeroMfgrcMonoEventListener interface {
	OnPending(MfgrcMono) error
	OnRevoke(MfgrcMono) error
	OnExecuting(MfgrcMono) error
	OnRetrying(MfgrcMono) error
	OnComplete(MfgrcMono) error
	OnFailed(MfgrcMono, string) error
}

type ZeroMfgrcGroupEventListener interface {
	OnPending(MfgrcGroup) error
	OnExecuting(MfgrcGroup) error
	OnComplete(MfgrcGroup) error
	OnFailed(MfgrcGroup, string) error
}

type ZeroMfgrcKeeperOpts interface {
	FetchUncompleteMonos() ([]MfgrcMono, error)
	DatebaseDatetime() (*time.Time, error)
	MonoStore() ZeroMfgrcMonoStore
}

type ZeroMfgrcGroupKeeperOpts interface {
	FetchUncompleteGroups() ([]MfgrcGroup, error)
	MonoStore() ZeroMfgrcMonoStore
	GroupStore() ZeroMfgrcGroupStore
}

type ZeroMfgrcGroupStore interface {
	UpdateGroup(MfgrcGroup) error
	DeleteGroup(MfgrcGroup) error
	AddGroupMono(MfgrcGroup, MfgrcMono) error
	MonoStore() ZeroMfgrcMonoStore
	NewSerial(string, string) (int64, error)
}

type ZeroMfgrcMonoActuator struct {
	Keeper  *ZeroMfgrcKeeper
	mono    MfgrcMono
	errchan chan error
}

func (act *ZeroMfgrcMonoActuator) Exec(mono MfgrcMono) chan error {
	act.mono = mono
	act.errchan = make(chan error, 1)

	act.mono.EventListener(act)
	err := act.Keeper.AddMono(act.mono)
	if err != nil {
		go func() {
			time.After(time.Millisecond * time.Duration(100))
			act.errchan <- err
			act.errchan = nil
		}()
	}
	return act.errchan
}

func (act *ZeroMfgrcMonoActuator) Mono() MfgrcMono {
	return act.mono
}

func (act *ZeroMfgrcMonoActuator) OnPending(MfgrcMono) error   { return nil }
func (act *ZeroMfgrcMonoActuator) OnExecuting(MfgrcMono) error { return nil }
func (act *ZeroMfgrcMonoActuator) OnRetrying(MfgrcMono) error  { return nil }
func (act *ZeroMfgrcMonoActuator) OnRevoke(mono MfgrcMono) error {
	if act.errchan != nil {
		act.errchan <- errors.New(fmt.Sprintf("mono `%s` is already revoke", act.mono.XmonoId()))
	}
	return nil
}
func (act *ZeroMfgrcMonoActuator) OnComplete(MfgrcMono) error {
	if act.errchan != nil {
		act.errchan <- nil
	}
	return nil
}
func (act *ZeroMfgrcMonoActuator) OnFailed(mono MfgrcMono, reason string) error {
	if act.errchan != nil {
		act.errchan <- errors.New(fmt.Sprintf("mono `%s` exec failed, reason: %s", act.mono.XmonoId(), reason))
	}
	return nil
}

type ZeroMfgrcGroupActuator struct {
	Keeper  *ZeroMfgrcGroupKeeper
	group   MfgrcGroup
	errchan chan error
}

func (act *ZeroMfgrcGroupActuator) Exec(group MfgrcGroup) chan error {
	act.group = group
	act.errchan = make(chan error, 1)

	act.group.EventListener(act)
	err := act.Keeper.AddGroup(act.group)
	if err != nil {
		go func() {
			time.After(time.Millisecond * time.Duration(100))
			act.errchan <- err
			act.errchan = nil
		}()
	}
	return act.errchan
}

func (act *ZeroMfgrcGroupActuator) Group() MfgrcGroup {
	return act.group
}

func (_ *ZeroMfgrcGroupActuator) OnPending(MfgrcGroup) error   { return nil }
func (_ *ZeroMfgrcGroupActuator) OnExecuting(MfgrcGroup) error { return nil }
func (act *ZeroMfgrcGroupActuator) OnComplete(MfgrcGroup) error {
	if act.errchan != nil {
		act.errchan <- nil
	}
	return nil
}
func (act *ZeroMfgrcGroupActuator) OnFailed(group MfgrcGroup, reason string) error {
	if act.errchan != nil {
		act.errchan <- errors.New(fmt.Sprintf("group `%s` exec failed, reason: %s", act.group.XgroupId(), reason))
	}
	return nil
}

type ZeroMfgrcMonoQueueActuator struct {
	Keeper  *ZeroMfgrcKeeper
	monos   []MfgrcMono
	errchan chan error

	counterLock sync.Mutex
	success     int
	failed      int
	result      map[string]string
}

func (act *ZeroMfgrcMonoQueueActuator) Exec(monos ...MfgrcMono) chan error {
	act.monos = monos
	act.errchan = make(chan error, 1)
	act.success = 0
	act.failed = 0
	act.result = make(map[string]string)

	for _, mono := range act.monos {
		mono.EventListener(act)
	}
	err := act.Keeper.AddMonosQueue(act.monos...)
	if err != nil {
		go func() {
			time.After(time.Millisecond * time.Duration(100))
			act.errchan <- err
			act.errchan = nil
		}()
	}

	return act.errchan
}

func (act *ZeroMfgrcMonoQueueActuator) OnPending(MfgrcMono) error   { return nil }
func (act *ZeroMfgrcMonoQueueActuator) OnExecuting(MfgrcMono) error { return nil }
func (act *ZeroMfgrcMonoQueueActuator) OnRetrying(MfgrcMono) error  { return nil }
func (act *ZeroMfgrcMonoQueueActuator) OnRevoke(mono MfgrcMono) error {
	if act.errchan == nil {
		return nil
	}

	act.counterLock.Lock()
	defer act.counterLock.Unlock()

	act.failed++
	act.result[mono.XmonoId()] = fmt.Sprintf("mono `%s` is already revoke", mono.XmonoId())
	act.check()
	return nil
}
func (act *ZeroMfgrcMonoQueueActuator) OnComplete(mono MfgrcMono) error {
	if act.errchan == nil {
		return nil
	}

	act.counterLock.Lock()
	defer act.counterLock.Unlock()

	act.success++
	act.result[mono.XmonoId()] = ""
	act.check()
	return nil
}
func (act *ZeroMfgrcMonoQueueActuator) OnFailed(mono MfgrcMono, reason string) error {
	if act.errchan == nil {
		return nil
	}
	act.counterLock.Lock()
	defer act.counterLock.Unlock()

	act.failed++
	act.result[mono.XmonoId()] = fmt.Sprintf("mono `%s` exec failed, reason: %s", mono.XmonoId(), reason)
	act.check()
	return nil
}

func (act *ZeroMfgrcMonoQueueActuator) Result() map[string]string {
	return act.result
}

func (act *ZeroMfgrcMonoQueueActuator) Success() int {
	return act.success
}

func (act *ZeroMfgrcMonoQueueActuator) Failed() int {
	return act.failed
}

func (act *ZeroMfgrcMonoQueueActuator) check() {
	if act.success+act.failed == len(act.monos) {
		if act.failed > 0 {
			act.errchan <- errors.New(fmt.Sprintf("queue exec failed, failed: %d, success: %d", act.failed, act.success))
		} else {
			act.errchan <- nil
		}
	}
	act.errchan = nil
}
