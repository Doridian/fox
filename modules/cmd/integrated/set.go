package integrated

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
)

type SetCmd struct {
}

var _ Cmd = &SetCmd{}

func (c *SetCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("set: missing variable name\n"))
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
	} else {
		varVal = varKey[eqPos+1:]
		varKey = varKey[:eqPos]
	}

	// TODO: Set some global table thingie
	log.Printf("VAR %s = %s", varKey, varVal)

	return 0, nil
}

func (c *SetCmd) SetContext(ctx context.Context) {

}
