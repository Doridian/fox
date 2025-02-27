package fs

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/file"
	"github.com/Doridian/fox/modules/fs/fileinfo"
	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/modules/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:fs"

type LuaModule struct {
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"stat":  doStat,
		"lstat": doLStat,

		"readDir": doReadDir,

		"open":     doOpen,
		"remove":   doRemove,
		"mkdir":    doMkdir,
		"mkdirAll": doMkdirAll,

		"glob":       doGlob,
		"hasGlob":    doHasGlob,
		"globEscape": doGlobEscape,
	})
	file.Load(L, mod)
	fileinfo.Load(L, mod)
	direntry.Load(L, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{time.LuaName}
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
