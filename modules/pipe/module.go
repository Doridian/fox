package pipe

import (
	"github.com/Doridian/fox/modules"
	luaio "github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/util"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:pipe"
const LuaTypeName = "Pipe"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mt := L.NewTypeMetatable(LuaType)

	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), util.MergeMaps(luaio.IndexFuncs(), map[string]lua.LGFunction{
		"isNull": pipeIsNull,
	})))
	mt.RawSetString("__tostring", L.NewFunction(pipeToString))

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": pipeNew,
	})

	mod.RawSetString("null", ToUserdata(L, &nullPipe))

	mod.RawSetString(LuaTypeName, mt)

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
