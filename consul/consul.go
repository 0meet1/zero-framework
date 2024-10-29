package consul

import (
	"github.com/0meet1/zero-framework/global"
)

const (
	DSC_LOCK_TRUNK = "zero.dsc.lock.trunk"
)

func RunDCSMutex() {
	dcsMutex := &xZeroDCSMutexTrunk{}
	err := dcsMutex.runDCSMutex()
	if err != nil {
		panic(err)
	}
	global.Key(DSC_LOCK_TRUNK, dcsMutex)
}

func RunServeRegistry() {
	serveRegistry := &xZeroServeRegistryTrunk{}
	err := serveRegistry.runServeRegistry()
	if err != nil {
		panic(err)
	}
	global.Key(REGISTRY_TRUNK, serveRegistry)
}

// func main() {
// 	client, err := api.NewClient(api.DefaultConfig())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			// 1. 获取锁
// 			// 锁具备默认的超时时间
// 			// 如果想自定义超时时间以及一些额外的配置可以用 LockOpts
// 			// client.LockOpts()
// 			lock, err := client.LockKey("lock")
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			// 2. 尝试锁住
// 			ch := make(chan struct{})
// 			defer close(ch)

// 			_, err = lock.Lock(ch)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			defer lock.Unlock()

// 			// 3. 操作需要的资源，如修改consul里面的数据
// 			// .. 此处的操作为原子性的.
// 			// 更新某个key.
// 			val, _, err := client.KV().Get("key", nil)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			var v string
// 			if val == nil {
// 				v = "0"
// 			} else {
// 				v = string(val.Value) + "0"
// 			}
// 			_, err = client.KV().Put(&api.KVPair{
// 				Key:   "key",
// 				Value: []byte(v),
// 			}, nil)
// 			if err != nil {
// 				log.Fatal(err)
// 				return
// 			}
// 		}()
// 	}

// 	wg.Wait()

// 	val, _, err := client.KV().Get("key", nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	log.Println(string(val.Value))
// }
