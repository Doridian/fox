package config

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.readline"
const LuaTypeName = "Config"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"prompt":                    cfgSetPrompt,
		"getPrompt":                 cfgGetPrompt,
		"vimMode":                   cfgSetVimMode,
		"getVimMpde":                cfgGetVimMode,
		"historyFile":               cfgSetHistoryFile,
		"getHistoryFile":            cfgGetHistoryFile,
		"historyLimit":              cfgSetHistoryLimit,
		"getHistoryLimit":           cfgGetHistoryLimit,
		"historySearchFold":         cfgSetHistorySearchFold,
		"getHistorySearchFold":      cfgGetHistorySearchFold,
		"disableAutoSaveHistory":    cfgSetDisableAutoSaveHistory,
		"getDisableAutoSaveHistory": cfgGetDisableAutoSaveHistory,
		"enableMask":                cfgSetEnableMask,
		"getEnableMask":             cfgGetEnableMask,
		"maskRune":                  cfgSetMaskRune,
		"getMaskRune":               cfgGetMaskRune,
	}))
	mt.RawSetString("__tostring", L.NewFunction(cfgToString))
	tbl.RawSetString(LuaTypeName, mt)
}
