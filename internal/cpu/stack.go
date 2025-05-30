package cpu

import "fmt"

func (cpu *SM83) push16(registerA *byte, registerB *byte) {
	cpu.push8(registerA)
	cpu.push8(registerB)
}

func (cpu *SM83) push8(register *byte) {
	err := cpu.memory.Write8(cpu.r_SP-1, *register)
	if err != nil {
		panic(fmt.Sprintf("Failed to write BC to stack: %v", err))
	}
	cpu.r_SP -= 1
}

func (cpu *SM83) pop16(registerA *byte, registerB *byte) {
	var err error
	*registerB, err = cpu.pop8()
	if err != nil {
		panic(fmt.Sprintf("Failed to pop B from stack: %v", err))
	}
	*registerA, err = cpu.pop8()
	if err != nil {
		panic(fmt.Sprintf("Failed to pop A from stack: %v", err))
	}
}

func (cpu *SM83) pop8() (byte, error) {
	data, err := cpu.memory.Read8(cpu.r_SP)
	if err != nil {
		return 0, fmt.Errorf("failed to read from stack: %w", err)
	}
	cpu.r_SP += 1
	return data, nil
}
