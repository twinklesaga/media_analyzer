package analyzer

import (
	"errors"
)

type analyseLv int

const(
	Lv1 analyseLv = 0		// 컨테이너만 분석
)


type Report struct {
	FileType 	string
	FileSize 	int64

	BitRate  	int
	SampleRate 	int

	Err 		error
}


type MediaAnalyser interface {
	Analyser( string, analyseLv) Report
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