package os

import (
	"os"

	"github.com/Doridian/fox/modules/pipe"
)

var osPipeCreator = pipe.FixedPipeCreator{
	Name: LuaName + ":os",
}

var stderrPipe = pipe.NewPipe(&osPipeCreator, "stderr", nil, os.Stderr, nil)
var stdoutPipe = pipe.NewPipe(&osPipeCreator, "stdout", nil, os.Stdout, nil)
var stdinPipe = pipe.NewPipe(&osPipeCreator, "stdin", os.Stdin, nil, nil)
