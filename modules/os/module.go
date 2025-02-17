package os

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:os"

type LuaModule struct {
}

func newLuaModule() modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"executable": osExecutable,
		"chdir":      osChdir,
		"getwd":      osGetwd,
	})
	mod.RawSetString("stdin", pipe.ToUserdata(L, stdinPipe))
	mod.RawSetString("stderr", pipe.ToUserdata(L, stderrPipe))
	mod.RawSetString("stdout", pipe.ToUserdata(L, stdoutPipe))
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
