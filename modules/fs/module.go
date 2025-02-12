package fs

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/file"
	"github.com/Doridian/fox/modules/fs/fileinfo"
	foxtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	modules.RequireDependencies(L, m)

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"stat":  doStat,
		"lstat": doLStat,

		"readDir": doReadDir,

		"open": doOpen,
	})
	file.Load(L, mod)
	fileinfo.Load(L, mod)
	direntry.Load(L, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{foxtime.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}
