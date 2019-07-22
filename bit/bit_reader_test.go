package bit_test

import (
	"fmt"
	"github.com/twinklesaga/media_analyzer/bit"
	"testing"
)


type TestReader struct {

}

func (r *TestReader)Read(b []byte)  (n int, err error) {

	l := len(b)
	for i := 0 ; i < l; i++ {
		b[i] = 0xF0
	}
	return l , nil
}



func TestBitReader_Read32(t *testing.T) {

	br := bit.NewBitReader(&TestReader{})

	fmt.Printf("%b\n",br.ReadBit32(0,3))
	fmt.Printf("%b\n",br.ReadBit32(0,3))
	fmt.Printf("%b\n",br.ReadBit32(0,4))
}
