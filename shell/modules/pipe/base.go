package pipe

import (
	"io"

	lua "github.com/yuin/gopher-lua"
)

type Pipe struct {
	creator interface{}

	forwardClose bool
	isNull       bool
	rc           io.ReadCloser
	wc           io.WriteCloser
}

func NewReadPipe(creator interface{}, rc io.ReadCloser) *Pipe {
	return &Pipe{
		creator:      creator,
		rc:           rc,
		forwardClose: true,
	}
}

func NewWritePipe(creator interface{}, wc io.WriteCloser) *Pipe {
	return &Pipe{
		creator:      creator,
		wc:           wc,
		forwardClose: true,
	}
}

func (p *Pipe) Close() {
	if !p.forwardClose {
		return
	}

	if p.rc != nil {
		p.rc.Close()
	}
	if p.wc != nil {
		p.wc.Close()
	}
}

func (p *Pipe) IsNull() bool {
	return p.isNull
}

func (p *Pipe) CanRead() bool {
	return p.rc != nil || p.isNull
}

func (p *Pipe) CanWrite() bool {
	return p.wc != nil || p.isNull
}

func (p *Pipe) GetReader() io.ReadCloser {
	return p.rc
}

func (p *Pipe) GetWriter() io.WriteCloser {
	return p.wc
}

func (p *Pipe) Creator() interface{} {
	return p.creator
}

func (p *Pipe) Push(L *lua.LState) int {
	return pushPipe(L, p)
}
