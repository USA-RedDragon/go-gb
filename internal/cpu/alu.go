package cpu

func incRegister(cpu *SM83, register *byte) {
	*register++
	cpu.SetFlag(ZeroFlag, *register == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, (*register&0x0F) == 0)
}

func decRegister(cpu *SM83, register *byte) {
	*register--
	cpu.SetFlag(ZeroFlag, *register == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, (*register&0x0F) == 0x0F)
}

func incCombRegister(_ *SM83, top *byte, bottom *byte) {
	if *bottom == 0xFF {
		*bottom = 0x00
		*top++
	} else {
		*bottom++
	}
}

func decCombRegister(_ *SM83, top *byte, bottom *byte) {
	if *bottom == 0x00 {
		*bottom = 0xFF
		*top--
	} else {
		*bottom--
	}
}

func inc16Register(cpu *SM83, reg *uint16) {
	*reg++
	cpu.SetFlag(ZeroFlag, false)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, (*reg&0x0FFF) == 0)
}

func addRegister(cpu *SM83, dst *byte, src *byte) {
	result := uint16(*dst) + uint16(*src)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(*src&0x0F)) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func rla(cpu *SM83) {
	var carry byte
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	} else {
		carry = 0
	}
	cpu.SetFlag(CarryFlag, (cpu.rA&0x80) != 0)
	cpu.rA = (cpu.rA << 1) | carry
	cpu.SetFlag(ZeroFlag, cpu.rA == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
}

func addCombRegisterCombRegister(cpu *SM83, dstTop *byte, dstBottom *byte, srcTop *byte, srcBottom *byte) {
	dstAddr := uint16(*dstTop)<<8 | uint16(*dstBottom)
	srcAddr := uint16(*srcTop)<<8 | uint16(*srcBottom)
	result := uint32(dstAddr) + uint32(srcAddr)
	cpu.SetFlag(ZeroFlag, result&0xFFFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((dstAddr&0x0FFF)+(srcAddr&0x0FFF)) > 0x0FFF)
	cpu.SetFlag(CarryFlag, result > 0xFFFF)
	*dstTop = byte((result >> 8) & 0xFF)
	*dstBottom = byte(result & 0xFF)
}

func daa(cpu *SM83) {
	correction := byte(0)

	if cpu.GetFlag(HalfCarryFlag) || (!cpu.GetFlag(NegativeFlag) && (cpu.rA&0xf) > 9) {
		correction |= 0x6
	}

	if cpu.GetFlag(CarryFlag) || (!cpu.GetFlag(NegativeFlag) && cpu.rA > 0x99) {
		correction |= 0x60
		cpu.SetFlag(CarryFlag, true)
	}

	if cpu.GetFlag(NegativeFlag) {
		cpu.rA -= correction
	} else {
		cpu.rA += correction
	}

	cpu.SetFlag(ZeroFlag, cpu.rA == 0)
	cpu.SetFlag(HalfCarryFlag, false)
}

func cpl(cpu *SM83) {
	cpu.rA = ^cpu.rA

	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, true)
}

func incMemComb(cpu *SM83, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	value++
	cpu.SetFlag(ZeroFlag, value == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, (value&0x0F) == 0)
	err = cpu.memory.Write8(addr, value)
	if err != nil {
		panic(err)
	}
}

func decMemComb(cpu *SM83, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	value--
	cpu.SetFlag(ZeroFlag, value == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, (value&0x0F) == 0x0F)
	err = cpu.memory.Write8(addr, value)
	if err != nil {
		panic(err)
	}
}

func addCombRegister16Register(cpu *SM83, dstTop *byte, dstBottom *byte, src *uint16) {
	dstAddr := uint16(*dstTop)<<8 | uint16(*dstBottom)
	result := uint32(dstAddr) + uint32(*src)
	cpu.SetFlag(ZeroFlag, result&0xFFFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((dstAddr&0x0FFF)+(*src&0x0FFF)) > 0x0FFF)
	cpu.SetFlag(CarryFlag, result > 0xFFFF)
	*dstTop = byte((result >> 8) & 0xFF)
	*dstBottom = byte(result & 0xFF)
}

func dec16Register(cpu *SM83, reg *uint16) {
	*reg--
	cpu.SetFlag(ZeroFlag, *reg == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, (*reg&0x0FFF) == 0x0FFF)
}

func addRegisterMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	result := uint16(*dst) + uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(value&0x0F)) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func subRegister(cpu *SM83, dst *byte, src *byte) {
	result := uint16(*dst) - uint16(*src)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-(*src&0x0F)) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func subMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	result := uint16(*dst) - uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func andRegister(cpu *SM83, dst *byte, src *byte) {
	*dst &= *src
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, true)
	cpu.SetFlag(CarryFlag, false)
}

func andMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	*dst &= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, true)
	cpu.SetFlag(CarryFlag, false)
}

func xorRegister(cpu *SM83, dst *byte, src *byte) {
	*dst ^= *src
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func xorMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	*dst ^= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func orRegister(cpu *SM83, dst *byte, src *byte) {
	*dst |= *src
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func orMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	*dst |= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func cpRegister(cpu *SM83, dst *byte, src *byte) {
	result := uint16(*dst) - uint16(*src)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-(*src&0x0F)) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
}

func cpMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	result := uint16(*dst) - uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
}

func cpImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	result := uint16(*dst) - uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
}

func addImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	result := uint16(*dst) + uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(value&0x0F)) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func subImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	result := uint16(*dst) - uint16(value)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func andImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	*dst &= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, true)
	cpu.SetFlag(CarryFlag, false)
}

func addSPImmediate(cpu *SM83) {
	offset, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	if offset&0x80 != 0 {
		offset = -((^offset + 1) & 0xFF) // Convert to signed
	}

	result := cpu.rSP + uint16(offset)
	cpu.SetFlag(ZeroFlag, false)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, (cpu.rSP&0x0FFF+uint16(offset)&0x0FFF) > 0x0FFF)
	cpu.SetFlag(CarryFlag, result < 0 || result > 0xFFFF)

	cpu.rSP = result & 0xFFFF
}

func xorImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	*dst ^= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func orImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	*dst |= value
	cpu.SetFlag(ZeroFlag, *dst == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, false)
	cpu.SetFlag(CarryFlag, false)
}

func sbcRegister(cpu *SM83, dst *byte, src *byte) {
	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) - uint16(*src) - uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-(*src&0x0F)-carry) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func adcRegister(cpu *SM83, dst *byte, src *byte) {
	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) + uint16(*src) + uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(*src&0x0F)+carry) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func sbcMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) - uint16(value) - uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F-carry) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func adcMemComb(cpu *SM83, dst *byte, addrTop *byte, addrBottom *byte) {
	addr := uint16(*addrTop)<<8 | uint16(*addrBottom)
	value, err := cpu.memory.Read8(addr)
	if err != nil {
		panic(err)
	}

	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) + uint16(value) + uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(value&0x0F)+carry) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func adcImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) + uint16(value) + uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, false)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)+(value&0x0F)+carry) > 0x0F)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}

func sbcImmediate(cpu *SM83, dst *byte) {
	value, err := cpu.memory.Read8(cpu.rPC)
	if err != nil {
		panic(err)
	}
	cpu.rPC++ // Increment program counter

	carry := byte(0)
	if cpu.GetFlag(CarryFlag) {
		carry = 1
	}
	result := uint16(*dst) - uint16(value) - uint16(carry)
	cpu.SetFlag(ZeroFlag, result&0xFF == 0)
	cpu.SetFlag(NegativeFlag, true)
	cpu.SetFlag(HalfCarryFlag, ((*dst&0x0F)-value&0x0F-carry) < 0)
	cpu.SetFlag(CarryFlag, result > 0xFF)
	*dst = byte(result & 0xFF)
}
