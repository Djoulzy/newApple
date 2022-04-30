package crtc

import (
	"image/color"
	"newApple/config"

	"github.com/Djoulzy/emutools/render"
)

var (
	Black      byte = 0
	White      byte = 1
	Red        byte = 2
	Cyan       byte = 3
	Violet     byte = 4
	Green      byte = 5
	Blue       byte = 6
	Yellow     byte = 7
	Orange     byte = 8
	Brown      byte = 9
	Lightred   byte = 10
	Darkgrey   byte = 11
	Grey       byte = 12
	Lightgreen byte = 13
	Lightblue  byte = 14
	Lightgrey  byte = 15
)

var Colors [16]color.Color = [16]color.Color{
	color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
	color.RGBA{R: 137, G: 78, B: 67, A: 255},   // Red
	color.RGBA{R: 146, G: 195, B: 203, A: 255}, // Cyan
	color.RGBA{R: 138, G: 87, B: 176, A: 255},  // Violet
	color.RGBA{R: 128, G: 174, B: 89, A: 255},  // Green
	color.RGBA{R: 68, G: 63, B: 164, A: 255},   // Blue
	color.RGBA{R: 215, G: 221, B: 137, A: 255}, // Yellow
	color.RGBA{R: 146, G: 106, B: 56, A: 255},  // Orange
	color.RGBA{R: 100, G: 82, B: 23, A: 255},   // Brown
	color.RGBA{R: 184, G: 132, B: 122, A: 255}, // Lightred
	color.RGBA{R: 96, G: 96, B: 96, A: 255},    // Darkgrey
	color.RGBA{R: 138, G: 138, B: 138, A: 255}, // Grey
	color.RGBA{R: 191, G: 233, B: 155, A: 255}, // Lightgreen
	color.RGBA{R: 131, G: 125, B: 216, A: 255}, // Lightblue
	color.RGBA{R: 179, G: 179, B: 179, A: 255}, // Lightgrey
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
	RAM          []byte
	screenWidth  int
	screenHeight int

	conf        *config.ConfigData
	BeamX       int
	BeamY       int
	RasterLine  byte
	RasterCount byte
	CCLK        byte

	VideoPage uint16
	VideoSize uint16

	graph *render.SDL2Driver
	MODE  byte

	videoRam []byte
	charRom  []byte
}

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
