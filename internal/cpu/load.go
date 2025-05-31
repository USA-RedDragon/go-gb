package cpu

func ldRegisterRegister(c *SM83, dst, src *byte) {
	*dst = *src
}

func ldRegisterImm(c *SM83, dst *byte) {
	var err error
	*dst, err = c.memory.Read8(c.r_PC)
	if err != nil {
		panic(err)
	}
	c.r_PC++
}

func ldCombRegister16Imm(c *SM83, dstTop *byte, dstBottom *byte) {
	combReg, err := c.memory.Read16(c.r_PC)
	if err != nil {
		panic(err)
	}
	c.r_PC += 2
	*dstTop = byte(combReg >> 8)
	*dstBottom = byte(combReg & 0xFF)
}

func ldRegisterMemComb(c *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	var err error
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
}

func ldMemCombRegister(c *SM83, addrTop *byte, addrBottom *byte, src *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	err := c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
}

func ldMem16Register(c *SM83, src *uint16) {
	addr, err := c.memory.Read16(c.r_PC)
	if err != nil {
		panic(err)
	}
	c.r_PC += 2
	*src, err = c.memory.Read16(addr)
	if err != nil {
		panic(err)
	}
}

func ldRegisterMemCombInc(c *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	var err error
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
	*addrBottom++
	if *addrBottom == 0x00 {
		*addrTop++
	}
}

func ldRegisterMemCombDec(c *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	var err error
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
	if *addrBottom == 0x00 {
		*addrBottom = 0xFF
		*addrTop--
	} else {
		*addrBottom--
	}
}

func ldMemCombRegisterInc(c *SM83, addrTop *byte, addrBottom *byte, src *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	err := c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
	*addrBottom++
	if *addrBottom == 0x00 {
		*addrTop++
	}
}

func ldMemCombRegisterDec(c *SM83, addrTop *byte, addrBottom *byte, src *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	err := c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
	if *addrBottom == 0x00 {
		*addrBottom = 0xFF
		*addrTop--
	} else {
		*addrBottom--
	}
}

func ld16Register16Imm(c *SM83, dst *uint16) {
	var err error
	*dst, err = c.memory.Read16(c.r_PC)
	if err != nil {
		panic(err)
	}
	c.r_PC += 2
}

func ld16RegCombRegister(c *SM83, dst *uint16, srcTop *byte, srcBottom *byte) {
	*dst = uint16(*srcTop)<<8 | uint16(*srcBottom)
}

func ldMemCombImm(c *SM83, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := c.memory.Read8(c.r_PC)
	if err != nil {
		panic(err)
	}
	c.r_PC++
	err = c.memory.Write8(addr, value)
	if err != nil {
		panic(err)
	}
}

func popRegisterPair(c *SM83, dstTop *byte, dstBottom *byte) {
	// Read the value from the stack
	addr, err := c.memory.Read16(c.r_SP)
	if err != nil {
		panic(err)
	}
	c.r_SP += 2 // Increment stack pointer

	// Set the destination registers
	*dstTop = byte(addr >> 8)
	*dstBottom = byte(addr & 0xFF)
}

func pushRegisterPair(c *SM83, srcTop *byte, srcBottom *byte) {
	// Combine the two registers into a single 16-bit value
	addr := uint16(*srcTop)<<8 | uint16(*srcBottom)

	// Write the value to the stack
	err := c.memory.Write16(c.r_SP-2, addr)
	if err != nil {
		panic(err)
	}
	c.r_SP -= 2 // Decrement stack pointer
}
