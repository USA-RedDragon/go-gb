package cpu

func scf(cpu *SM83) {
	cpu.SetFlag(CarryFlag, true)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(NegativeFlag, false)
}
