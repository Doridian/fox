package info

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:info"

var infoTable map[string]lua.LValue

func init() {
	infoTable = make(map[string]lua.LValue)
	infoTable["version"] = lua.LString(version)
	infoTable["gitrev"] = lua.LString(gitrev)
}

type LuaModule struct {
}

func newLuaModule() modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"__index": infoIndex,
		"__call":  infoCall,
	})
	L.SetMetatable(mod, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}

func Register(key string, val lua.LValue) {
	if infoTable[key] != nil {
		panic("Info key already registered: " + key)
	}
	infoTable[key] = val
}
