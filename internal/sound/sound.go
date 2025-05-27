package sound

type Sound struct {
	NR11 byte // Channel 1 length timer and duty cycle register
	NR12 byte // Channel 1 volume and envelope register
	NR50 byte // Master volume and VIN panning register
	NR51 byte // Sound panning register
	NR52 byte // Audio master control register
}

func NewSound() *Sound {
	snd := &Sound{}
	snd.Reset()
	return snd
}

func (s *Sound) Reset() {
	s.NR11 = 0x00
	s.NR12 = 0x00
	s.NR50 = 0x00
	s.NR51 = 0x00
	s.NR52 = 0x00
}
