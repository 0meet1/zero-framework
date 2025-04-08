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
