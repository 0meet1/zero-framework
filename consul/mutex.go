package consul

import (
	"strings"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/hashicorp/consul/api"
)

const (
	DEFAULT_LOCK_WAIT_TIMEOUT = 30
)

type ZeroDCSMutex struct {
	KeyName string
	Lock    *api.Lock
	Lockc   chan struct{}
}

type ZeroDCSMutexTrunk interface {
	Key(string, string) error
	Value(string) (string, error)
	Del(string) error

	Acquire(string) (bool, *api.WriteMeta, error)
	Release(string) (bool, *api.WriteMeta, error)

	Lock(string, string, ...int) (*ZeroDCSMutex, error)
	Unlock(*ZeroDCSMutex) error
}

type xZeroDCSMutexTrunk struct {
	apiConfig *api.Config
	apiClient *api.Client
}

func (mtx *xZeroDCSMutexTrunk) Key(keyName, value string) error {
	_, err := mtx.apiClient.KV().Put(&api.KVPair{
		Key:   keyName,
		Value: []byte(value),
	}, &api.WriteOptions{})
	return err
}

func (mtx *xZeroDCSMutexTrunk) Value(keyName string) (string, error) {
	val, _, err := mtx.apiClient.KV().Get(keyName, &api.QueryOptions{})
	if err != nil {
		return "", err
	}
	if val != nil && val.Value != nil {
		return string(val.Value), nil
	}
	return "", nil
}

func (mtx *xZeroDCSMutexTrunk) Del(keyName string) error {
	_, err := mtx.apiClient.KV().Delete(keyName, &api.WriteOptions{})
	return err
}

func (mtx *xZeroDCSMutexTrunk) Acquire(keyName string) (bool, *api.WriteMeta, error) {
	return mtx.apiClient.KV().Acquire(&api.KVPair{
		Key: keyName,
	}, &api.WriteOptions{})
}

func (mtx *xZeroDCSMutexTrunk) Release(keyName string) (bool, *api.WriteMeta, error) {
	return mtx.apiClient.KV().Release(&api.KVPair{
		Key: keyName,
	}, &api.WriteOptions{})
}

func (mtx *xZeroDCSMutexTrunk) Lock(keyName, operator string, timeout ...int) (*ZeroDCSMutex, error) {
	_timeout := DEFAULT_LOCK_WAIT_TIMEOUT
	if len(timeout) > 0 && timeout[0] > 0 {
		_timeout = timeout[0]
	}
	_xLock, err := mtx.apiClient.LockOpts(&api.LockOptions{
		Key:          keyName,
		Value:        []byte(operator),
		LockWaitTime: time.Duration(_timeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	_c := make(chan struct{})
	_, err = _xLock.Lock(_c)
	if err != nil {
		return nil, err
	}

	return &ZeroDCSMutex{
		Lock:  _xLock,
		Lockc: _c,
	}, nil
}

func (mtx *xZeroDCSMutexTrunk) Unlock(_mutex *ZeroDCSMutex) error {
	_mutex.Lock.Unlock()
	close(_mutex.Lockc)
	return _mutex.Lock.Destroy()
}

func (mtx *xZeroDCSMutexTrunk) runDCSMutex() error {
	mtx.apiConfig = api.DefaultConfig()

	if strings.TrimSpace(global.StringValue("zero.consul.serverAddr")) != "" {
		mtx.apiConfig.Address = global.StringValue("zero.consul.serverAddr")
	}

	client, err := api.NewClient(mtx.apiConfig)
	if err != nil {
		return err
	}
	mtx.apiClient = client
	return nil
}
