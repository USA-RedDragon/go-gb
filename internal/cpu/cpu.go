package cpu

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/consts"
	"github.com/USA-RedDragon/go-gb/internal/impls"
	"github.com/USA-RedDragon/go-gb/internal/input"
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
	Input     *input.Input

	halted bool
	exit   bool

	RAM  [consts.RAMSize]byte  // 8KB of RAM
	HRAM [consts.HRAMSize]byte // 127 bytes of HRAM

	ime             bool  // Interrupt Master Enable flag
	interruptFlag   byte  // Interrupt Flag register, used to check which interrupts are pending
	interruptEnable byte  // Interrupt Enable register, used to enable/disable interrupts
	serialData      byte  // SB, serial data register
	serialControl   byte  // SC, serial control register
	bank            byte  // 0xFF50, used to disable BIOS
	TMA             uint8 // Timer Modulo register, used for the timer
	TAC             byte  // Timer Control register, used to control the timer operation
	OAMDMA          byte  // Object Attribute Memory DMA register, used for sprite data transfer

	rA byte // A, accumulator register
	rF byte // F, flags register
	rB byte // B
	rC byte // C
	rD byte // D
	rE byte // E
	rH byte // H
	rL byte // L

	// 16-bit registers
	rPC uint16 // PC (Program Counter)
	rSP uint16 // SP (Stack Pointer)
}

func NewSM83(config *config.Config, cartridge *cartridge.Cartridge) *SM83 {
	cpu := &SM83{
		config:    config,
		memory:    memory.MMIO{},
		cartridge: cartridge,
		Sound:     sound.NewSound(),
		Input:     input.NewInput(),
	}
	cpu.PPU = ppu.NewPPU(cpu)

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
	if c.cartridge != nil {
		c.cartridge.Reset()
	}
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
		c.memory.AddMMIO(biosData, 0x0000, consts.BIOSSize, true)
		if c.cartridge != nil {
			c.memory.AddMMIO(c.cartridge.ROMBank0[consts.BIOSSize:], 0x0100, consts.ROMBankSize-consts.BIOSSize, true)
		} else {
			c.memory.AddMMIO(bytes.Repeat([]byte{0xff}, consts.ROMBankSize-consts.BIOSSize), 0x0100, consts.ROMBankSize-consts.BIOSSize, true)
		}
	} else {
		if c.cartridge != nil {
			c.memory.AddMMIO(c.cartridge.ROMBank0[:], 0x0, consts.ROMBankSize, true)
		} else {
			c.memory.AddMMIO(bytes.Repeat([]byte{0xff}, consts.ROMBankSize), 0x0, consts.ROMBankSize, true)
		}
	}
	if c.cartridge != nil && len(c.cartridge.AdditionalROMBanks) > 0 {
		c.memory.AddMMIO(c.cartridge.AdditionalROMBanks[0][:], 0x4000, consts.ROMBankSize, true)
	}
	c.memory.AddMMIO(c.PPU.VRAM[:], 0x8000, consts.VRAMSize, false)
	if c.cartridge != nil && c.cartridge.RAMSize.Bytes() > 0 {
		c.memory.AddMMIO(c.cartridge.CartridgeRAMBanks[0], 0xA000, consts.CartridgeRAMBankSize, false)
	} else {
		c.memory.AddMMIO(bytes.Repeat([]byte{0xff}, consts.CartridgeRAMBankSize), 0xA000, consts.CartridgeRAMBankSize, false)
	}
	c.memory.AddMMIO(c.RAM[:], 0xC000, consts.RAMSize, false)
	c.memory.AddMMIO(c.PPU.OAM[:], 0xFE00, consts.OAMSize, false)
	c.memory.AddMMIO(make([]byte, consts.ProhibitedSize), 0xFEA0, consts.ProhibitedSize, false)
	c.memory.AddMMIOByte(&c.Input.JOYP, 0xFF00, true)
	c.memory.AddMMIOByte(&c.serialData, 0xFF01, false)
	c.memory.AddMMIOByte(&c.serialControl, 0xFF02, false)
	c.memory.AddMMIOByte(&c.TMA, 0xFF06, false)
	c.memory.AddMMIOByte(&c.TAC, 0xFF07, false)
	c.memory.AddMMIOByte(&c.interruptFlag, 0xFF0F, false)
	c.memory.AddMMIOByte(&c.Sound.NR10, 0xFF10, false)
	c.memory.AddMMIOByte(&c.Sound.NR11, 0xFF11, false)
	c.memory.AddMMIOByte(&c.Sound.NR12, 0xFF12, false)
	c.memory.AddMMIOByte(&c.Sound.NR13, 0xFF13, false)
	c.memory.AddMMIOByte(&c.Sound.NR14, 0xFF14, false)
	c.memory.AddMMIOByte(&c.Sound.NR22, 0xFF17, false)
	c.memory.AddMMIOByte(&c.Sound.NR24, 0xFF19, false)
	c.memory.AddMMIOByte(&c.Sound.NR30, 0xFF1A, false)
	c.memory.AddMMIOByte(&c.Sound.NR42, 0xFF21, false)
	c.memory.AddMMIOByte(&c.Sound.NR44, 0xFF23, false)
	c.memory.AddMMIOByte(&c.Sound.NR50, 0xFF24, false)
	c.memory.AddMMIOByte(&c.Sound.NR51, 0xFF25, false)
	c.memory.AddMMIOByte(&c.Sound.NR52, 0xFF26, false)
	c.memory.AddMMIOByte(&c.PPU.LCDControl, 0xFF40, false)
	c.memory.AddMMIOByte(&c.PPU.LCDStatus, 0xFF41, false)
	c.memory.AddMMIOByte(&c.PPU.SCY, 0xFF42, false)
	c.memory.AddMMIOByte(&c.PPU.SCX, 0xFF43, false)
	c.memory.AddMMIOByte(&c.PPU.LY, 0xFF44, true)
	c.memory.AddMMIOByte(&c.PPU.LYC, 0xFF45, false)
	c.memory.AddMMIOByte(&c.OAMDMA, 0xFF46, false)
	c.memory.AddMMIOByte(&c.PPU.BGP, 0xFF47, false)
	c.memory.AddMMIOByte(&c.PPU.OBP0, 0xFF48, false)
	c.memory.AddMMIOByte(&c.PPU.OBP1, 0xFF49, false)
	c.memory.AddMMIOByte(&c.PPU.WY, 0xFF4A, false)
	c.memory.AddMMIOByte(&c.PPU.WX, 0xFF4B, false)
	c.memory.AddMMIOByte(&c.bank, 0xFF50, false)
	byt := byte(0x00)
	c.memory.AddMMIOByte(&byt, 0xFF7F, false) // Unused
	c.memory.AddMMIO(c.HRAM[:], 0xFF80, consts.HRAMSize, false)
	c.memory.AddMMIOByte(&c.interruptEnable, 0xFFFF, false)

	if c.config.BIOS != "" {
		c.ime = false
		c.rA = 0
		c.rF = 0
		c.rB = 0
		c.rC = 0
		c.rD = 0
		c.rE = 0
		c.rH = 0
		c.rL = 0
		c.rPC = 0x0
		c.rSP = 0x0
	} else {
		c.ime = false
		c.rA = 0x01
		c.rF = byte(ZeroFlag)
		if c.cartridge.ROMBank0[0x014D] == 0x00 {
			// Header checksum, if zero the Carry and Half Carry flags are set
			c.rF |= byte(CarryFlag) | byte(HalfCarryFlag)
		}
		c.rB = 0
		c.rC = 0x13
		c.rD = 0
		c.rE = 0xD8
		c.rH = 0x01
		c.rL = 0x4D
		c.rPC = 0x0100 // Program Counter starts at 0x0100
		c.rSP = 0xFFFE // Stack Pointer starts at 0xFFFE
	}
	c.halted = false
	c.exit = false
}

func (c *SM83) GetPC() uint16 {
	return c.rPC
}

func (c *SM83) Step() int {
	if !c.halted {
		if c.ime && c.interruptFlag != 0 {
			// Handle interrupts if IME is set and there are pending interrupts
			var interrupt impls.Interrupt
			if c.GetInterruptEnableFlag(impls.JoypadInterrupt) && c.GetInterruptFlag(impls.JoypadInterrupt) {
				slog.Info("Joypad interrupt triggered")
				interrupt = impls.JoypadInterrupt
			} else if c.GetInterruptEnableFlag(impls.SerialInterrupt) && c.GetInterruptFlag(impls.SerialInterrupt) {
				slog.Info("Serial interrupt triggered")
				interrupt = impls.SerialInterrupt
			} else if c.GetInterruptEnableFlag(impls.TimerInterrupt) && c.GetInterruptFlag(impls.TimerInterrupt) {
				slog.Info("Timer interrupt triggered")
				interrupt = impls.TimerInterrupt
			} else if c.GetInterruptEnableFlag(impls.LCDInterrupt) && c.GetInterruptFlag(impls.LCDInterrupt) {
				slog.Info("LCD interrupt triggered")
				interrupt = impls.LCDInterrupt
			} else if c.GetInterruptEnableFlag(impls.VBlankInterrupt) && c.GetInterruptFlag(impls.VBlankInterrupt) {
				slog.Info("VBlank interrupt triggered")
				interrupt = impls.VBlankInterrupt
			}

			if interrupt != 0 {
				c.SetFlag(NegativeFlag, false)
				c.SetFlag(HalfCarryFlag, false)
				c.SetFlag(CarryFlag, false)
				c.SetFlag(ZeroFlag, false)

				c.ime = false // Disable IME to prevent re-entrancy

				// Push PC onto stack
				c.rSP -= 2
				err := c.memory.Write16(c.rSP, c.rPC)
				if err != nil {
					panic(fmt.Sprintf("Failed to push PC onto stack: %v", err))
				}

				switch interrupt {
				case impls.JoypadInterrupt:
					c.rPC = 0x0060 // Joypad interrupt vector
				case impls.SerialInterrupt:
					c.rPC = 0x0058 // Serial interrupt vector
				case impls.TimerInterrupt:
					c.rPC = 0x0050 // Timer interrupt vector
				case impls.LCDInterrupt:
					c.rPC = 0x0048 // LCD interrupt vector
				case impls.VBlankInterrupt:
					c.rPC = 0x0040 // VBlank interrupt vector
				}
				c.SetInterruptFlag(interrupt, false) // Clear the interrupt flag
			}
		}

		instruction := c.fetch()

		// c.DebugRegisters() is expensive in the hot path
		if c.config.LogLevel == config.LogLevelDebug {
			slog.Debug(c.DebugRegisters())
			slog.Debug("Instruction", "instruction", instruction)
		}

		if instruction == nil {
			mem, err := c.memory.Read8(c.rPC - 1)
			if err != nil {
				panic(fmt.Sprintf("Failed to read memory at PC 0x%04X: %v", c.rPC-1, err))
			}
			panic(fmt.Sprintf("Unknown instruction at PC 0x%04X: 0x%02X", c.rPC-1, mem))
		}

		preBank := c.bank
		instruction.Exec(c)
		if c.bank != preBank && c.bank != 0 && c.config.BIOS != "" {
			// Boot rom disabled
			err := c.memory.RemoveMMIO(0x0000, consts.BIOSSize)
			if err != nil {
				panic(fmt.Sprintf("Failed to remove BIOS MMIO: %v", err))
			}
			c.memory.AddMMIO(c.cartridge.ROMBank0[:consts.BIOSSize], 0x0000, consts.BIOSSize, true)
		}

		return int(instruction.Cycles)
	}
	return 1
}

func (c *SM83) fetch() *OpCode {
	// Fetch the next instruction from memory at the current PC
	instruction, err := c.memory.Read8(c.rPC)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch instruction at PC 0x%04X: %v", c.rPC, err))
	}
	c.rPC++
	return opcodes[instruction]
}

func (c *SM83) DebugRegisters() string {
	var ret = "\n"
	ret += fmt.Sprintf("IE: 0x%02X\n", c.interruptEnable)
	ret += fmt.Sprintf(" A: 0x%02X\t  F: 0x%02X\n", c.rA, c.rF)
	ret += fmt.Sprintf(" B: 0x%02X\t  C: 0x%02X\n", c.rB, c.rC)
	ret += fmt.Sprintf(" D: 0x%02X\t  E: 0x%02X\n", c.rD, c.rE)
	ret += fmt.Sprintf(" H: 0x%02X\t  L: 0x%02X\n", c.rH, c.rL)
	ret += fmt.Sprintf("PC: 0x%04X\t SP: 0x%04X\n", c.rPC, c.rSP)
	ret += fmt.Sprintf("Flags: C: %t, H: %t, N: %t, Z: %t\n",
		c.GetFlag(CarryFlag),
		c.GetFlag(HalfCarryFlag),
		c.GetFlag(NegativeFlag),
		c.GetFlag(ZeroFlag),
	)
	ret += fmt.Sprintf("Interrupts: Joy: %t, Serial: %t, Timer: %t, LCD: %t, VBlank: %t\n",
		c.GetInterruptEnableFlag(impls.JoypadInterrupt),
		c.GetInterruptEnableFlag(impls.SerialInterrupt),
		c.GetInterruptEnableFlag(impls.TimerInterrupt),
		c.GetInterruptEnableFlag(impls.LCDInterrupt),
		c.GetInterruptEnableFlag(impls.VBlankInterrupt),
	)

	return ret
}

func (c *SM83) RunUntilFrame() [23040]byte {
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
		c.rF |= byte(flag)
	} else {
		c.rF &^= byte(flag)
	}
}

func (c *SM83) GetFlag(flag Flag) bool {
	return c.rF&byte(flag) != 0
}

func (c *SM83) Halt() {
	c.halted = true
}

func (c *SM83) IsHalted() bool {
	return c.halted
}

func (c *SM83) Resume() {
	c.halted = false
}

func (c *SM83) Quit() {
	c.exit = true
}
