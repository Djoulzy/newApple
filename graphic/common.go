package graphic

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	Xadjust        = 200
	Yadjust        = 100
	fontWidth      = 7
	fontHeight     = 9
	mnemonicWidth  = fontWidth * 3
	mnemonicHeight = fontHeight

	setFPS = 50
	throttleFPS = 1000/setFPS
)

type KEYPressed struct {
	KeyCode uint
	Mode    uint
}

var fps, frameCount, lastFrame, lastTime, timerFPS uint32
var debug *sdl.Texture
