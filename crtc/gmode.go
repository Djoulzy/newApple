package crtc

var (
	CharOutCount uint16 = 0
	CharBlock    uint16 = 0
)

var videoBase = [8]uint16{0x0000, 0x0080, 0x0100, 0x0180, 0x0200, 0x0280, 0x0300, 0x0380}

func (C *CRTC) StandardTextMode(X int, Y int) {
	charPack := C.RasterLine / 3

	// charPos := uint16(C.RasterLine)*uint16(C.Reg[R1]) + uint16(C.CCLK)
	charPos := (C.RasterLine-(byte(charPack)*3))*C.Reg[R1] + C.CCLK

	screenChar := C.videoRam[videoBase[charPack]+uint16(charPos)]
	pixelData := C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	// fmt.Printf("%04X - %02X\n", charPos, pixelData)

	// fmt.Printf("(%d %d): %d - %d\n", C.RasterLine, C.CCLK, charPack, charPos)
	for column := 0; column < 7; column++ {
		bit := byte(0b01000000 >> column)
		if pixelData&bit > 0 {
			C.graph.DrawPixel(X+column, Y, Colors[Green])
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}
