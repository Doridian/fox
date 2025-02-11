package shellcmd

import (
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func getSetErrorPropagation(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		c.ErrorPropagation = L.CheckBool(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LBool(c.ErrorPropagation))
	return 1
}

func getSetPath(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		c.Gocmd.Path = L.CheckString(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.Gocmd.Path))
	return 1
}

func getSetDir(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		c.Gocmd.Dir = L.CheckString(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.Gocmd.Dir))
	return 1
}

func getSetArgs(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		argsL := L.CheckTable(2)
		if argsL == nil {
			return 0
		}
		argsLLen := argsL.MaxN()
		args := make([]string, 0, argsLLen)
		for i := 1; i <= argsLLen; i++ {
			args = append(args, lua.LVAsString(argsL.RawGetInt(i)))
		}
		c.Gocmd.Args = args
		L.Push(ud)
		return 1
	}

	ret := L.NewTable()
	for _, arg := range c.Gocmd.Args {
		ret.Append(lua.LString(arg))
	}
	L.Push(ret)
	return 1
}

func getSetEnv(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		envL := L.CheckTable(2)
		if envL == nil {
			return 0
		}
		env := make([]string, 0)
		envK, envV := envL.Next(lua.LNil)
		for envK != lua.LNil {
			env = append(env, fmt.Sprintf("%s=%s", lua.LVAsString(envK), lua.LVAsString(envV)))
			envK, envV = envL.Next(envK)
		}
		c.Gocmd.Env = env
		L.Push(ud)
		return 1
	}

	ret := L.NewTable()
	for _, arg := range c.Gocmd.Env {
		envSplit := strings.SplitN(arg, "=", 2)
		if len(envSplit) < 2 {
			continue
		}
		ret.RawSetString(envSplit[0], lua.LString(envSplit[1]))
	}
	L.Push(ret)
	return 1
}
