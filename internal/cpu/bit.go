package cpu

func bitRegister(cpu *SM83, bit uint8, reg *byte) {
	mask := byte(1 << bit)
	cpu.SetFlag(HalfCarryFlag, true)
	cpu.SetFlag(NegativeFlag, false)
	if *reg&mask == 0 {
		cpu.SetFlag(ZeroFlag, true)
	} else {
		cpu.SetFlag(ZeroFlag, false)
	}
}

func swapRegister(cpu *SM83, reg *byte) {
	*reg = (*reg >> 4) | (*reg << 4)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(ZeroFlag, *reg == 0)
}

func resRegister(cpu *SM83, bit uint8, reg *byte) {
	mask := byte(1 << bit)
	*reg &^= mask // Clear the specified bit
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(ZeroFlag, *reg == 0)
}

func rlRegister(cpu *SM83, reg *byte) {
	carry := (*reg & 0x80) >> 7 // Get the carry bit
	*reg = (*reg << 1) | carry  // Shift left and set the carry bit
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(CarryFlag, carry != 0)
	cpu.SetFlag(ZeroFlag, *reg == 0)
}
