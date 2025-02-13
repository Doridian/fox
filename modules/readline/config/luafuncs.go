package config

import (
	lua "github.com/yuin/gopher-lua"
)

func cfgSetPrompt(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.Prompt = L.CheckString(2)
	L.Push(ud)
	return 1
}

func cfgGetPrompt(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LString(cfg.Prompt))
	return 1
}

func cfgGetHistoryFile(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LString(cfg.HistoryFile))
	return 1
}

func cfgSetHistoryFile(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.HistoryFile = L.CheckString(2)
	L.Push(ud)
	return 1
}

func cfgGetHistoryLimit(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LNumber(cfg.HistoryLimit))
	return 1
}

func cfgSetHistoryLimit(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.HistoryLimit = int(L.CheckNumber(2))
	L.Push(ud)
	return 1
}

func cfgSetVimMode(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.VimMode = L.CheckBool(2)
	L.Push(ud)
	return 1
}

func cfgGetVimMode(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LBool(cfg.VimMode))
	return 1
}

func cfgSetDisableAutoSaveHistory(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.DisableAutoSaveHistory = L.CheckBool(2)
	L.Push(ud)
	return 1
}

func cfgGetDisableAutoSaveHistory(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LBool(cfg.DisableAutoSaveHistory))
	return 1
}

func cfgSetHistorySearchFold(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.HistorySearchFold = L.CheckBool(2)
	L.Push(ud)
	return 1
}

func cfgGetHistorySearchFold(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LBool(cfg.HistorySearchFold))
	return 1
}

func cfgSetEnableMask(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg.EnableMask = L.CheckBool(2)
	L.Push(ud)
	return 1
}

func cfgGetEnableMask(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LBool(cfg.EnableMask))
	return 1
}

func cfgSetMaskRune(L *lua.LState) int {
	cfg, ud := Check(L, 1)
	if cfg == nil {
		return 0
	}

	str := L.CheckString(2)
	runes := []rune(str)
	if len(runes) != 1 {
		L.ArgError(2, "expect 1 rune")
		return 0
	}

	cfg.MaskRune = runes[0]
	L.Push(ud)
	return 1
}

func cfgGetMaskRune(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}

	L.Push(lua.LString(cfg.MaskRune))
	return 1
}

func cfgToString(L *lua.LState) int {
	cfg, _ := Check(L, 1)
	if cfg == nil {
		return 0
	}
	L.Push(lua.LString(ToString(cfg)))
	return 1
}
