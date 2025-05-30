package memory

import (
	"fmt"
	"log/slog"
	"sort"
)

type mmioType int

const (
	MMIOTypeByte mmioType = iota
	MMIOTypeByteArray
)

type mmioMapping struct {
	address  uint16
	size     uint16
	readOnly bool

	mmioType mmioType

	data     []byte
	byteData *byte // For single byte MMIO mappings
}

type MMIO struct {
	mmios []mmioMapping
}

func (h *MMIO) mapMemory(addr uint16) uint16 {
	return addr
}

// findMMIO finds the MMIO device index that contains the given address
func (h *MMIO) findMMIOIndex(addr *uint16) (int, error) {
	*addr = h.mapMemory(*addr)

	for i, mapping := range h.mmios {
		if *addr >= mapping.address && *addr <= mapping.address+mapping.size-1 {
			// Found the MMIO mapping that contains the address
			return i, nil
		}
	}

	return 0, fmt.Errorf("MMIO address %04x not found", *addr)
}

func (h *MMIO) RemoveMMIO(address uint16, size uint16) error {
	for i, mapping := range h.mmios {
		if mapping.address == address && mapping.size == size {
			h.mmios = append(h.mmios[:i], h.mmios[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("MMIO address %04x with size %d not found", address, size)
}

func (h *MMIO) AddMMIO(data []byte, address uint16, size uint16, readOnly bool) {
	// Add the MMIO, but ensure that the entries are sorted by address.
	// This is required for the MMIO handler to work properly.

	mapping := mmioMapping{address, size, readOnly, MMIOTypeByteArray, data, nil}
	h.mmios = append(h.mmios, mapping)

	sort.Slice(h.mmios, func(i, j int) bool {
		return h.mmios[i].address < h.mmios[j].address
	})
}

func (h *MMIO) AddMMIOByte(data *byte, address uint16, readOnly bool) {
	// Add a single byte MMIO mapping.
	// This is useful for registers that are not larger than 1 byte.
	mapping := mmioMapping{address, 1, readOnly, MMIOTypeByte, []byte{}, data}
	h.mmios = append(h.mmios, mapping)

	sort.Slice(h.mmios, func(i, j int) bool {
		return h.mmios[i].address < h.mmios[j].address
	})
}

// Read8 reads a 8-bit value from the MMIO address space and returns it.
func (h *MMIO) Read8(addr uint16) (uint8, error) {
	index, err := h.findMMIOIndex(&addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %04x not found", addr)
	}
	if h.mmios[index].mmioType == MMIOTypeByte {
		return *h.mmios[index].byteData, nil
	} else if h.mmios[index].mmioType != MMIOTypeByteArray {
		return 0, fmt.Errorf("MMIO address %04x is not a byte array", addr)
	}
	return h.mmios[index].data[nonMapped], nil
}

// Write8 writes a 8-bit value to the MMIO address space.
func (h *MMIO) Write8(addr uint16, data uint8) error {
	index, err := h.findMMIOIndex(&addr)
	if err != nil {
		return err
	}
	slog.Debug("MMIO Write8", "addr", fmt.Sprintf("%04x", addr), "data", fmt.Sprintf("%02x", data))
	if h.mmios[index].readOnly {
		return nil
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %04x not found", addr)
	}
	if h.mmios[index].mmioType == MMIOTypeByte {
		*h.mmios[index].byteData = data
		return nil
	} else if h.mmios[index].mmioType != MMIOTypeByteArray {
		return fmt.Errorf("MMIO address %04x is not a byte array", addr)
	}
	h.mmios[index].data[nonMapped] = data
	return nil
}

// Read16 reads a 16-bit value from the MMIO address space and returns it.
func (h *MMIO) Read16(addr uint16) (uint16, error) {
	index, err := h.findMMIOIndex(&addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %04x not found", addr)
	}
	if h.mmios[index].mmioType != MMIOTypeByteArray {
		return 0, fmt.Errorf("MMIO address %04x is not a byte array", addr)
	}
	dataBytes := h.mmios[index].data[nonMapped : nonMapped+2]
	return uint16(dataBytes[0]) | uint16(dataBytes[1])<<8, nil
}

// Write16 writes a 16-bit value to the MMIO address space.
func (h *MMIO) Write16(addr uint16, data uint16) error {
	index, err := h.findMMIOIndex(&addr)
	if err != nil {
		return err
	}
	if h.mmios[index].readOnly {
		return nil
	}
	slog.Debug("MMIO Write16", "addr", fmt.Sprintf("%04x", addr), "data", fmt.Sprintf("%04x", data))
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %04x not found", addr)
	}
	if h.mmios[index].mmioType != MMIOTypeByteArray {
		return fmt.Errorf("MMIO address %04x is not a byte array", addr)
	}
	h.mmios[index].data[nonMapped] = byte(data)
	h.mmios[index].data[nonMapped+1] = byte(data >> 8)
	return nil
}
