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

	isNull bool
	cl     io.Closer
	rd     io.Reader
	wr     io.Writer
}

func NewReadPipe(creator PipeCreator, description string, rc io.ReadCloser) *Pipe {
	return &Pipe{
		creator:     creator,
		description: description,
		rd:          rc,
		cl:          rc,
	}
}

func NewPipe(creator PipeCreator, description string, rd io.Reader, wr io.Writer, cl io.Closer) *Pipe {
	return &Pipe{
		creator:     creator,
		description: description,
		rd:          rd,
		wr:          wr,
		cl:          cl,
	}
}

func NewWritePipe(creator PipeCreator, description string, wc io.WriteCloser) *Pipe {
	return &Pipe{
		creator:     creator,
		description: description,
		wr:          wc,
		cl:          wc,
	}
}

func (p *Pipe) Close() {
	if p.cl != nil {
		p.cl.Close()
	}
}

func (p *Pipe) IsNull() bool {
	return p.isNull
}

func (p *Pipe) CanRead() bool {
	return p.rd != nil || p.isNull
}

func (p *Pipe) CanWrite() bool {
	return p.wr != nil || p.isNull
}

func (p *Pipe) GetReader() io.Reader {
	return p.rd
}

func (p *Pipe) GetWriter() io.Writer {
	return p.wr
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
	if p.rd != nil {
		mode += "r"
	}
	if p.wr != nil {
		mode += "w"
	}

	creatorStr := "nil"
	if p.creator != nil {
		creatorStr = p.creator.ToString()
	}

	return fmt.Sprintf("%s{%s, %s, %s}", LuaType, mode, p.description, creatorStr)
}
