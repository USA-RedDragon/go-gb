package consts

const (
	ROMBankSize          = 16384 // 16KB ROM bank size
	VRAMSize             = 8192  // 8KB of VRAM
	BackgroundMapOffset  = 6144  // 6KB into the VRAM for background map data
	CartridgeRAMBankSize = 8192  // 8KB of cartridge RAM bank
	RAMSize              = 8192  // 8KB of RAM
	HRAMSize             = 127   // 127 bytes of HRAM (0xFF80 - 0xFFFE)
	BIOSSize             = 256   // 256 bytes of BIOS (0x0000 - 0x00FF)
	ProhibitedSize       = 96    // 96 bytes of prohibited memory (0xFEA0 - 0xFEFF)
	OAMSize              = 160   // 160 bytes of OAM (0xFE00 - 0xFE9F)
)
