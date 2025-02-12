package shell

import (
	"errors"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/prompt"
	"github.com/Doridian/fox/shell/modules/shellcmd"
	lua "github.com/yuin/gopher-lua"
)

//go:embed init.lua
var initCode string

type ShellManager struct {
	l *lua.LState
}

func NewShellManager() *ShellManager {
	s := &ShellManager{
		l: lua.NewState(),
	}
	s.init()
	return s
}

func (s *ShellManager) init() {
	mod := shellcmd.NewLuaModule()
	mod.Init(s.l)
	err := s.l.DoString(initCode)
	if err != nil {
		log.Fatalf("Error initializing shell: %v", err)
	}
}

// if .. then .. end
// while .. do .. end
// repeat .. until ..
// for .. do .. end
// .. = ..
// ( .. )

var ErrNeedMore = errors.New("Need more input")

func (s *ShellManager) CommandToLua(cmd string) (string, error) {
	lua := strings.Builder{}
	lua.WriteString(cmd)
	return lua.String(), nil
}

// Return true to exit shell
func (s *ShellManager) Run(p *prompt.PromptManager, cmd string) int {
	var luaCode string
	var err error
	for {
		luaCode, err = s.CommandToLua(cmd)
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return 0
		}
		cmdAdd, err := p.Prompt("f.x> ")
		if err != nil {
			log.Printf("Prompt aborted: %v", err)
			os.Exit(0)
			return 0
		}
		cmd += cmdAdd
	}

	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))
	err = s.l.DoString(luaCode)
	exitCode := int(lua.LVAsNumber(s.l.GetGlobal("_LAST_EXIT_CODE")))

	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}
		log.Printf("Error running command (code %d): %v", exitCode, err)
		return exitCode
	}

	retC := s.l.GetTop()
	if retC > 0 {
		s.l.Pop(retC)
	}

	return exitCode
}
