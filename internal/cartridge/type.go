package cartridge

type CartridgeType byte

const (
	CartridgeTypeROMOnly                    CartridgeType = 0x00
	CartridgeTypeMBC1                       CartridgeType = 0x01
	CartridgeTypeMBC1RAM                    CartridgeType = 0x02
	CartridgeTypeMBC1RAMBattery             CartridgeType = 0x03
	CartridgeTypeMBC2                       CartridgeType = 0x05
	CartridgeTypeMBC2Battery                CartridgeType = 0x06
	CartridgeTypeROMRAM                     CartridgeType = 0x08
	CartridgeTypeROMRAMBattery              CartridgeType = 0x09
	CartridgeTypeMMM01                      CartridgeType = 0x0B
	CartridgeTypeMMM01RAM                   CartridgeType = 0x0C
	CartridgeTypeMMM01RAMBattery            CartridgeType = 0x0D
	CartridgeTypeMBC3TimerRAM               CartridgeType = 0x0F
	CartridgeTypeMBC3TimerRAMBattery        CartridgeType = 0x10
	CartridgeTypeMBC3                       CartridgeType = 0x11
	CartridgeTypeMBC3RAM                    CartridgeType = 0x12
	CartridgeTypeMBC3RAMBattery             CartridgeType = 0x13
	CartridgeTypeMBC5                       CartridgeType = 0x19
	CartridgeTypeMBC5RAM                    CartridgeType = 0x1A
	CartridgeTypeMBC5RAMBattery             CartridgeType = 0x1B
	CartridgeTypeMBC5Rumble                 CartridgeType = 0x1C
	CartridgeTypeMBC5RumbleRAM              CartridgeType = 0x1D
	CartridgeTypeMBC5RumbleRAMBattery       CartridgeType = 0x1E
	CartridgeTypeMBC6                       CartridgeType = 0x20
	CartridgeTypeMBC7SensorRumbleRAMBattery CartridgeType = 0x22
	CartridgeTypePocketCamera               CartridgeType = 0xFC
	CartridgeTypeBandaiTAMA5                CartridgeType = 0xFD
	CartridgeTypeHuC3                       CartridgeType = 0xFE
	CartridgeTypeHuC1RAMBattery             CartridgeType = 0xFF
)

func (c CartridgeType) String() string {
	switch c {
	case CartridgeTypeROMOnly:
		return "ROM Only"
	case CartridgeTypeMBC1:
		return "MBC1"
	case CartridgeTypeMBC1RAM:
		return "MBC1 + RAM"
	case CartridgeTypeMBC1RAMBattery:
		return "MBC1 + RAM + Battery"
	case CartridgeTypeMBC2:
		return "MBC2"
	case CartridgeTypeMBC2Battery:
		return "MBC2 + Battery"
	case CartridgeTypeROMRAM:
		return "ROM + RAM"
	case CartridgeTypeROMRAMBattery:
		return "ROM + RAM + Battery"
	case CartridgeTypeMMM01:
		return "MMM01"
	case CartridgeTypeMMM01RAM:
		return "MMM01 + RAM"
	case CartridgeTypeMMM01RAMBattery:
		return "MMM01 + RAM + Battery"
	case CartridgeTypeMBC3TimerRAM:
		return "MBC3 + Timer + RAM"
	case CartridgeTypeMBC3TimerRAMBattery:
		return "MBC3 + Timer + RAM + Battery"
	case CartridgeTypeMBC3:
		return "MBC3"
	case CartridgeTypeMBC3RAM:
		return "MBC3 + RAM"
	case CartridgeTypeMBC3RAMBattery:
		return "MBC3 + RAM + Battery"
	case CartridgeTypeMBC5:
		return "MBC5"
	case CartridgeTypeMBC5RAM:
		return "MBC5 + RAM"
	case CartridgeTypeMBC5RAMBattery:
		return "MBC5 + RAM + Battery"
	case CartridgeTypeMBC5Rumble:
		return "MBC5 + Rumble"
	case CartridgeTypeMBC5RumbleRAM:
		return "MBC5 + Rumble + RAM"
	case CartridgeTypeMBC5RumbleRAMBattery:
		return "MBC5 + Rumble + RAM + Battery"
	case CartridgeTypeMBC6:
		return "MBC6"
	case CartridgeTypeMBC7SensorRumbleRAMBattery:
		return "MBC7 + Sensor + Rumble + RAM + Battery"
	case CartridgeTypePocketCamera:
		return "Pocket Camera"
	case CartridgeTypeBandaiTAMA5:
		return "Bandai TAMA5"
	case CartridgeTypeHuC3:
		return "HuC3"
	case CartridgeTypeHuC1RAMBattery:
		return "HuC1 + RAM + Battery"
	default:
		return "Unknown Cartridge Type"
	}
}
