package pipe

import lua "github.com/yuin/gopher-lua"

const luaPipeType = "go://fox/pipe/Pipe"

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

	mt := L.NewTypeMetatable(luaPipeType)
	L.SetGlobal("pipe", mt)
	L.SetField(mt, "null", L.NewFunction(newNullPipe))
	L.SetField(mt, "stdin", L.NewFunction(newStdinPipe))
	L.SetField(mt, "stderr", L.NewFunction(newStderrPipe))
	L.SetField(mt, "stdout", L.NewFunction(newStdoutPipe))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}
