package pipe

import (
	"fmt"
	"io"

	lua "github.com/yuin/gopher-lua"
)

type PipeCreator interface {
	ToString() string
}

type Pipe struct {
	creator     PipeCreator
	description string

	forwardClose bool
	isNull       bool
	rc           io.ReadCloser
	wc           io.WriteCloser
}

func NewReadPipe(creator PipeCreator, description string, rc io.ReadCloser) *Pipe {
	return &Pipe{
		creator:      creator,
		description:  description,
		rc:           rc,
		forwardClose: true,
	}
}

func NewWritePipe(creator PipeCreator, description string, wc io.WriteCloser) *Pipe {
	return &Pipe{
		creator:      creator,
		description:  description,
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

func (p *Pipe) PushNew(L *lua.LState) int {
	return PushNew(L, p)
}

func (p *Pipe) ToString() string {
	if p.isNull {
		return fmt.Sprintf("%s{null}", LuaType)
	}

	mode := ""
	if p.rc != nil {
		mode += "r"
	}
	if p.wc != nil {
		mode += "w"
	}

	creatorStr := "nil"
	if p.creator != nil {
		creatorStr = p.creator.ToString()
	}

	return fmt.Sprintf("%s{%s, %s, %s}", LuaType, mode, p.description, creatorStr)
}
