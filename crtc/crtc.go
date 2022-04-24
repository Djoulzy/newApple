package crtc

import (
	"newApple/config"
	"newApple/graphic"
)

const (
	// Durée d'un cycle PAL : 1.015uS
	screenWidthPAL  = 504
	screenHeightPAL = 312
	rasterWidthPAL  = 403
	rasterHeightPAL = 284
	cyclesPerLine   = 63

	rasterTime = 1                  // Nb of cycle to put 1 byte on a line
	rasterLine = rasterWidthPAL / 8 // Nb of cycle to draw a full line
	fullRaster = rasterLine * rasterHeightPAL

	winWidth      = screenWidthPAL
	winHeight     = screenHeightPAL
	visibleWidth  = 320
	visibleHeight = 200

	firstVBlankLine  = 300
	lastVBlankLine   = 15
	firstDisplayLine = 51
	lastDisplayLine  = 250
)

func (C *CRTC) Init(ram []byte, io []byte, chargen []byte, video interface{}, conf *config.ConfigData) {
	C.Reg[R0] = 63
	C.Reg[R1] = 40
	C.Reg[R2] = 50
	C.Reg[R3] = 0b10001000
	C.Reg[R4] = 32
	C.Reg[R5] = 16
	C.Reg[R6] = 24
	C.Reg[R7] = 29
	C.Reg[R9] = 8
	C.Reg[R12] = 0
	C.Reg[R13] = 0

	C.graph = video.(graphic.Driver)
	C.graph.Init(winWidth, winHeight)
	C.conf = conf

	C.videoRam = ram[screenStart : screenStart+screenSize]
	C.charRom = chargen

	C.BeamX = 0
	C.BeamY = 0
	C.RasterLine = 0
	C.RasterCount = 0
	C.CCLK = 0
}

// func (C *CRTC) Disassemble() string {
// 	var buf string
// 	buf = fmt.Sprintf("BeamX: %03d - BeamY: %03d - CCLK: %02d - RasterLine: %02d", C.BeamX, C.BeamY, C.CCLK, C.RasterLine)
// 	return buf
// }

// func (C *CRTC) saveRasterPos(val int) {
// 	C.Reg[REG_RASTER] = byte(val)
// 	mask := byte(val>>1) & RST8
// 	res := C.Reg[REG_CTRL1] & 0b01111111
// 	C.Reg[REG_CTRL1] = res | mask
// }

// func (C *CRTC) readVideoMatrix() {
// 	C.ColorBuffer[C.VMLI] = C.color[C.VC] & 0b00001111
// 	C.CharBuffer[C.VMLI] = C.bankMem.Read(C.ScreenBase + C.VC)
// 	// fmt.Printf("VMLI: %02X - VC: %02X - Screen Code: %d - Color: %04X\n", C.VMLI, C.VC, C.CharBuffer[C.VMLI], C.ColorBuffer[C.VMLI])
// }

func (C *CRTC) drawChar(X int, Y int) {
	// if C.drawArea && (C.Reg[REG_CTRL1]&DEN > 0) {
	if C.visibleArea {
		C.StandardTextMode(X, Y)
	} else if C.syncArea {
		for column := 0; column < 8; column++ {
			C.graph.DrawPixel(X+column, Y, Colors[Blue])
		}
	} else {
		for column := 0; column < 8; column++ {
			C.graph.DrawPixel(X+column, Y, Colors[Red])
		}
	}
}

func (C *CRTC) Run(debug bool) bool {
	C.visibleArea = (C.CCLK < C.Reg[R1]) && (C.RasterLine < C.Reg[R6])
	SyncAreaH := (C.CCLK >= C.Reg[R2]) && (C.CCLK <= C.Reg[R2]+(C.Reg[R3]&0b00001111))
	SyncAreaV := (C.RasterLine >= C.Reg[R7]) && (C.RasterLine <= C.Reg[R7]+(C.Reg[R3]>>4))
	C.syncArea = SyncAreaH || SyncAreaV
	C.BeamX = int(C.CCLK) * 8

	// log.Printf("BeamX: %d - BeamY: %d - CCLK: %02d - RasterLine: %02d", C.BeamX, C.BeamY, C.CCLK, C.RasterLine)

	C.drawChar(C.BeamX, C.BeamY)

	C.CCLK++
	if C.CCLK >= C.Reg[R0] {
		C.CCLK = 0
		C.BeamY++
		if C.BeamY >= screenHeightPAL {
			C.BeamY = 0
			C.RasterCount = 0
			C.RasterLine = 0
			C.graph.UpdateFrame()
		} else {
			C.RasterCount++
			if C.RasterCount == C.Reg[R9] {
				C.RasterLine++
				C.RasterCount = 0
			}
		}
	}
	if debug {
		C.graph.UpdateFrame()
	}
	return true
}

// func (C *CRTC) Dump(addr uint16) {
// 	fmt.Printf("Bank: %d - VideoBase: %04X - CharBase: %04X", C.BankSel, C.ScreenBase, C.CharBase)
// 	C.bankMem.Show()
// 	C.bankMem.Dump(addr)
// }

// func (C *CRTC) Stats() {
// 	banks := [4]uint16{BankStart0, BankStart1, BankStart2, BankStart3}

// 	fmt.Printf("VIC:\n")
// 	fmt.Printf("Bank: %d - VideoBase: %04X (%04X) - CharBase: %04X (%04X)\n", C.BankSel, C.ScreenBase, banks[C.BankSel]+C.ScreenBase, C.CharBase, banks[C.BankSel]+C.CharBase)
// 	fmt.Printf("RstX: %04X - RstY: %04X - RC: %02d - VC: %03X - VCBase: %03X - VMLI: %02d\n", C.BeamX, C.BeamY, C.RC, C.VC, C.VCBASE, C.VMLI)
// 	fmt.Printf("IRQ Line: ")
// 	if C.Reg[REG_IRQ]&0b10000000 > 0 {
// 		fmt.Printf("On")
// 	} else {
// 		fmt.Printf("Off")
// 	}
// 	fmt.Printf(" - IRQ Enabled: ")
// 	if C.Reg[REG_IRQ_ENABLED]&0b00001111 > 0 {
// 		fmt.Printf("Yes")
// 	} else {
// 		fmt.Printf("None")
// 	}
// 	fmt.Printf(" - Raster IRQ: %04X\n", C.RasterIRQ)
// }
