package cpu

func jr(cpu *SM83) {
	offset, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++
	if offset&0x80 != 0 {
		offset = -((^offset + 1) & 0xFF) // Convert to signed
	}
	cpu.rPC += uint16(offset)
}

func jrCond(cpu *SM83, condition bool) {
	if condition {
		offset, err := cpu.memory.Read8(cpu.rPC)
		if err != nil {
			panic(err)
		}
		cpu.rPC++
		if offset&0x80 != 0 {
			offset = -((^offset + 1) & 0xFF) // Convert to signed
		}
		cpu.rPC += uint16(offset)
	} else {
		cpu.rPC++ // Just skip the offset
	}
}

func ret(cpu *SM83) {
	// Read the return address from the stack
	addr, err := cpu.memory.Read16(cpu.rSP)
	if err != nil {
		panic(err)
	}
	cpu.rSP += 2   // Increment stack pointer
	cpu.rPC = addr // Set program counter to return address
}

func retCond(cpu *SM83, condition bool) {
	if condition {
		// Read the return address from the stack
		addr, err := cpu.memory.Read16(cpu.rSP)
		if err != nil {
			panic(err)
		}
		cpu.rSP += 2   // Increment stack pointer
		cpu.rPC = addr // Set program counter to return address
	} else {
		cpu.rSP += 2 // Just skip the return address
	}
}

func call(cpu *SM83) {
	// Read the call address
	addr, err := cpu.memory.Read16(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC += 2 // Increment program counter

	// Push the current program counter onto the stack
	err = cpu.memory.Write16(cpu.rSP-2, cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rSP -= 2 // Decrement stack pointer

	cpu.rPC = addr // Set program counter to call address
}

func callCond(cpu *SM83, condition bool) {
	if condition {
		// Read the call address
		addr, err := cpu.memory.Read16(cpu.rPC)
		if err != nil {
			panic(err)
		}
		cpu.rPC += 2 // Increment program counter

		// Push the current program counter onto the stack
		err = cpu.memory.Write16(cpu.rSP-2, cpu.rPC)
		if err != nil {
			panic(err)
		}
		cpu.rSP -= 2 // Decrement stack pointer

		cpu.rPC = addr // Set program counter to call address
	} else {
		cpu.rPC += 2 // Just skip the call address
	}
}

func rst(cpu *SM83, vector byte) {
	// Push the current program counter onto the stack
	err := cpu.memory.Write16(cpu.rSP-2, cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rSP -= 2 // Decrement stack pointer

	// Set program counter to the reset vector
	cpu.rPC = uint16(vector)
}

func reti(cpu *SM83) {
	// Read the return address from the stack
	addr, err := cpu.memory.Read16(cpu.rSP)
	if err != nil {
		panic(err)
	}
	cpu.rSP += 2   // Increment stack pointer
	cpu.rPC = addr // Set program counter to return address

	// Enable interrupts
	cpu.ime = true
}

func jp(cpu *SM83) {
	// Read the jump address
	addr, err := cpu.memory.Read16(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC += 2 // Increment program counter

	cpu.rPC = addr // Set program counter to jump address
}

func jpCond(cpu *SM83, condition bool) {
	if condition {
		// Read the jump address
		addr, err := cpu.memory.Read16(cpu.rPC)
		if err != nil {
			panic(err)
		}
		cpu.rPC += 2 // Increment program counter

		cpu.rPC = addr // Set program counter to jump address
	} else {
		cpu.rPC += 2 // Just skip the jump address
	}
}

func jpMemComb(cpu *SM83, addrTop *byte, addrBottom *byte) {
	// Combine the two registers into a single 16-bit address
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)

	cpu.rPC = addr // Set program counter to the combined address
}
