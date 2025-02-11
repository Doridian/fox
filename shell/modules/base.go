package modules

import lua "github.com/yuin/gopher-lua"

type Module interface {
	Init(l *lua.LState)
}
