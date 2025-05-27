package ppu

import "github.com/USA-RedDragon/go-gb/internal/consts"

type PPU struct {
	VRAM       [consts.VRAMSize]byte // 8KB of VRAM
	LCDControl byte                  // LCDC, LCD control register
	LCDStatus  byte                  // STAT, LCD status register
	SCX        byte                  // SCX, scroll X register
	SCY        byte                  // SCY, scroll Y register
	LY         byte                  // LY, LCD Y coordinate register, read-only
	LYC        byte                  // LYC, LY Compare register
	BGP        byte                  // BGP, background palette data
}

func NewPPU() *PPU {
	return &PPU{
		VRAM: [consts.VRAMSize]byte{},
	}
}

func (ppu *PPU) Reset() {
	ppu.VRAM = [consts.VRAMSize]byte{}
	ppu.LCDControl = 0x00
	ppu.LCDStatus = 0x00
	ppu.SCX = 0x00
	ppu.SCY = 0x00
	ppu.LY = 0x00
	ppu.LYC = 0x00
	ppu.BGP = 0x00
}
