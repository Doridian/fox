package shell

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/modules/loader"
	"github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

//go:embed init.lua
var initCode string

var ErrNeedMore = errors.New("Need more input")

// TODO: Handle SIGTERM

func New() *Shell {
	rl, err := readline.New("?fox?> ")
	if err != nil {
		log.Panicf("Error initializing readline: %v", err)
	}

	s := &Shell{
		l: lua.NewState(lua.Options{
			SkipOpenLibs:        true,
			IncludeGoStackTrace: true,
		}),
		rl: rl,
	}
	s.init()

	return s
}

func (s *Shell) sendInterrupt() {
	s.mainMod.Interrupt()

	cancelCtx := s.cancelCtx
	if cancelCtx != nil {
		cancelCtx()
	}
}

func (s *Shell) signalInit() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for range signals {
			s.sendInterrupt()
		}
	}()
}

func (s *Shell) Loader(L *lua.LState) int {
	mod := s.l.SetFuncs(s.l.NewTable(), map[string]lua.LGFunction{
		"exit":              luaExit,
		"readlineConfig":    s.luaSetReadlineConfig,
		"getReadlineConfig": s.luaGetReadlineConfig,

		"defaultShellParser":  luaDefaultShellParser,
		"defaultRenderPrompt": luaDefaultRenderPrompt,
	})
	s.mod = mod
	L.Push(mod)
	return 1
}

func (s *Shell) init() {
	s.l.Pop(lua.OpenBase(s.l))
	s.l.Pop(lua.OpenPackage(s.l))
	s.l.Pop(lua.OpenTable(s.l))
	// s.l.Pop(lua.OpenIo(s.l))
	// s.l.Pop(lua.OpenOs(s.l))
	s.l.Pop(lua.OpenString(s.l))
	s.l.Pop(lua.OpenMath(s.l))
	s.l.Pop(lua.OpenDebug(s.l))
	s.l.Pop(lua.OpenChannel(s.l))
	s.l.Pop(lua.OpenCoroutine(s.l))

	s.print = s.l.GetGlobal("print").(*lua.LFunction)

	mainMod := loader.NewLuaModule()
	err := mainMod.ManualRegisterModuleDefault(s)
	if err != nil {
		log.Fatalf("Error registering shell as module: %v", err)
	}
	mainMod.Load(s.l)
	s.mainMod = mainMod

	s.signalInit()

	s.startLuaLock()
	defer s.endLuaLock(false, nil)
	err = s.l.DoString(initCode)
	if err != nil {
		log.Fatalf("Error initializing shell: %v", err)
	}

	if s.l.GetTop() > 0 {
		log.Fatalf("luaInit %d left stack frames!", s.l.GetTop())
	}
}

func defaultShellParser(cmd string) (string, error) {
	if cmd[len(cmd)-2:] == "\\\n" {
		return "", ErrNeedMore
	}
	return strings.ReplaceAll(cmd, "\\\n", "\n"), nil
}

func luaDefaultShellParser(L *lua.LState) int {
	cmd := L.CheckString(1)
	parsed, err := defaultShellParser(cmd)
	if err != nil {
		if errors.Is(err, ErrNeedMore) {
			L.Push(lua.LTrue)
			return 1
		}
		L.RaiseError("Error parsing command: %v", err)
		return 0
	}
	L.Push(lua.LString(parsed))
	return 1
}

func (s *Shell) shellParser(cmd string) (string, error) {
	if s.mod == nil {
		return defaultShellParser(cmd)
	}

	shellParser := s.mod.RawGetString("parser")
	if shellParser == nil || shellParser == lua.LNil {
		return defaultShellParser(cmd)
	}

	s.startLuaLock()
	defer s.endLuaLock(false, nil)
	s.l.Push(shellParser)
	s.l.Push(lua.LString(cmd))
	err := s.l.PCall(1, 1, nil)
	if err != nil {
		log.Printf("Error in Lua shell.parser: %v", err)
		return defaultShellParser(cmd)
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)
	if parseRet == nil || parseRet == lua.LNil || parseRet == lua.LFalse {
		return defaultShellParser(cmd)
	} else if parseRet == lua.LTrue {
		return "", ErrNeedMore
	}
	return lua.LVAsString(parseRet), nil
}

func defaultRenderPrompt(lineNo int) string {
	if lineNo < 2 {
		return "fox> "
	}
	return "fo+> "
}

func luaDefaultRenderPrompt(L *lua.LState) int {
	lineNo := L.CheckInt(1)
	L.Push(lua.LString(defaultRenderPrompt(lineNo)))
	return 1
}

func (s *Shell) renderPrompt(lineNo int) string {
	if s.mod == nil {
		return defaultRenderPrompt(lineNo)
	}

	renderPrompt := s.mod.RawGetString("renderPrompt")
	if renderPrompt == nil || renderPrompt == lua.LNil {
		return defaultRenderPrompt(lineNo)
	}

	s.startLuaLock()
	defer s.endLuaLock(false, nil)
	s.l.Push(renderPrompt)
	s.l.Push(lua.LNumber(lineNo))
	err := s.l.PCall(1, 1, nil)
	if err != nil {
		log.Printf("Error in Lua shell.renderPrompt: %v", err)
		return defaultRenderPrompt(lineNo)
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)

	if parseRet == nil || parseRet == lua.LNil || parseRet == lua.LFalse {
		return defaultRenderPrompt(lineNo)
	}
	return lua.LVAsString(parseRet)
}

func (s *Shell) readLine(disp string) (string, error) {
	s.rlLock.Lock()
	defer s.rlLock.Unlock()

	s.rl.SetPrompt(disp)
	return s.rl.ReadLine()
}

func (s *Shell) RunFile(file string) error {
	s.startLuaLock()
	err := s.l.DoFile(file)
	s.endLuaLock(err == nil, err)
	return err
}

func (s *Shell) RunString(code string) error {
	if code == "" || code == "\n" {
		return nil
	}

	s.startLuaLock()
	err := s.l.DoString(code)
	s.endLuaLock(err == nil, err)
	return err
}

func (s *Shell) RunPrompt() error {
	var err error
	running := true
	for running {
		running, err = s.runPromptOne()
	}
	return err
}

func (s *Shell) startLuaLock() {
	s.lLock.Lock()
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	s.l.SetContext(s.ctx)
	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))
}

func (s *Shell) endLuaLock(printStack bool, err error) {
	defer s.lLock.Unlock()

	if printStack {
		retC := s.l.GetTop()
		if retC > 0 {
			s.l.Insert(s.print, 0)
			s.l.Call(retC, 0)
		}
	}

	cancelCtx := s.cancelCtx
	if cancelCtx != nil {
		cancelCtx()
	}
	s.ctx = nil
	s.cancelCtx = nil

	if err != nil {
		log.Printf("Lua error: %v", err)
	}
}

func (s *Shell) runPromptOne() (bool, error) {
	var luaCode string

	s.mainMod.PrePrompt()

	cmdBuilder := strings.Builder{}
	lineNo := 1
	for {
		cmdAdd, err := s.readLine(s.renderPrompt(lineNo))
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				return true, err
			}
			log.Printf("Prompt aborted: %v", err)
			return false, err
		}
		cmdBuilder.WriteString(cmdAdd)
		cmdBuilder.WriteRune('\n')
		lineNo++

		luaCode, err = s.shellParser(cmdBuilder.String())
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Internal error parsing command: %v", err)
			return false, err
		}
	}

	return true, s.RunString(luaCode)
}
