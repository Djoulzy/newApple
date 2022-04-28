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
	romSize      = 2048
	ioSize       = 2048
	chargenSize  = 2048
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
	ROM_AID  []byte
	ROM_D0   []byte
	ROM_D8   []byte
	ROM_E0   []byte
	ROM_E8   []byte
	ROM_F0   []byte
	ROM_F8   []byte
	IO       []byte
	KEYB     []byte
	CHARGEN  []byte
	BLANK    []byte
	MEM      mem.BANK
	IOAccess mem.MEMAccess

	InputLine    graphic.KEYPressed
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
	// BLANK = make([]byte, blanckSize)
	ROM_D0 = mem.LoadROM(romSize, "assets/roms/II/3410011D0.bin")
	ROM_D8 = mem.LoadROM(romSize, "assets/roms/II/3410012D8.bin")
	ROM_E0 = mem.LoadROM(romSize, "assets/roms/II/3410013E0.bin")
	ROM_E8 = mem.LoadROM(romSize, "assets/roms/II/3410014E8.bin")
	ROM_F0 = mem.LoadROM(romSize, "assets/roms/II/3410015F0.bin")
	ROM_F8 = mem.LoadROM(romSize, "assets/roms/II/3410020F8.bin")
	ROM_AID = mem.LoadROM(romSize, "assets/roms/II/3410016.bin")
	// KEYB = mem.LoadROM(keyboardSize, "assets/roms/Keyb.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/II/3410036.bin")

	mem.Clear(RAM)
	// mem.DisplayCharRom(CHARGEN, 1, 8, 16)

	// RAM[0x0001] = 0x00
	// MEM = mem.InitBanks(nbMemLayout, &RAM[0x0001])
	var test byte = 0
	MEM = mem.InitBanks(nbMemLayout, &test)
	IOAccess = &accessor{}

	// MEM Setup

	memLayouts()

	outputDriver = &graphic.SDL2Driver{}
	initKeyboard()
	outputDriver.SetKeyboardLine(&InputLine)
	CRTC.Init(RAM, IO, CHARGEN, outputDriver, conf)

	// CPU Setup
	cpu.Init(&MEM, conf)
	outputDriver.SetCodeList(cpu.DumpRom(0xD000))
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
			fmt.Printf("\nFill Screen")
			cpt := 0
			for i := 0x0400; i < 0x0450; i++ {
				MEM.Write(uint16(i), byte(cpt))
				cpt++
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
	fmt.Printf("%s\n", cpu.Trace())
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Now().Sub(start)
	log.Printf("%s took %s", name, elapsed)
}

func RunEmulation() {
	var key byte
	// defer timeTrack(time.Now(), "RunEmulation")
	CRTC.Run(!run)
	if cpu.State == mos6510.ReadInstruction && !run {
		execInst.Lock()
	}

	if MEM.Read(0xC000) == 0 {
		key = keyMap[InputLine.KeyCode]
		if InputLine.Mode == 1073742048 {
			key -= 0x40
		}
		MEM.Write(0xC000, key)
		InputLine.KeyCode = 0
		InputLine.Mode = 0
	}

	cpu.NextCycle()
	if cpu.State == mos6510.ReadInstruction {
		// go outputDriver.ShowCode(&cpu.PC)
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
