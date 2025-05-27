package cpu

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/memory"
)

const (
	ROMBankSize          = 16384 // 16KB ROM bank size
	VRAMSize             = 8192  // 8KB of VRAM
	CartridgeRAMBankSize = 8192  // 8KB of cartridge RAM bank
	RAMSize              = 8192  // 8KB of RAM
)

type SM83 struct {
	config *config.Config
	memory memory.MMIO // Memory-mapped I/O
	halted bool
	exit   bool

	ROMBank0          [ROMBankSize]byte          // 16KB of ROM bank 0
	ROMSwitchableBank [ROMBankSize]byte          // 16KB of ROM, can be expanded for more banks
	VRAM              [VRAMSize]byte             // 8KB of VRAM
	CartridgeRAM      [CartridgeRAMBankSize]byte // 8KB of cartridge RAM, if present
	RAM               [RAMSize]byte              // 8KB of RAM

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

func NewSM83(config *config.Config) *SM83 {
	cpu := &SM83{
		config: config,
		memory: memory.MMIO{},
	}

	cpu.Reset()

	cpu.memory.AddMMIO(cpu.ROMBank0[:], 0x0, ROMBankSize)
	cpu.memory.AddMMIO(cpu.ROMSwitchableBank[:], 0x4000, ROMBankSize)
	cpu.memory.AddMMIO(cpu.VRAM[:], 0x8000, VRAMSize)
	cpu.memory.AddMMIO(cpu.CartridgeRAM[:], 0xA000, CartridgeRAMBankSize)
	cpu.memory.AddMMIO(cpu.RAM[:], 0xC000, RAMSize)

	if config.ROM != "" {
		cpu.loadROM()
	}

	return cpu
}

func (c *SM83) Reset() {
	c.RAM = [RAMSize]byte{}
	c.ROMBank0 = [ROMBankSize]byte{}
	c.ROMSwitchableBank = [ROMBankSize]byte{}
	c.VRAM = [VRAMSize]byte{}
	c.CartridgeRAM = [CartridgeRAMBankSize]byte{}
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
	c.r_PC = 0x0100 // Program Counter starts at 0x0100
	c.r_SP = 0
	c.halted = false
	c.exit = false
}

func (c *SM83) GetTitle() string {
	str := []byte{}
	for i := 0x134; i < 0x144; i++ {
		mem, err := c.memory.Read8(uint16(i))
		if err != nil {
			panic(fmt.Sprintf("Failed to read memory at 0x%04X: %v", i, err))
		}
		if mem == 0 {
			break
		}
		str = append(str, mem)
	}
	return string(str)
}

func (c *SM83) Step() {
	if !c.halted {
		c.step()
	}
}

func (c *SM83) step() {
	instruction := c.fetch()

	slog.Debug("Instruction", "instruction", fmt.Sprintf("0x%02X", instruction))
	slog.Debug(c.DebugRegisters())

	c.execute(instruction)
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

func (c *SM83) loadROM() {
	rom, err := os.ReadFile(c.config.ROM)
	if err != nil {
		panic(fmt.Sprintf("Failed to load rom: %v", err))
	}
	switch len(rom) {
	case 0:
		panic("ROM file is empty")
	case ROMBankSize:
		slog.Debug("Loaded ROM bank 0", "size", len(rom))
		copy(c.ROMBank0[:], rom)
	case ROMBankSize * 2:
		slog.Debug("Loaded ROM bank 0 and switchable bank", "size", len(rom))
		copy(c.ROMBank0[:], rom[:ROMBankSize])
		copy(c.ROMSwitchableBank[:], rom[ROMBankSize:])
	default:
		panic(fmt.Sprintf("ROM size is %d bytes, expected 16KB or 32KB", len(rom)))
	}

	slog.Debug("ROM loaded", "size", len(rom), "title", c.GetTitle())
}

func (c *SM83) DebugRegisters() string {
	var ret = "\n"
	ret += fmt.Sprintf("IR: 0x%02X\t IE: 0x%02X\n", c.r_IR, c.r_IE)
	ret += fmt.Sprintf(" A: 0x%02X\t  F: 0x%02X\n", c.r_A, c.r_F)
	ret += fmt.Sprintf(" B: 0x%02X\t  C: 0x%02X\n", c.r_B, c.r_C)
	ret += fmt.Sprintf(" D: 0x%02X\t  E: 0x%02X\n", c.r_D, c.r_E)
	ret += fmt.Sprintf(" H: 0x%02X\t  L: 0x%02X\n", c.r_H, c.r_L)
	ret += fmt.Sprintf("PC: 0x%04X\n", c.r_PC)
	ret += fmt.Sprintf("SP: 0x%04X\n", c.r_SP)

	return ret
}

func (c *SM83) Run() {
	cycleTime := time.Second / 4194304
	prevTime := time.Now()
	for !c.exit {
		c.Step()
		time.Sleep(cycleTime - time.Since(prevTime))
		prevTime = time.Now()
	}
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
