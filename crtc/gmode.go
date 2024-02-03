package crtc

import (
	"image/color"
)

var (
	screenChar  byte     = 0
	pixelData   byte     = 0
	hiresPixels [14]byte = [14]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

var screenLine = [24]uint16{
	0x0000, 0x0080, 0x0100, 0x0180, 0x0200, 0x0280, 0x0300, 0x0380,
	0x0028, 0x00A8, 0x0128, 0x01A8, 0x0228, 0x02A8, 0x0328, 0x03A8,
	0x0050, 0x00D0, 0x0150, 0x01D0, 0x0250, 0x02D0, 0x0350, 0x03D0,
}

var boxLine = [8]uint16{0x0000, 0x0400, 0x0800, 0x0C00, 0x1000, 0x1400, 0x1800, 0x1C00}

// ////////////////////////////////////////////////////////////////////
//
//	Pour Apple II Original                      //
//
// ////////////////////////////////////////////////////////////////////
func (C *CRTC) StandardTextModeA2(X int, Y int) {
	C.videoRam = C.VideoMEM[Set_MEM][0][Set_PAGE]
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
		if pixelData&bit == bit {
			C.graph.DrawPixel(X+column, Y, C.TextColor)
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}

//////////////////////////////////////////////////////////////////////
//                       Pour Apple II+ / IIe                       //
//////////////////////////////////////////////////////////////////////

func (C *CRTC) Standard80ColTextMode(X int, Y int) {
	screenChar = C.videoAux[screenLine[C.RasterLine]+uint16(C.CCLK/2)]
	pixelData = C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	pixelData = ^pixelData

	for column := 0; column < 7; column++ {
		bit := byte(0b00000001 << column)
		if pixelData&bit == bit {
			C.graph.DrawPixel(X+column, Y, C.TextColor)
			C.graph.DrawPixel(X+column, Y+1, C.TextColor)
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
			C.graph.DrawPixel(X+column, Y+1, Colors[Black])
		}
	}

	screenChar = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK/2)]
	pixelData = C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	pixelData = ^pixelData

	for column := 0; column < 7; column++ {
		bit := byte(0b00000001 << column)
		if pixelData&bit == bit {
			C.graph.DrawPixel(X+7+column, Y, C.TextColor)
			C.graph.DrawPixel(X+7+column, Y+1, C.TextColor)
		} else {
			C.graph.DrawPixel(X+7+column, Y, Colors[Black])
			C.graph.DrawPixel(X+7+column, Y+1, Colors[Black])
		}
	}
}

func (C *CRTC) StandardTextModeA2E(X int, Y int) {
	C.videoRam = C.VideoMEM[Set_MEM][0][Set_PAGE]
	screenChar = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)]
	pixelData = C.charRom[uint16(screenChar)<<3+uint16(C.RasterCount)]
	pixelData = ^pixelData

	for column := 0; column < 7; column++ {
		bit := byte(0b00000001 << column)
		if pixelData&bit == bit {
			C.graph.DrawPixel(X+column, Y, C.TextColor)
		} else {
			C.graph.DrawPixel(X+column, Y, Colors[Black])
		}
	}
}

func (C *CRTC) LoResMode(X int, Y int) {
	var color byte

	if Set_MIXED == 1 && C.RasterLine >= 20 {
		if C.conf.Model == "Apple2" {
			C.StandardTextModeA2(X, Y)
		} else {
			C.StandardTextModeA2E(X, Y)
		}
	} else {
		C.videoRam = C.VideoMEM[Set_MEM][0][Set_PAGE]
		screenChar = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)]
		// fmt.Printf("%d ", C.RasterCount)
		if C.RasterCount < 4 {
			// if screenLine[C.RasterLine]&0x01 == 0x01 {
			color = screenChar & 0b00001111
			// color = screenChar >> 4
		} else {
			// color = screenChar & 0b00001111
			color = screenChar >> 4
		}
		for column := 0; column < 7; column++ {
			C.graph.DrawPixel(X+column, Y, Colors[color])
		}
	}
}

var hiresColor [2][4]color.Color = [2][4]color.Color{
	{Colors[Black], Colors[Purple], Colors[LightGreen], Colors[White]},
	{Colors[Black], Colors[MediumBlue], Colors[Orange], Colors[White]},
}

func (C *CRTC) HiResMode(X int, Y int) {

	if Set_MIXED == 1 && C.RasterLine >= 20 {
		if C.conf.Model == "Apple2" {
			C.StandardTextModeA2(X, Y)
		} else {
			C.StandardTextModeA2E(X, Y)
		}
	} else {
		C.videoRam = C.VideoMEM[Set_MEM][Set_MODE][Set_PAGE]
		if C.conf.ColorDisplay {
			if C.CCLK%2 == 0 {
				line := boxLine[Y%8]
				pixelData = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)+line]

				colMode := (pixelData & 0b10000000) >> 7
				hiresPixels[0] = (pixelData & 0b00000011)
				hiresPixels[1] = (pixelData & 0b00000110) >> 1
				hiresPixels[2] = (pixelData & 0b00001100) >> 2
				hiresPixels[3] = (pixelData & 0b00011000) >> 3
				hiresPixels[4] = (pixelData & 0b00110000) >> 4
				hiresPixels[5] = (pixelData & 0b01100000) >> 5
				hiresPixels[6] = (pixelData & 0b01000000) >> 6
				for i := 0; i < 3; i++ {
					C.graph.DrawPixel(X+i*2, Y, hiresColor[colMode][hiresPixels[i*2]])
					// C.graph.DrawPixel(X+i*2+1, Y, hiresColor[colMode][hiresPixels[i*2+1]])
					if (hiresPixels[i*2+1] != 0x00) && (hiresPixels[i*2+1] != 0b00000011) {
						C.graph.DrawPixel(X+(i*2)+1, Y, hiresColor[colMode][hiresPixels[i*2]])
					} else {
						C.graph.DrawPixel(X+(i*2)+1, Y, hiresColor[colMode][hiresPixels[i*2+1]])
					}
				}

				pixelData = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK+1)+line]
				colMode = (pixelData & 0b10000000) >> 7
				hiresPixels[6] += (pixelData & 0b00000001) << 1
				hiresPixels[7] = (pixelData & 0b00000011)
				hiresPixels[8] = (pixelData & 0b00000110) >> 1
				hiresPixels[9] = (pixelData & 0b00001100) >> 2
				hiresPixels[10] = (pixelData & 0b00011000) >> 3
				hiresPixels[11] = (pixelData & 0b00110000) >> 4
				hiresPixels[12] = (pixelData & 0b01100000) >> 5
				hiresPixels[13] = (pixelData & 0b01100000) >> 5

				for i := 3; i < 7; i++ {
					C.graph.DrawPixel(X+i*2, Y, hiresColor[colMode][hiresPixels[i*2]])
					// C.graph.DrawPixel(X+i*2+1, Y, hiresColor[colMode][hiresPixels[i*2+1]])
					if (hiresPixels[i*2+1] != 0x00) && (hiresPixels[i*2+1] != 0b00000011) {
						C.graph.DrawPixel(X+(i*2)+1, Y, hiresColor[colMode][hiresPixels[i*2]])
					} else {
						C.graph.DrawPixel(X+(i*2)+1, Y, hiresColor[colMode][hiresPixels[i*2+1]])
					}
				}
			}
		} else {
			line := boxLine[Y%8]
			pixelData = C.videoRam[screenLine[C.RasterLine]+uint16(C.CCLK)+line]
			for column := 0; column < 7; column++ {
				bit := byte(0b00000001 << column)
				if pixelData&bit == bit {
					C.graph.DrawPixel(X+column, Y, C.TextColor)
				} else {
					C.graph.DrawPixel(X+column, Y, Colors[Black])
				}
			}
		}
	}
}
