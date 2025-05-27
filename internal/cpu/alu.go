package cpu

func (cpu *SM83) incRegister(register *byte) {
	*register++
	cpu.SetFlag(ZeroFlag, *register == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, (*register&0x0F) == 0)
}

func (cpu *SM83) decRegister(register *byte) {
	*register--
	cpu.SetFlag(ZeroFlag, *register == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, (*register&0x0F) == 0x0F)
}
