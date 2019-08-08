package analyzer_test

import (
	"fmt"
	"github.com/twinklesaga/media_analyzer/analyzer"
	"testing"
)

func TestWAVAnalyzer(t *testing.T) {

	files := []string{
		"/Volumes/Work/mcp/400172377.wav",
		"/Volumes/Work/mcp/31513321.wav",
	}

	for _,f :=range files {

		fmt.Println(f)
		a := analyzer.NewWavAnalyzer()

		r:= a.Analyze(f, analyzer.Lv1)

		if r.Err != nil {
			fmt.Println(r.Err)
		} else {
			fmt.Print(r)
		}
		fmt.Println("-------------------")
	}
}

func TestMP3Analyzer(t *testing.T) {
	files := []string{

		"/Volumes/Work/mcp/427244641/427245163_5d4b876f.mp3",
		"/Volumes/Work/mcp/427244641/427244641_5d4b6503.mp3",
		"/Volumes/Work/mcp/427205215/320k.mp3",
		"/Volumes/Work/mcp/427205206/320k.mp3",
		"/Users/1100117/Downloads/NewArea - Space Bound/01 SPACE BOUND.mp3",
		"/Users/1100117/Downloads/427204637_5d3fe218.mp3",
	}
	/*	"/Volumes/Work/mcp/err/416/405328928.mp3",
		"/Volumes/Work/mcp/err/416/418933583.mp3",
		"/Volumes/Work/mcp/err/416/419739376.mp3",
		"/Volumes/Work/mcp/err/416/419934920.mp3",
		"/Volumes/Work/mcp/err/416/420075481.mp3",

		"/Volumes/Work/mcp/11116.mp3",
		"/Volumes/Work/mcp/31513321.mp3",
		"/Volumes/Work/mcp/400172377.mp3",
		"/Volumes/Work/mcp/err/22401046.mp3",
	}*/

	for _,f :=range files {

		fmt.Println(f)
		a := analyzer.NewMP3Analyzer()

		r:= a.Analyze(f, analyzer.Lv1)

		if r.Err != nil {
			fmt.Println(r.Err)
		} else {
			fmt.Print(r)
		}
		fmt.Println("-------------------")
	}
}


func TestFLACAnalyzer(t *testing.T) {
	files := []string{

		"/Volumes/Work/_mcp_work/423734259/423734259_5ce2397b.flac",
		"/Volumes/Work/mcp/test.flac",
		"/Volumes/Work/mcp/0000985746.flac",

	}

	for _,f :=range files {

		fmt.Println(f)
		a := analyzer.NewFALCAnalyzer()

		r:= a.Analyze(f, analyzer.Lv1)

		if r.Err != nil {
			fmt.Println(r.Err)
		} else {
			fmt.Print(r)
		}
		fmt.Println("-------------------")
	}
}

func TestM4AAnalyzer(t *testing.T) {
	files := []string{
		"/Volumes/Work/mcp/0514/420329934_128k.aac",
		"/Volumes/Work/mcp/420329924/128k.aac",
		"/Volumes/Work/mcp/420329924/256k.aac",
	}

	for _,f :=range files {

		fmt.Println(f)
		a := analyzer.NewM4AAnalyzer()

		r:= a.Analyze(f, analyzer.Lv1)

		if r.Err != nil {
			fmt.Println(r.Err)
		} else {
			fmt.Print(r)
		}
		fmt.Println("-------------------")
	}
}