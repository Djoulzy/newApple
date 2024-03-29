package crtc

import (
	"fmt"
	"newApple/config"
	"time"

	"github.com/Djoulzy/emutools/render"
	"github.com/Djoulzy/mmu"
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

func (C *CRTC) Init(mem *mmu.MMU, chargen *mmu.ROM, video *render.SDL2Driver, conf *config.ConfigData) {
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

	C.screenWidth = int(C.Reg[R1]) * 7
	C.screenHeight = int(C.Reg[R6]) * 8 // * 2

	C.graph = video
	C.graph.Init(C.screenWidth, C.screenHeight, 2, "Go Apple II", true, false)
	C.conf = conf
	C.VideoPages[0] = [2]uint16{0x0400, 0x2000}
	C.VideoPages[1] = [2]uint16{0x0400, 0x2000}
	C.VideoPages[2] = [2]uint16{0x0800, 0x4000}

	C.charRom = chargen.Buff

	C.VideoMEM[0][0][0] = mem.GetChipMem("MN_TXT") // [main][textmode][page1]
	C.VideoMEM[0][0][1] = mem.GetChipMem("MN___2") // [main][textmode][page2]
	C.VideoMEM[0][1][0] = mem.GetChipMem("MN_HGR") // [main][hires][page1]
	C.VideoMEM[0][1][1] = mem.GetChipMem("MN___3") // [main][hires][page2]

	if C.conf.Model == "Apple2e" {
		C.VideoMEM[1][0][0] = mem.GetChipMem("AX_TXT") // [aux][textmode][page1]
		C.VideoMEM[1][0][1] = mem.GetChipMem("AX___2") // [aux][textmode][page2]
		C.VideoMEM[1][1][0] = mem.GetChipMem("AX_HGR") // [aux][hires][page1]
		C.VideoMEM[1][1][1] = mem.GetChipMem("AX___3") // [aux][hires][page2]
	}

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

	C.SetTexMode()

	if C.conf.Model == "Apple2" {
		go NE5555()
	}
}

func (C *CRTC) SetTexMode() {
	Set_MODE = 0
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	if C.conf.Model == "Apple2" {
		C.videoMode = (*CRTC).StandardTextModeA2
	} else {
		if Set_80COL == 1 {
			C.videoAux = C.VideoMEM[1][Set_MODE][Set_PAGE]
			C.videoMode = (*CRTC).Standard80ColTextMode
		} else {
			C.videoMode = (*CRTC).StandardTextModeA2E
		}
	}
}

func (C *CRTC) Set40Cols() {
	Set_80COL = 0
	// C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	// if C.conf.Model == "Apple2" {
	// 	C.videoMode = (*CRTC).StandardTextModeA2
	// } else {
	// 	C.videoMode = (*CRTC).StandardTextModeA2E
	// }
}

func (C *CRTC) Set80Cols() {
	Set_80COL = 1
	// C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	// if C.conf.Model == "Apple2e" {
	// 	C.videoAux = C.VideoMEM[1][Set_MODE][Set_PAGE]
	// 	C.videoMode = (*CRTC).Standard80ColTextMode
	// }
}

func (C *CRTC) SetGraphMode() {
	Set_MODE = 1
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	if Set_HIRES == 1 {
		C.videoMode = (*CRTC).HiResMode
	} else {
		C.videoMode = (*CRTC).LoResMode
	}
}

func (C *CRTC) SetMixedMode() {
	Set_MIXED = 1
}

func (C *CRTC) SetFullMode() {
	Set_MIXED = 0
}

func (C *CRTC) SetLoResMode() {
	Set_HIRES = 0
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	if Set_MODE == 1 {
		C.videoMode = (*CRTC).LoResMode
	}
}

func (C *CRTC) SetHiResMode() {
	Set_HIRES = 1
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	if Set_MODE == 1 {
		C.videoMode = (*CRTC).HiResMode
	}
}

func (C *CRTC) SetPage1() {
	Set_PAGE = 0
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	C.videoAux = C.VideoMEM[1][Set_MODE][Set_PAGE]
}

func (C *CRTC) SetPage2() {
	Set_PAGE = 1
	C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
	C.videoAux = C.VideoMEM[1][Set_MODE][Set_PAGE]
}

func (C *CRTC) ToggleMonitorColor() {
	C.conf.ColorDisplay = !C.conf.ColorDisplay
	if C.conf.ColorDisplay {
		C.TextColor = Colors[White]
	} else {
		C.TextColor = Colors[LightGreen]
	}
}

// func (C *CRTC) UpdateGraphMode() {
// 	C.UpdateVideoRam()
// 	if Is_TEXTMODE {
// 		if C.conf.Model == "Apple2" {
// 			C.videoMode = (*CRTC).StandardTextModeA2
// 		} else {
// 			if Is_80COL {
// 				C.videoMode = (*CRTC).Standard80ColTextMode
// 			} else {
// 				C.videoMode = (*CRTC).StandardTextModeA2E
// 			}
// 		}
// 	} else {
// 		if Is_HIRESMODE {
// 			C.videoMode = (*CRTC).HiResMode
// 		} else {
// 			C.videoMode = (*CRTC).LoResMode
// 		}
// 	}
// }

// func (C *CRTC) UpdateVideoRam() {
// 	var page byte
// 	if Is_PAGE2 {
// 		page = 2
// 	} else {
// 		page = 1
// 	}
// 	if Is_TEXTMODE {
// 		C.videoBase = C.VideoPages[page][0]
// 		C.pageSize = C.VideoPages[0][0]
// 		C.videoRam = C.RAM[C.videoBase : C.videoBase+C.pageSize]
// 		// C.videoAux = C.AUX[C.videoBase : C.videoBase+C.pageSize]
// 	} else {
// 		C.videoBase = C.VideoPages[page][1]
// 		C.pageSize = C.VideoPages[0][1]
// 		C.videoRam = C.RAM[C.videoBase : C.videoBase+C.pageSize]
// 		// C.videoAux = C.AUX[C.videoBase : C.videoBase+C.pageSize]
// 	}
// }

func (C *CRTC) DumpMode() {
	var mode string
	if Set_MODE == 0 {
		mode = "TEXT"
	} else {
		if Set_HIRES == 1 {
			mode = "HiRES"
		} else {
			mode = "LoRES"
		}
		if Set_MIXED == 1 {
			mode = mode + " Mixed"
		}
	}
	fmt.Printf("Mode: %s - Page: %d - 80Cols: %d - Mem: %d\n", mode, Set_PAGE, Set_80COL, Set_MEM)
}

func (C *CRTC) drawChar(X int, Y int) {
	// if C.drawArea && (C.Reg[REG_CTRL1]&DEN > 0) {
	C.videoMode(C, X, Y)
}

func (C *CRTC) Run() bool {
	C.BeamX = int(C.CCLK) * 7

	// log.Printf("BeamX: %d - BeamY: %d - CCLK: %02d - RasterLine: %02d", C.BeamX, C.BeamY, C.CCLK, C.RasterLine)

	if C.CCLK < (C.Reg[R1]) {
		C.drawChar(C.BeamX, C.BeamY)
	}

	C.CCLK++ // += 2
	// C.CCLK += 2
	if C.CCLK == C.Reg[R0] {
		C.CCLK = 0
		C.BeamY++ // += 2
		// C.BeamY += 2
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
