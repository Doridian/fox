package shell

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/modules/cmd"
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
	err := mainMod.ManualRegisterModule(s, nil)
	if err != nil {
		log.Fatalf("Error registering shell as module: %v", err)
	}
	mainMod.Load(s.l)
	s.mainMod = mainMod

	s.signalInit()

	s.startLuaLock()
	err = s.l.DoString(initCode)
	if err != nil {
		log.Fatalf("Error initializing shell: %v", err)
	}

	if s.l.GetTop() > 0 {
		log.Fatalf("luaInit %d left stack frames!", s.l.GetTop())
	}
	s.endLuaLock(nil)
}

func defaultShellParser(cmd string) (string, error) {
	if cmd[len(cmd)-1] == '\\' {
		return "", ErrNeedMore
	}
	return strings.ReplaceAll(cmd, "\\\n", "\n"), nil
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
	defer s.endLuaLock(nil)
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

func (s *Shell) renderPrompt(lineNo int) string {
	if s.mod == nil {
		return defaultRenderPrompt(lineNo)
	}

	renderPrompt := s.mod.RawGetString("renderPrompt")
	if renderPrompt == nil || renderPrompt == lua.LNil {
		return defaultRenderPrompt(lineNo)
	}

	s.startLuaLock()
	defer s.endLuaLock(nil)
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

func (s *Shell) RunScript(file string) int {
	s.startLuaLock()
	err := s.l.DoFile(file)
	return s.endLuaLock(err)
}

func (s *Shell) RunString(code string) int {
	if code == "" || code == "\n" {
		return 0
	}

	s.startLuaLock()
	err := s.l.DoString(code)
	return s.endLuaLock(err)
}

func (s *Shell) RunPrompt() bool {
	s.mainMod.PrePrompt()

	res, err := s.readLine(s.renderPrompt(1))
	if err != nil {
		if errors.Is(err, readline.ErrInterrupt) {
			return true
		}
		log.Printf("Prompt aborted: %v", err)
		return false
	}
	if res != "" {
		exitCode := s.runPromptInt(res)
		if exitCode != 0 {
			log.Printf("Exit code: %v", exitCode)
		}
	}
	return true
}

func (s *Shell) startLuaLock() {
	s.lLock.Lock()
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	s.l.SetContext(s.ctx)
	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))
}

func (s *Shell) endLuaLock(err error) int {
	exitCode := int(lua.LVAsNumber(s.l.GetGlobal("_LAST_EXIT_CODE")))

	if err == nil {
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
	s.lLock.Unlock()

	if err != nil {
		if exitCode == 0 {
			exitCode = cmd.ExitCodeInternalShellError
		}
		log.Printf("Lua error: %v", err)
		return exitCode
	}

	return exitCode
}

func (s *Shell) runPromptInt(firstLine string) int {
	var luaCode string
	var err error

	cmdBuilder := strings.Builder{}
	cmdBuilder.WriteString(firstLine)

	lineNo := 1

	for {
		cmdBuilder.WriteRune('\n')

		luaCode, err = s.shellParser(cmdBuilder.String())
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return 0
		}

		lineNo++
		cmdAdd, err := s.readLine(s.renderPrompt(lineNo))
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				return 0
			}
			log.Printf("Prompt aborted: %v", err)
			os.Exit(0)
			return 0
		}
		cmdBuilder.WriteString(cmdAdd)
	}

	return s.RunString(luaCode)
}
