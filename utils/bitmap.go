package utils

import "encoding/binary"

type BitMap struct {
	data []uint64
}

func NewBitMap(len int64) BitMap {
	bitMap := BitMap{
		data: make([]uint64, (len+63)/64),
	}
	bitMap.Clear()
	return bitMap
}

func (b BitMap) Clear() {
	for i := range b.data {
		b.data[i] = 0
	}
}

func (b BitMap) Set(pos int64, value bool) {
	index := pos / 64
	offset := pos % 64
	if value {
		b.data[index] |= (uint64(1)) << offset
	} else {
		b.data[index] &= ^((uint64(1)) << offset)
	}
}

func (b BitMap) At(pos uint64) bool {
	index := pos / 64
	offset := pos % 64
	return b.data[index]&((uint64(1))<<offset) != 0
}

func (b BitMap) ToBytes() []byte {
	var res []byte
	for _, v := range b.data {
		res = binary.LittleEndian.AppendUint64(res, v)
	}
	return res
}

func (b BitMap) FromBytes(data []byte) {
	if len(data)%8 != 0 {
		panic("")
	}
	b.data = make([]uint64, len(data)/8)
	for i := 0; i < len(data)/8; i++ {
		b.data[i] = binary.LittleEndian.Uint64(data[i*8 : i*8+64])
	}
}
