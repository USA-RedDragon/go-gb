package cpu

func jr(cpu *SM83) {
	offset, err := cpu.memory.Read8(cpu.r_PC)
	if err != nil {
		panic(err)
	}
	cpu.r_PC++
	if offset&0x80 != 0 {
		offset = -((^offset + 1) & 0xFF) // Convert to signed
	}
	cpu.r_PC += uint16(offset)
}

func jrCond(cpu *SM83, condition bool) {
	if condition {
		offset, err := cpu.memory.Read8(cpu.r_PC)
		if err != nil {
			panic(err)
		}
		cpu.r_PC++
		if offset&0x80 != 0 {
			offset = -((^offset + 1) & 0xFF) // Convert to signed
		}
		cpu.r_PC += uint16(offset)
	} else {
		cpu.r_PC++ // Just skip the offset
	}
}

func ret(cpu *SM83) {
	// Read the return address from the stack
	addr, err := cpu.memory.Read16(cpu.r_SP)
	if err != nil {
		panic(err)
	}
	cpu.r_SP += 2   // Increment stack pointer
	cpu.r_PC = addr // Set program counter to return address
}

func retCond(cpu *SM83, condition bool) {
	if condition {
		// Read the return address from the stack
		addr, err := cpu.memory.Read16(cpu.r_SP)
		if err != nil {
			panic(err)
		}
		cpu.r_SP += 2   // Increment stack pointer
		cpu.r_PC = addr // Set program counter to return address
	} else {
		cpu.r_SP += 2 // Just skip the return address
	}
}

func call(cpu *SM83) {
	// Read the call address
	addr, err := cpu.memory.Read16(cpu.r_PC)
	if err != nil {
		panic(err)
	}
	cpu.r_PC += 2 // Increment program counter

	// Push the current program counter onto the stack
	err = cpu.memory.Write16(cpu.r_SP-2, cpu.r_PC)
	if err != nil {
		panic(err)
	}
	cpu.r_SP -= 2 // Decrement stack pointer

	cpu.r_PC = addr // Set program counter to call address
}

func callCond(cpu *SM83, condition bool) {
	if condition {
		// Read the call address
		addr, err := cpu.memory.Read16(cpu.r_PC)
		if err != nil {
			panic(err)
		}
		cpu.r_PC += 2 // Increment program counter

		// Push the current program counter onto the stack
		err = cpu.memory.Write16(cpu.r_SP-2, cpu.r_PC)
		if err != nil {
			panic(err)
		}
		cpu.r_SP -= 2 // Decrement stack pointer

		cpu.r_PC = addr // Set program counter to call address
	} else {
		cpu.r_PC += 2 // Just skip the call address
	}
}

func rst(cpu *SM83, vector byte) {
	// Push the current program counter onto the stack
	err := cpu.memory.Write16(cpu.r_SP-2, cpu.r_PC)
	if err != nil {
		panic(err)
	}
	cpu.r_SP -= 2 // Decrement stack pointer

	// Set program counter to the reset vector
	cpu.r_PC = uint16(vector)
}

func reti(cpu *SM83) {
	// Read the return address from the stack
	addr, err := cpu.memory.Read16(cpu.r_SP)
	if err != nil {
		panic(err)
	}
	cpu.r_SP += 2   // Increment stack pointer
	cpu.r_PC = addr // Set program counter to return address

	// Enable interrupts
	cpu.ime = true
}

func jp(cpu *SM83) {
	// Read the jump address
	addr, err := cpu.memory.Read16(cpu.r_PC)
	if err != nil {
		panic(err)
	}
	cpu.r_PC += 2 // Increment program counter

	cpu.r_PC = addr // Set program counter to jump address
}

func jpCond(cpu *SM83, condition bool) {
	if condition {
		// Read the jump address
		addr, err := cpu.memory.Read16(cpu.r_PC)
		if err != nil {
			panic(err)
		}
		cpu.r_PC += 2 // Increment program counter

		cpu.r_PC = addr // Set program counter to jump address
	} else {
		cpu.r_PC += 2 // Just skip the jump address
	}
}

func jpMemComb(cpu *SM83, addrTop *byte, addrBottom *byte) {
	// Combine the two registers into a single 16-bit address
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)

	cpu.r_PC = addr // Set program counter to the combined address
}
