package config

import (
	"fmt"

	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:readline"
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

func ToString(cfg *goreadline.Config) string {
	if cfg == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s{\"%s\"}", LuaType, cfg.Prompt)
}
