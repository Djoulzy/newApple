package graphic

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var bitmap *sdl.Surface

func buildNumeric() {
	numeric, err := sdl.CreateRGBSurface(0, int32(fontWidth*256), int32(fontHeight), 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	for num := 0x00; num < 0x100; num++ {
		val := []rune(fmt.Sprintf("%02X", num))
		for i, r := range val {
			bitmap.Blit(getGlyph(r), numeric, &sdl.Rect{int32(num*fontWidth*2 + i*fontWidth), 0, fontWidth, fontHeight})
		}
	}
	err = numeric.SaveBMP("graphic/assets/num.bmp")
	if err != nil {
		panic(err)
	}
}

func buildMnemonic(inst []string) {
	mne, err := sdl.CreateRGBSurface(0, int32(mnemonicWidth*16), int32(mnemonicHeight*16), 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			val := []rune(fmt.Sprintf("%s", inst[y*16+x]))
			for i, r := range val {
				bitmap.Blit(getGlyph(r), mne, &sdl.Rect{int32((x*mnemonicWidth) +(i*fontWidth)), int32(y*mnemonicHeight), mnemonicWidth, mnemonicHeight})
			}
		}
	}
	err = mne.SaveBMP("graphic/assets/mnemonic.bmp")
	if err != nil {
		panic(err)
	}
}

func MakeAsset(inst []string) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	bitmap, err = sdl.LoadBMP("graphic/assets/white.bmp")
	if err != nil {
		panic(err)
	}

	buildNumeric()
	buildMnemonic(inst)
}
