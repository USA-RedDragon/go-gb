package impls

type Interrupt uint8

const (
	JoypadInterrupt Interrupt = 1 << 4
	SerialInterrupt Interrupt = 1 << 3
	TimerInterrupt  Interrupt = 1 << 2
	LCDInterrupt    Interrupt = 1 << 1
	VBlankInterrupt Interrupt = 1 << 0
)

type CPU interface {
	SetInterruptFlag(flag Interrupt, val bool)
}
