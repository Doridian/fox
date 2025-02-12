package pipe

import lua "github.com/yuin/gopher-lua"

const LuaType = "go://fox/pipe/Pipe"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Init(L *lua.LState) {
	funcs := map[string]lua.LGFunction{
		"read":     pipeRead,
		"canRead":  pipeCanRead,
		"write":    pipeWrite,
		"canWrite": pipeCanWrite,
		"close":    pipeClose,
	}

	mt := L.NewTypeMetatable(LuaType)
	L.SetGlobal("pipe", mt)
	L.SetField(mt, "null", makePipe(L, &nullPipe))
	L.SetField(mt, "stdin", makePipe(L, &stdinPipe))
	L.SetField(mt, "stderr", makePipe(L, &stderrPipe))
	L.SetField(mt, "stdout", makePipe(L, &stdoutPipe))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}
