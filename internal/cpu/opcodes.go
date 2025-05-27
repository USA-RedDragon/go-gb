package cpu

import (
	"fmt"
	"log/slog"
)

func (c *SM83) execute(instruction byte) (cycles int) {
	switch instruction {
	case 0x00: // NOP
		cycles += 1
	case 0x05: // DEC B
		// Decrement B and set flags accordingly
		c.r_B--
		slog.Debug("Executing DEC B", "value", fmt.Sprintf("0x%02X", c.r_B))
		c.SetFlag(ZeroFlag, c.r_B == 0)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_B&0x0F) == 0x0F)
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
	case 0x0d: // DEC C
		// Decrement C and set flags accordingly
		c.r_C--
		slog.Debug("Executing DEC C", "value", fmt.Sprintf("0x%02X", c.r_C))
		c.SetFlag(ZeroFlag, c.r_C == 0)
		c.SetFlag(NegativeFlag, true)
		c.SetFlag(HalfCarryFlag, (c.r_C&0x0F) == 0x0F)
		cycles += 1
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
	case 0xAF: // XOR A
		c.r_A = 0
		c.SetFlag(ZeroFlag, true)
		c.SetFlag(NegativeFlag, false)
		c.SetFlag(HalfCarryFlag, false)
		c.SetFlag(CarryFlag, false)
		cycles += 1
	case 0xC3: // JP nn
		newloc, err := c.memory.Read16(c.r_PC)
		if err != nil {
			panic(fmt.Sprintf("Failed to read address for JP: %v", err))
		}
		slog.Debug("Executing JP nn", "address", fmt.Sprintf("0x%04X", newloc))
		c.r_PC = newloc
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
	case 0xF3: // DI
		// Disable interrupts
		slog.Debug("Executing DI")
		c.ime = false
		cycles += 1
	default:
		panic("Unknown instruction: " + fmt.Sprintf("0x%02X", instruction))
	}

	return
}
