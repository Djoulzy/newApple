package crtc

import (
	"image/color"
	"newApple/config"

	"github.com/Djoulzy/emutools/render"
)

var (
	Black      byte = 0
	Red        byte = 1
	DarkBlue   byte = 2
	Purple     byte = 3
	DarkGreen  byte = 4
	Gray       byte = 5
	MediumBlue byte = 6
	LightBlue  byte = 7
	Brown      byte = 8
	Orange     byte = 9
	Grey       byte = 10
	Pink       byte = 11
	LightGreen byte = 12
	Yellow     byte = 13
	Aqua       byte = 14
	White      byte = 15
)

var Colors [16]color.Color = [16]color.Color{
	color.RGBA{R: 0, G: 0, B: 0, A: 255},       //black
	color.RGBA{R: 72, G: 58, B: 211, A: 255},   //red
	color.RGBA{R: 163, G: 30, B: 9, A: 255},    //dk blue
	color.RGBA{R: 221, G: 84, B: 213, A: 255},  //purple
	color.RGBA{R: 57, G: 133, B: 54, A: 255},   //dk green
	color.RGBA{R: 104, G: 104, B: 104, A: 255}, //gray
	color.RGBA{R: 246, G: 68, B: 51, A: 255},   //med blue
	color.RGBA{R: 249, G: 185, B: 134, A: 255}, //lt blue
	color.RGBA{R: 33, G: 106, B: 147, A: 255},  //brown
	color.RGBA{R: 49, G: 131, B: 240, A: 255},  //orange
	color.RGBA{R: 184, G: 184, B: 184, A: 255}, //grey
	color.RGBA{R: 157, G: 175, B: 244, A: 255}, //pink
	color.RGBA{R: 64, G: 219, B: 97, A: 255},   //lt green
	color.RGBA{R: 82, G: 251, B: 254, A: 255},  //yellow
	color.RGBA{R: 210, G: 247, B: 134, A: 255}, //aqua
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, //white
}

// 	{R: 0, G: 0, B: 0},
// 	{R: 255, G: 255, B: 255}, // White
// 	{R: 137, G: 78, B: 67},   // Red
// 	{R: 146, G: 195, B: 203}, // Cyan
// 	{R: 138, G: 87, B: 176},  // Violet
// 	{R: 128, G: 174, B: 89},  // Green
// 	{R: 68, G: 63, B: 164},   // Blue
// 	{R: 215, G: 221, B: 137}, // Yellow
// 	{R: 146, G: 106, B: 56},  // Orange
// 	{R: 100, G: 82, B: 23},   // Brown
// 	{R: 184, G: 132, B: 122}, // Lightred
// 	{R: 96, G: 96, B: 96},    // Darkgrey
// 	{R: 138, G: 138, B: 138}, // Grey
// 	{R: 191, G: 233, B: 155}, // Lightgreen
// 	{R: 131, G: 125, B: 216}, // Lightblue
// 	{R: 179, G: 179, B: 179}, // Lightgrey
// }

// VIC :
type CRTC struct {
	Reg          [18]byte
	screenWidth  int
	screenHeight int

	conf        *config.ConfigData
	BeamX       int
	BeamY       int
	RasterLine  byte
	RasterCount byte
	CCLK        byte
	TextColor   color.Color

	VideoPages [3][2]uint16
	videoBase  uint16
	pageSize   uint16

	graph *render.SDL2Driver

	videoRam  []byte
	videoAux  []byte
	charRom   []byte
	videoMode func(*CRTC, int, int)
}

const (
	TXTPGSIZE  = 1024
	GRPGSIZE   = 8192
	TEXTPAGE1  = 0x0400
	TEXTPAGE2  = 0x0800
	HIRESPAGE1 = 0x2000
	HIRESPAGE2 = 0x4000
)

const (
	R0 byte = iota // Longueur d'une ligne (displayed + sync)
	R1             // Nb de characteres par ligne
	R2             // Pos du sync start par apport au debut de la ligne
	R3             // Sync control (0-3: Horizontal, 4-7: Vertical)
	R4             // Nb total de lignes
	R5             // Nb de scanlines à ajouter pour compléter l'ecran
	R6             // Nb de lignes visibles affichées
	R7             // Pos du vertical sync
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
	R16
	R17
)

var (
	Is_TEXTMODE  bool = true
	Is_MIXEDMODE bool = true
	Is_HIRESMODE bool = false
	Is_PAGE2     bool = false
	Is_80COL     bool = false
)
