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

const (
	// Bit 0 - BG/Window Display/Priority     (0=Off, 1=On)
	LCDCBGDisplay uint8 = 1 << iota
	// Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	LCDCSpriteDisplayEnable
	// Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	LCDCSpriteSize
	// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	LCDCBGTileMapDisplaySelect
	// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	LCDCBGWindowTileDataSelect
	// Bit 5 - Window Display Enable          (0=Off, 1=On)
	LCDCWindowDisplayEnable
	// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	LCDCWindowTileMapDisplayeSelect
	// Bit 7 - LCD Display Enable             (0=Off, 1=On)
	LCDCDisplayEnable
)

type PPU struct {
	VRAM       [consts.VRAMSize]byte // 8KB of VRAM
	LCDControl byte                  // LCDC, LCD control register
	LCDStatus  byte                  // STAT, LCD status register
	SCX        byte                  // SCX, scroll X register
	SCY        byte                  // SCY, scroll Y register
	LY         byte                  // LY, LCD Y coordinate register, read-only
	LYC        byte                  // LYC, LY Compare register
	WX         byte                  // WX, window X coordinate register
	WY         byte                  // WY, window Y coordinate register
	BGP        byte                  // BGP, background palette data
	OBP0       byte                  // OBP0, object palette 0 data
	OBP1       byte                  // OBP1, object palette 1 data
	OAM        [consts.OAMSize]byte  // OAM, Object Attribute Memory
	Fetcher    *Fetcher              // Fetcher for pixel data

	HaveFrame     bool
	FrameBuffer_A [consts.FrameBufferSize]byte
	FrameBuffer_B [consts.FrameBufferSize]byte // Frame buffer for double buffering

	state    ppuState // Current state of the PPU
	ticks    uint16
	x        byte // Current X coordinate in the pixel transfer state
	disabled bool // Indicates if the PPU is disabled
}

func NewPPU() *PPU {
	ppu := &PPU{}
	ppu.Fetcher = NewFetcher(ppu)
	ppu.Reset()
	return ppu
}

func (ppu *PPU) GetFrame() [consts.FrameBufferSize]byte {
	ppu.HaveFrame = false
	return ppu.FrameBuffer_B
}

func (ppu *PPU) Reset() {
	ppu.VRAM = [consts.VRAMSize]byte{}
	ppu.OAM = [consts.OAMSize]byte{}
	ppu.HaveFrame = false
	ppu.Fetcher.Reset()
	ppu.LCDControl = 0x00
	ppu.LCDStatus = 0x00
	ppu.SCX = 0x00
	ppu.SCY = 0x00
	ppu.LY = 0x00
	ppu.LYC = 0x00
	ppu.BGP = 0x00
	ppu.OBP0 = 0x00
	ppu.OBP1 = 0x00
	ppu.WX = 0x00
	ppu.WY = 0x00
	ppu.disabled = true
}

func (ppu *PPU) GetPalleteColor(idx byte) byte {
	// ppu.BGP is the background palette data
	// Bits 7 and 6 are the color for index 3
	// Bits 5 and 4 are the color for index 2
	// Bits 3 and 2 are the color for index 1
	// Bits 1 and 0 are the color for index 0
	var color byte
	switch idx {
	case 0:
		color = (ppu.BGP & 0x03) // Get bits 1 and 0
	case 1:
		color = (ppu.BGP & 0x0C) >> 2 // Get bits 3 and 2
	case 2:
		color = (ppu.BGP & 0x30) >> 4 // Get bits 5 and 4
	case 3:
		color = (ppu.BGP & 0xC0) >> 6 // Get bits 7 and 6
	default:
		slog.Error("PPU: Invalid palette index", "index", idx)
		color = 0x00 // Default to black if index is invalid
	}
	return color
}

func (ppu *PPU) Step() {
	if ppu.disabled {
		if ppu.LCDControl&LCDCDisplayEnable != 0 {
			ppu.disabled = false
			ppu.state = ppuStateOAMSearch
		} else {
			return
		}
	} else {
		if ppu.LCDControl&LCDCDisplayEnable == 0 {
			// Turn screen off and reset PPU state machine.
			ppu.LY = 0
			ppu.x = 0
			ppu.disabled = true
			return
		}
	}

	if ppu.LCDControl&LCDCBGWindowTileDataSelect == 1 {
		slog.Error("PPU: BG/Window Tile Data Select is set to 1, this is not supported in DMG mode")
	}

	switch ppu.state {
	case ppuStateOAMSearch:
		// TODO: find sprites
		if ppu.ticks == 80 {
			ppu.x = 0
			ppu.Fetcher.Reset()
			slog.Debug("PPU: OAM Search complete", "LY", ppu.LY, "ticks", ppu.ticks, "x", ppu.x)
			ppu.state = ppuStatePixelTransfer
		}
	case ppuStatePixelTransfer:
		if ppu.ticks%2 == 0 {
			ppu.Fetcher.Step()
		}
		if ppu.Fetcher.PixelFIFO.Size() <= 8 {
			break
		}

		// Put a pixel from the FIFO on screen.
		pixelIdx, ok := ppu.Fetcher.PixelFIFO.Pop()
		if !ok {
			slog.Error("PPU: Pixel FIFO is empty, cannot transfer pixel", "LY", ppu.LY, "ticks", ppu.ticks, "x", ppu.x)
			break
		}

		fbOffset := (uint16(ppu.LY) * 160) + uint16(ppu.x)
		ppu.FrameBuffer_A[fbOffset] = ppu.GetPalleteColor(pixelIdx)

		ppu.x++
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
				ppu.FrameBuffer_B = ppu.FrameBuffer_A
				ppu.FrameBuffer_A = [consts.FrameBufferSize]byte{}
				ppu.HaveFrame = true
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
	ppu.ticks++
}
