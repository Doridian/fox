package shell

import (
	"context"
	"errors"
	"fmt"
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

var ErrNeedMore = errors.New("need more input")
var ErrShellNotInited = errors.New("shell not initialized")

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
		rl:         rl,
		ShowErrors: true,
	}

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
	if s.args != nil {
		argsL := s.l.NewTable()
		for _, arg := range s.args {
			argsL.Append(lua.LString(arg))
		}
		s.mod.RawSetString("args", argsL)
	}
	L.Push(mod)
	return 1
}

func (s *Shell) MustInit(args []string) {
	err := s.Init(args)
	if err != nil {
		log.Fatalf("Error initializing shell: %v", err)
	}
}

func (s *Shell) Init(args []string) error {
	if s.args != nil {
		return errors.New("shell already initialized")
	}
	s.args = args

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
		return fmt.Errorf("error registering shell as module: %w", err)
	}
	mainMod.Load(s.l)
	s.mainMod = mainMod

	s.signalInit()

	s.startLuaLock()
	defer s.endLuaLock(false, nil)
	err = s.l.DoString(initCode)
	if err != nil {
		return fmt.Errorf("error initializing shell: %w", err)
	}

	if s.l.GetTop() > 0 {
		return fmt.Errorf("luaInit %d left stack frames", s.l.GetTop())
	}

	return nil
}

func defaultShellParser(cmdAdd string, _ int, prev lua.LValue) (lua.LValue, bool, *string) {
	if cmdAdd == "" {
		return prev, false, nil
	}

	prevStr := ""
	if prev != nil && prev != lua.LNil {
		prevStr = lua.LVAsString(prev)
	}

	if strings.HasSuffix(cmdAdd, "\\") {
		cmd := prevStr + cmdAdd[:len(cmdAdd)-1] + "\n"
		return lua.LString(cmd), true, nil
	}

	cmd := prevStr + cmdAdd + "\n"
	if strings.HasPrefix(cmd, "--") {
		return lua.LString(cmd), true, nil
	}

	if cmd[0] == '=' {
		return lua.LString("return " + cmd[1:]), false, nil
	}
	return lua.LString(cmd), false, nil
}

func luaDefaultShellParser(L *lua.LState) int {
	cmd := L.CheckString(1)
	lineNo := L.OptInt(2, 1)
	prev := L.Get(3)
	parsed, needMore, promptOverride := defaultShellParser(cmd, lineNo, prev)
	L.Push(parsed)
	L.Push(lua.LBool(needMore))
	if promptOverride == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(*promptOverride))
	}
	return 3
}

func (s *Shell) shellParser(cmdAdd string, lineNo int, prev lua.LValue) (lua.LValue, bool, *string) {
	hasPrev := prev != nil && prev != lua.LNil
	if s.mod == nil || (hasPrev && prev.Type() == lua.LTString) {
		return defaultShellParser(cmdAdd, lineNo, prev)
	}

	if !hasPrev && strings.HasPrefix(cmdAdd, "--[[DEFAULT]]") {
		return defaultShellParser(cmdAdd, lineNo, prev)
	}

	shellParser := s.mod.RawGetString("parser")
	if shellParser == nil || shellParser == lua.LNil {
		return defaultShellParser(cmdAdd, lineNo, prev)
	}

	s.startLuaLock()
	defer s.endLuaLock(false, nil)
	s.l.Push(shellParser)
	s.l.Push(lua.LString(cmdAdd))
	s.l.Push(lua.LNumber(lineNo))
	s.l.Push(prev)
	err := s.l.PCall(3, 3, nil)
	if err != nil {
		if s.ShowErrors {
			log.Printf("Error in Lua shell.parser: %v", err)
		}
		return lua.LNil, false, nil
	}
	parseRet := s.l.Get(-3)
	needMoreL := s.l.OptBool(-2, false)
	promptOverride := s.l.Get(-1)
	s.l.Pop(3)

	if needMoreL {
		if parseRet == lua.LNil || parseRet == nil {
			s.l.RaiseError("needMore is true but parser returned nil, this is not allowed")
		} else if parseRet.Type() == lua.LTString {
			s.l.RaiseError("needMore is true but parser returned a string, this is not allowed")
		}
	}

	var promptOverrideRes *string
	if promptOverride != nil && promptOverride != lua.LNil {
		str := lua.LVAsString(promptOverride)
		promptOverrideRes = &str
	}

	if parseRet == nil || parseRet == lua.LNil || parseRet == lua.LFalse {
		cmd, needMore, _ := defaultShellParser(cmdAdd, lineNo, prev)
		return cmd, needMore, promptOverrideRes
	}
	return parseRet, needMoreL, promptOverrideRes
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
		if s.ShowErrors {
			log.Printf("Error in Lua shell.renderPrompt: %v", err)
		}
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
	if s.args == nil {
		return ErrShellNotInited
	}

	s.startLuaLock()
	err := s.l.DoFile(file)
	s.endLuaLock(err == nil, err)
	return err
}

func (s *Shell) RunString(code string) error {
	if s.args == nil {
		return ErrShellNotInited
	}

	if code == "" || code == "\n" {
		return nil
	}

	s.startLuaLock()
	err := s.l.DoString(code)
	s.endLuaLock(err == nil, err)
	return err
}

func (s *Shell) RunCommand(cmd string) error {
	if s.args == nil {
		return ErrShellNotInited
	}

	s.startLuaLock()
	s.l.Push(s.mod.RawGetString("runCommand"))
	s.l.Push(lua.LString(cmd))
	err := s.l.PCall(1, 0, nil)
	s.endLuaLock(false, err)

	return err
}

func (s *Shell) RunFunc(luaFunc *lua.LFunction) error {
	if s.args == nil {
		return ErrShellNotInited
	}

	s.startLuaLock()
	s.l.Push(luaFunc)
	err := s.l.PCall(0, 0, nil)
	s.endLuaLock(false, err)

	return err
}

func (s *Shell) RunPrompt() error {
	if s.args == nil {
		return ErrShellNotInited
	}

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

	if s.ShowErrors && err != nil {
		log.Printf("Lua error: %v", err)
	}
}

func (s *Shell) runPromptOne() (bool, error) {
	var parseRes lua.LValue = lua.LNil

	s.mainMod.PrePrompt()

	lineNo := 1
	needMore := true

	var nextPromptOverride *string
	var nextPrompt string
	for needMore {
		if nextPromptOverride != nil {
			nextPrompt = *nextPromptOverride
		} else {
			nextPrompt = s.renderPrompt(lineNo)
		}
		cmdAdd, err := s.readLine(nextPrompt)
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				return true, err
			}
			err = fmt.Errorf("prompt aborted: %w", err)
			if s.ShowErrors {
				log.Println(err.Error())
			}
			return false, err
		}

		parseRes, needMore, nextPromptOverride = s.shellParser(cmdAdd, lineNo, parseRes)
		lineNo++
	}

	if parseRes == nil || parseRes == lua.LNil {
		return true, nil
	}

	switch typedLuaCode := parseRes.(type) {
	case lua.LString:
		return true, s.RunString(typedLuaCode.String())
	case *lua.LFunction:
		return true, s.RunFunc(typedLuaCode)
	}

	err := fmt.Errorf("unknown shell.parser return type: %v", parseRes.Type().String())
	if s.ShowErrors {
		log.Println(err.Error())
	}
	return true, err
}
