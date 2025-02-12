package pipe

import lua "github.com/yuin/gopher-lua"

const LuaType = "go://fox/pipe/Pipe"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mt := L.NewTypeMetatable(LuaType)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"read":     pipeRead,
		"canRead":  pipeCanRead,
		"write":    pipeWrite,
		"canWrite": pipeCanWrite,
		"close":    pipeClose,
	}))

	mod := L.NewTable()
	L.SetField(mod, "null", makePipe(L, &nullPipe))
	L.SetField(mod, "stdin", makePipe(L, &stdinPipe))
	L.SetField(mod, "stderr", makePipe(L, &stderrPipe))
	L.SetField(mod, "stdout", makePipe(L, &stdoutPipe))
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return "pipe"
}
