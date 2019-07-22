package analyzer

import (
	"fmt"
	"github.com/mewkiz/flac"
	"io"
	"os"
	"runtime"
	"time"
)

type FLACAnalyzer struct {


}
func NewFALCAnalyzer() MediaAnalyzer{
	return &FLACAnalyzer{

	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (a *FLACAnalyzer)Analyze(filePath string , lv AnalyzeLV)(r Report) {
	var err error
	var f *os.File
	var fi os.FileInfo

	fi, err = os.Stat(filePath)
	if err != nil {
		r.Err = err
		return
	}
	r.FileSize = fi.Size()

	f, err = os.Open(filePath)
	if err != nil {
		r.Err = err
		return
	}
	defer f.Close()
	var flacStream *flac.Stream
	flacStream, err = flac.Parse(f)
	if err != nil {
		r.Err = err
		return
	}
	defer flacStream.Close()

	r.FileType = "FLAC"
	for {
		_, err := flacStream.ParseNext()
		if err != nil {
			if err != io.EOF {
				r.Err = err
			}
			break
		}
	}

	r.SampleRate = int(flacStream.Info.SampleRate)
	r.Channel = int(flacStream.Info.NChannels)
	r.SampleBit = int(flacStream.Info.BitsPerSample)
	r.Duration = time.Duration(int64(flacStream.Info.NSamples*1000)/int64(flacStream.Info.SampleRate)) * time.Millisecond
	r.BitRate = int(r.FileSize * 8 / int64(r.Duration.Seconds()) / 1000)
	return
}