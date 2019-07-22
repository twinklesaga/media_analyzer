package analyzer

import (
	"errors"
	"fmt"
	"time"
)

type AnalyzeLV int

const(
	Lv1 AnalyzeLV = 0 // 컨테이너만 분석
)


type Report struct {
	FileType 	string
	FileSize 	int64

	SampleRate 	int
	BitRate  	int

	Channel 	int
	SampleBit  	int

	Duration    time.Duration

	Err 		error
	SubErr		error
}

const ReportFormat = "FileType : %s\n" +
					 "FileSize : %d\n" +
					 "SampleRate : %d\n" +
					 "BitRate : %d\n" +
					 "Channel : %d\n" +
					 "SampleBit : %d\n" +
	                 "Duration : %v\n"


var (
	NotSupportFileFormat 	= errors.New("not support file format")
	NotSupportData			= errors.New("not support data")
	MismatchContainerFormat = errors.New("mismatch container format")
)


func (r Report)String() string{
	return fmt.Sprintf(ReportFormat , r.FileType , r.FileSize , r.SampleRate, r.BitRate , r.Channel , r.SampleBit , r.Duration)
}

type MediaAnalyzer interface {
	Analyze( string, AnalyzeLV) Report
}

func SliceToUint32(s []byte) (uint32 , error ) {
	if len(s) >=4 {
		return  uint32(s[0]) << 24 | (uint32(s[1]) << 16) | (uint32(s[2]) << 8) | uint32(s[3]) , nil
	}
	return 0, errors.New("slice size small")
}


func SliceToInt32BigEndian(s []byte) (int , error ) {
	if len(s) >= 4 {
		return int(s[3])<<24 | (int(s[2]) << 16) | (int(s[1]) << 8) | int(s[0]), nil
	}
	return 0, errors.New("slice size small")
}