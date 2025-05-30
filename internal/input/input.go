package input

type Input struct {
	JOYP byte // Joypad register
}

func NewInput() *Input {
	input := &Input{}
	input.Reset()
	return input
}

func (s *Input) Reset() {
	s.JOYP = 0xCF
}
