package pipe

import lua "github.com/yuin/gopher-lua"

const LuaName = "fox.pipe"
const LuaTypeName = "Pipe"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"read":     pipeRead,
		"canRead":  pipeCanRead,
		"write":    pipeWrite,
		"canWrite": pipeCanWrite,
		"close":    pipeClose,
	}))

	mod := L.NewTable()

	mod.RawSetString("null", makePipe(L, &nullPipe))
	mod.RawSetString("stdin", makePipe(L, &stdinPipe))
	mod.RawSetString("stderr", makePipe(L, &stderrPipe))
	mod.RawSetString("stdout", makePipe(L, &stdoutPipe))

	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
