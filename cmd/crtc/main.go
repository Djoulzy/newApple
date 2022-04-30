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
	chargenSize = 4096
	ioSize      = 4096
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
	mem.Clear(RAM)
	IO = make([]byte, ioSize)
	mem.Clear(IO)
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/characters-2.901447-10")
	CRTC.Init(RAM, IO, CHARGEN, &outputDriver, &conf)
}

func main() {
	start()
	for {
		CRTC.Run(false)
		// vic.Stats()
	}
}
