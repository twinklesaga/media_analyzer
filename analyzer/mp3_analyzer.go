package analyzer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	ID3 				uint32 = 0x49443300
	ID3Mask 			uint32 = 0xFFFFFF00

	AAU 				uint32 = 0xFFE00000
	AAUMask 			uint32 = 0xFFE00000

	TAG 				uint32 = 0x54414700
	TAGMask 			uint32 = 0xFFFFFF00
)


const (
	AAUMpegAudioVersionMask 	uint32 = 0x00180000
	AAUMpegAudioVersionOffset 	uint32 = 19

	AAUMpegLayerMask 			uint32 = 0x00060000
	AAUMpegLayerOffset 			uint32 = 17

	AAUChecksumMask				uint32 = 0x00010000
	AAUChecksumOffset			uint32 = 16

	AAUBitRateMask	 	uint32 = 0x0000F000
	AAUBitRateOffset 	uint32 = 12

	AAUFrequencyMask 	uint32 = 0x00000C00
	AAUFrequencyOffset 	uint32 = 10

	AAUPaddingBitMask 	uint32 = 0x00000200

	AAUChannelMask		uint32 = 0x000000C0
	AAUChannelOffset	uint32 = 6

)

const(
	Layer1 int = 0
	Layer2 int = 1
	Layer3 int = 2
	LayerReserved int = 3

	MPEG1 int = 0
	MPEG2 int = 1
	MPEG25 int = 2
	MPEGReserved int = 3
)

type MP3Context struct {
	err 		error
	fileSize 	int64

	mpeg 		int
	layer 		int

	tag 		bool
	aauCount 	int
	totalAAUSize int64

	crc 		bool

	Frequency 	int
	BitRate 	int

}

func (c *MP3Context)IsMP3() bool{
	if c.mpeg == MPEG1 && c.layer == Layer3 {
		return true
	}
	return false
}

type MP3Analyzer struct {
	BitRateTable [][]int
	FrequencyTable [][]int

}

func NewMP3Analyzer() MediaAnalyser{
	mp3 := new(MP3Analyzer)

	mp3.BitRateTable = [][]int{
		{0, 32000, 64000, 96000, 128000, 160000, 192000, 224000,256000, 288000, 320000, 352000, 384000, 416000, 448000 , 0},
		{0, 32000, 48000, 56000, 64000, 80000, 96000, 112000, 128000, 160000, 192000, 224000, 256000, 320000, 384000 , 0},
		{0, 32000, 40000, 48000, 56000, 64000, 80000, 96000,112000, 128000, 160000, 192000, 224000, 256000, 320000 ,0},
	}

	mp3.FrequencyTable = [][]int{
		{44100,	48000,	32000 , 0},
		{22050,	24000,	16000 , 0},
		{11025,	12000,	8000 , 0},
	}

	return mp3
}

func (a *MP3Analyzer)Analyser(filePath string , lv analyseLv) Report{

	ctx := MP3Context{}
	var fi os.FileInfo
	var f  *os.File

	fi , ctx.err = os.Stat(filePath)
	if ctx.err == nil {
		ctx.fileSize = fi.Size()
	}
	f , ctx.err = os.Open(filePath)
	if  ctx.err == nil {
		defer f.Close()

		var pos int64 = 0
		eof := false
		for {
			header := make([]byte, 0x4)

			_ , err := f.Read(header)
			if err != nil {
				if err == io.EOF {
					eof = true
				}
				break
			}
			pos += 4

			code, err := SliceToUint32(header)

			if (code & ID3Mask) == ID3 {
				id3Header := make([]byte, 10)
				copy(id3Header, header)
				f.Read(id3Header[4:])
				pos += 6
				//fmt.Println(id3Header)

				size, err := SliceToUint32(id3Header[6:])
				if err != nil {

				}

				unSyncSize := Unsynchsafe(size)
				//	fmt.Println(size , unSyncSize)
				id3Data := make([]byte , unSyncSize)

				pos += int64(unSyncSize)
				f.Read(id3Data)
			} else if (code & AAUMask) == AAU {

				mpegVersion := (code & AAUMpegAudioVersionMask) >> AAUMpegAudioVersionOffset
				switch mpegVersion {
					case 3:ctx.mpeg = MPEG1
					case 2:ctx.mpeg = MPEG2
					case 1:ctx.mpeg = MPEGReserved
					case 0:ctx.mpeg = MPEG25
				}
				mpegLayer := (code & AAUMpegLayerMask) >> AAUMpegLayerOffset
				switch mpegLayer {
					case 3:ctx.layer = Layer1
					case 2:ctx.layer = Layer2
					case 1:ctx.layer = Layer3
					case 0:ctx.layer = LayerReserved
				}

				checksum := (code & AAUChecksumMask) >> AAUChecksumOffset
				if checksum > 0 {
					ctx.crc = true
				}

				if !ctx.IsMP3() {
					fmt.Println("AAU is not mpeg1 layer3" , ctx.aauCount)
				}


				bitRateIndex := (code & AAUBitRateMask) >> AAUBitRateOffset
				bitRate := a.BitRateTable[ctx.layer][bitRateIndex]

				frequencyIndex := (code & AAUFrequencyMask) >> AAUFrequencyOffset
				frequency := a.FrequencyTable[ctx.mpeg][frequencyIndex]

				if ctx.aauCount == 0 {
					ctx.Frequency = frequency
					ctx.BitRate = bitRate
				}else if ctx.Frequency != frequency|| ctx.BitRate != bitRate {
					ctx.err = errors.New("mismatch aau header")
				}

				paddingByte := 0
				if code & AAUPaddingBitMask > 0 {
					paddingByte = 1
				}
				AAUSize := 144 * ctx.BitRate/ctx.Frequency + paddingByte - 4

				buf := make([]byte, AAUSize)
				n, err := f.Read(buf)
				if err != nil || n != AAUSize{
					fmt.Println("AAU Read fail" , err , n , AAUSize , hex.Dump(header))
					fmt.Println(hex.Dump(buf))
				}
				pos += int64(n)

				ctx.aauCount++

				if err != nil {
					fmt.Println(err)
				}
				ctx.totalAAUSize  += int64(AAUSize + 4)

			} else if (code & TAGMask) == TAG {
				f.Seek(124 , 1)
				pos += 124
				fmt.Println("TAG")

			} else {
				fmt.Println(hex.Dump(header))
				ctx.err = errors.New(fmt.Sprintf("unknown header %v %d" , header , pos ))
				break
			}
		}

		fmt.Println(ctx.fileSize , pos)
		if eof {
			return Report{
				FileType:"MP3",
				FileSize:ctx.fileSize,
				BitRate:ctx.BitRate,
				SampleRate:ctx.Frequency,
				Err:ctx.err,
			}
		}
	}

	return Report{Err:ctx.err}

}

func Unsynchsafe(in uint32) uint32{
	var out uint32  = 0
	var mask uint32 = 0x7F000000

	for mask > 0 {
		out >>= 1
		out |= in & mask
		mask >>= 8
	}

	return out
}
