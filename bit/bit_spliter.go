package bit

import (
	"errors"
)


func Split(data []byte ,ns ...int) ([]int64 , error){

	l := len(data) * 8

	readBit := 0
	for _,n :=range ns {
		if n > 64 {
			return nil , errors.New("out of range 64")
		}
		readBit += n
	}
	if l < readBit {
		return nil , errors.New("out of range data len")
	}

	pos := 0

	var err error

	result := make([]int64 , len(ns))

	for i , n := range ns {
		var r int64 = 0
		var cur int64 = 0
		var num = n
		for num > 0 {
			bytePos := 0
			if pos > 0 {
				bytePos =  pos / 8
			}
			bitPos := pos % 8
			//fmt.Println(bytePos , bitPos , data[bytePos])
			if bitPos + num >= 8 {
				mask , offset := posMask(bitPos)
				r = r<<offset | int64(data[bytePos]) & int64(mask)
				bitPos += int(offset)
				num-=int(offset)
			//	fmt.Println(mask , offset)
				pos += int(offset)
			}else {
				if data[bytePos]& (1<< uint8(7 - bitPos)) > 0 {
					cur = 1
				} else {
					cur = 0
				}
				pos++
				num--
				r = r<<1 | cur
			}
		}
		result[i] = r
	}
	return result , err
}


func posMask(bitPos int) (int , uint32) {
	switch bitPos {
	case 0: return  0xFF , 8
	case 1: return  0x7F , 7
	case 2: return  0x3F , 6
	case 3: return  0x1F , 5
	case 4: return  0x0F , 4
	case 5: return  0x07 , 3
	case 6: return  0x03 , 2
	case 7: return  0x01 , 1
	}
	panic("")
}