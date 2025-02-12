package file

import (
	"github.com/Doridian/fox/luautil"
	foxio "github.com/Doridian/fox/modules/io"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs"
const LuaTypeName = "File"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), luautil.MergeFuncMaps(foxio.IndexFuncs(), map[string]lua.LGFunction{
		"stat": fileStat,
	})))
	mt.RawSetString("__tostring", L.NewFunction(fileToString))
	tbl.RawSetString(LuaTypeName, mt)
}
