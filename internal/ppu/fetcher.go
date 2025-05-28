package ppu

import "github.com/USA-RedDragon/go-gb/internal/consts"

type fetcherState uint8

const (
	fetcherStateTileNumber fetcherState = iota
	fetcherStateTileDataLow
	fetcherStateTileDataHigh
	fetcherStatePushToFIFO
)

type Fetcher struct {
	state     fetcherState
	PixelFIFO *FIFO
	PPU       *PPU // Reference to the PPU for accessing VRAM and other registers

	tileNum     uint8   // Tile number being fetched
	tileIndex   byte    // Index of the current tile in the line
	pixelBuffer [8]byte // Buffer for the pixel data of the current tile line
}

func NewFetcher(ppu *PPU) *Fetcher {
	fetcher := &Fetcher{
		PixelFIFO: NewFIFO(),
		PPU:       ppu,
	}
	fetcher.Reset()
	return fetcher
}

func (f *Fetcher) Reset() {
	f.PixelFIFO.Reset()
	f.state = fetcherStateTileNumber
	f.pixelBuffer = [8]byte{}
	f.tileIndex = 0
	f.tileNum = 0
}

func (f *Fetcher) Step() {
	tileLine := f.PPU.LY % 8
	switch f.state {
	case fetcherStateTileNumber:
		// Fetch the tile number from VRAM
		mapAddr := consts.BackgroundMapOffset + (uint16(tileLine) * 32)
		f.tileNum = f.PPU.VRAM[uint16(f.tileIndex)+mapAddr]
		f.state = fetcherStateTileDataLow
	case fetcherStateTileDataLow:
		// Fetch the low byte of the tile data
		offset := uint16(f.tileNum) * 16 // Each tile is 16 bytes (8 lines, 2 bytes per line)
		addr := offset + (uint16(tileLine) * 2)
		lowByte := f.PPU.VRAM[addr]
		for i := range 8 {
			f.pixelBuffer[i] = (lowByte >> i) & 0x01
		}
		f.state = fetcherStateTileDataHigh
	case fetcherStateTileDataHigh:
		// Fetch the high byte of the tile data
		offset := uint16(f.tileNum) * 16 // Each tile is 16 bytes (8 lines, 2 bytes per line)
		addr := offset + (uint16(tileLine) * 2) + 1
		highByte := f.PPU.VRAM[addr]
		for i := range 8 {
			// Combine low and high bytes to form the pixel data
			f.pixelBuffer[i] |= ((highByte >> i) & 0x01) << 1
		}
		f.state = fetcherStatePushToFIFO
	case fetcherStatePushToFIFO:
		// Push the fetched tile data to the FIFO
		if f.PixelFIFO.Size() <= 8 {
			for i := 7; i >= 0; i-- {
				f.PixelFIFO.Push(f.pixelBuffer[i])
			}
		}
		f.tileIndex++
		f.state = fetcherStateTileNumber
	}
}
