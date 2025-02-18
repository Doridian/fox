package integrated

import (
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/Doridian/fox/modules/vars"
	lua "github.com/yuin/gopher-lua"
)

type SetCmd struct {
}

var _ Cmd = &SetCmd{}

func (c *SetCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("missing variable name\n"))
		return 1, nil
	}

	varKey := gocmd.Args[1]
	var varVal string

	eqPos := strings.IndexRune(varKey, '=')
	if eqPos < 0 {
		varB, err := io.ReadAll(gocmd.Stdin)
		if err != nil {
			return 1, err
		}
		varVal = string(varB)
		varVal = strings.TrimSuffix(varVal, "\n")
	} else {
		varVal = varKey[eqPos+1:]
		varKey = varKey[:eqPos]
	}

	vars.Set(varKey, lua.LString(varVal))

	return 0, nil
}

func (c *SetCmd) SetContext(ctx context.Context) {

}
