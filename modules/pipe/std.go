package pipe

import (
	"os"
)

var stderrPipe = Pipe{
	wc:           os.Stderr,
	forwardClose: false,
}

var stdoutPipe = Pipe{
	wc:           os.Stdout,
	forwardClose: false,
}

var stdinPipe = Pipe{
	rc:           os.Stdin,
	forwardClose: false,
}

var nullPipe = Pipe{
	isNull:       true,
	forwardClose: false,
}
