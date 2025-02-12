package pipe

import (
	"os"
)

type FixedPipeCreator struct {
	Name string
}

func (s *FixedPipeCreator) ToString() string {
	return s.Name
}

var osPipeCreator = FixedPipeCreator{
	Name: LuaName + ":os",
}

var stderrPipe = Pipe{
	creator:      &osPipeCreator,
	wc:           os.Stderr,
	description:  "stderr",
	forwardClose: false,
}

var stdoutPipe = Pipe{
	creator:      &osPipeCreator,
	wc:           os.Stdout,
	description:  "stdout",
	forwardClose: false,
}

var stdinPipe = Pipe{
	creator:      &osPipeCreator,
	rc:           os.Stdin,
	description:  "stdin",
	forwardClose: false,
}

var nullPipe = Pipe{
	isNull:       true,
	forwardClose: false,
}
