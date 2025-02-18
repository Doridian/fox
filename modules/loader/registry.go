package loader

import (
	"sync"

	"github.com/Doridian/fox/modules"
)

var ctors []*ModuleCtor
var ctorLock sync.Mutex

func AddModuleDefault(ctor func(loader *LuaModule) modules.LuaModule) {
	AddModule(ctor, ModuleConfig{})
}

func AddModule(ctor func(loader *LuaModule) modules.LuaModule, cfg ModuleConfig) {
	ctorLock.Lock()
	defer ctorLock.Unlock()

	inst := &ModuleCtor{
		ctor: ctor,
		cfg:  cfg,
	}
	ctors = append(ctors, inst)
}

func (l *LuaModule) GetModule(name string) modules.LuaModule {
	ctorLock.Lock()
	defer ctorLock.Unlock()

	if inst, ok := l.gomods[name]; ok {
		return inst.mod
	}
	return nil
}
