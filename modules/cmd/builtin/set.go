package builtin

import (
	"context"
	"errors"
	"flag"
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
	flags := flag.NewFlagSet("set", flag.ContinueOnError)
	flags.SetOutput(gocmd.Stderr)
	rawPtr := flags.Bool("r", false, "raw (do not strip trailing newline)")
	err := flags.Parse(gocmd.Args[1:])
	if err != nil {
		return 1, err
	}

	if flags.NArg() < 1 {
		return 1, errors.New("missing variable name")
	}

	varKey := flags.Arg(0)
	var varVal string

	eqPos := strings.IndexRune(varKey, '=')
	if eqPos < 0 {
		varB, err := io.ReadAll(gocmd.Stdin)
		if err != nil {
			return 1, err
		}
		varVal = string(varB)
	} else {
		varVal = varKey[eqPos+1:]
		varKey = varKey[:eqPos]
	}

	if !*rawPtr {
		varVal = strings.TrimSuffix(varVal, "\n")
	}
	vars.Set(varKey, lua.LString(varVal))

	return 0, nil
}

func (c *SetCmd) SetContext(ctx context.Context) {

}

func init() {
	Register("set", func() Cmd { return &SetCmd{} })
}
