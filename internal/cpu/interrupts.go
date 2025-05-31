package cpu

import "github.com/USA-RedDragon/go-gb/internal/impls"

func (c *SM83) GetInterruptEnableFlag(flag impls.Interrupt) bool {
	return c.interruptEnable&byte(flag) != 0
}

func (c *SM83) GetInterruptFlag(flag impls.Interrupt) bool {
	return c.interruptFlag&byte(flag) != 0
}

func (c *SM83) SetInterruptFlag(flag impls.Interrupt, val bool) {
	if c.ime && val {
		c.interruptFlag |= byte(flag)
	} else {
		c.interruptFlag &^= byte(flag)
	}
}
