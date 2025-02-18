package readline

import (
	"errors"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, rl *goreadline.Instance) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = rl
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func PushNew(L *lua.LState, rl *goreadline.Instance) int {
	if rl == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(ToUserdata(L, rl))
	return 1
}

func Check(L *lua.LState, i int) (*goreadline.Instance, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*goreadline.Instance); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}

func rlResultHandler(L *lua.LState, val string, err error) int {
	if err != nil {
		if errors.Is(err, goreadline.ErrInterrupt) {
			L.Push(lua.LNil)
			return 1
		}

		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(val))
	return 1
}

func fixConfig(loader *loader.LuaModule, cfg *goreadline.Config) *goreadline.Config {
	shellMod := loader.GetModule(LuaName).(*shell.Shell)
	cfg.Stdout = shellMod.Stdout()
	cfg.Stderr = shellMod.Stderr()
	cfg.Stdin = shellMod.Stdin()
	return cfg
}
