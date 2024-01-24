package global

import (
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


	 /**  :: Zero Framewrok For Golang ::  **********   **********   **********   **********  ( v1.7.10.RELEASE )  **/

`

var (
	_appName string

	_map    map[string]interface{}
	_wMap   map[string]interface{}
	_once   sync.Once
	_rwLock sync.RWMutex
	_wLock  sync.Mutex
)

func copyMap(src map[string]interface{}) (map[string]interface{}, error) {
	dest := make(map[string]interface{})
	for key, value := range src {
		dest[key] = value
	}
	return dest, nil
}

func shared() map[string]interface{} {
	_once.Do(func() {
		_wMap = make(map[string]interface{})
	})
	return _wMap
}

func synchronize() error {
	_rwLock.Lock()
	defer _rwLock.Unlock()
	history := _map
	dist, err := copyMap(shared())
	if err != nil {
		_map = history
		return err
	}
	_map = dist
	return nil
}

func Key(key string, value interface{}) error {
	_wLock.Lock()
	shared()[key] = value
	_wLock.Unlock()
	return synchronize()
}

func Pop(key string) error {
	_wLock.Lock()
	delete(shared(), key)
	_wLock.Unlock()
	return synchronize()
}

func Value(key string) interface{} {
	_rwLock.RLock()
	defer _rwLock.RUnlock()
	return _map[key]
}

var (
	channel chan os.Signal
)

func RunServer() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}

func systemAbsPath() string {
	_zeroFrameworkHome := os.Getenv("ZERO_HOME")
	if len(_zeroFrameworkHome) > 0 {
		return _zeroFrameworkHome
	}
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		dir, file := path.Split(filename)
		if !strings.HasPrefix(file, _appName) {
			panic("global context must be initialized in `main func` and appname must same as `main package filename` .")
		}
		return dir
	}
	return ""
}

func AppName() string {
	return _appName
}

func InitGlobalContext(appName string) {
	if len(cfg.ServerAbsPath()) == 0 {
		_appName = appName
		cfg.NewConfigs(systemAbsPath())
		Key("zero.system.logger", log.InitLogger())
		Logger().Info(ZERO_FRAMEWORK_BANNER)
	}
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
