package pipe

type FixedPipeCreator struct {
	Name string
}

func (s *FixedPipeCreator) ToString() string {
	return s.Name
}

var nullPipe = Pipe{
	isNull: true,
}
