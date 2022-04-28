package graphic

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SDL2Driver struct {
	winHeight int
	winWidth  int
	emuHeight int
	emuWidth  int
	window    *sdl.Window
	w_surf    *sdl.Surface
	emul      *sdl.Surface
	debug     [2]*sdl.Surface
	debugShow int
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	bitmap    *sdl.Surface
	mnemo     *sdl.Surface
	num       *sdl.Surface
	cmd       *sdl.Surface
	keybLine  *KEYPressed
	codeList  [][]byte

	LastPC   uint16
	dumpCode chan bool

	font *ttf.Font
}

func (S *SDL2Driver) DirectDrawPixel(x, y int, c color.Color) {
	// S.renderer.SetDrawColor(byte(color.R), byte(color.G), byte(color.B), 255)
	// S.renderer.DrawPoint(int32(x), int32(y))
	S.emul.Set(x, y, c)
}

func (S *SDL2Driver) DrawPixel(x, y int, c color.Color) {
	S.emul.Set(x, y, c)
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
	S.dumpCode = make(chan bool)

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	// S.window, S.renderer, err = sdl.CreateWindowAndRenderer(int32(S.winWidth*2), int32(S.winHeight*2), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	// S.window.SetTitle(title)
	S.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(S.winWidth), int32(S.winHeight), sdl.WINDOW_SHOWN)
	// S.renderer, err = sdl.CreateRenderer(S.window, -1, sdl.RENDERER_ACCELERATED)
	// if err != nil {
	// 	panic(err)
	// }
	S.w_surf, err = S.window.GetSurface()
	S.emul, err = sdl.CreateRGBSurface(0, int32(S.emuWidth), int32(S.emuHeight), 32, 0, 0, 0, 0)
	S.emul.SetRLE(true)
	if err != nil {
		panic(err)
	}

	S.debug[0], err = sdl.CreateRGBSurface(0, int32(Xadjust), int32(S.winHeight), 32, 0, 0, 0, 0)
	err = S.debug[0].SetRLE(true)
	if err != nil {
		panic(err)
	}
	S.debug[1], err = sdl.CreateRGBSurface(0, int32(Xadjust), int32(S.winHeight), 32, 0, 0, 0, 0)
	S.debug[1].SetRLE(true)
	if err != nil {
		panic(err)
	}
	S.debugShow = 1

	if err != nil {
		panic(err)
	}

	S.bitmap, err = sdl.LoadBMP("graphic/assets/white.bmp")
	if err != nil {
		panic(err)
	}
	S.mnemo, err = sdl.LoadBMP("graphic/assets/mnemonic.bmp")
	if err != nil {
		panic(err)
	}
	S.num, err = sdl.LoadBMP("graphic/assets/num.bmp")
	if err != nil {
		panic(err)
	}
	S.cmd, err = sdl.LoadBMP("graphic/assets/cmd.bmp")
	if err != nil {
		panic(err)
	}
	ttf.Init()
	S.font, err = ttf.OpenFont("assets/ttf/PetMe.ttf", 8)
	if err != nil {
		panic(err)
	}
}

func (S *SDL2Driver) SetKeyboardLine(line *KEYPressed) {
	S.keybLine = line
}

func (S *SDL2Driver) SetCodeList(list [][]byte) {
	S.codeList = list
}

func getGlyph(char rune) *sdl.Rect {
	pos := int32(char - 32)
	// posy := int32(pos / 18)
	// posx := pos - int32(pos / 18)*18
	// fmt.Printf("r: %c ASCII: %d - abs: %d - x: %d - y: %d\n", char, char, pos, posx, posy)
	return &sdl.Rect{pos*7 - int32(pos/18)*126, int32(pos/18) * 9, 7, 9}
}

func (S *SDL2Driver) throttleFPS(showFps bool) {
	timerFPS = sdl.GetTicks() - lastFrame
	if timerFPS < throttleFPS {
		sdl.Delay(throttleFPS - timerFPS)
	}
	lastFrame = sdl.GetTicks()

	if showFps {
		if lastFrame >= (lastTime + 1000) {
			lastTime = lastFrame
			fps = frameCount
			frameCount = 0
		}
		runes := []rune(fmt.Sprintf("%d", fps))
		for i, r := range runes {
			S.bitmap.Blit(getGlyph(r), S.emul, &sdl.Rect{int32(S.emuWidth - 21 + i*7), 2, 7, 9})
		}
		frameCount++
	}
}

func (S *SDL2Driver) ShowCode(pc_done uint16, inst string) {
	var debugHide int

	debugHide = int(math.Abs(float64(S.debugShow - 1)))
	S.debug[debugHide].FillRect(&sdl.Rect{0, 0, Xadjust, int32(S.winHeight)}, 16)
	S.debug[S.debugShow].Blit(&sdl.Rect{0, fontHeight, Xadjust, int32(S.winHeight - fontHeight)}, S.debug[debugHide], nil)

	surf, _ := S.font.RenderUTF8Solid(fmt.Sprintf("%04X: %s", pc_done, inst), sdl.Color(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
	surf.Blit(nil, S.debug[debugHide], &sdl.Rect{5, int32(S.winHeight - fontHeight*2), mnemonicWidth, 8})
	surf.Free()
	S.debugShow = debugHide
}

func (S *SDL2Driver) UpdateFrame() {
	S.throttleFPS(true)

	// S.texture, _ = S.renderer.CreateTextureFromSurface(S.emul)
	// S.renderer.Copy(S.texture, nil, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	// S.texture, _ = S.renderer.CreateTextureFromSurface(S.debug[S.debugShow])
	// S.renderer.Copy(S.texture, nil, &sdl.Rect{0, 0, int32(Xadjust), int32(S.winHeight)})
	// S.renderer.Present()

	S.emul.BlitScaled(nil, S.w_surf, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	S.debug[S.debugShow].Blit(nil, S.w_surf, nil)
	S.window.UpdateSurface()
}

func (S *SDL2Driver) Run() {
	for {
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
	}
}

func (S *SDL2Driver) IOEvents() *KEYPressed {
	return S.keybLine
}
