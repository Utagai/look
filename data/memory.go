package data

// MemoryData is a data that lives entirely in memory.
type MemoryData struct {
	data []Datum
}

func NewMemoryData(data []Datum) *MemoryData {
	return &MemoryData{
		data: data,
	}
}

var _ Data = (*MemoryData)(nil)

func (md *MemoryData) Find(q string) Data {
	return md
}

func (md *MemoryData) At(index int) (Datum, bool) {
	if index >= len(md.data) || index < 0 {
		return nil, false
	}

	return md.data[index], true
}

func (md *MemoryData) Length() int {
	return len(md.data)
}
