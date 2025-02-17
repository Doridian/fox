package integrated

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ExportCmd struct {
}

func (c *ExportCmd) RunAs(gocmd *exec.Cmd) (int, error) {
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

	return 0, os.Setenv(varKey, varVal)
}

func (c *ExportCmd) SetContext(ctx context.Context) {

}

var _ Cmd = &ExportCmd{}
