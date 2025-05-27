package cpu

import (
	"fmt"
)

func (c *SM83) execute(instruction byte) {
	switch instruction {
	case 0x00: // NOP
	default:
		panic("Unknown instruction: " + fmt.Sprintf("0x%02X", instruction))
	}
}
