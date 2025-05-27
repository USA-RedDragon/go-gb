package ppu

import (
	"log/slog"

	"github.com/USA-RedDragon/go-gb/internal/consts"
)

type ppuState uint8

const (
	ppuStateOAMSearch     ppuState = iota // OAM Search state
	ppuStatePixelTransfer                 // Pixel Transfer state
	ppuStateHBlank                        // H-Blank state
	ppuStateVBlank                        // V-Blank state
)

type PPU struct {
	VRAM       [consts.VRAMSize]byte // 8KB of VRAM
	LCDControl byte                  // LCDC, LCD control register
	LCDStatus  byte                  // STAT, LCD status register
	SCX        byte                  // SCX, scroll X register
	SCY        byte                  // SCY, scroll Y register
	LY         byte                  // LY, LCD Y coordinate register, read-only
	LYC        byte                  // LYC, LY Compare register
	BGP        byte                  // BGP, background palette data
	OBP0       byte                  // OBP0, object palette 0 data
	OBP1       byte                  // OBP1, object palette 1 data
	OAM        [consts.OAMSize]byte  // OAM, Object Attribute Memory

	state ppuState // Current state of the PPU
	ticks uint16
	x     byte // Current X coordinate in the pixel transfer state
}

func NewPPU() *PPU {
	ppu := &PPU{}
	ppu.Reset()
	return ppu
}

func (ppu *PPU) Reset() {
	ppu.VRAM = [consts.VRAMSize]byte{}
	ppu.OAM = [consts.OAMSize]byte{}
	ppu.LCDControl = 0x00
	ppu.LCDStatus = 0x00
	ppu.SCX = 0x00
	ppu.SCY = 0x00
	ppu.LY = 0x00
	ppu.LYC = 0x00
	ppu.BGP = 0x00
	ppu.OBP0 = 0x00
	ppu.OBP1 = 0x00
}

func (ppu *PPU) Step() {
	ppu.ticks++
	switch ppu.state {
	case ppuStateOAMSearch:
		// TODO: find sprites
		if ppu.ticks == 80 {
			slog.Debug("PPU: OAM Search complete", "LY", ppu.LY, "ticks", ppu.ticks, "x", ppu.x)
			ppu.state = ppuStatePixelTransfer
		}
	case ppuStatePixelTransfer:
		ppu.x++
		// TODO: transfer pixels to the screen
		if ppu.x == 160 {
			slog.Debug("PPU: Pixel Transfer complete", "LY", ppu.LY, "ticks", ppu.ticks, "x", ppu.x)
			ppu.state = ppuStateHBlank
		}
	case ppuStateHBlank:
		// no-ops
		if ppu.ticks == 456 {
			ppu.ticks = 0
			ppu.LY++
			if ppu.LY == 144 {
				slog.Debug("PPU: VBlank started", "LY", ppu.LY, "ticks", ppu.ticks)
				ppu.state = ppuStateVBlank
			} else {
				slog.Debug("PPU: HBlank complete", "LY", ppu.LY, "ticks", ppu.ticks)
				ppu.state = ppuStateOAMSearch
			}
		}
	case ppuStateVBlank:
		// no-ops
		if ppu.ticks == 456 {
			ppu.ticks = 0
			ppu.LY++
			if ppu.LY == 153 {
				ppu.LY = 0
				slog.Debug("PPU: VBlank ended, LY reset", "LY", ppu.LY, "ticks", ppu.ticks)
				ppu.state = ppuStateOAMSearch
			}
		}
		ppu.state = ppuStateOAMSearch
	default:
		slog.Error("Unknown PPU state encountered", "state", ppu.state)
	}
}
