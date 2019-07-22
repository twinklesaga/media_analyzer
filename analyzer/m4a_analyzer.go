package analyzer

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type m4aContext struct {
	Version int
}

type AtomProcess struct {
	HasChild bool
	Process  func(*m4aContext, []byte)error
}

type M4AAnalyzer struct {
	atomMap map[atomType]AtomProcess
}

func NewM4AAnalyzer() MediaAnalyzer{
	m4a := new(M4AAnalyzer)
	m4a.makeAtomMap()

	return m4a
}

func (a *M4AAnalyzer)Analyze(filePath string , lv AnalyzeLV) Report{
	f , err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return Report{}
	}
	defer f.Close()



	ctx := m4aContext{}
	for {
		AtomBuf := make([]byte, 8)
		n, err := f.Read(AtomBuf)
		if err != nil || n != 8 {
			if err == io.EOF {
				fmt.Println("EOF")
			}
			break
		}
		aSize , err := SliceToUint32(AtomBuf)
		aType , err := SliceToUint32(AtomBuf[4:])
		p , ok := a.atomMap[atomType(aType)]

		fmt.Printf("%s (%d)\n" , string(AtomBuf[4:]) , aSize)
		if ok {
			if p.HasChild {

			}else if p.Process != nil {
				buf := make([]byte , aSize - 8)
				n , err = f.Read(buf)
				if err != nil {
					fmt.Println(err)
				}else if n != int(aSize - 8) {
					//fmt.Println("------------------------------------missmatch read size")
				}

				p.Process(&ctx , buf)
			}else {
				f.Seek(int64(aSize - 8) , 1)
			}
		} else {
			break
		}
	}
	return Report{}
}


type atomType uint32
const(
	ftyp atomType = 0x66747970
	moov atomType = 0x6d6f6f76
	mvhd atomType = 0x6d766864
	trak atomType = 0x7472616b
	tkhd atomType = 0x746b6864
	edts atomType = 0x65647473
	elst atomType = 0x656c7374
	mdia atomType = 0x6d646961
	mdhd atomType = 0x6d646864
	hdlr atomType = 0x68646c72
	minf atomType = 0x6d696e66
	smhd atomType = 0x736d6864
	dinf atomType = 0x64696e66
	stbl atomType = 0x7374626c
	stsd atomType = 0x73747364
	stts atomType = 0x73747473
	stsc atomType = 0x73747363
	stsz atomType = 0x7374737a
	stco atomType = 0x7374636f
	sgpd atomType = 0x73677064
	sbgp atomType = 0x73626770
	udta atomType = 0x75647461
	free atomType = 0x66726565
	mdat atomType = 0x6d646174
)

func (a *M4AAnalyzer)makeAtomMap(){
	a.atomMap = make(map[atomType]AtomProcess)
	a.atomMap[ftyp] = AtomProcess{HasChild:false ,  Process:a.checkFTYP}
	a.atomMap[moov] = AtomProcess{HasChild:true , Process:a.checkDUMP}
	a.atomMap[mvhd] = AtomProcess{HasChild:false , Process:a.checkDUMP}
	a.atomMap[trak] = AtomProcess{HasChild:true ,  Process:a.checkDUMP}
	a.atomMap[tkhd] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[edts] = AtomProcess{HasChild:true ,  Process:a.checkDUMP}
	a.atomMap[elst] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[mdia] = AtomProcess{HasChild:true ,  Process:a.checkDUMP}
	a.atomMap[mdhd] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[hdlr] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[minf] = AtomProcess{HasChild:true ,  Process:a.checkMINF}
	a.atomMap[smhd] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[dinf] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[stbl] = AtomProcess{HasChild:true ,  Process:a.checkDUMP}
	a.atomMap[stsd] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[stts] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[stsc] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[stsz] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[stco] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[sgpd] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[sbgp] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[udta] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[free] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
	a.atomMap[mdat] = AtomProcess{HasChild:false ,  Process:a.checkDUMP}
}

func (a *M4AAnalyzer)checkFTYP(ctx *m4aContext, data []byte) error {
	fmt.Println(hex.Dump(data))
	return nil
}

func (a *M4AAnalyzer)checkMVHD(ctx *m4aContext, data []byte ) error {

	return nil
}

func (a *M4AAnalyzer)checkMINF(ctx *m4aContext, data []byte ) error {

	return nil
}

func (a *M4AAnalyzer)checkDUMP(ctx *m4aContext, data []byte) error {
	if len(data) > 1024 {
		fmt.Println(hex.Dump(data[:1024]))
	}else {
		fmt.Println(hex.Dump(data))
	}
	return nil
}