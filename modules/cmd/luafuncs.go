package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func getRaiseForBadExit(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LBool(c.RaiseForBadExit)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setRaiseForBadExit(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckBool(2)
	c.lock.Lock()
	c.RaiseForBadExit = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getAutoLookPath(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LBool(c.AutoLookPath)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setAutoLookPath(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckBool(2)
	c.lock.Lock()
	c.AutoLookPath = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getDir(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LString(c.gocmd.Dir)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setDir(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckString(2)
	c.lock.Lock()
	c.gocmd.Dir = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getArgs(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	ret := L.NewTable()
	for _, arg := range c.gocmd.Args {
		ret.Append(lua.LString(arg))
	}
	c.lock.RUnlock()
	L.Push(ret)
	return 1
}

func setArgs(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	argsL := L.CheckTable(2)
	if argsL == nil {
		return 0
	}
	err := c.setArgs(argsL)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(ud)
	return 1
}

func (c *Cmd) setArgs(argsL *lua.LTable) error {
	if argsL == nil || argsL == lua.LNil {
		c.lock.Lock()
		c.gocmd.Args = nil
		c.lock.Unlock()
		return nil
	}

	argsLLen := argsL.MaxN()
	if argsLLen < 1 {
		return errors.New("cmd must have at least one argument (the process binary)")
	}

	args := make([]string, 0, argsLLen)
	for i := 1; i <= argsLLen; i++ {
		args = append(args, lua.LVAsString(argsL.RawGetInt(i)))
	}

	c.lock.Lock()
	c.gocmd.Args = args
	c.lock.Unlock()
	return nil
}

func getEnv(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	ret := L.NewTable()
	for _, arg := range c.gocmd.Env {
		envSplit := strings.SplitN(arg, "=", 2)
		if len(envSplit) < 2 {
			continue
		}
		ret.RawSetString(envSplit[0], lua.LString(envSplit[1]))
	}
	c.lock.RUnlock()
	L.Push(ret)
	return 1
}

func setEnv(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	envL := L.CheckTable(2)
	if envL == nil {
		return 0
	}

	c.setEnv(envL)

	L.Push(ud)
	return 1
}

func (c *Cmd) setEnv(envL *lua.LTable) {
	if envL == nil || envL == lua.LNil {
		c.lock.Lock()
		c.gocmd.Env = nil
		c.lock.Unlock()
		return
	}

	env := make([]string, 0)
	envK, envV := envL.Next(lua.LNil)
	for envK != nil && envK != lua.LNil {
		env = append(env, fmt.Sprintf("%s=%s", lua.LVAsString(envK), lua.LVAsString(envV)))
		envK, envV = envL.Next(envK)
	}

	c.lock.Lock()
	c.gocmd.Env = env
	c.lock.Unlock()
}

func cmdToString(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}
	L.Push(lua.LString(c.ToString()))
	return 1
}

func lookPath(L *lua.LState) int {
	cmd := L.CheckString(1)
	path, err := exec.LookPath(cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(path))
	return 1
}

func (m *LuaModule) newCmdInt(L *lua.LState) (*Cmd, *lua.LUserData) {
	c := &Cmd{
		mod:             m,
		gocmd:           &exec.Cmd{},
		closeQueue:      make([]io.Closer, 0),
		AutoLookPath:    true,
		RaiseForBadExit: false,
	}

	// new|run|start([args, [dir, [env]]])
	err := c.setArgs(L.OptTable(1, nil))
	if err != nil {
		L.RaiseError("%v", err)
		return nil, nil
	}
	c.gocmd.Dir = L.OptString(2, "")
	c.setEnv(L.OptTable(3, nil))

	return c, ToUserdata(L, c)
}

func (m *LuaModule) newCmd(L *lua.LState) int {
	_, ud := m.newCmdInt(L)
	L.Push(ud)
	return 1
}

func (m *LuaModule) runCmd(L *lua.LState) int {
	c, ud := m.newCmdInt(L)
	L.Push(ud)
	return c.doRun(L) + 1
}

func (m *LuaModule) startCmd(L *lua.LState) int {
	c, ud := m.newCmdInt(L)
	L.Push(ud)
	return c.doStart(L) + 1
}

func (m *LuaModule) getRunning(L *lua.LState) int {
	res := L.NewTable()
	for job := range allCmds {
		res.Append(ToUserdata(L, job))
	}
	L.Push(res)
	return 1
}
