package cpu

import (
	"fmt"
	"log/slog"
)

func (c *SM83) executeCB(instruction byte) (cycles int) {
	switch instruction {
	case 0x11: // RL C
		// Rotate C left. Current carry flag goes to LSB, rotated MSB goes to carry flag.
		carry := c.GetFlag(CarryFlag)
		topBit := c.r_C & 0x80 // Get the top bit (7th bit) of C
		c.r_C = (c.r_C << 1)
		if carry {
			c.r_C |= 0x01 // Set LSB if carry was set
		} else {
			c.r_C &= 0xFE // Clear LSB if carry was not set
		}
		if topBit != 0 {
			c.SetFlag(CarryFlag, true) // Set carry flag if the top bit was 1
		} else {
			c.SetFlag(CarryFlag, false) // Clear carry flag if the top bit was 0
		}

		c.SetFlag(ZeroFlag, c.r_C == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)

		slog.Debug("Executing RL C", "C", fmt.Sprintf("0x%02X", c.r_C))
		cycles += 2
	case 0x37: // SWAP A
		// Swap the nibbles of A
		c.r_A = (c.r_A << 4) | (c.r_A >> 4)
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 2
	case 0x7C: // BIT 7,H
		// Check if the 7th bit of H is set
		c.SetFlag(ZeroFlag, (c.r_H&0x80) == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, true)
		slog.Debug("Executing BIT 7,H", "H", fmt.Sprintf("0x%02X", c.r_H))
		cycles += 2
	case 0x87: // RES 0,A
		// Reset the 0th bit of A
		c.r_A &= 0xFE // Clear the 0th bit
		cycles += 2
	default:
		panic("Unknown CB-prefixed instruction: " + fmt.Sprintf("0x%02X", instruction))
	}
	return
}

func (c *SM83) execute(instruction byte) (cycles int) {
	switch instruction {
	case 0x00: // NOP
		cycles += 1
	case 0x01: // LD BC,nn
		// Load the next two bytes as a 16-bit value into BC
		bc, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read BC value: %v", err))
		}
		slog.Debug("Executing LD BC,nn", "value", fmt.Sprintf("0x%04X", bc))
		c.r_B = byte(bc >> 8)
		c.r_C = byte(bc & 0xFF)
		c.r_PC += 2
		cycles += 3
	case 0x04: // INC B
		c.incRegister(&c.r_B)
		cycles += 1
		slog.Debug("Executing INC B", "value", fmt.Sprintf("0x%02X", c.r_B))
	case 0x05: // DEC B
		// Decrement B and set flags accordingly
		c.decRegister(&c.r_B)
		cycles += 1
		slog.Debug("Executing DEC B", "value", fmt.Sprintf("0x%02X", c.r_B))
	case 0x06: // LD B,n
		// Load the next byte as an immediate value into B
		breg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read B value: %v", err))
		}
		slog.Debug("Executing LD B,n", "value", fmt.Sprintf("0x%02X", breg))
		c.r_B = breg
		c.r_PC++
		cycles += 2
	case 0x0B: // DEC BC
		// Decrement BC
		c.r_C--
		if c.r_C == 0xFF {
			c.r_B-- // Decrement B if C wraps around
		}
		slog.Debug("Executing DEC BC", "BC", fmt.Sprintf("0x%04X", uint16(c.r_B)<<8|uint16(c.r_C)))
		cycles += 2
	case 0x0C: // INC C
		c.incRegister(&c.r_C)
		cycles += 1
		slog.Debug("Executing INC C", "value", fmt.Sprintf("0x%02X", c.r_C))
	case 0x0d: // DEC C
		c.decRegister(&c.r_C)
		cycles += 1
		slog.Debug("Executing DEC C", "value", fmt.Sprintf("0x%02X", c.r_C))
	case 0x0e: // LD C,n
		// Load the next byte as an immediate value into C
		creg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read C value: %v", err))
		}
		slog.Debug("Executing LD C,n", "value", fmt.Sprintf("0x%02X", creg))
		c.r_C = creg
		c.r_PC++
		cycles += 2
	case 0x11: // LD DE,nn
		// Load the next two bytes as a 16-bit value into DE
		de, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read DE value: %v", err))
		}
		slog.Debug("Executing LD DE,nn", "value", fmt.Sprintf("0x%04X", de))
		c.r_D = byte(de >> 8)
		c.r_E = byte(de & 0xFF)
		c.r_PC += 2
		cycles += 3
	case 0x12: // LD (DE),A
		// Store the value of A at the address pointed to by DE
		addr := uint16(c.r_D)<<8 | uint16(c.r_E)
		slog.Debug("Executing LD (DE),A", "address", fmt.Sprintf("0x%04X", addr))
		err := c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (DE): %v", err))
		}
		slog.Debug("Stored A in (DE)", "value", fmt.Sprintf("0x%02X", c.r_A))
		cycles += 2
	case 0x13: // INC DE
		// Increment DE
		c.r_E++
		if c.r_E == 0x00 {
			c.r_D++ // Increment D if E wraps around
		}
		slog.Debug("Executing INC DE", "DE", fmt.Sprintf("0x%04X", uint16(c.r_D)<<8|uint16(c.r_E)))
		cycles += 2
	case 0x15: // DEC D
		c.decRegister(&c.r_D)
		cycles += 1
		slog.Debug("Executing DEC D", "value", fmt.Sprintf("0x%02X", c.r_D))
	case 0x16: // LD D,n
		// Load the next byte as an immediate value into D
		dreg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read D value: %v", err))
		}
		slog.Debug("Executing LD D,n", "value", fmt.Sprintf("0x%02X", dreg))
		c.r_D = dreg
		c.r_PC++
		cycles += 2
	case 0x17: // RLA
		// Rotate A left
		previousCarry := c.GetFlag(CarryFlag)
		topBit := c.r_A & 0x80 // Get the top bit (7th bit) of A
		c.r_A = (c.r_A << 1)
		if topBit != 0 {
			c.SetFlag(CarryFlag, true) // Set carry flag if the top bit was 1
		} else {
			c.SetFlag(CarryFlag, false) // Clear carry flag if the top bit was 0
		}
		if previousCarry {
			c.r_A |= 0x01 // Set LSB if previous carry was set
		} else {
			c.r_A &= 0xFE // Clear LSB if previous carry was not set
		}
		c.SetFlag(ZeroFlag, false)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		cycles += 1
	case 0x18: // JR n
		// Read the next byte as a signed offset
		offset, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read offset for JR: %v", err))
		}
		c.r_PC++
		slog.Debug("Executing JR n", "offset", fmt.Sprintf("0x%02X", offset))
		// Jump to the new location
		c.r_PC += uint16(int8(offset)) // Sign-extend the offset
		slog.Debug("Jumping to new PC", "new_PC", fmt.Sprintf("0x%04X", c.r_PC))
		cycles += 3
	case 0x19: // ADD HL,DE
		// Add DE to HL
		hl := uint16(c.r_H)<<8 | uint16(c.r_L)
		de := uint16(c.r_D)<<8 | uint16(c.r_E)
		slog.Debug("Executing ADD HL,DE", "HL", fmt.Sprintf("0x%04X", hl), "DE", fmt.Sprintf("0x%04X", de))
		result := hl + de
		c.r_H = byte(result >> 8)
		c.r_L = byte(result & 0xFF)
		c.SetFlag(CarryFlag, result > 0xFFFF) // Set carry flag if result exceeds 16 bits
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, (hl&0x0FFF)+(de&0x0FFF) > 0x0FFF) // Check for half carry
		cycles += 2
	case 0x1A: // LD A,(DE)
		// Load the value at address DE into A
		addr := uint16(c.r_D)<<8 | uint16(c.r_E)
		slog.Debug("Executing LD A,(DE)", "address", fmt.Sprintf("0x%04X", addr))
		areg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read A from (DE): %v", err))
		}
		slog.Debug("Loaded A from (DE)", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		cycles += 2
	case 0x1C: // INC E
		c.incRegister(&c.r_E)
		cycles += 1
	case 0x1D: // DEC E
		c.decRegister(&c.r_E)
		cycles += 1
		slog.Debug("Executing DEC E", "value", fmt.Sprintf("0x%02X", c.r_E))
	case 0x1E: // LD E,n
		// Load the next byte as an immediate value into E
		ereg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read E value: %v", err))
		}
		slog.Debug("Executing LD E,n", "value", fmt.Sprintf("0x%02X", ereg))
		c.r_E = ereg
		c.r_PC++
		cycles += 2
	case 0x20: // JR NZ,nn
		// Read the next byte as a signed offset
		offset, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read offset for JR NZ: %v", err))
		}
		c.r_PC++
		slog.Debug("Executing JR NZ,nn", "offset", fmt.Sprintf("0x%02X", offset))
		if !c.GetFlag(ZeroFlag) {
			// If Z flag is not set, jump to the new location
			c.r_PC += uint16(int8(offset)) // Sign-extend the offset
			slog.Debug("Jumping to new PC", "new_PC", fmt.Sprintf("0x%04X", c.r_PC))
			cycles += 3
		} else {
			slog.Debug("Not jumping")
			cycles += 2
		}
	case 0x21: // LD HL,nn
		// Load the next two bytes as a 16-bit value into HL
		hl, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read HL value: %v", err))
		}
		slog.Debug("Executing LD HL,nn", "value", fmt.Sprintf("0x%04X", hl))
		c.r_H = byte(hl >> 8)
		c.r_L = byte(hl & 0xFF)
		c.r_PC += 2
		cycles += 3
	case 0x22: // LD (HL+),A
		// Store the value of A at the address pointed to by HL, then increment HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD (HL+),A", "address", fmt.Sprintf("0x%04X", addr))
		err := c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (HL+): %v", err))
		}
		c.r_L++
		if c.r_L == 0x00 {
			c.r_H++ // Increment H if L wraps around
		}
		cycles += 2
	case 0x23: // INC HL
		// Increment HL
		c.r_L++
		if c.r_L == 0x00 {
			c.r_H++ // Increment H if L wraps around
		}
		slog.Debug("Executing INC HL", "HL", fmt.Sprintf("0x%04X", uint16(c.r_H)<<8|uint16(c.r_L)))
		cycles += 1
	case 0x24: // INC H
		c.incRegister(&c.r_H)
		cycles += 1
		slog.Debug("Executing INC H", "value", fmt.Sprintf("0x%02X", c.r_H))
	case 0x27: // DAA
		correction := byte(0)

		if c.GetFlag(HalfCarryFlag) || (!c.GetFlag(NegativeFlag) && (c.r_A&0xf) > 9) {
			correction |= 0x6
		}

		if c.GetFlag(CarryFlag) || (!c.GetFlag(NegativeFlag) && c.r_A > 0x99) {
			correction |= 0x60
			c.SetFlag(CarryFlag, true)
		}

		if c.GetFlag(NegativeFlag) {
			c.r_A -= correction
		} else {
			c.r_A += correction
		}

		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(HalfCarryFlag, false)
		cycles += 1
	case 0x28: // JR Z,nn
		// Read the next byte as a signed offset
		offset, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read offset for JR Z: %v", err))
		}
		c.r_PC++
		slog.Debug("Executing JR Z,nn", "offset", fmt.Sprintf("0x%02X", offset))
		if c.GetFlag(ZeroFlag) {
			// If Z flag is set, jump to the new location
			c.r_PC += uint16(int8(offset)) // Sign-extend the offset
			slog.Debug("Jumping to new PC", "new_PC", fmt.Sprintf("0x%04X", c.r_PC))
			cycles += 3
		} else {
			slog.Debug("Not jumping")
			cycles += 2
		}
	case 0x2A: // LD A,(HL+)
		// Load the value at address HL into A, then increment HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD A,(HL+)", "address", fmt.Sprintf("0x%04X", addr))
		areg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read A from (HL+): %v", err))
		}
		slog.Debug("Loaded A from (HL+)", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		c.r_L++
		if c.r_L == 0x00 {
			c.r_H++ // Increment H if L wraps around
		}
		cycles += 2
	case 0x2C: // INC L
		c.incRegister(&c.r_L)
		cycles += 1
	case 0x2E: // LD L,n
		// Load the next byte as an immediate value into L
		lreg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read L value: %v", err))
		}
		slog.Debug("Executing LD L,n", "value", fmt.Sprintf("0x%02X", lreg))
		c.r_L = lreg
		c.r_PC++
		cycles += 2
	case 0x2F: // CPL
		// Complement A (invert all bits)
		slog.Debug("Executing CPL")
		c.r_A = ^c.r_A

		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, true)
		cycles += 1
	case 0x31: // LD SP,nn
		// Load the next two bytes as a 16-bit value into SP
		sp, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read SP value: %v", err))
		}
		slog.Debug("Executing LD SP,nn", "value", fmt.Sprintf("0x%04X", sp))
		c.r_SP = sp
		c.r_PC += 2
		cycles += 3
	case 0x32: // LD (HL-),A
		// Store the value of A at the address pointed to by HL, then decrement HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD (HL-),A", "address", fmt.Sprintf("0x%04X", addr))
		err := c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (HL-): %v", err))
		}
		c.r_L--
		if c.r_L == 0xFF {
			c.r_H-- // Decrement H if L wraps around
		}
		cycles += 2
	case 0x35: // DEC (HL)
		// Decrement the value at address HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing DEC (HL)", "address", fmt.Sprintf("0x%04X", addr))
		value, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read value from (HL): %v", err))
		}
		slog.Debug("Value before DEC", "value", fmt.Sprintf("0x%02X", value))
		value--
		if value == 0xFF {
			c.SetFlag(ZeroFlag, true)
		} else {
			c.SetFlag(ZeroFlag, false)
		}
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (value&0x0F) == 0x0F)
		slog.Debug("Value after DEC", "value", fmt.Sprintf("0x%02X", value))
		err = c.memory.Write8(addr, value)
		if err != nil {
			panic(fmt.Sprintf("Failed to write decremented value to (HL): %v", err))
		}
		cycles += 3
	case 0x36: // LD (HL),n
		// Load the next byte as an immediate value into the address pointed to by HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		n, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read immediate value for LD (HL),n: %v", err))
		}
		slog.Debug("Executing LD (HL),n", "address", fmt.Sprintf("0x%04X", addr), "value", fmt.Sprintf("0x%02X", n))
		err = c.memory.Write8(addr, n)
		if err != nil {
			panic(fmt.Sprintf("Failed to write value to (HL): %v", err))
		}
		c.r_PC++
		cycles += 3
	case 0x37: // SCF
		// Set the carry flag
		slog.Debug("Executing SCF")
		c.SetFlag(CarryFlag, true)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		cycles += 1
	case 0x3D: // DEC A
		c.decRegister(&c.r_A)
		cycles += 1
		slog.Debug("Executing DEC A", "value", fmt.Sprintf("0x%02X", c.r_A))
	case 0x3E: // LD A,n
		// Load the next byte as an immediate value into A
		areg, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read A value: %v", err))
		}
		slog.Debug("Executing LD A,n", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		c.r_PC++
		cycles += 2
	case 0x47: // LD B,A
		// Load the value of A into B
		slog.Debug("Executing LD B,A", "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_B = c.r_A
		cycles += 1
	case 0x4F: // LD C,A
		// Load the value of A into C
		slog.Debug("Executing LD C,A", "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_C = c.r_A
		cycles += 1
	case 0x56: // LD D,(HL)
		// Load the value at address HL into D
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD D,(HL)", "address", fmt.Sprintf("0x%04X", addr))
		dreg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read D from (HL): %v", err))
		}
		slog.Debug("Loaded D from (HL)", "value", fmt.Sprintf("0x%02X", dreg))
		c.r_D = dreg
		cycles += 2
	case 0x57: // LD D,A
		// Load the value of A into D
		slog.Debug("Executing LD D,A", "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_D = c.r_A
		cycles += 1
	case 0x5E: // LD E,(HL)
		// Load the value at address HL into E
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD E,(HL)", "address", fmt.Sprintf("0x%04X", addr))
		ereg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read E from (HL): %v", err))
		}
		slog.Debug("Loaded E from (HL)", "value", fmt.Sprintf("0x%02X", ereg))
		c.r_E = ereg
		cycles += 2
	case 0x5F: // LD E,A
		// Load the value of A into E
		slog.Debug("Executing LD E,A", "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_E = c.r_A
		cycles += 1
	case 0x67: // LD H,A
		// Load the value of A into H
		slog.Debug("Executing LD H,A", "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_H = c.r_A
		cycles += 1
	case 0x77: // LD (HL),A
		// Store the value of A at the address pointed to by HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD (HL),A", "address", fmt.Sprintf("0x%04X", addr), "value", fmt.Sprintf("0x%02X", c.r_A))
		err := c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (HL): %v", err))
		}
		cycles += 2
	case 0x78: // LD A,B
		// Load the value of B into A
		slog.Debug("Executing LD A,B", "value", fmt.Sprintf("0x%02X", c.r_B))
		c.r_A = c.r_B
		cycles += 1
	case 0x79: // LD A,C
		// Load the value of C into A
		slog.Debug("Executing LD A,C", "value", fmt.Sprintf("0x%02X", c.r_C))
		c.r_A = c.r_C
		cycles += 1
	case 0x7B: // LD A,E
		// Load the value of E into A
		slog.Debug("Executing LD A,E", "value", fmt.Sprintf("0x%02X", c.r_E))
		c.r_A = c.r_E
		cycles += 1
	case 0x7C: // LD A,H
		// Load the value of H into A
		slog.Debug("Executing LD A,H", "value", fmt.Sprintf("0x%02X", c.r_H))
		c.r_A = c.r_H
		cycles += 1
	case 0x7D: // LD A,L
		// Load the value of L into A
		slog.Debug("Executing LD A,L", "value", fmt.Sprintf("0x%02X", c.r_L))
		c.r_A = c.r_L
		cycles += 1
	case 0x7E: // LD A,(HL)
		// Load the value at address HL into A
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing LD A,(HL)", "address", fmt.Sprintf("0x%04X", addr))
		areg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read A from (HL): %v", err))
		}
		slog.Debug("Loaded A from (HL)", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		cycles += 2
	case 0x86: // Add A,(HL)
		// Add the value at address HL to A and set flags
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing ADD A,(HL)", "address", fmt.Sprintf("0x%04X", addr))
		hlValue, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read value from (HL): %v", err))
		}
		slog.Debug("Value at (HL)", "value", fmt.Sprintf("0x%02X", hlValue))
		c.r_A += hlValue
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (hlValue&0x0F))
		c.SetFlag(CarryFlag, c.r_A < hlValue)
		cycles += 2
	case 0x87: // ADD A,A
		// Add A to itself and set flags
		slog.Debug("Executing ADD A,A", "A", fmt.Sprintf("0x%02X", c.r_A))
		c.r_A += c.r_A
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (c.r_A&0x0F))
		c.SetFlag(CarryFlag, c.r_A < c.r_A)
		cycles += 1
	case 0x90: // SUB B
		// Subtract B from A and set flags
		slog.Debug("Executing SUB B", "B", fmt.Sprintf("0x%02X", c.r_B))
		c.SetFlag(ZeroFlag, c.r_A == c.r_B)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (c.r_B&0x0F))
		c.SetFlag(CarryFlag, c.r_A < c.r_B)
		c.r_A -= c.r_B
		cycles += 1
	case 0x95: // SUB L
		// Subtract L from A and set flags
		slog.Debug("Executing SUB L", "L", fmt.Sprintf("0x%02X", c.r_L))
		c.SetFlag(ZeroFlag, c.r_A == c.r_L)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (c.r_L&0x0F))
		c.SetFlag(CarryFlag, c.r_A < c.r_L)
		c.r_A -= c.r_L
		cycles += 1
	case 0x96: // SUB (HL)
		// Subtract the value at address HL from A and set flags
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing SUB (HL)", "address", fmt.Sprintf("0x%04X", addr))
		hlValue, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read value from (HL): %v", err))
		}
		slog.Debug("Value at (HL)", "value", fmt.Sprintf("0x%02X", hlValue))
		c.SetFlag(ZeroFlag, c.r_A == hlValue)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (hlValue&0x0F))
		c.SetFlag(CarryFlag, c.r_A < hlValue)
		c.r_A -= hlValue
		cycles += 2
	case 0xA1: // AND C
		// Perform bitwise AND between A and C
		slog.Debug("Executing AND C", "C", fmt.Sprintf("0x%02X", c.r_C))
		c.r_A &= c.r_C
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, true)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xA7: // AND A
		// Perform bitwise AND between A and itself (effectively clearing the flags)
		slog.Debug("Executing AND A", "A", fmt.Sprintf("0x%02X", c.r_A))
		c.r_A &= c.r_A
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, true)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xA9: // XOR C
		// Perform bitwise XOR between A and C
		slog.Debug("Executing XOR C", "C", fmt.Sprintf("0x%02X", c.r_C))
		c.r_A ^= c.r_C
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xAF: // XOR A
		c.r_A = 0
		c.SetFlag(ZeroFlag, true)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xB0: // OR B
		// Perform bitwise OR between A and B
		slog.Debug("Executing OR B", "B", fmt.Sprintf("0x%02X", c.r_B))
		c.r_A |= c.r_B
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xB1: // OR C
		// Perform bitwise OR between A and C
		slog.Debug("Executing OR C", "C", fmt.Sprintf("0x%02X", c.r_C))
		c.r_A |= c.r_C
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xBE: // CP (HL)
		// Compare the value at address HL with A and set flags
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing CP (HL)", "address", fmt.Sprintf("0x%04X", addr))
		hlValue, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read value from (HL): %v", err))
		}
		slog.Debug("Value at (HL)", "value", fmt.Sprintf("0x%02X", hlValue))
		c.SetFlag(ZeroFlag, c.r_A == hlValue)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (hlValue&0x0F))
		c.SetFlag(CarryFlag, c.r_A < hlValue)
		cycles += 2
	case 0xC1: // POP BC
		c.pop16(&c.r_B, &c.r_C)
		slog.Debug("Executing POP BC")
		cycles += 3
	case 0xC3: // JP nn
		newloc, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for JP: %v", err))
		}
		slog.Debug("Executing JP nn", "address", fmt.Sprintf("0x%04X", newloc))
		c.r_PC = newloc
		cycles += 4
	case 0xC5: // PUSH BC
		// Push the current BC register onto the stack
		c.push16(&c.r_B, &c.r_C)
		cycles += 4
	case 0xC8: // RET Z
		// Return from subroutine if Z flag is set
		slog.Debug("Executing RET Z")
		if c.GetFlag(ZeroFlag) {
			addr, err := c.memory.Read16(c.r_SP)
			if err != nil {
				panic(fmt.Sprintf("Failed to read return address from stack: %v", err))
			}
			slog.Debug("Popped return address from stack", "address", fmt.Sprintf("0x%04X", addr))
			c.r_SP += 2
			c.r_PC = addr
			slog.Debug("Jumping to new PC", "new_PC", fmt.Sprintf("0x%04X", c.r_PC))
			cycles += 5
		} else {
			slog.Debug("Not returning, Z flag not set")
			cycles += 2
		}
	case 0xC9: // RET
		// Return from subroutine
		slog.Debug("Executing RET")
		addr, err := c.memory.Read16(c.r_SP)
		if err != nil {
			panic(fmt.Sprintf("Failed to read return address from stack: %v", err))
		}
		slog.Debug("Popped return address from stack", "address", fmt.Sprintf("0x%04X", addr))
		c.r_SP += 2
		c.r_PC = addr
		cycles += 4
	case 0xCA: // JP Z,nn
		// Read the next two bytes as a 16-bit address
		newloc, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for JP Z: %v", err))
		}
		c.r_PC += 2
		slog.Debug("Executing JP Z,nn", "address", fmt.Sprintf("0x%04X", newloc))
		if c.GetFlag(ZeroFlag) {
			// If Z flag is set, jump to the new location
			c.r_PC = newloc
			slog.Debug("Jumping to new PC", "new_PC", fmt.Sprintf("0x%04X", c.r_PC))
			cycles += 4
		} else {
			slog.Debug("Not jumping")
			cycles += 3
		}
	case 0xCB: // CB-prefixed instructions
		// Handle CB-prefixed instructions
		instruction, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read CB-prefixed instruction: %v", err))
		}
		c.r_PC++
		cycles += 1
		cycles += c.executeCB(instruction)
	case 0xCD: // CALL nn
		newloc, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for CALL: %v", err))
		}
		c.r_PC += 2
		slog.Debug("Executing CALL nn", "address", fmt.Sprintf("0x%04X", newloc))
		// Push the current PC onto the stack
		err = c.memory.Write16(c.r_SP-2, c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to write PC to stack for CALL: %v", err))
		}
		c.r_SP -= 2
		// Set the PC to the new location
		c.r_PC = newloc
		cycles += 6
	case 0xD1: // POP DE
		c.pop16(&c.r_D, &c.r_E)
		slog.Debug("Executing POP DE")
		cycles += 3
	case 0xD5: // PUSH DE
		c.push16(&c.r_D, &c.r_E)
		slog.Debug("Executing PUSH DE", "value", fmt.Sprintf("0x%04X", uint16(c.r_D)<<8|uint16(c.r_E)))
		cycles += 4
	case 0xE0: // LDH (n),A
		// Load the value of A into the address 0xFF00 + n
		n, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read n for LDH: %v", err))
		}
		slog.Debug("Executing LDH (n),A", "address", fmt.Sprintf("0x%02X", 0xFF00+uint16(n)), "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_PC++
		err = c.memory.Write8(0xFF00+uint16(n), c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (n): %v", err))
		}
		cycles += 3
	case 0xE1: // POP HL
		c.pop16(&c.r_H, &c.r_L)
		slog.Debug("Executing POP HL")
		cycles += 3
	case 0xE2: // LD (C),A
		// Load the value of A into the address pointed to by C
		addr := uint16(0xFF00) | uint16(c.r_C)
		slog.Debug("Executing LD (C),A", "address", fmt.Sprintf("0x%04X", addr), "value", fmt.Sprintf("0x%02X", c.r_A))
		err := c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (C): %v", err))
		}
		cycles += 2
	case 0xEA: // LD (nn),A
		// Load the value of A into the address nn
		addr, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for LD (nn),A: %v", err))
		}
		slog.Debug("Executing LD (nn),A", "address", fmt.Sprintf("0x%04X", addr), "value", fmt.Sprintf("0x%02X", c.r_A))
		c.r_PC += 2
		err = c.memory.Write8(addr, c.r_A)
		if err != nil {
			panic(fmt.Sprintf("Failed to write A to (nn): %v", err))
		}
		cycles += 4
	case 0xE5: // PUSH HL
		c.push16(&c.r_H, &c.r_L)
		slog.Debug("Executing PUSH HL", "value", fmt.Sprintf("0x%04X", uint16(c.r_H)<<8|uint16(c.r_L)))
		cycles += 4
	case 0xE6: // AND n
		// Perform bitwise AND between A and the next byte
		n, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read n for AND: %v", err))
		}
		slog.Debug("Executing AND n", "value", fmt.Sprintf("0x%02X", n))
		c.r_PC++
		c.r_A &= n
		c.SetFlag(ZeroFlag, c.r_A == 0)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, true)
		c.SetFlag(CarryFlag, false)
		cycles += 2
	case 0xE9: // JP (HL)
		// Jump to the address in HL
		addr := uint16(c.r_H)<<8 | uint16(c.r_L)
		slog.Debug("Executing JP (HL)", "address", fmt.Sprintf("0x%04X", addr))
		c.r_PC = addr
		cycles += 1
	case 0xEF: // RST 0x28
		// Push the current PC onto the stack and set PC to 0x28
		slog.Debug("Executing RST 0x28")
		err := c.memory.Write16(c.r_SP-2, c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to write PC to stack for RST 0x28: %v", err))
		}
		c.r_SP -= 2
		c.r_PC = 0x28
		cycles += 4
	case 0xF0: // LDH A,(n)
		// Load the value from the address 0xFF00 + n into A
		n, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read n for LDH: %v", err))
		}
		slog.Debug("Executing LDH A,(n)", "address", fmt.Sprintf("0x%02X", 0xFF00+uint16(n)))
		c.r_PC++
		areg, err := c.memory.Read8(0xFF00 + uint16(n))
		if err != nil {
			panic(fmt.Sprintf("Failed to read A from (n): %v", err))
		}
		slog.Debug("Loaded A from (n)", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		cycles += 3
	case 0xF1: // POP AF
		// Pop the top of the stack into AF
		c.pop16(&c.r_A, &c.r_F)
		slog.Debug("Executing POP AF", "AF", fmt.Sprintf("0x%04X", uint16(c.r_A)<<8|uint16(c.r_F)))
		cycles += 3
	case 0xF3: // DI
		// Disable interrupts
		slog.Debug("Executing DI")
		c.ime = false
		cycles += 1
	case 0xF5: // PUSH AF
		// Push the current AF register onto the stack
		c.push16(&c.r_A, &c.r_F)
		slog.Debug("Executing PUSH AF", "value", fmt.Sprintf("0x%04X", uint16(c.r_A)<<8|uint16(c.r_F)))
		cycles += 4
	case 0xFA: // LD A,(nn)
		// Load the value at address nn into A
		addr, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for LD A,(nn): %v", err))
		}
		slog.Debug("Executing LD A,(nn)", "address", fmt.Sprintf("0x%04X", addr))
		c.r_PC += 2
		areg, err := c.memory.Read8(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to read A from (nn): %v", err))
		}
		slog.Debug("Loaded A from (nn)", "value", fmt.Sprintf("0x%02X", areg))
		c.r_A = areg
		cycles += 4
	case 0xFB: // EI
		// Enable interrupts
		slog.Debug("Executing EI")
		c.ime = true
		cycles += 1
	case 0xFE: // CP n
		// Subtract the next byte from A and set flags
		n, err := c.memory.Read8(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read n for CP: %v", err))
		}
		slog.Debug("Executing CP n", "value", fmt.Sprintf("0x%02X", n))
		c.r_PC++
		c.SetFlag(ZeroFlag, c.r_A == n)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_A&0x0F) < (n&0x0F))
		c.SetFlag(CarryFlag, c.r_A < n)
		cycles += 2
	default:
		panic("Unknown instruction: " + fmt.Sprintf("0x%02X", instruction))
	}

	return
}
