package main

import (
	"fmt"
	"log"
	"newApple/config"
	"newApple/crtc"
	"newApple/graphic"
	"newApple/mem"
	"newApple/mos6510"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/mattn/go-tty"
)

const (
	ramSize      = 65536
	romSize      = 8192
	ioSize       = 4096
	chargenSize  = 4096
	keyboardSize = 2048
	blanckSize   = 12288

	nbMemLayout = 1

	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &config.ConfigData{}

	cpu mos6510.CPU

	RAM      []byte
	ROM_CD   []byte
	ROM_EF   []byte
	IO       []byte
	KEYB     []byte
	CHARGEN  []byte
	BLANK    []byte
	MEM      mem.BANK
	IOAccess mem.MEMAccess

	outputDriver graphic.Driver
	CRTC         crtc.CRTC
	cpuTurn      bool
	run          bool
	execInst     sync.Mutex
)

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

func setup() {
	// ROMs & RAM Setup
	RAM = make([]byte, ramSize)
	IO = make([]byte, ioSize)
	BLANK = make([]byte, blanckSize)
	ROM_CD = mem.LoadROM(romSize, "assets/roms/CD.bin")
	ROM_EF = mem.LoadROM(romSize, "assets/roms/EF.bin")
	KEYB = mem.LoadROM(keyboardSize, "assets/roms/Keyb.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/Video_US.bin")

	mem.Clear(RAM)
	// mem.DisplayCharRom(CHARGEN, 1, 8, 16)

	// RAM[0x0001] = 0x00
	// MEM = mem.InitBanks(nbMemLayout, &RAM[0x0001])
	var test byte = 0
	MEM = mem.InitBanks(nbMemLayout, &test)
	IOAccess = &accessor{}

	// MEM Setup

	memLayouts()

	outputDriver = &graphic.SDLDriver{}
	CRTC.Init(RAM, IO, CHARGEN, outputDriver, conf)

	// CPU Setup
	cpu.Init(&MEM, conf)
}

func input() {
	dumpAddr := ""
	var keyb *tty.TTY
	keyb, _ = tty.Open()

	for {
		r, _ := keyb.ReadRune()
		switch r {
		case 's':
			Disassamble()
			MEM.DumpStack(cpu.SP)
		case 'z':
			Disassamble()
			MEM.Dump(0)
		case 'x':
			// DumpMem(&pla, "memDump.bin")
		case 'r':
			conf.Disassamble = false
			run = true
			execInst.Unlock()
		case 'l':
			// LoadPRG(&pla, "./prg/GARDEN.prg")
			LoadPRG(&MEM, conf.LoadPRG)
			// addr, _ := LoadPRG(mem.Val, conf.LoadPRG)
			// cpu.GoTo(addr)
		case ' ':
			if run {
				conf.Disassamble = true
				run = false
			} else {
				execInst.Unlock()
			}
			// fmt.Printf("\n(s) Stack Dump - (z) Zero Page - (r) Run - (sp) Pause / unpause > ")
		case 'w':
			fmt.Printf("\nFill Color RAM")
			for i := 0xD800; i < 0xDC00; i++ {
				MEM.Write(uint16(i), 0)
			}
			// for i := 0x0800; i < 0x0C00; i++ {
			// 	IO[uint16(i)] = 0
			// }
		case 'q':
			cpu.DumpStats()
			os.Exit(0)
		default:
			dumpAddr += string(r)
			fmt.Printf("%c", r)
			if len(dumpAddr) == 4 {
				hx, _ := strconv.ParseInt(dumpAddr, 16, 64)
				MEM.Dump(uint16(hx))
				dumpAddr = ""
			}
		}

	}
}

func Disassamble() {
	// fmt.Printf("\n%s %s", vic.Disassemble(), cpu.Disassemble())
	fmt.Printf("%s\n", cpu.Disassemble())
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Now().Sub(start)
	log.Printf("%s took %s", name, elapsed)
}

func RunEmulation() {
	// defer timeTrack(time.Now(), "RunEmulation")
	CRTC.Run(!run)
	if cpu.State == mos6510.ReadInstruction && !run {
		execInst.Lock()
	}

	cpu.NextCycle()
	if cpu.State == mos6510.ReadInstruction {
		if conf.Breakpoint == cpu.InstStart {
			conf.Disassamble = true
			run = false
		}
	}

	if cpu.State == mos6510.ReadInstruction {
		if !run || conf.Disassamble {
			Disassamble()
		}
	}
}

func main() {
	// var exit chan bool
	// exit = make(chan bool)

	confload.Load("config.ini", conf)

	clog.LogLevel = conf.LogLevel
	clog.StartLogging = conf.StartLogging
	if conf.FileLog != "" {
		clog.EnableFileLog(conf.FileLog)
	}

	// f, err := os.Create("newC64.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	setup()
	go input()

	run = true
	cpuTurn = true
	// go func() {
	for {
		RunEmulation()
	}
	// }()

	// outputDriver.Run()

	// cpu.DumpStats()
	// <-exit
}
