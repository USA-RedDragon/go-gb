package emulator

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/USA-RedDragon/go-gb/internal/cartridge"
	"github.com/USA-RedDragon/go-gb/internal/config"
	"github.com/USA-RedDragon/go-gb/internal/cpu"
	"github.com/USA-RedDragon/go-gb/internal/impls"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	return upscaledImage.Pix
}

func (e *Emulator) convertToScreen(frame [23040]byte) []byte {
	originalRender := image.NewRGBA(image.Rect(0, 0, 160, 144))

	// The input frame is expected to be in 2BPP format, where each pixel is represented by 2 bits.
	for i := 0; i < len(frame); i++ {
		pixel := frame[i]
		switch pixel {
		case 0x00:
			// White
			originalRender.Pix[i*4+0] = 0xFF // R
			originalRender.Pix[i*4+1] = 0xFF // G
			originalRender.Pix[i*4+2] = 0xFF // B
			originalRender.Pix[i*4+3] = 0xFF // A
		case 0x01:
			// 33% Gray
			originalRender.Pix[i*4+0] = 0x55 // R
			originalRender.Pix[i*4+1] = 0x55 // G
			originalRender.Pix[i*4+2] = 0x55 // B
			originalRender.Pix[i*4+3] = 0xFF // A
		case 0x02:
			// 66% Gray
			originalRender.Pix[i*4+0] = 0xAA // R
			originalRender.Pix[i*4+1] = 0xAA // G
			originalRender.Pix[i*4+2] = 0xAA // B
			originalRender.Pix[i*4+3] = 0xFF // A
		case 0x03:
			// Black
			originalRender.Pix[i*4+0] = 0x00 // R
			originalRender.Pix[i*4+1] = 0x00 // G
			originalRender.Pix[i*4+2] = 0x00 // B
			originalRender.Pix[i*4+3] = 0xFF // A
		default:
			// Invalid pixel value, default to red
			originalRender.Pix[i*4+0] = 0xFF // R
			originalRender.Pix[i*4+1] = 0x00 // G
			originalRender.Pix[i*4+2] = 0x00 // B
			originalRender.Pix[i*4+3] = 0xFF // A
		}
	}

	return e.upscale(originalRender)
}

func (e *Emulator) updateFrame() {
	e.frame = e.convertToScreen(e.cpu.RunUntilFrame())
}

func (e *Emulator) Update() error {
	start := time.Now()

	// Frame stepping
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if e.cpu.IsHalted() {
			e.cpu.Resume()
		}
		e.updateFrame()
		e.cpu.Halt()
	} else if inpututil.KeyPressDuration(ebiten.KeyF) > 30 {
		if e.cpu.IsHalted() {
			e.cpu.Resume()
		}
		e.updateFrame()
		e.cpu.Halt()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		e.cpu.Resume()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		e.cpu.Halt()
	}

	if !e.cpu.IsHalted() {
		e.updateFrame()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		halted := e.cpu.IsHalted()
		e.cpu.Reset()
		if halted {
			e.cpu.Halt()
		}
	}

	e.frametime = int(time.Since(start).Milliseconds())
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.stopped {
		return
	}
	screen.WritePixels(e.frame)
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"FPS: %0.2f\nFrame Time: %dms\nTPS: %0.2f\nPC: 0x%04X\n%s",
			1000.0/float64(e.frametime),
			e.frametime,
			ebiten.ActualTPS(),
			e.cpu.GetPC(),
			fmt.Sprintf("Interrupts:\n\tJoy: %t, Serial: %t, Timer: %t, LCD: %t, VBlank: %t\n",
				e.cpu.GetInterruptEnableFlag(impls.JoypadInterrupt),
				e.cpu.GetInterruptEnableFlag(impls.SerialInterrupt),
				e.cpu.GetInterruptEnableFlag(impls.TimerInterrupt),
				e.cpu.GetInterruptEnableFlag(impls.LCDInterrupt),
				e.cpu.GetInterruptEnableFlag(impls.VBlankInterrupt),
			),
		),
	)
}

func (e *Emulator) Layout(_, _ int) (int, int) {
	return int(e.config.Scale * 160), int(e.config.Scale * 144)
}

func (e *Emulator) Stop() {
	e.stopped = true
	e.cpu.Halt()
	os.Exit(0)
}
