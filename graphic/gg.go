package graphic

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"github.com/golang/freetype"
	"github.com/veandco/go-sdl2/sdl"
)

type GGDriver struct {
	winHeight    int
	winWidth     int
	emuHeight    int
	emuWidth     int
	window       *sdl.Window
	w_surf       *sdl.Surface
	emul         *image.RGBA
	emul_s       *sdl.Surface
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	keybLine     *KEYPressed
	codeList     []string
	nextCodeLine int

	font         *freetype.Context
	Update       chan bool
	debugBGColor *color.RGBA
}

func (S *GGDriver) DrawPixel(x, y int, c color.Color) {
	S.emul.Set(x+Xadjust, y, c)
}

func (S *GGDriver) CloseAll() {
	S.window.Destroy()
	sdl.Quit()
}

func (S *GGDriver) Init(width, height int, title string) {
	S.emuHeight = height
	S.emuWidth = width + Xadjust
	S.winHeight = S.emuHeight * 2
	S.winWidth = S.emuWidth * 2

	S.codeList = make([]string, nbCodeLines)
	S.nextCodeLine = 0
	S.Update = make(chan bool)

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	S.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(S.winWidth), int32(S.winHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	S.renderer, err = sdl.CreateRenderer(S.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	S.w_surf, err = S.window.GetSurface()
	S.w_surf.SetRLE(true)

	S.emul = image.NewRGBA(image.Rect(0, 0, S.emuWidth, S.emuHeight))
	S.emul_s, _ = sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&S.emul.Pix[0]), int32(S.emuWidth), int32(S.emuHeight), 32, 4*S.emuWidth, 0, 0, 0, 0)
	S.emul_s.SetRLE(true)

	fontBytes, err := ioutil.ReadFile("assets/ttf/PetMe.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	fg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	S.font = freetype.NewContext()
	S.font.SetDPI(72)
	S.font.SetFont(f)
	S.font.SetFontSize(fontWidth)
	S.font.SetClip(S.emul.Bounds())
	S.font.SetDst(S.emul)
	S.font.SetSrc(fg)

	S.debugBGColor = &color.RGBA{50, 50, 50, 255}
}

func (S *GGDriver) SetKeyboardLine(line *KEYPressed) {
	S.keybLine = line
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
		pt := freetype.Pt((S.emuWidth - fontWidth*3), fontHeight)
		S.font.DrawString(fmt.Sprintf("%d", fps), pt)
	}
}

func (S *GGDriver) DumpCode(inst string) {
	S.codeList[S.nextCodeLine] = inst
	S.nextCodeLine++
	if S.nextCodeLine == nbCodeLines {
		S.nextCodeLine = 0
	}
}

func (S *GGDriver) ShowCode() {
	b := image.Rect(0, 0, Xadjust, S.emuHeight)
	draw.Draw(S.emul, b, &image.Uniform{S.debugBGColor}, image.ZP, draw.Src)
	base := (S.emuHeight - fontHeight)
	cpt := S.nextCodeLine - 1
	for i := 0; i < nbCodeLines; i++ {
		if cpt < 0 {
			cpt = nbCodeLines - 1
		}
		pt := freetype.Pt(0, base-fontHeight*i)
		S.font.DrawString(fmt.Sprintf("%s\n", S.codeList[cpt]), pt)
		cpt--
	}
}

func (S *GGDriver) UpdateFrame() {
	rect := sdl.Rect{0, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2}

	S.throttleFPS(true)
	S.ShowCode()

	// SDL2 Texture + Render
	S.texture, _ = S.renderer.CreateTextureFromSurface(S.emul_s)
	S.renderer.Copy(S.texture, nil, &rect)
	S.renderer.Present()

	// SDL2 Surface
	// S.emul_s.BlitScaled(nil, S.w_surf, &sdl.Rect{0, 0, int32(S.emuWidth) * 2, int32(S.emuHeight) * 2})
	// S.window.UpdateSurface()

	frameCount++

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
		sdl.Delay(10)
	}
}

func (S *GGDriver) IOEvents() *KEYPressed {
	return S.keybLine
}
