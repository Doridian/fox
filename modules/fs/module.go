package fs

import (
	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/fileinfo"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	fileinfo.Load(L)
	direntry.Load(L)

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
