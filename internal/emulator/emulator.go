package emulator

import (
	"fmt"
	"image"
	"os"
	"runtime/pprof"
	"time"

	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/cpu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/draw"
)

type Emulator struct {
	config    *config.Config
	cpu       *cpu.SM83
	stopped   bool
	frametime int
	frame     []byte
}

func New(config *config.Config, cartridge *cartridge.Cartridge) *Emulator {
	emu := &Emulator{
		config: config,
		cpu:    cpu.NewSM83(config, cartridge),
	}
	return emu
}

func (e *Emulator) upscale(render *image.RGBA) []byte {
	targetWidth := int(160 * e.config.Scale)
	targetHeight := int(144 * e.config.Scale)

	// Create a new blank image with the target dimensions
	upscaledImage := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Perform bicubic interpolation
	draw.CatmullRom.Scale(upscaledImage, upscaledImage.Bounds(), render, render.Bounds(), draw.Src, nil)

	// Convert the upscaled image to a byte array for the renderer
	upscaled := make([]byte, targetWidth*targetHeight*4)
	copy(upscaled, upscaledImage.Pix)

	return upscaled
}

func (e *Emulator) convertToScreen(frame []byte) []byte {
	// The input frame is expected to be in 2BPP format, where each pixel is represented by 2 bits.
	screen := make([]byte, len(frame)*4) // 2BPP to 4BPP conversion
	for i := 0; i < len(frame); i++ {
		pixel := frame[i]
		switch pixel {
		case 0x00:
			// White
			screen[i*4+0] = 0xFF // R
			screen[i*4+1] = 0xFF // G
			screen[i*4+2] = 0xFF // B
			screen[i*4+3] = 0xFF // A
		case 0x01:
			// 66% Gray
			screen[i*4+0] = 0xAA // R
			screen[i*4+1] = 0xAA // G
			screen[i*4+2] = 0xAA // B
			screen[i*4+3] = 0xFF // A
		case 0x02:
			// 33% Gray
			screen[i*4+0] = 0x55 // R
			screen[i*4+1] = 0x55 // G
			screen[i*4+2] = 0x55 // B
			screen[i*4+3] = 0xFF // A
		case 0x03:
			// Black
			screen[i*4+0] = 0x00 // R
			screen[i*4+1] = 0x00 // G
			screen[i*4+2] = 0x00 // B
			screen[i*4+3] = 0xFF // A
		default:
			// Invalid pixel value, default to red
			screen[i*4+0] = 0xFF // R
			screen[i*4+1] = 0x00 // G
			screen[i*4+2] = 0x00 // B
			screen[i*4+3] = 0xFF // A
		}
	}

	var originalRender *image.RGBA
	originalRender = image.NewRGBA(image.Rect(0, 0, 160, 144))
	copy(originalRender.Pix, screen)

	return e.upscale(originalRender)
}

func (e *Emulator) Update() error {
	start := time.Now()
	e.frame = e.convertToScreen(e.cpu.RunUntilFrame())
	e.frametime = int(time.Since(start).Milliseconds())
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.stopped {
		return
	}
	screen.Clear()
	screen.WritePixels(e.frame)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nFrame Time: %dms\nTPS: %0.2f", 1000.0/float64(e.frametime), e.frametime, ebiten.ActualTPS()))
}

func (e *Emulator) Layout(_, _ int) (int, int) {
	return int(e.config.Scale * 160), int(e.config.Scale * 144)
}

func (e *Emulator) Stop() {
	e.stopped = true
	pprof.StopCPUProfile()
	e.cpu.Halt()
	os.Exit(0)
}
