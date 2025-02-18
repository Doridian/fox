package cmd

import (
	"errors"
	"io"
)

func ioCanRead(pipeIo interface{}) bool {
	if pipeIo == nil {
		return false
	}

	_, ok := pipeIo.(io.Reader)
	return ok
}

func ioCanWrite(pipeIo interface{}) bool {
	if pipeIo == nil {
		return false
	}

	_, ok := pipeIo.(io.Writer)
	return ok
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
