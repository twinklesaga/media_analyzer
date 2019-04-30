package analyzer

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type FLACAnalyzer struct {
	metaBlocks []*MetadataBlock
	frames []*Frame
}


func NewFALCAnalyzer() MediaAnalyser {

	return &FLACAnalyzer{
		metaBlocks:make([]*MetadataBlock , 0),
	}
}

type blockType byte
const (
	STREAMINFO 		blockType = 0
	PADDING 		blockType = 1
	APPLICATION		blockType = 2
	SEEKTABLE		blockType = 3
	VORBIS_COMMENT 	blockType = 4
	CUESHEET		blockType = 5
	PICTURE			blockType = 6
	RESERVED		blockType = 126
	INVALID 		blockType = 127
)
type MetadataBlock struct {
	bType 	blockType
	size 	int32


	data []byte
}




func (b *MetadataBlock) Read(reader io.Reader ) (error ,bool){

	headerBuf := make([]byte , 4)

	_,err := reader.Read(headerBuf)
	lastMetaBlock := false
	if err == nil {
		header , _ := SliceToUint32(headerBuf)

		if (header & 0x80000000) > 0 {
			lastMetaBlock = true
		}

		bType := header & 0x7F000000 >> 24
		switch bType {
			case 0:b.bType=STREAMINFO
			case 1:b.bType=PADDING
			case 2:b.bType=APPLICATION
			case 3:b.bType=SEEKTABLE
			case 4:b.bType=VORBIS_COMMENT
			case 5:b.bType=CUESHEET
			case 6:b.bType=PICTURE
			case 127:b.bType=INVALID
			default:b.bType=RESERVED
		}

		b.size = int32(header&0x00FFFFFF)

		b.data = make([]byte , b.size)
		_,err =reader.Read(b.data)
	}
	return err,lastMetaBlock
}

type Frame struct {

}


func (f *Frame)Read(reader io.Reader) error {
	headerBuf := make([]byte , 4)

	_,err := reader.Read(headerBuf)
	if err == nil {
	//	header, _ := SliceToUint32(headerBuf)

	//	header & 0xFFC00000 = b11111111111110
	}

	return nil
}


func (a *FLACAnalyzer)Analyser(filePath string , lv analyseLv) Report {
	f , err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return Report{}
	}
	defer f.Close()

	fLaC := make([]byte , 4)
	f.Read(fLaC)

	fLaCType , _:= SliceToUint32(fLaC)

	if fLaCType == 0x664c6143 {

		lastMetaBlock := false

		for !lastMetaBlock {
			block := MetadataBlock{}
			err, lastMetaBlock = block.Read(f)
			if err == nil {
				a.metaBlocks = append(a.metaBlocks, &block)
			}

			fmt.Println(block.bType , block.size )
			if block.bType == STREAMINFO {
				fmt.Println(hex.Dump(block.data))
			}
		}


	}else {
		fmt.Printf("is not flac : %s : %s" , filePath , hex.Dump(fLaC))
	}


	return Report{

	}
}