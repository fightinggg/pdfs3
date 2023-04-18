package pdfs_blocks

import (
	"encoding/binary"
	"errors"
	"go/types"
	"pdfs3/pdfs/pdfs_err"
)

type ArrayBaseBlock struct {
	head  map[int]int
	datas [][]byte
}

func ReadAll(block []byte) (*ArrayBaseBlock, error) {
	if block[0] != 0 {
		return nil, types.Error{Msg: "Valid Block Head Byte"}
	}

	b := &ArrayBaseBlock{
		head:  map[int]int{},
		datas: [][]byte{},
	}
	// array base block begin at a block head just like:
	// [index1,index1MapTo,index2,index2MapTo,index3,index3MapTo,..., end]
	// in this format , every number is 32bit, only the end can equal to 0
	// "indexiMapTo" means which position the file with index "readIndex" is stored.
	// so the first step is reading this map
	readIndex := 1
	for ; block[readIndex] != 0 && readIndex < len(block); readIndex += 2 * 4 {
		key := binary.LittleEndian.Uint32(block[readIndex : readIndex+4])
		value := binary.LittleEndian.Uint32(block[readIndex+4 : readIndex+8])
		b.head[int(key)-1] = int(value)
	}

	// skip this bytes, now block[readIndex] = 0
	readIndex++

	// the second area in array base block just like:
	// [len1,data1,len2,data2,len3,data3,...]
	// every leni is 32bit, every datai's len equal leni
	for fileI := 0; readIndex < len(block); fileI++ {
		lenI := int(binary.LittleEndian.Uint32(block[readIndex : readIndex+4]))
		b.datas = append(b.datas, block[readIndex+4:readIndex+4+lenI])
		readIndex += 4 + lenI
	}

	return b, nil
}

func NewArrayBaseBlock(block []byte) *ArrayBaseBlock {
	all, _ := ReadAll([]byte{0, 0})
	return all
}

func (b *ArrayBaseBlock) WriteAll() ([]byte, error) {
	block := []byte{0}

	// array base block begin at a block head just like:
	// [index1,index1MapTo,index2,index2MapTo,index3,index3MapTo,..., end]
	// in this format , every number is 32bit, only the end can equal to 0
	// "indexiMapTo" means which position the file with index "writeIndex" is stored.
	// so the first step is reading this map
	for k, v := range b.head {
		if k < 0 {
			return nil, errors.New("ERROR")
		}

		block = binary.LittleEndian.AppendUint32(block, uint32(k)+1)
		block = binary.LittleEndian.AppendUint32(block, uint32(v))
	}

	block = append(block, 0)

	// the second area in array base block just like:
	// [len1,data1,len2,data2,len3,data3,...]
	// every leni is 32bit, every datai's len equal leni
	for _, data := range b.datas {
		block = binary.LittleEndian.AppendUint32(block, uint32(len(data)))
		block = append(block, data...)
	}

	return block, nil

}

func (b *ArrayBaseBlock) Read(index int) ([]byte, error) {
	blockIndex, ok := b.head[index]
	if !ok || blockIndex > len(b.datas) {
		return nil, pdfs_err.PdfsNotFoundError()
	}
	return b.datas[blockIndex], nil
}

func (b *ArrayBaseBlock) Write(index int, data []byte) (int, int) {
	if index == -1 {
		index = len(b.head)
		b.head[index] = len(b.datas)
		b.datas = append(b.datas, data)
	} else {
		blockIndex, ok := b.head[index]
		if !ok || blockIndex > len(b.datas) {
			panic("E")
		}
		b.datas[blockIndex] = data
	}

	byteSizeSum := 1                   // block start head
	byteSizeSum += len(b.head) * 2 * 4 // head size
	byteSizeSum += 1                   // head end
	for _, item := range b.datas {
		byteSizeSum += len(item) // data size
	}
	return index, byteSizeSum

}
