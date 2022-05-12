package main

import (
	"io/ioutil"
	"newApple/config"
	"newApple/crtc"
	"runtime"

	"github.com/Djoulzy/emutools/mem"
	"github.com/Djoulzy/emutools/render"
)

const (
	ramSize     = 65536
	chargenSize = 2048
	ioSize      = 4096
	screenStart = 0x0400
	screenSize  = 1024
)

var (
	conf             config.ConfigData
	RAM, IO, CHARGEN []byte
	CRTC             crtc.CRTC
	outputDriver     render.SDL2Driver
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func LoadData(mem []byte, file string, memStart uint16) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	for i, val := range content {
		mem[memStart+uint16(i)] = val
	}
	return nil
}

func start() {
	conf.Disassamble = false

	RAM = make([]byte, ramSize)
	mem.Clear(RAM, 0x1000, 0xFF)
	IO = make([]byte, ioSize)
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/II/3410036.bin")
	CRTC.Init(RAM, IO, CHARGEN, &outputDriver, &conf)

	cpt := 0
	for i := screenStart; i < screenStart+screenSize; i++ {
		RAM[uint16(i)] = byte(cpt)
		cpt++
	}
}

func main() {
	start()
	go CRTC.Run(false)
	outputDriver.Run()
}
