package reader

import (
	"errors"
	"io"
)


type ReadMode int

const (
	IdleMode 		ReadMode = iota
	BitReadMode
	ByteReadMode
)

type DetailReader struct {
	reader 		io.Reader
	mode 		ReadMode

	buf 		[]byte
	pos 		int32
	len 		int32
}

func NewDetailReader(reader io.Reader) *DetailReader {
	return &DetailReader{
		reader:reader,

		mode:IdleMode,
		buf:make([]byte , 64),
		pos:0,
		len:0,
	}
}

func BitPosMask(bitPos int32) (int32 , uint32) {
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

func (b *DetailReader)ReadBit32(n int32) (int32 , error){
	if n >0 && n <= 32 {
		b.mode = BitReadMode

		var result int32 = 0
		var cur int32 = 0
		for n > 0 {
			if b.len == 0 {
				readByte := (n / 8) + 1
				io.ReadFull(b.reader, b.buf[0:readByte])
				b.len = int32(readByte * 8)
				b.pos = 0
			}

			bytePos := 8 / b.pos
			bitPos := 8 % b.pos

			if bitPos + int32(n) >= 8 {
				mask , offset := BitPosMask(bitPos)
				result = result<<offset | int32(b.buf[bytePos]) & mask
				bitPos += int32(offset)
				n-=int32(offset)
			}else {
				if b.buf[bytePos]& (1<< uint8(7 - bitPos)) > 0 {
					cur = 1
				} else {
					cur = 0
				}
				b.pos++
				n--
				result = result<<1 | cur
			}
		}
		if b.len == b.pos {
			b.len = 0
			b.pos = 0
			b.mode = IdleMode
		}
		return result , nil
	}
	return 0 , errors.New("out of range 32")
}

func (b *DetailReader)ReadBit64(n int32) (int64 , error){
	if n >0 && n <= 64 {
		b.mode = BitReadMode
		var result int64 = 0
		var cur int64 = 0
		for n > 0 {
			if b.len == 0 {
				readByte := (n / 8) + 1
				io.ReadFull(b.reader, b.buf[0:readByte])
				b.len = int32(readByte * 8)
				b.pos = 0
			}

			bytePos := 8 / b.pos
			bitPos := 8 % b.pos

			if bitPos + int32(n) >= 8 {
				mask , offset := BitPosMask(bitPos)
				result = result<<offset | int64(b.buf[bytePos]) & int64(mask)
				bitPos += int32(offset)
				n-=int32(offset)
			}else {
				if b.buf[bytePos]& (1<< uint8(7 - bitPos)) > 0 {
					cur = 1
				} else {
					cur = 0
				}
				b.pos++
				n--
				result = result<<1 | cur
			}
		}
		if b.len == b.pos {
			b.len = 0
			b.pos = 0
			b.mode = IdleMode
		}
		return result , nil
	}

	return 0 , errors.New("out of range 32")
}


func (b *DetailReader)ReadBitMultiple32(ns ...int32) ([]int32 , error){

	result := make([]int32, len(ns))
	if b.mode == IdleMode{
		var sum int32 = 0
		for _, n := range ns {
			sum += int32(n)
		}

		readByte := sum / 8
		if sum % 8 > 0{
			readByte ++
		}
		if readByte <= 64 {
			io.ReadFull(b.reader, b.buf[0:readByte])
		}
		b.mode = BitReadMode
		b.pos = 0
		b.len = sum
	}

	for i, n := range ns {
		v , err := b.ReadBit32(n)
		if err != nil {
			return nil , err
		}

		result[i] = v
	}
	return result , nil
}

func (b *DetailReader)ReadBitMultiple64(ns ...int32) ([]int64 , error){

	result := make([]int64, len(ns))
	if b.mode == IdleMode{
		var sum int32 = 0
		for _, n := range ns {
			sum += int32(n)
		}

		readByte := sum / 8
		if sum % 8 > 0{
			readByte ++
		}
		if readByte <= 64 {
			io.ReadFull(b.reader, b.buf[0:readByte])
		}
		b.mode = BitReadMode
		b.pos = 0
		b.len = sum
	}

	for i, n := range ns {
		v , err := b.ReadBit64(n)
		if err != nil {
			return nil , err
		}

		result[i] = v
	}
	return result , nil
}

func (d *DetailReader)Read(p []byte) (int , error){
	return d.reader.Read(p)
}