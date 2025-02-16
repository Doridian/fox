package env

import (
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// __index(t, k)
func envIndex(L *lua.LState) int {
	k := L.CheckString(2)

	v, ok := os.LookupEnv(k)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(v))
	return 1
}

// __newindex(t, k, v)
func envNewIndex(L *lua.LState) int {
	k := L.CheckString(2)
	v := L.CheckString(3)

	err := os.Setenv(k, v)
	if err != nil {
		L.RaiseError("%v", err)
	}
	return 0
}

// __call()
func envCall(L *lua.LState) int {
	ret := L.NewTable()

	for _, env := range os.Environ() {
		spl := strings.SplitN(env, "=", 2)
		k := spl[0]
		v := ""
		if len(spl) > 1 {
			v = spl[1]
		}

		ret.RawSetString(k, lua.LString(v))
	}

	L.Push(ret)
	return 1
}
