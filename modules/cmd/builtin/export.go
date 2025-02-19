package builtin

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Doridian/fox/modules/loader"
)

type ExportCmd struct {
}

var _ Cmd = &ExportCmd{}

func (c *ExportCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
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

	varKey := gocmd.Args[1]
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
	return 0, os.Setenv(varKey, varVal)
}

func init() {
	Register("export", func() Cmd { return &ExportCmd{} })
}
