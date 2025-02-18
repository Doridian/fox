package integrated

import (
	"context"
	"os/exec"
	"strconv"
	"time"
)

type SleepCmd struct {
	ctx context.Context
}

var _ Cmd = &SleepCmd{}

func (c *SleepCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("missing duration\n"))
		return 1, nil
	}
	dur, err := time.ParseDuration(gocmd.Args[1])
	if err != nil {
		durF, errF := strconv.ParseFloat(gocmd.Args[1], 64)
		if errF == nil {
			dur = time.Duration(durF * float64(time.Second))
		} else {
			_, _ = gocmd.Stderr.Write([]byte(err.Error()))
			_, _ = gocmd.Stderr.Write([]byte("\n"))
			_, _ = gocmd.Stderr.Write([]byte(errF.Error()))
			_, _ = gocmd.Stderr.Write([]byte("\n"))
			return 1, nil
		}
	}

	select {
	case <-c.ctx.Done():
		return 1, nil
	case <-time.After(dur):
		return 0, nil
	}
}

func (c *SleepCmd) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func init() {
	Register("sleep", func() Cmd { return &SleepCmd{} })
}
