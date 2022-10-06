package crtc

import (
	"fmt"
	"newApple/config"
	"time"

	"github.com/Djoulzy/emutools/render"
)

var blink bool = false

func NE5555() {
	ticker := time.NewTicker(time.Millisecond * 200)
	defer func() {
		ticker.Stop()
	}()

	for {
		<-ticker.C
		blink = !blink
	}
}

func (C *CRTC) Init(ram []byte, aux []byte, io []byte, chargen []byte, video *render.SDL2Driver, conf *config.ConfigData) {
	C.Reg[R0] = 63 // 126
	C.Reg[R1] = 40 // 80 colonnes
	C.Reg[R2] = 50
	C.Reg[R3] = 0b10001000
	C.Reg[R4] = 32
	C.Reg[R5] = 16
	C.Reg[R6] = 24
	C.Reg[R7] = 29
	C.Reg[R9] = 8
	C.Reg[R12] = 0
	C.Reg[R13] = 0

	C.RAM = ram
	C.AUX = aux
	C.screenWidth = int(C.Reg[R1]) * 7
	C.screenHeight = int(C.Reg[R6]) * 8 // * 2

	C.graph = video
	C.graph.Init(C.screenWidth, C.screenHeight, 2, "Go Apple II", false, conf.Disassamble)
	C.conf = conf
	C.VideoPages[0] = [2]uint16{0x0400, 0x2000}
	C.VideoPages[1] = [2]uint16{0x0400, 0x2000}
	C.VideoPages[2] = [2]uint16{0x0800, 0x4000}

	C.charRom = chargen

	C.BeamX = 0
	C.BeamY = 0
	C.RasterLine = 0
	C.RasterCount = 0
	C.CCLK = 0

	if conf.ColorDisplay {
		C.TextColor = Colors[White]
	} else {
		C.TextColor = Colors[LightGreen]
	}

	C.UpdateGraphMode()

	if C.conf.Model == "2" {
		go NE5555()
	}
}

func (C *CRTC) ToggleMonitorColor() {
	C.conf.ColorDisplay = !C.conf.ColorDisplay
	if C.conf.ColorDisplay {
		C.TextColor = Colors[White]
	} else {
		C.TextColor = Colors[LightGreen]
	}
}

func (C *CRTC) UpdateGraphMode() {
	C.UpdateVideoRam()
	if Is_TEXTMODE {
		if C.conf.Model == "2" {
			C.videoMode = (*CRTC).StandardTextModeA2
		} else {
			if Is_80COL {
				C.videoMode = (*CRTC).Standard80ColTextMode
			} else {
				C.videoMode = (*CRTC).StandardTextModeA2E
			}
		}
	} else {
		if Is_HIRESMODE {
			C.videoMode = (*CRTC).HiResMode
		} else {
			C.videoMode = (*CRTC).LoResMode
		}
	}
}

func (C *CRTC) UpdateVideoRam() {
	var page byte
	if Is_PAGE2 {
		page = 2
	} else {
		page = 1
	}
	if Is_TEXTMODE {
		C.videoBase = C.VideoPages[page][0]
		C.pageSize = C.VideoPages[0][0]
		C.videoRam = C.RAM[C.videoBase : C.videoBase+C.pageSize]
		C.videoAux = C.AUX[C.videoBase : C.videoBase+C.pageSize]
	} else {
		C.videoBase = C.VideoPages[page][1]
		C.pageSize = C.VideoPages[0][1]
		C.videoRam = C.RAM[C.videoBase : C.videoBase+C.pageSize]
		C.videoAux = C.AUX[C.videoBase : C.videoBase+C.pageSize]
	}
}

func (C *CRTC) DumpMode() {
	if Is_TEXTMODE {
		fmt.Printf("TEXT ")
	} else if Is_HIRESMODE {
		fmt.Printf("HIRES ")
	} else {
		fmt.Printf("LORES ")
	}
	fmt.Printf("VideoRam: %04X Size: %04X\n", C.videoBase, C.pageSize)
}

func (C *CRTC) drawChar(X int, Y int) {
	// if C.drawArea && (C.Reg[REG_CTRL1]&DEN > 0) {
	C.videoMode(C, X, Y)
}

func (C *CRTC) Run(debug bool) bool {
	C.BeamX = int(C.CCLK) * 7

	// log.Printf("BeamX: %d - BeamY: %d - CCLK: %02d - RasterLine: %02d", C.BeamX, C.BeamY, C.CCLK, C.RasterLine)

	if C.CCLK < (C.Reg[R1]) {
		C.drawChar(C.BeamX, C.BeamY)
	}

	C.CCLK++ // += 2
	if C.CCLK == C.Reg[R0] {
		C.CCLK = 0
		C.BeamY++ // += 2
		if C.BeamY >= C.screenHeight {
			C.BeamY = 0
			C.RasterCount = 0
			C.RasterLine = 0
			// C.graph.UpdateFrame()
		} else {
			C.RasterCount++
			if C.RasterCount == C.Reg[R9] {
				C.RasterLine++
				C.RasterCount = 0
			}
		}
	}
	return true
}
