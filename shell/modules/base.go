package modules

import lua "github.com/yuin/gopher-lua"

type LuaModule interface {
	Init(l *lua.LState)
}
