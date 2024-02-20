package global

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"
	"syscall"

	cfg "github.com/0meet1/zero-framework/config"
	"github.com/0meet1/zero-framework/log"
)

const ZERO_FRAMEWORK_BANNER = `



	███████ ███████ ██████   ██████      ███████ ██████   █████  ███    ███ ███████ ██     ██  ██████  ██████  ██   ██ 
	   ███  ██      ██   ██ ██    ██     ██      ██   ██ ██   ██ ████  ████ ██      ██     ██ ██    ██ ██   ██ ██  ██  
	  ███   █████   ██████  ██    ██     █████   ██████  ███████ ██ ████ ██ █████   ██  █  ██ ██    ██ ██████  █████   
	 ███    ██      ██   ██ ██    ██     ██      ██   ██ ██   ██ ██  ██  ██ ██      ██ ███ ██ ██    ██ ██   ██ ██  ██  
	███████ ███████ ██   ██  ██████      ██      ██   ██ ██   ██ ██      ██ ███████  ███ ███   ██████  ██   ██ ██   ██


	 /**  :: Zero Framewrok For Golang ::  **********   **********   **********   **********  ( v1.10.0.RELEASE )  **/

`

type ZeroGlobalEventsObserver interface {
	Shutdown() error
}

type ZeroGlobalInitiator func() error

var (
	_appName string

	_once sync.Once

	_map    map[string]interface{}
	_rwLock sync.RWMutex

	_wMap  map[string]interface{}
	_wLock sync.Mutex

	_observers map[string]ZeroGlobalEventsObserver
	_oLock     sync.Mutex
)

func copyMap(src map[string]interface{}) map[string]interface{} {
	dest := make(map[string]interface{})
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

func shared() map[string]interface{} {
	_once.Do(func() {
		_wMap = make(map[string]interface{})
	})
	return _wMap
}

func synchronize() {
	_rwLock.Lock()
	_wLock.Lock()
	defer _rwLock.Unlock()
	defer _wLock.Unlock()
	_map = copyMap(shared())
}

func Key(key string, value interface{}) {
	if _observers == nil {
		panic("global context not initialized")
	}

	if _, ok := shared()[key]; ok {
		panic(fmt.Sprintf("key `%s` already exists", key))
	}

	_wLock.Lock()
	shared()[key] = value
	_wLock.Unlock()
	synchronize()
}

func Pop(key string) {
	if _observers == nil {
		panic("global context not initialized")
	}

	_wLock.Lock()
	delete(shared(), key)
	_wLock.Unlock()
	synchronize()
}

func Value(key string) interface{} {
	if _observers == nil {
		panic("global context not initialized")
	}

	_rwLock.RLock()
	defer _rwLock.RUnlock()
	return _map[key]
}

func Contains(key string) bool {
	if _observers == nil {
		panic("global context not initialized")
	}

	_rwLock.RLock()
	defer _rwLock.RUnlock()
	_, ok := _map[key]
	return ok
}

var (
	channel chan os.Signal
)

func RunServer() {
	if _observers == nil {
		panic("global context not initialized")
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	if _observers != nil {
		_oLock.Lock()
		for _, observer := range _observers {
			observer.Shutdown()
		}
		_oLock.Unlock()
	}
}

func ListenEvents(name string, observer ZeroGlobalEventsObserver) {
	if _observers == nil {
		panic("global context not initialized")
	}

	_oLock.Lock()
	_observers[name] = observer
	_oLock.Unlock()
}

func LeaveEventsObserver(name string) {
	if _observers == nil {
		panic("global context not initialized")
	}

	_oLock.Lock()
	delete(_observers, name)
	_oLock.Unlock()
}

func systemAbsPath() string {
	_zeroFrameworkHome := os.Getenv("ZERO_HOME")
	if len(_zeroFrameworkHome) > 0 {
		return _zeroFrameworkHome
	}
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		dir, file := path.Split(filename)
		if !strings.HasPrefix(file, _appName) && !strings.HasPrefix(file, "main") {
			panic("global context must be initialized in `main func` and appname must same as main package filename or 'main' .")
		}
		return dir
	}
	return ""
}

func AppName() string {
	if _observers == nil {
		panic("global context not initialized")
	}
	return _appName
}

func GlobalContext(appName string) {
	if _observers == nil {
		_observers = make(map[string]ZeroGlobalEventsObserver)
	}
	if len(cfg.ServerAbsPath()) == 0 {
		_appName = appName
		cfg.NewConfigs(systemAbsPath())
		Key("zero.system.logger", log.InitLogger())
		Logger().Info(ZERO_FRAMEWORK_BANNER)
	}
}

func Run(appName string, initiators ...ZeroGlobalInitiator) {
	GlobalContext(appName)
	if initiators != nil {
		for _, initiator := range initiators {
			if err := initiator(); err != nil {
				panic(err)
			}
		}
	}
	RunServer()
}

func ServerAbsPath() string {
	return cfg.ServerAbsPath()
}

func StringValue(cfgName string) string {
	return cfg.StringValue(cfgName)
}

func IntValue(cfgName string) int {
	return cfg.IntValue(cfgName)
}

func SliceStringValue(cfgName string) []string {
	return cfg.SliceStringValue(cfgName)
}

func Logger() *log.ZeroLogger {
	return Value("zero.system.logger").(*log.ZeroLogger)
}
