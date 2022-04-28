package graphic

import "image/color"

type Driver interface {
	Init(int, int, string)
	DirectDrawPixel(int, int, color.Color)
	DrawPixel(int, int, color.Color)
	UpdateFrame()
	// Run()
	IOEvents() *KEYPressed
	SetKeyboardLine(*KEYPressed)
	SetCodeList([][]byte)
	ShowCode(*uint16)
	CloseAll()
}
