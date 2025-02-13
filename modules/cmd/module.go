package cmd

import (
	"sync"

	"github.com/Doridian/fox/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.cmd"
const LuaTypeName = "Cmd"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
	awaitedCmds    map[*Cmd]bool
	awaitedCmdLock sync.Mutex
}

func NewLuaModule() *LuaModule {
	return &LuaModule{
		awaitedCmds: make(map[*Cmd]bool),
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"getDir": getDir,
		"dir":    setDir,
		"getCmd": getCmd,
		"cmd":    setCmd,
		"getEnv": getCmd,
		"env":    setEnv,

		"stdout":     setStdout,
		"getStdout":  getStdout,
		"stdoutPipe": acquireStdoutPipe,
		"stderr":     setStderr,
		"getStderr":  getStderr,
		"stderrPipe": acquireStderrPipe,
		"stdin":      setStdin,
		"getStdin":   getStdin,
		"stdinPipe":  acquireStdinPipe,

		"run":   doRun,
		"start": doStart,
		"wait":  doWait,

		"getErrorPropagation": getErrorPropagation,
		"errorPropagation":    setErrorPropagation,
		"getAutoLookPath":     getAutoLookPath,
		"autoLookPath":        setAutoLookPath,
	}))
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__tostring": cmdToString,
		"__call":     doRun,
	})

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new":      m.newCmd,
		"run":      m.runCmd,
		"start":    m.startCmd,
		"lookPath": lookPath,
	})

	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{pipe.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt(all bool) bool {
	m.awaitedCmdLock.Lock()
	defer m.awaitedCmdLock.Unlock()

	triedKill := false
	toDelete := make([]*Cmd, 0)
	for cmd := range m.awaitedCmds {
		if cmd.gocmd.Process != nil {
			cmd.gocmd.Process.Kill()
			triedKill = true
		}

		if cmd.gocmd.ProcessState != nil {
			toDelete = append(toDelete, cmd)
		}
	}

	for _, cmd := range toDelete {
		delete(m.awaitedCmds, cmd)
	}

	return triedKill
}

func (m *LuaModule) AwaitCmd(cmd *Cmd) {
	m.awaitedCmdLock.Lock()
	m.awaitedCmds[cmd] = true
	m.awaitedCmdLock.Unlock()
}

func (m *LuaModule) StopAwaitCmd(cmd *Cmd) {
	m.awaitedCmdLock.Lock()
	delete(m.awaitedCmds, cmd)
	m.awaitedCmdLock.Unlock()
}
