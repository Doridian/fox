package duration

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func durationParse(L *lua.LState) int {
	dStr := L.CheckString(1)
	d, err := time.ParseDuration(dStr)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return Push(L, d)
}
