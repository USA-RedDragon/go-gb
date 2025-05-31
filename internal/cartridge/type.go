package cartridge

type Type byte

const (
	TypeROMOnly                    Type = 0x00
	TypeMBC1                       Type = 0x01
	TypeMBC1RAM                    Type = 0x02
	TypeMBC1RAMBattery             Type = 0x03
	TypeMBC2                       Type = 0x05
	TypeMBC2Battery                Type = 0x06
	TypeROMRAM                     Type = 0x08
	TypeROMRAMBattery              Type = 0x09
	TypeMMM01                      Type = 0x0B
	TypeMMM01RAM                   Type = 0x0C
	TypeMMM01RAMBattery            Type = 0x0D
	TypeMBC3TimerRAM               Type = 0x0F
	TypeMBC3TimerRAMBattery        Type = 0x10
	TypeMBC3                       Type = 0x11
	TypeMBC3RAM                    Type = 0x12
	TypeMBC3RAMBattery             Type = 0x13
	TypeMBC5                       Type = 0x19
	TypeMBC5RAM                    Type = 0x1A
	TypeMBC5RAMBattery             Type = 0x1B
	TypeMBC5Rumble                 Type = 0x1C
	TypeMBC5RumbleRAM              Type = 0x1D
	TypeMBC5RumbleRAMBattery       Type = 0x1E
	TypeMBC6                       Type = 0x20
	TypeMBC7SensorRumbleRAMBattery Type = 0x22
	TypePocketCamera               Type = 0xFC
	TypeBandaiTAMA5                Type = 0xFD
	TypeHuC3                       Type = 0xFE
	TypeHuC1RAMBattery             Type = 0xFF
)

func (c Type) String() string {
	switch c {
	case TypeROMOnly:
		return "ROM Only"
	case TypeMBC1:
		return "MBC1"
	case TypeMBC1RAM:
		return "MBC1 + RAM"
	case TypeMBC1RAMBattery:
		return "MBC1 + RAM + Battery"
	case TypeMBC2:
		return "MBC2"
	case TypeMBC2Battery:
		return "MBC2 + Battery"
	case TypeROMRAM:
		return "ROM + RAM"
	case TypeROMRAMBattery:
		return "ROM + RAM + Battery"
	case TypeMMM01:
		return "MMM01"
	case TypeMMM01RAM:
		return "MMM01 + RAM"
	case TypeMMM01RAMBattery:
		return "MMM01 + RAM + Battery"
	case TypeMBC3TimerRAM:
		return "MBC3 + Timer + RAM"
	case TypeMBC3TimerRAMBattery:
		return "MBC3 + Timer + RAM + Battery"
	case TypeMBC3:
		return "MBC3"
	case TypeMBC3RAM:
		return "MBC3 + RAM"
	case TypeMBC3RAMBattery:
		return "MBC3 + RAM + Battery"
	case TypeMBC5:
		return "MBC5"
	case TypeMBC5RAM:
		return "MBC5 + RAM"
	case TypeMBC5RAMBattery:
		return "MBC5 + RAM + Battery"
	case TypeMBC5Rumble:
		return "MBC5 + Rumble"
	case TypeMBC5RumbleRAM:
		return "MBC5 + Rumble + RAM"
	case TypeMBC5RumbleRAMBattery:
		return "MBC5 + Rumble + RAM + Battery"
	case TypeMBC6:
		return "MBC6"
	case TypeMBC7SensorRumbleRAMBattery:
		return "MBC7 + Sensor + Rumble + RAM + Battery"
	case TypePocketCamera:
		return "Pocket Camera"
	case TypeBandaiTAMA5:
		return "Bandai TAMA5"
	case TypeHuC3:
		return "HuC3"
	case TypeHuC1RAMBattery:
		return "HuC1 + RAM + Battery"
	default:
		return "Unknown Cartridge Type"
	}
}
