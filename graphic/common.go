package graphic

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	Xadjust = 200
	Yadjust = 100
)

type KEYPressed struct {
	KeyCode uint
	Mode    uint
}

var setFPS uint32 = 50
var fps, frameCount, lastFrame, lastTime, timerFPS uint32
var font *ttf.Font
var fpsDisp *sdl.Surface
var debug *sdl.Texture
