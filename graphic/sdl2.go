package graphic

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type SDL2Driver struct {
	winHeight int
	winWidth  int
	emuHeight int
	emuWidth  int
	window    *sdl.Window
	surface   *sdl.Surface
	debug     *sdl.Surface
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	bitmap    *sdl.Surface
	screen    []byte
	keybLine  *KEYPressed
	codeList  map[int][]rune
}

func (S *SDL2Driver) DirectDrawPixel(x, y int, c color.Color) {
	// S.renderer.SetDrawColor(byte(color.R), byte(color.G), byte(color.B), 255)
	// S.renderer.DrawPoint(int32(x), int32(y))
	S.surface.Set(x, y, c)
}

func (S *SDL2Driver) DrawPixel(x, y int, c color.Color) {
	S.surface.Set(x, y, c)
}

func (S *SDL2Driver) CloseAll() {
	S.window.Destroy()
	sdl.Quit()
}

func (S *SDL2Driver) Init(winWidth, winHeight int, title string) {
	S.emuHeight = winHeight
	S.emuWidth = winWidth
	S.winHeight = S.emuHeight * 2
	S.winWidth = S.emuWidth*2 + Xadjust

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	// S.window, S.renderer, err = sdl.CreateWindowAndRenderer(int32(S.winWidth*2), int32(S.winHeight*2), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	// S.window.SetTitle(title)
	S.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(S.winWidth), int32(S.winHeight), sdl.WINDOW_SHOWN)
	S.renderer, err = sdl.CreateRenderer(S.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	S.surface, err = sdl.CreateRGBSurface(0, int32(S.emuWidth), int32(S.emuHeight), 32, 0, 0, 0, 0)
	S.debug, err = sdl.CreateRGBSurface(0, int32(Xadjust), int32(S.winHeight), 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}

	S.bitmap, err = sdl.LoadBMP("graphic/assets/white.bmp")
	if err != nil {
		panic(err)
	}

	S.screen = make([]byte, S.winWidth*S.winHeight*3)
}

func (S *SDL2Driver) SetKeyboardLine(line *KEYPressed) {
	S.keybLine = line
}

func (S *SDL2Driver) SetCodeList(list map[int][]rune) {
	S.codeList = list
}

func getGlyph(char rune) *sdl.Rect {
	pos := int32(char - 32)
	// posy := int32(pos / 18)
	// posx := pos - int32(pos / 18)*18
	// fmt.Printf("r: %c ASCII: %d - abs: %d - x: %d - y: %d\n", char, char, pos, posx, posy)
	return &sdl.Rect{pos*7 - int32(pos/18)*126, int32(pos/18) * 9, 7, 9}
}

func (S *SDL2Driver) getFPS() {
	lastFrame = sdl.GetTicks()
	if lastFrame >= (lastTime + 1000) {
		lastTime = lastFrame
		fps = frameCount
		frameCount = 0
	}
	runes := []rune(fmt.Sprintf("%d", fps))
	for i, r := range runes {
		S.bitmap.Blit(getGlyph(r), S.surface, &sdl.Rect{int32(S.emuWidth - 21 + i*7), 2, 7, 9})
	}
}

func (S *SDL2Driver) GenCodeImage() {
	tmpSurface, err := sdl.CreateRGBSurface(0, int32(Xadjust), int32(10*7), 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	for pc, runes := range S.codeList {
		for i, r := range runes {
			err = S.bitmap.Blit(getGlyph(r), tmpSurface, &sdl.Rect{int32(i * 7), int32((pc-0xD000) * 9), 7, 9})
			if err != nil {
				panic(err)
			}
		}
	}

	err = tmpSurface.SaveBMP("code.bmp")
	if err != nil {
		panic(err)
	}
	tmpSurface.Free()
}

func (S *SDL2Driver) ShowCode(pc int) {
	S.debug.FillRect(&sdl.Rect{0, 0, Xadjust, int32(S.winHeight)}, 16)

	for i := -10; i < 10; i++ {
		for ri, r := range S.codeList[pc+i] {
			S.bitmap.Blit(getGlyph(r), S.debug, &sdl.Rect{int32(ri * 7), int32(i+10) * 9, 7, 9})
		}
	}
}

func (S *SDL2Driver) UpdateFrame() {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(1)
		case *sdl.KeyboardEvent:
			switch t.Type {
			case sdl.KEYDOWN:
				S.keybLine.KeyCode = uint(t.Keysym.Sym)
				S.keybLine.Mode = 0
				switch t.Keysym.Mod {
				case 1:
					if S.keybLine.KeyCode != sdl.K_LSHIFT {
						S.keybLine.Mode = sdl.K_LSHIFT
					}
				case 2:
					if S.keybLine.KeyCode != sdl.K_RSHIFT {
						S.keybLine.Mode = sdl.K_RSHIFT
					}
				case 64:
					if S.keybLine.KeyCode != sdl.K_LCTRL {
						S.keybLine.Mode = sdl.K_LCTRL
					}
				default:
					S.keybLine.Mode = S.keybLine.KeyCode
				}
				log.Printf("KEY DOWN : %d - %d %d", t.Keysym.Mod, S.keybLine.KeyCode, S.keybLine.Mode)
			case sdl.KEYUP:
				// *S.keybLine = 1073742049
				S.keybLine.KeyCode = 0
				S.keybLine.Mode = 0
			}
		default:
			// buffer = 0
		}
	}

	// timerFPS = sdl.GetTicks() - lastFrame
	// if timerFPS < (1000 / setFPS) {
	// 	sdl.Delay((1000 / setFPS) - timerFPS)
	// 	// return
	// }
	frameCount++
	// S.renderer.Clear()
	// S.texture.Update(nil, S.screen, S.winWidth*3)
	// S.renderer.Copy(S.texture, nil, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	S.getFPS()

	S.texture, _ = S.renderer.CreateTextureFromSurface(S.surface)
	S.renderer.Copy(S.texture, nil, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	S.texture, _ = S.renderer.CreateTextureFromSurface(S.debug)
	S.renderer.Copy(S.texture, nil, &sdl.Rect{0, 0, int32(Xadjust), int32(S.winHeight)})
	S.renderer.Present()

	// S.window.UpdateSurface()
}

func (S *SDL2Driver) IOEvents() *KEYPressed {
	return S.keybLine
}
