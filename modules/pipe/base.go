package pipe

import (
	"errors"
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

var _ io.ReadWriteCloser = (*Pipe)(nil)

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
	if p.isNull {
		return p
	}
	return p.rd
}

func (p *Pipe) GetWriter() io.Writer {
	if p.isNull {
		return p
	}
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
	if p.cl != nil {
		mode += "c"
	}

	creatorStr := "nil"
	if p.creator != nil {
		creatorStr = p.creator.ToString()
	}

	return fmt.Sprintf("%s{%s, %s, %s}", LuaType, mode, p.description, creatorStr)
}

func (p *Pipe) Close() error {
	if p.cl != nil {
		return p.cl.Close()
	}
	return errors.New("cannot close pipe")
}

func (p *Pipe) Read(data []byte) (n int, err error) {
	if p.isNull {
		return 0, io.EOF
	}
	if p.rd != nil {
		return p.rd.Read(data)
	}
	return 0, errors.New("cannot read from pipe")
}

func (p *Pipe) Write(data []byte) (n int, err error) {
	if p.isNull {
		return len(data), nil
	}
	if p.wr != nil {
		return p.wr.Write(data)
	}
	return 0, errors.New("cannot write to pipe")
}
