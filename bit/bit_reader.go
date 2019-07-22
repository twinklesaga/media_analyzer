package bit

import "io"

type BitReader struct {
	reader 	io.Reader
	buf 		[]byte
	data 		[]byte
	residual 	uint
}


func NewBitReader(reader io.Reader) *BitReader {

	return &BitReader{
		reader:reader,
		buf : make([]byte , 1),
		data : make([]byte , 0),
		residual: 0,
	}
}

func (r *BitReader)ReadBit32(prev int32, bit uint) (result int32){

	result = prev
	if r.residual >= bit {
		mask := GetBitMask(bit)

		shift := uint(r.residual - bit)
		val := (r.buf[0]  & (mask <<  shift)) >> shift
		result = int32(result << bit) | int32(val)

		r.residual -= bit
	}else {
		if r.residual > 0 {
			mask := GetBitMask(r.residual)
			result = int32(result << r.residual) | int32(r.buf[0] & mask)

			bit -= r.residual
			r.residual = 0

		}
		r.reader.Read(r.buf)
		r.residual = 8

		r.data = append(r.data , r.buf[0])
		result = r.ReadBit32(result , bit)
	}
	return
}

func (r *BitReader)GetResidual() uint {
	return r.residual
}


func (r *BitReader)ReadUnaryUnsigned() (val int) {
	val = 1
	for {
		bit := r.ReadBit32(0, 1)
		if bit == 1 {
			return
		}else{
			val++
		}
	}
}

func GetBitMask(bit uint) byte {

	switch bit {
	case 1: return 0x01
	case 2: return 0x03
	case 3: return 0x07
	case 4: return 0x0F
	case 5: return 0x1F
	case 6: return 0x3F
	case 7: return 0x7F
	case 8: return 0xFF
	}
	return 0
}