package pipe

import (
	"fmt"
	"os"

	lua "github.com/yuin/gopher-lua"
)

func pipeToString(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}

	L.Push(lua.LString(p.ToString()))
	return 1
}

func pipeIsNull(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.IsNull()))
	return 1
}

func pipeNew(L *lua.LState) int {
	description := fmt.Sprintf("lua(%s)", L.OptString(1, ""))

	r, w, err := os.Pipe()
	if err != nil {
		L.RaiseError("failed to create pipe: %v", err)
		return 0
	}

	p := NewPipe(nil, description, r, w, w)
	L.Push(ToUserdata(L, p))
	return 1
}
