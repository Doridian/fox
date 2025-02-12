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
		c.gocmd.Path = L.CheckString(2)
		if c.gocmd.Args == nil {
			c.gocmd.Args = []string{c.gocmd.Path}
		} else {
			c.gocmd.Args[0] = c.gocmd.Path
		}
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.gocmd.Path))
	return 1
}

func getSetDir(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		c.gocmd.Dir = L.CheckString(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.gocmd.Dir))
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
		args := make([]string, 1, argsLLen)
		args[0] = c.gocmd.Path
		for i := 1; i <= argsLLen; i++ {
			args = append(args, lua.LVAsString(argsL.RawGetInt(i)))
		}
		c.gocmd.Args = args
		L.Push(ud)
		return 1
	}

	ret := L.NewTable()
	for i := 1; i < len(c.gocmd.Args); i++ {
		ret.Append(lua.LString(c.gocmd.Args[i]))
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
		c.gocmd.Env = env
		L.Push(ud)
		return 1
	}

	ret := L.NewTable()
	for _, arg := range c.gocmd.Env {
		envSplit := strings.SplitN(arg, "=", 2)
		if len(envSplit) < 2 {
			continue
		}
		ret.RawSetString(envSplit[0], lua.LString(envSplit[1]))
	}
	L.Push(ret)
	return 1
}
