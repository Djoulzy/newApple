package graphic

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GGDriver struct {
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

	test *image.RGBA
}

func (S *GGDriver) DrawPixel(x, y int, c color.Color) {
	S.test.Set(x, y, c)
}

func (S *GGDriver) CloseAll() {
	S.window.Destroy()
	sdl.Quit()
}

func (S *GGDriver) Init(winWidth, winHeight int, title string) {
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

	S.debug[0], err = sdl.CreateRGBSurface(0, int32(Xadjust/2), int32(S.emuHeight), 32, 0, 0, 0, 0)
	err = S.debug[0].SetRLE(true)
	if err != nil {
		panic(err)
	}
	S.debug[1], err = sdl.CreateRGBSurface(0, int32(Xadjust/2), int32(S.emuHeight), 32, 0, 0, 0, 0)
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

	S.test = image.NewRGBA(image.Rect(0, 0, S.emuWidth, S.emuHeight))
}

func (S *GGDriver) SetKeyboardLine(line *KEYPressed) {
	S.keybLine = line
}

func (S *GGDriver) SetCodeList(list [][]byte) {
	S.codeList = list
}

func (S *GGDriver) throttleFPS(showFps bool) {
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
		// runes := []rune(fmt.Sprintf("%d", fps))
		// for i, r := range runes {
		// 	S.bitmap.Blit(getGlyph(r), S.emul, &sdl.Rect{int32(S.emuWidth - 21 + i*7), 2, 7, 9})
		// }
		fmt.Printf("%d\n",fps)
		frameCount++
	}
}

func (S *GGDriver) ShowCode(pc_done uint16, inst string) {
	// var debugHide int

	// debugHide = int(math.Abs(float64(S.debugShow - 1)))
	// S.debug[debugHide].FillRect(&sdl.Rect{0, 0, Xadjust / 2, int32(S.emuHeight)}, 16)
	// S.debug[S.debugShow].Blit(&sdl.Rect{0, fontHeight, Xadjust / 2, int32(S.emuHeight - fontHeight)}, S.debug[debugHide], nil)

	// // surf, _ := S.font.RenderUTF8Solid(fmt.Sprintf("%04X: %s", pc_done, inst), sdl.Color(color.RGBA{R: 255, G: 255, B: 255, A: 255}))
	// // surf.Blit(nil, S.debug[debugHide], &sdl.Rect{5, int32(S.winHeight - fontHeight*2), mnemonicWidth, 8})
	// // surf.Free()
	// y := int32(S.emuHeight - fontHeight*2)
	// runes := []rune(inst)
	// for i, r := range runes {
	// 	S.bitmap.Blit(getGlyph(r), S.debug[debugHide], &sdl.Rect{int32(5 + (i * 7)), y, 7, 9})
	// }

	// S.debugShow = debugHide
}

func (S *GGDriver) UpdateFrame() {
	S.throttleFPS(true)

	// S.texture, _ = S.renderer.CreateTextureFromSurface(S.emul)
	// S.renderer.Copy(S.texture, nil, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	// S.texture, _ = S.renderer.CreateTextureFromSurface(S.debug[S.debugShow])
	// S.renderer.Copy(S.texture, nil, &sdl.Rect{0, 0, int32(Xadjust), int32(S.winHeight)})
	// S.renderer.Present()

	// S.emul.BlitScaled(nil, S.w_surf, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	tmp, _ := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&S.test.Pix[0]), int32(S.emuWidth), int32(S.emuHeight), 32, 4*S.emuWidth, 0,0,0,0)
	tmp.BlitScaled(nil, S.w_surf, &sdl.Rect{Xadjust, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})

	S.debug[S.debugShow].BlitScaled(nil, S.w_surf, &sdl.Rect{0, 0, int32(Xadjust), int32(S.emuHeight) * 2})
	S.window.UpdateSurface()
}

func (S *GGDriver) Run() {
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

func (S *GGDriver) IOEvents() *KEYPressed {
	return S.keybLine
}
