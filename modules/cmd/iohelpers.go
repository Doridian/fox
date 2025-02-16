package cmd

import (
	"errors"
	"io"

	"github.com/Doridian/fox/modules/pipe"
)

func ioCanRead(pipeIo interface{}) bool {
	if pipeIo == nil {
		return false
	}

	p, ok := pipeIo.(*pipe.Pipe)
	if ok {
		return p.CanRead()
	}
	_, ok = pipeIo.(io.Reader)
	return ok
}

func ioCanWrite(pipeIo interface{}) bool {
	if pipeIo == nil {
		return false
	}

	p, ok := pipeIo.(*pipe.Pipe)
	if ok {
		return p.CanWrite()
	}
	_, ok = pipeIo.(io.Writer)
	return ok
}

func ioIsNull(pipeIo interface{}) bool {
	if pipeIo == nil {
		return false
	}

	p, ok := pipeIo.(*pipe.Pipe)
	if ok {
		return p.IsNull()
	}
	return false
}

func ioClose(pipeIo interface{}) error {
	if pipeIo == nil {
		return nil
	}

	p, ok := pipeIo.(io.Closer)
	if ok {
		return p.Close()
	}
	return errors.New("not closable")
}
