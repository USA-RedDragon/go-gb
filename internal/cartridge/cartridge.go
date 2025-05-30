package cartridge

import (
	"fmt"
	"os"

	"github.com/USA-RedDragon/go-gb/internal/consts"
)

type RAMSize byte

func (r RAMSize) String() string {
	switch r {
	case 0x00:
		return "No RAM"
	case 0x01:
		return "2KB"
	case 0x02:
		return "8KB"
	case 0x03:
		return "32KB"
	case 0x04:
		return "128KB"
	case 0x05:
		return "64KB"
	default:
		return fmt.Sprintf("Unknown RAM size: %d", r)
	}
}

func (r RAMSize) NumberOfBanks() int {
	switch r {
	case 0x00:
		return 0
	case 0x01:
		return -1
	case 0x02:
		return 1
	case 0x03:
		return 4
	case 0x04:
		return 16
	case 0x05:
		return 8
	default:
		return -1 // Unknown size
	}
}

func (r RAMSize) Bytes() int {
	switch r {
	case 0x00:
		return 0
	case 0x01:
		return 2 * 1024
	case 0x02:
		return 8 * 1024
	case 0x03:
		return 32 * 1024
	case 0x04:
		return 128 * 1024
	case 0x05:
		return 64 * 1024
	default:
		return -1 // Unknown size
	}
}

type ROMSize byte

func (r ROMSize) String() string {
	switch r {
	case 0x00:
		return "32KB"
	case 0x01:
		return "64KB"
	case 0x02:
		return "128KB"
	case 0x03:
		return "256KB"
	case 0x04:
		return "512KB"
	case 0x05:
		return "1MB"
	case 0x06:
		return "2MB"
	case 0x07:
		return "4MB"
	case 0x08:
		return "8MB"
	case 0x52:
		return "1.1MB"
	case 0x53:
		return "1.2MB"
	case 0x54:
		return "1.5MB"
	default:
		return fmt.Sprintf("Unknown ROM size: %d", r)
	}
}

func (r ROMSize) NumberOfBanks() int {
	switch r {
	case 0x00:
		return 2
	case 0x01:
		return 4
	case 0x02:
		return 8
	case 0x03:
		return 16
	case 0x04:
		return 32
	case 0x05:
		return 64
	case 0x06:
		return 128
	case 0x07:
		return 256
	case 0x08:
		return 512
	case 0x52:
		return 72
	case 0x53:
		return 80
	case 0x54:
		return 96
	default:
		return -1 // Unknown size
	}
}

func (r ROMSize) Bytes() int {
	switch r {
	case 0x00:
		return 32 * 1024
	case 0x01:
		return 64 * 1024
	case 0x02:
		return 128 * 1024
	case 0x03:
		return 256 * 1024
	case 0x04:
		return 512 * 1024
	case 0x05:
		return 1 * 1024 * 1024
	case 0x06:
		return 2 * 1024 * 1024
	case 0x07:
		return 4 * 1024 * 1024
	case 0x08:
		return 8 * 1024 * 1024
	case 0x52:
		return 1024*1024 + 128
	case 0x53:
		return 1024*1024 + 256
	case 0x54:
		return 1024*1024 + 512
	default:
		return -1 // Unknown size
	}
}

type Cartridge struct {
	ROMBank0           [consts.ROMBankSize]byte
	AdditionalROMBanks [][consts.ROMBankSize]byte // Additional ROM banks if the ROM is larger than 16KB
	CartridgeRAMBanks  [][]byte
	RAMSize            RAMSize // Size of the cartridge RAM, if present
	ROMSize            ROMSize
	OldPublisher       OldPublisher
	NewPublisher       NewPublisher
	Publisher          string
	Title              string
	SGB                bool
	CGB                bool
	CGBOnly            bool
	Version            uint8
	CartridgeType      CartridgeType
	Japanese           bool
}

func NewCartridge(romPath string) (*Cartridge, error) {
	c := &Cartridge{}

	// Load the ROM file
	romData, err := os.ReadFile(romPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ROM file: %w", err)
	}
	if len(romData) < consts.ROMBankSize {
		return nil, fmt.Errorf("ROM file is too small, must be at least %d bytes", consts.ROMBankSize)
	}

	c.ROMSize = ROMSize(romData[0x148])
	c.RAMSize = RAMSize(romData[0x149])

	// Copy the first 16KB to ROMBank0
	copy(c.ROMBank0[:], romData[:consts.ROMBankSize])
	c.AdditionalROMBanks = make([][consts.ROMBankSize]byte, c.ROMSize.NumberOfBanks()-1)

	for i := 1; i <= c.ROMSize.NumberOfBanks()-1; i++ {
		start := i * consts.ROMBankSize
		end := start + consts.ROMBankSize
		if end > len(romData) {
			end = len(romData)
		}
		copy(c.AdditionalROMBanks[i-1][:], romData[start:end])
	}

	if c.RAMSize.Bytes() > 0 {
		c.CartridgeRAMBanks = make([][]byte, c.RAMSize.NumberOfBanks())
		for i := range c.RAMSize.NumberOfBanks() {
			c.CartridgeRAMBanks[i] = make([]byte, consts.CartridgeRAMBankSize)
		}
	} else {
		c.CartridgeRAMBanks = [][]byte{}
	}

	c.OldPublisher = OldPublisher(romData[0x014B])
	c.Publisher = c.OldPublisher.String()
	if c.OldPublisher == OldPublisherSeeOther {
		c.NewPublisher = NewPublisher(string([]byte{c.ROMBank0[0x0144], c.ROMBank0[0x0145]}))
		c.Publisher = c.NewPublisher.String()
	}

	c.Title = c.getTitle()

	c.SGB = c.ROMBank0[0x146] == 0x03
	c.CGBOnly = c.ROMBank0[0x143] == 0xC0
	c.CGB = c.ROMBank0[0x143] == 0x80 || c.CGBOnly

	c.Version = c.ROMBank0[0x014C]

	c.Japanese = c.ROMBank0[0x014A] == 0x00

	c.CartridgeType = CartridgeType(c.ROMBank0[0x0147])

	return c, nil
}
func (c *Cartridge) Reset() {
	if c.RAMSize.Bytes() > 0 {
		c.CartridgeRAMBanks = make([][]byte, c.RAMSize.NumberOfBanks())
		for i := range c.RAMSize.NumberOfBanks() {
			c.CartridgeRAMBanks[i] = make([]byte, consts.CartridgeRAMBankSize)
		}
	} else {
		c.CartridgeRAMBanks = [][]byte{}
	}
}

func (c *Cartridge) NintendoLogoValid() bool {
	// The Nintendo logo is a specific sequence of bytes that must be present
	// at the start of the ROM. This is a simplified check.
	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}
	if len(c.ROMBank0) < len(logo) {
		return false
	}
	for i := 0x104; i < 0x134; i++ {
		if c.ROMBank0[i] != logo[i] {
			return false
		}
	}
	return true
}

func (c *Cartridge) getTitle() string {
	str := []byte{}
	for i := 0x134; i < 0x144; i++ {
		if i >= len(c.ROMBank0) || c.ROMBank0[i] == 0 {
			break
		}
		str = append(str, c.ROMBank0[i])
	}
	return string(str)
}

func (c *Cartridge) String() string {
	return fmt.Sprintf("Cartridge: Title: %s, Publisher: %s, RAM size: %s, ROM size: %s, SGB: %t, Version: %d, Type: %s, Japanese: %t, CGB: %t, CGBOnly: %t",
		c.Title, c.Publisher, c.RAMSize, c.ROMSize, c.SGB, c.Version, c.CartridgeType, c.Japanese, c.CGB, c.CGBOnly)
}
