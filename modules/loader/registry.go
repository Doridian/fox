package loader

import (
	"sync"

	"github.com/Doridian/fox/modules"
)

var ctors []*ModuleCtor
var ctorLock sync.Mutex

func AddModuleDefault(ctor func() modules.LuaModule) {
	AddModule(ctor, DefaultConfig())
}

func AddModule(ctor func() modules.LuaModule, cfg ModuleConfig) {
	ctorLock.Lock()
	defer ctorLock.Unlock()

	inst := &ModuleCtor{
		ctor: ctor,
		cfg:  cfg,
	}
	ctors = append(ctors, inst)
}
