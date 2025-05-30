package cpu

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/consts"
	"github.com/USA-RedDragon/go-gb/internal/memory"
	"github.com/USA-RedDragon/go-gb/internal/ppu"
	"github.com/USA-RedDragon/go-gb/internal/sound"
)

type SM83 struct {
	config    *config.Config
	memory    memory.MMIO // Memory-mapped I/O
	PPU       *ppu.PPU
	Sound     *sound.Sound
	cartridge *cartridge.Cartridge

	halted bool
	exit   bool

	RAM  [consts.RAMSize]byte  // 8KB of RAM
	HRAM [consts.HRAMSize]byte // 127 bytes of HRAM

	ime             bool // Interrupt Master Enable flag
	interruptFlag   byte // Interrupt Flag register, used to check which interrupts are pending
	interruptEnable byte // Interrupt Enable register, used to enable/disable interrupts
	serialData      byte // SB, serial data register
	serialControl   byte // SC, serial control register
	bank            byte // 0xFF50, used to disable BIOS

	// Registers
	r_IR byte // IR, instruction register
	r_IE byte // IE, interrupt enable register

	r_A byte // A, accumulator register
	r_F byte // F, flags register
	r_B byte // B
	r_C byte // C
	r_D byte // D
	r_E byte // E
	r_H byte // H
	r_L byte // L

	// 16-bit registers
	r_PC uint16 // PC (Program Counter)
	r_SP uint16 // SP (Stack Pointer)
}

func NewSM83(config *config.Config, cartridge *cartridge.Cartridge) *SM83 {
	cpu := &SM83{
		config:    config,
		memory:    memory.MMIO{},
		cartridge: cartridge,
		PPU:       ppu.NewPPU(),
		Sound:     sound.NewSound(),
	}

	cpu.Reset()

	cpu.PPU.Step()
	cpu.PPU.Step()
	cpu.PPU.Step()
	cpu.PPU.Step()

	slog.Debug("ROM loaded", "rom", cartridge)

	return cpu
}

func (c *SM83) Reset() {
	c.RAM = [consts.RAMSize]byte{}
	c.PPU.Reset()
	c.cartridge.Reset()
	c.HRAM = [consts.HRAMSize]byte{}
	c.interruptFlag = 0
	c.interruptEnable = 0
	c.serialData = 0
	c.serialControl = 0
	c.bank = 0x00

	c.memory = memory.MMIO{}

	if c.config.BIOS != "" {
		biosData, err := os.ReadFile(c.config.BIOS)
		if err != nil {
			slog.Error("Failed to load BIOS", "error", err)
			os.Exit(1)
		}
		if len(biosData) != consts.BIOSSize {
			slog.Error("Invalid BIOS size", "expected", consts.BIOSSize, "got", len(biosData))
			os.Exit(1)
		}
		// We need to add the BIOS to the memory map at 0x0000, and only add cartridge ROM banks after that
		c.memory.AddMMIO(biosData, 0x0000, consts.BIOSSize)
		c.memory.AddMMIO(c.cartridge.ROMBank0[consts.BIOSSize:], 0x0100, consts.ROMBankSize-consts.BIOSSize)
	} else {
		c.memory.AddMMIO(c.cartridge.ROMBank0[:], 0x0, consts.ROMBankSize)
	}
	if len(c.cartridge.AdditionalROMBanks) > 0 {
		c.memory.AddMMIO(c.cartridge.AdditionalROMBanks[0][:], 0x4000, consts.ROMBankSize)
	}
	c.memory.AddMMIO(c.PPU.VRAM[:], 0x8000, consts.VRAMSize)
	if c.cartridge.RAMSize.Bytes() > 0 {
		c.memory.AddMMIO(c.cartridge.CartridgeRAMBanks[0][:], 0xA000, consts.CartridgeRAMBankSize)
	}
	c.memory.AddMMIO(c.RAM[:], 0xC000, consts.RAMSize)
	c.memory.AddMMIO(c.PPU.OAM[:], 0xFE00, consts.OAMSize)
	c.memory.AddMMIO(make([]byte, consts.ProhibitedSize)[:], 0xFEA0, consts.ProhibitedSize)
	c.memory.AddMMIOByte(&c.serialData, 0xFF01)
	c.memory.AddMMIOByte(&c.serialControl, 0xFF02)
	c.memory.AddMMIOByte(&c.interruptFlag, 0xFF0F)
	c.memory.AddMMIOByte(&c.Sound.NR11, 0xFF11)
	c.memory.AddMMIOByte(&c.Sound.NR12, 0xFF12)
	c.memory.AddMMIOByte(&c.Sound.NR13, 0xFF13)
	c.memory.AddMMIOByte(&c.Sound.NR14, 0xFF14)
	c.memory.AddMMIOByte(&c.Sound.NR50, 0xFF24)
	c.memory.AddMMIOByte(&c.Sound.NR51, 0xFF25)
	c.memory.AddMMIOByte(&c.Sound.NR52, 0xFF26)
	c.memory.AddMMIOByte(&c.PPU.LCDControl, 0xFF40)
	c.memory.AddMMIOByte(&c.PPU.LCDStatus, 0xFF41)
	c.memory.AddMMIOByte(&c.PPU.SCY, 0xFF42)
	c.memory.AddMMIOByte(&c.PPU.SCX, 0xFF43)
	c.memory.AddMMIOByte(&c.PPU.LY, 0xFF44)
	c.memory.AddMMIOByte(&c.PPU.LYC, 0xFF45)
	c.memory.AddMMIOByte(&c.PPU.BGP, 0xFF47)
	c.memory.AddMMIOByte(&c.PPU.OBP0, 0xFF48)
	c.memory.AddMMIOByte(&c.PPU.OBP1, 0xFF49)
	c.memory.AddMMIOByte(&c.bank, 0xFF50)
	byt := byte(0x00)
	c.memory.AddMMIOByte(&byt, 0xFF7F) // Unused
	c.memory.AddMMIO(c.HRAM[:], 0xFF80, consts.HRAMSize)
	c.memory.AddMMIOByte(&c.interruptEnable, 0xFFFF)

	if c.config.BIOS != "" {
		c.ime = false
		c.r_IR = 0
		c.r_IE = 0
		c.r_A = 0
		c.r_F = 0
		c.r_B = 0
		c.r_C = 0
		c.r_D = 0
		c.r_E = 0
		c.r_H = 0
		c.r_L = 0
		c.r_PC = 0x0
	} else {
		c.ime = false
		c.r_IR = 0
		c.r_IE = 0
		c.r_A = 0x01
		c.r_F = byte(ZeroFlag)
		if c.cartridge.ROMBank0[0x014D] == 0x00 {
			// Header checksum, if zero the Carry and Half Carry flags are set
			c.r_F |= byte(CarryFlag) | byte(HalfCarryFlag)
		}
		c.r_B = 0
		c.r_C = 0x13
		c.r_D = 0
		c.r_E = 0xD8
		c.r_H = 0x01
		c.r_L = 0x4D
		c.r_PC = 0x0100 // Program Counter starts at 0x0100
	}
	c.r_SP = 0xFFFE
	c.halted = false
	c.exit = false
}

func (c *SM83) GetPC() uint16 {
	return c.r_PC
}

func (c *SM83) Step() int {
	if !c.halted {
		instruction := c.fetch()
		c.r_IR = instruction

		slog.Debug("Instruction", "instruction", fmt.Sprintf("0x%02X", instruction))
		slog.Debug(c.DebugRegisters())

		for range 4 {
			c.PPU.Step()
		}

		preBank := c.bank
		steps := c.execute(instruction)
		if c.bank != preBank && c.bank != 0 && c.config.BIOS != "" {
			// Boot rom disabled
			err := c.memory.RemoveMMIO(0x0000, consts.BIOSSize)
			if err != nil {
				panic(fmt.Sprintf("Failed to remove BIOS MMIO: %v", err))
			}
			c.memory.AddMMIO(c.cartridge.ROMBank0[:consts.BIOSSize], 0x0000, consts.BIOSSize)
		}

		return steps
	}
	return 1
}

func (c *SM83) fetch() byte {
	// Fetch the next instruction from memory at the current PC
	instruction, err := c.memory.Read8(c.r_PC)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch instruction at PC 0x%04X: %v", c.r_PC, err))
	}
	c.r_PC++
	return instruction
}

func (c *SM83) DebugRegisters() string {
	var ret = "\n"
	ret += fmt.Sprintf("IR: 0x%02X\t IE: 0x%02X\n", c.r_IR, c.r_IE)
	ret += fmt.Sprintf(" A: 0x%02X\t  F: 0x%02X\n", c.r_A, c.r_F)
	ret += fmt.Sprintf(" B: 0x%02X\t  C: 0x%02X\n", c.r_B, c.r_C)
	ret += fmt.Sprintf(" D: 0x%02X\t  E: 0x%02X\n", c.r_D, c.r_E)
	ret += fmt.Sprintf(" H: 0x%02X\t  L: 0x%02X\n", c.r_H, c.r_L)
	ret += fmt.Sprintf("PC: 0x%04X\t SP: 0x%04X\n", c.r_PC, c.r_SP)
	ret += fmt.Sprintf("Flags: C: %t, H: %t, N: %t, Z: %t\n",
		c.GetFlag(CarryFlag),
		c.GetFlag(HalfCarryFlag),
		c.GetFlag(NegativeFlag),
		c.GetFlag(ZeroFlag),
	)

	return ret
}

func (c *SM83) RunUntilFrame() []byte {
	cycleTime := time.Second / 4194304 / 4 // 4.194304 MHz, divided by 4 (1.048576 MHz) to count machine cycles
	for !c.PPU.HaveFrame {
		prevTime := time.Now()
		cycles := c.Step()
		for range cycles {
			c.PPU.Step()
			c.PPU.Step()
			c.PPU.Step()
			c.PPU.Step()
			time.Sleep(cycleTime - time.Since(prevTime))
			prevTime = time.Now()
		}
	}

	return c.PPU.GetFrame()
}

func (c *SM83) Run() {
	cycleTime := time.Second / 4194304 / 4 // 4.194304 MHz, divided by 4 (1.048576 MHz) to count machine cycles
	time.Sleep(cycleTime)                  // Simulate the initial delay from reading the first instruction
	for !c.exit {
		prevTime := time.Now()
		cycles := c.Step()
		for range cycles {
			c.PPU.Step()
			c.PPU.Step()
			c.PPU.Step()
			c.PPU.Step()
			time.Sleep(cycleTime - time.Since(prevTime))
			prevTime = time.Now()
		}
	}
}

type Flag uint8

const (
	CarryFlag     Flag = 1 << 4
	HalfCarryFlag Flag = 1 << 5
	NegativeFlag  Flag = 1 << 6
	ZeroFlag      Flag = 1 << 7
)

func (c *SM83) SetFlag(flag Flag, val bool) {
	if val {
		c.r_F |= byte(flag)
	} else {
		c.r_F &^= byte(flag)
	}
}

func (c *SM83) GetFlag(flag Flag) bool {
	return c.r_F&byte(flag) != 0
}

func (c *SM83) Halt() {
	c.halted = true
}

func (c *SM83) Unhalt() {
	c.halted = false
}

func (c *SM83) Quit() {
	c.exit = true
}
