package cpu

func ldRegisterRegister(_ *SM83, dst, src *byte) {
	*dst = *src
}

func ldRegisterImm(c *SM83, dst *byte) {
	var err error
	*dst, err = c.memory.Read8(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC++
}

func ldCombRegister16Imm(c *SM83, dstTop *byte, dstBottom *byte) {
	combReg, err := c.memory.Read16(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC += 2
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
	addr, err := c.memory.Read16(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC += 2
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
	*dst, err = c.memory.Read16(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC += 2
}

func ld16RegCombRegister(_ *SM83, dst *uint16, srcTop *byte, srcBottom *byte) {
	*dst = uint16(*srcTop)<<8 | uint16(*srcBottom)
}

func ldMemCombImm(c *SM83, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := c.memory.Read8(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC++
	err = c.memory.Write8(addr, value)
	if err != nil {
		panic(err)
	}
}

func popRegisterPair(c *SM83, dstTop *byte, dstBottom *byte) {
	// Read the value from the stack
	addr, err := c.memory.Read16(c.rSP)
	if err != nil {
		panic(err)
	}
	c.rSP += 2 // Increment stack pointer

	// Set the destination registers
	*dstTop = byte(addr >> 8)
	*dstBottom = byte(addr & 0xFF)
}

func pushRegisterPair(c *SM83, srcTop *byte, srcBottom *byte) {
	// Write the value to the stack
	err := c.memory.Write8(c.rSP-1, *srcTop)
	if err != nil {
		panic(err)
	}
	c.rSP--
	err = c.memory.Write8(c.rSP-1, *srcBottom)
	if err != nil {
		panic(err)
	}
	c.rSP--
}

func ldMemRegisterRegister(c *SM83, addrRegister *byte, src *byte) {
	addr := 0xFF00 + uint16(*addrRegister)
	err := c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
}

func ldRegisterMemRegister(c *SM83, dst *byte, addrRegister *byte) {
	addr := 0xFF00 + uint16(*addrRegister)
	var err error
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
}

func ld16ImmMemRegister(c *SM83, src *byte) {
	addr, err := c.memory.Read16(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC += 2
	err = c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
}

func ldh8ImmMemRegister(c *SM83, src *byte) {
	// Read the immediate value from the memory
	value, err := c.memory.Read8(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC++

	// The address is always 0xFF00 + value
	addr := 0xFF00 + uint16(value)

	err = c.memory.Write8(addr, *src)
	if err != nil {
		panic(err)
	}
}

func ldhRegisterMemImm(c *SM83, dst *byte) {
	// Read the immediate value from the memory
	value, err := c.memory.Read8(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC++

	// The address is always 0xFF00 + value
	addr := 0xFF00 + uint16(value)

	// Read the value from the memory at the specified address
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
}

func ldRegisterMem16Imm(c *SM83, dst *byte) {
	// Read the immediate value from the memory
	addr, err := c.memory.Read16(c.rPC)
	if err != nil {
		panic(err)
	}
	c.rPC += 2

	// Read the value from the memory at the specified address
	*dst, err = c.memory.Read8(addr)
	if err != nil {
		panic(err)
	}
}
