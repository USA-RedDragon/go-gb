package ppu

type FIFO struct {
	size int
	data []byte
}

func NewFIFO() *FIFO {
	fifo := &FIFO{}
	fifo.Reset()
	return fifo
}

func (f *FIFO) Reset() {
	f.size = 0
	f.data = []byte{}
}

func (f *FIFO) Push(value byte) {
	// slog.Warn("Pushing to FIFO", "size", f.size, "value", value)
	if f.size < len(f.data) {
		f.data[f.size] = value
	} else {
		f.data = append(f.data, value)
	}
	f.size++
}

func (f *FIFO) Pop() (byte, bool) {
	if f.size == 0 {
		return 0, false // FIFO is empty
	}
	value := f.data[0]
	f.data = f.data[1:]
	f.size--
	// slog.Warn("Popping from FIFO", "size", f.size, "value", value)
	return value, true
}

func (f *FIFO) Size() int {
	return f.size
}
