package mfgrc

import (
	"errors"
	"fmt"
	"time"
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
	XmonoId() string
	XuniqueCode() string
	Xoption() string
	Xprogress() int

	State() string
	Response() interface{}

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

	Do()
	Export() (map[string]interface{}, error)

	Store(ZeroMfgrcGroupStore)
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

type ZeroMfgrcMonoProgressListener interface {
	OnProgress(MfgrcMono) error
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
	keeper  *ZeroMfgrcKeeper
	mono    MfgrcMono
	errchan chan error
}

func (act *ZeroMfgrcMonoActuator) Exec(mono MfgrcMono) chan error {
	act.mono = mono
	act.errchan = make(chan error, 1)

	act.mono.EventListener(act)
	err := act.keeper.AddMono(act.mono)
	if err != nil {
		go func() {
			time.After(time.Millisecond * time.Duration(100))
			act.errchan <- err
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
	act.errchan <- errors.New(fmt.Sprintf("mono `%s` is already revoke", act.mono.XmonoId()))
	return nil
}
func (act *ZeroMfgrcMonoActuator) OnComplete(MfgrcMono) error {
	act.errchan <- nil
	return nil
}
func (act *ZeroMfgrcMonoActuator) OnFailed(mono MfgrcMono, reason string) error {
	act.errchan <- errors.New(fmt.Sprintf("mono `%s` exec failed, reason: %s", act.mono.XmonoId(), reason))
	return nil
}
