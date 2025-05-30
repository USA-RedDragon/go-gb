package sound

type Sound struct {
	NR10 byte // Channel 1 sweep register
	NR11 byte // Channel 1 length timer and duty cycle register
	NR12 byte // Channel 1 volume and envelope register
	NR13 byte // Channel 1 period low register, write-only
	NR14 byte // Channel 1 period high and control register
	NR22 byte // Channel 2 volume and envelope register
	NR24 byte // Channel 2 period high and control register
	NR30 byte // Channel 3 DAC enable register
	NR42 byte // Channel 4 volume and envelope register
	NR44 byte // Channel 4 control register
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
	s.NR10 = 0x00
	s.NR11 = 0x00
	s.NR12 = 0x00
	s.NR13 = 0x00
	s.NR14 = 0x00
	s.NR22 = 0x00
	s.NR24 = 0x00
	s.NR30 = 0x00
	s.NR42 = 0x00
	s.NR44 = 0x00
	s.NR50 = 0x00
	s.NR51 = 0x00
	s.NR52 = 0x00
}
