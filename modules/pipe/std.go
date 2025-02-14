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
	creator:     &osPipeCreator,
	wr:          os.Stderr,
	description: "stderr",
}

var stdoutPipe = Pipe{
	creator:     &osPipeCreator,
	wr:          os.Stdout,
	description: "stdout",
}

var stdinPipe = Pipe{
	creator:     &osPipeCreator,
	rd:          os.Stdin,
	description: "stdin",
}

var nullPipe = Pipe{
	isNull: true,
}
