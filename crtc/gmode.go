package crtc

func (C *CRTC) StandardTextMode(X int, Y int) {
	screenChar := C.videoRam[uint16(C.RasterLine)*uint16(C.Reg[R1])+uint16(C.CCLK)]
	pixelData := C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]

	for column := 0; column < 8; column++ {
		bit := byte(0b10000000 >> column)
		if pixelData&bit > 0 {
			C.graph.DrawPixel(X+column, Y, Colors[Green])
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}
