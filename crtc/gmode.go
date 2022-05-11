package crtc

var (
	screenChar byte = 0
	pixelData  byte = 0
)

var screenLine = [24]uint16{
	0x0000, 0x0080, 0x0100, 0x0180, 0x0200, 0x0280, 0x0300, 0x0380,
	0x0028, 0x00A8, 0x0128, 0x01A8, 0x0228, 0x02A8, 0x0328, 0x03A8,
	0x0050, 0x00D0, 0x0150, 0x01D0, 0x0250, 0x02D0, 0x0350, 0x03D0,
}

//////////////////////////////////////////////////////////////////////
//                      Pour Apple II Original                      //
//////////////////////////////////////////////////////////////////////
func (C *CRTC) StandardTextModeA2(X int, Y int) {
	screenChar = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)]
	pixelData = C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	switch screenChar & 0b11000000 {
	case 0:
		pixelData = ^pixelData
	case 0b01000000:
		if blink {
			pixelData = ^pixelData
		}
	}

	for column := 0; column < 7; column++ {
		bit := byte(0b01000000 >> column)
		if pixelData&bit > 1 {
			C.graph.DrawPixel(X+column, Y, Colors[Green])
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}

//////////////////////////////////////////////////////////////////////
//                       Pour Apple II+ / IIe                       //
//////////////////////////////////////////////////////////////////////
func (C *CRTC) StandardTextModeA2E(X int, Y int) {
	screenChar = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)]
	pixelData = C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	pixelData = ^pixelData

	for column := 0; column < 7; column++ {
		bit := byte(0b00000001 << column)
		if pixelData&bit > 1 {
			C.graph.DrawPixel(X+column, Y, Colors[Green])
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}
