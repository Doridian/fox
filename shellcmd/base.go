package shellcmd

import (
	"errors"
	"log"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/prompt"
	lua "github.com/yuin/gopher-lua"
)

//go:embed shell.lua
var shell string

type ShellManager struct {
	L *lua.LState
}

func NewShell() *ShellManager {
	s := &ShellManager{
		L: lua.NewState(),
	}
	s.L.DoString(shell)
	return s
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
func (s *ShellManager) Run(p *prompt.PromptManager, cmd string) (bool, int) {
	var luaCode string
	var err error
	for {
		luaCode, err = s.CommandToLua(cmd)
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return false, 0
		}
		cmdAdd, err := p.Prompt("f.x> ")
		if err != nil {
			log.Printf("Prompt aborted: %v", err)
			return true, 0
		}
		cmd += cmdAdd
	}

	s.L.SetGlobal("_LAST_DO_EXIT", lua.LBool(false))
	s.L.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))

	err = s.L.DoString(luaCode)
	if err != nil {
		log.Printf("Error running command: %v", err)
		return false, 1
	}
	retC := s.L.GetTop()
	if retC > 0 {
		s.L.Pop(retC)
	}

	doExit := lua.LVAsBool(s.L.GetGlobal("_LAST_DO_EXIT"))
	exitCode := lua.LVAsNumber(s.L.GetGlobal("_LAST_EXIT_CODE"))
	return doExit, int(exitCode)
}
