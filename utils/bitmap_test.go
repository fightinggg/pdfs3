package utils

import (
	"math/rand"
	"testing"
)

func TestName(t *testing.T) {
	for i := 0; i < 1000; i++ {
		testLen := rand.Int31() % 256
		bm := NewBitMap(int64(testLen))
		compate := make([]bool, testLen)

		for j := 0; j < int(10*testLen); j++ {
			pos := rand.Int31() % testLen
			value := rand.Int31()%2 == 0
			bm.Set(int64(pos), value)
			compate[pos] = value
		}

		for j := 0; j < int(testLen); j++ {
			if bm.At(uint64(j)) != compate[j] {
				panic("")
			}
		}
		//
		//for j := 0; j < int(testLen); j++ {
		//	if bm.At(uint64(j)) {
		//		fmt.Printf("*")
		//	} else {
		//		fmt.Printf(" ")
		//	}
		//
		//}
		//fmt.Println()
		//
		//for j := 0; j < int(testLen); j++ {
		//	if compate[j] {
		//		fmt.Printf("*")
		//	} else {
		//		fmt.Printf(" ")
		//	}
		//}
		//fmt.Println()
	}
}
