package main

import (
	"bytes"
	"fmt"
	"log"
	"newApple/config"
	"newApple/crtc"
	"newApple/disk"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/Djoulzy/emutools/mos6510"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/Djoulzy/emutools/mem"
	"github.com/Djoulzy/emutools/render"
	"github.com/mattn/go-tty"
)

const (
	ramSize      = 65536
	romSize      = 2048
	softSwitches = 256
	chargenSize  = 2048
	keyboardSize = 2048
	blanckSize   = 12288
	slot_roms    = 256

	nbMemLayout = 4

	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &config.ConfigData{}

	cpu     mos6510.CPU
	MODEL   int
	BankSel byte

	RAM     []byte
	ROM_AID []byte
	ROM_D0  []byte
	ROM_D8  []byte
	ROM_E0  []byte
	ROM_E8  []byte
	ROM_F0  []byte
	ROM_F8  []byte

	ROM_CD []byte
	ROM_EF []byte

	IO       []byte
	SLOT1    []byte
	SLOT2    []byte
	SLOT3    []byte
	SLOT4    []byte
	SLOT5    []byte
	SLOT6    []byte
	SLOT7    []byte
	KEYB     []byte
	CHARGEN  []byte
	BLANK    []byte
	MEM      mem.BANK
	IOAccess mem.MEMAccess

	InputLine    render.KEYPressed
	outputDriver render.SDL2Driver
	CRTC         crtc.CRTC
	cpuTurn      bool
	run          bool
	execInst     sync.Mutex
	lastPC       uint16
)

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

func apple2_Roms() {
	ROM_D0 = mem.LoadROM(romSize, "assets/roms/II/3410011D0.bin")
	ROM_D8 = mem.LoadROM(romSize, "assets/roms/II/3410012D8.bin")
	ROM_E0 = mem.LoadROM(romSize, "assets/roms/II/3410013E0.bin")
	ROM_E8 = mem.LoadROM(romSize, "assets/roms/II/3410014E8.bin")
	ROM_F0 = mem.LoadROM(romSize, "assets/roms/II/3410015F0.bin")
	ROM_F8 = mem.LoadROM(romSize, "assets/roms/II/3410020F8.bin")
	// ROM_AID = mem.LoadROM(romSize, "assets/roms/II/3410016.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/II/3410036.bin")
}

func apple2e_Roms() {
	ROM_CD = mem.LoadROM(romSize*4, "assets/roms/IIe/CD.bin")
	ROM_EF = mem.LoadROM(romSize*4, "assets/roms/IIe/EF.bin")
	CHARGEN = mem.LoadROM(chargenSize*2, "assets/roms/IIe/Video_US.bin")
}

func setup() {
	BankSel = 1
	MEM = mem.InitBanks(nbMemLayout, &BankSel)

	// Common Setup
	RAM = make([]byte, ramSize)
	mem.Clear(RAM, 0x1000, 0xFF)
	IO = make([]byte, softSwitches)
	mem.Clear(IO, 0, 0x00)

	SLOT1 = make([]byte, slot_roms)
	SLOT2 = make([]byte, slot_roms)
	SLOT3 = make([]byte, slot_roms)
	SLOT4 = make([]byte, slot_roms)
	SLOT5 = make([]byte, slot_roms)
	SLOT6 = mem.LoadROM(slot_roms, "assets/roms/slot_disk2_cx00.bin")
	SLOT7 = make([]byte, slot_roms)
	DiskDrive := disk.Attach()
	DiskDrive.LoadDiskImage("woz/DOS33.woz")

	IOAccess = &io_access{Disk: DiskDrive}

	if MODEL == 1 {
		apple2_Roms()
	} else {
		apple2e_Roms()
	}

	// mem.DisplayCharRom(CHARGEN, 1, 8, 16)

	// MEM Setup

	memLayouts(MODEL)

	outputDriver = render.SDL2Driver{}
	initKeyboard()
	CRTC.Init(RAM, IO, CHARGEN, &outputDriver, conf)
	outputDriver.SetKeyboardLine(&InputLine)

	// CPU Setup
	cpu.Init(&MEM)
}

func input() {
	dumpAddr := ""
	var keyb *tty.TTY
	keyb, _ = tty.Open()

	for {
		r, _ := keyb.ReadRune()
		switch r {
		case 's':
			MEM.DumpStack(cpu.SP)
			cpu.DumpStackDebug()
		case 'z':
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
			for i := 0x0400; i < 0x0800; i++ {
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
				fmt.Printf("\n")
				MEM.Dump(uint16(hx))
				dumpAddr = ""
			}
		}

	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Now().Sub(start)
	log.Printf("%s took %s", name, elapsed)
}

func RunEmulation() {
	var key byte
	// defer timeTrack(time.Now(), "RunEmulation")
	for {
		CRTC.Run(!run)
		if cpu.CycleCount == 1 && !run {
			execInst.Lock()
		}

		if InputLine.KeyCode != 0 {
			key = keyMap[InputLine.KeyCode]
			if InputLine.Mode == 1073742048 {
				key -= 0x40
			}
			MEM.Write(0xC000, key|0b10000000)
			InputLine.KeyCode = 0
			InputLine.Mode = 0
		}

		cpu.NextCycle()
		// if cpu.State == mos6510.ReadInstruction {
		// 	outputDriver.DumpCode(cpu.FullInst)
		// 	if conf.Breakpoint == cpu.InstStart {
		// 		conf.Disassamble = true
		// 		run = false
		// 	}
		// }

		if cpu.CycleCount == 1 {
			// outputDriver.DumpCode(cpu.FullInst)
			if conf.Breakpoint == cpu.InstStart {
				conf.Disassamble = true
				run = false
			}
			if !run || conf.Disassamble {
				fmt.Printf("%s\n", cpu.FullDebug)
			}
		}
	}
}

func packageName(v interface{}) string {
	if v == nil {
		return ""
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return val.Elem().Type().PkgPath()
	}
	return val.Type().PkgPath()
}

func main() {
	var b bytes.Buffer
	fmt.Println(reflect.TypeOf(b).PkgPath())
	// var exit chan bool
	// exit = make(chan bool)

	confload.Load("config.ini", conf)

	clog.LogLevel = conf.LogLevel
	clog.StartLogging = conf.StartLogging
	if conf.FileLog != "" {
		clog.EnableFileLog(conf.FileLog)
	}
	switch conf.Model {
	case "2":
		MODEL = 1
	case "2e":
		fallthrough
	default:
		MODEL = 2
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
	// outputDriver.ShowCode = true
	outputDriver.ShowFps = true

	go RunEmulation()
	// }()

	outputDriver.Run()

	// cpu.DumpStats()
	// <-exit
}
