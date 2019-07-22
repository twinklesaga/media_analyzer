package analyzer

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

type WAVAnalyzer struct {

}

const (
	RIFF uint32 = 0x52494646
	FMT  uint32 = 0x666d7420
	LIST uint32 = 0x4C495354
	DATA uint32 = 0x64617461
	FACT uint32 = 0x66616374

	WAVE uint32 = 0x57415645
)

type AudioFormat uint16
const(
	WAVE_FORMAT_UNKNOWN 		AudioFormat = 0X0000
	WAVE_FORMAT_PCM 			AudioFormat = 0X0001
	WAVE_FORMAT_MS_ADPCM 		AudioFormat = 0X0002
	WAVE_FORMAT_IEEE_FLOAT 		AudioFormat = 0X0003
	WAVE_FORMAT_ALAW 			AudioFormat = 0X0006
	WAVE_FORMAT_MULAW 			AudioFormat = 0X0007
	WAVE_FORMAT_IMA_ADPCM 		AudioFormat = 0X0011
	WAVE_FORMAT_YAMAHA_ADPCM 	AudioFormat = 0X0016
	WAVE_FORMAT_GSM 			AudioFormat = 0X0031
	WAVE_FORMAT_ITU_ADPCM 		AudioFormat = 0X0040
	WAVE_FORMAT_MPEG 			AudioFormat = 0X0050
	WAVE_FORMAT_EXTENSIBLE 		AudioFormat = 0XFFFE
)

func (f AudioFormat)String() string {
	switch f {
	case WAVE_FORMAT_UNKNOWN: 		return "UNKNOWN"
	case WAVE_FORMAT_PCM: 			return "PCM"
	case WAVE_FORMAT_MS_ADPCM: 		return "MS_ADPCM"
	case WAVE_FORMAT_IEEE_FLOAT:	return "IEE_FLOAT"
	case WAVE_FORMAT_ALAW: 			return "ALAW"
	case WAVE_FORMAT_MULAW: 		return "MULAW"
	case WAVE_FORMAT_IMA_ADPCM: 	return "IMA_ADPCM"
	case WAVE_FORMAT_YAMAHA_ADPCM:	return "YAMAHA_ADPCM"
	case WAVE_FORMAT_GSM: 			return "GSM"
	case WAVE_FORMAT_ITU_ADPCM:		return "ITU_ADPCM"
	case WAVE_FORMAT_MPEG:			return "MPEG"
	case WAVE_FORMAT_EXTENSIBLE:	return "EXTENSIBLE"
	default:
		return "ERROR"
	}
}

func NewWavAnalyzer()MediaAnalyzer  {
	return &WAVAnalyzer{}
}


func (a *WAVAnalyzer)Analyze(filePath string , lv AnalyzeLV) Report {
	var err error
	var f *os.File

	f, err = os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return Report{}
	}
	defer f.Close()

	chunkHeaderSize := 8
	chunkHeaderBuf := make([]byte , chunkHeaderSize)
	var n int
	var n64 int64
	var report Report
	var wavInfo WAVInfo
	var FileSize int64 = 0
	for err == nil {
		n, err = f.Read(chunkHeaderBuf)
		FileSize += int64(n)

		if err == io.EOF {
			if n == 0 {
				err = nil
			}
			break
		}

		if n == chunkHeaderSize && err == nil {
			var chunkId uint32
			var size int
			chunkId, err = SliceToUint32(chunkHeaderBuf)
			size, err = SliceToInt32BigEndian(chunkHeaderBuf[4:])

			if chunkId == RIFF {
				Riff := make([]byte, 4)
				n, err = f.Read(Riff)
				FileSize += int64(n)
				if n == 4 {
					var wave uint32
					wave ,err = SliceToUint32(Riff)
					if wave == WAVE {
						report.FileType = "WAVE"
					}else {
						fmt.Println(hex.Dump(Riff))
					}
				}
			} else if chunkId == FMT {
				fmtChunk := make([]byte, size)
				n, err = f.Read(fmtChunk)
				FileSize += int64(n)

				if err == nil && size == n {
					wavInfo = GetWaveInfo(fmtChunk)
					if wavInfo.Format == WAVE_FORMAT_PCM {
						report.SampleRate = wavInfo.SampleRate
						report.BitRate = int(wavInfo.BytePerSec)

					}else{
						fmt.Println(wavInfo)
						err = NotSupportData
						break
					}
				} else {
					//err = MismatchContainerFormat
					//break
				}
			}else if chunkId == LIST{
				n64 , err = f.Seek(int64(size), 1)
				if err != nil {

					fmt.Printf("LIST error : %d , %d : %v\n" , size , n64 , err)
					err = MismatchContainerFormat
					report.SubErr = err
					break
				}
				FileSize += int64(size)
			}else if chunkId == FACT{
				n64 , err = f.Seek(int64(size), 1)
				if err != nil {
					err = MismatchContainerFormat
					report.SubErr = err
					fmt.Printf("FACT error : %d , %d : %v" , size , n64 , err)
					break
				}
				FileSize += int64(size)
			} else if chunkId == DATA {
				wavInfo.DataBlockSize = uint32(size)
				if wavInfo.Format == WAVE_FORMAT_PCM {
					wavInfo.NumSamples = int(size / int(wavInfo.BytePerSample))
					wavInfo.Duration = time.Duration(float64(wavInfo.NumSamples)*1000/float64(wavInfo.SampleRate)) * time.Millisecond
				}



				n64 , err = f.Seek(int64(size), 1)
				if err == nil {
					FileSize += int64(size)
				}
			} else {
				fmt.Println("Unknown chunk " , hex.Dump(chunkHeaderBuf))
				break
			}

		} else {

		}
	}
	//fmt.Println(wavInfo)
	if err != nil {
		report.Err = err
	}
	report.FileSize = FileSize
	report.Channel = int(wavInfo.Channel)
	report.Duration = wavInfo.Duration
	report.SampleBit = wavInfo.SampleBit

	return report

}

func GetWaveInfo(fmt []byte) WAVInfo{
	var info WAVInfo
	info.ParseFMT(fmt)

	info.SampleBit = int(info.ByteAlign / info.Channel * 8)
	return info
}

type WAVInfo struct {
	Format 			AudioFormat
	Channel 		uint16
	SampleRate 		int
	SampleBit		int
	BytePerSec  	uint32
	ByteAlign   	uint16
	BytePerSample 	uint16
	SamplePerBit 	uint16

	DataBlockSize 	uint32
	NumSamples 		int
	Duration 		time.Duration
}


func (w *WAVInfo)ParseFMT(fmt []byte){
	w.Format 	 	= AudioFormat(fmt[1] << 8 | fmt[0])
	w.Channel 	 	= uint16(fmt[3] << 8 | fmt[2])
	w.SampleRate 	= int(fmt[7]) << 24 | int(fmt[6]) << 16 | int(fmt[5]) << 8 | int(fmt[4])
	w.BytePerSec 	= uint32(fmt[11]) << 24 | uint32(fmt[10]) << 16 | uint32(fmt[9]) << 8 | uint32(fmt[8])
	w.ByteAlign  	= uint16(fmt[13]) << 8 | uint16(fmt[12])
	w.BytePerSample  = uint16(fmt[15]) << 8 | uint16(fmt[14])

	if len(fmt) > 16 {

	}
}


const PrintFmt string = "Format : %s(%d)\n" +
						"Channel : %d\n" +
						"SampleRate : %d\n" +
						"BytePerSec : %d\n" +
						"ByteAlign : %d\n" +
						"BytePerSample :%d\n" +
						"DataBlockSize : %d\n" +
						"NumSamples: %d\n" +
						"Duration : %v\n"


func (w WAVInfo)String()string{
	return fmt.Sprintf(
		PrintFmt,
	 	w.Format , w.Format,
		w.Channel ,
		w.SampleRate,
		w.BytePerSec,
		w.ByteAlign,
		w.BytePerSec,
		w.DataBlockSize,
		w.NumSamples,
		w.Duration)
}

