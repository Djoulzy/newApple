package main

import (
	"bytes"
	"fmt"
	"log"
	"newApple/config"
	"newApple/crtc"
	"newApple/disk"
	woz "newApple/goWoz"
	"os"
	"reflect"
	"runtime"
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
	romSize      = 4096
	softSwitches = 256
	chargenSize  = 2048
	keyboardSize = 2048
	blanckSize   = 12288
	slot_roms    = 256

	nbMemLayout = 1

	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &config.ConfigData{}

	cpu     mos6510.CPU
	MODEL   int
	BankSel byte

	RAM    []byte
	BANK1  []byte
	BANK2  []byte
	ROM_C  []byte
	ROM_D  []byte
	ROM_EF []byte

	IO      []byte
	SLOT1   []byte
	SLOT2   []byte
	SLOT3   []byte
	SLOT4   []byte
	SLOT5   []byte
	SLOT6   []byte
	SLOT7   []byte
	KEYB    []byte
	CHARGEN []byte

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

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func apple2_Roms() {
	ROM_D = mem.LoadROM(romSize, "assets/roms/II/D.bin")
	ROM_EF = mem.LoadROM(romSize*2, "assets/roms/II/EF.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/II/3410036.bin")
}

func apple2e_Roms() {
	ROM_C = mem.LoadROM(romSize, "assets/roms/IIe/C.bin")
	ROM_D = mem.LoadROM(romSize, "assets/roms/IIe/D.bin")
	ROM_EF = mem.LoadROM(romSize*2, "assets/roms/IIe/EF.bin")
	CHARGEN = mem.LoadROM(chargenSize*2, "assets/roms/IIe/Video_US.bin")
}

func setup() {
	BankSel = 0
	MEM = mem.InitBanks(nbMemLayout, &BankSel)

	// Common Setup
	RAM = make([]byte, ramSize)
	mem.Clear(RAM, 0x1000, 0xFF)
	BANK1 = make([]byte, romSize)
	mem.Clear(BANK1, 0x1000, 0xFF)
	BANK2 = make([]byte, romSize*3)
	mem.Clear(BANK2, 0x1000, 0xFF)
	IO = make([]byte, softSwitches)
	mem.Clear(IO, 0, 0x00)

	SLOT1 = make([]byte, slot_roms)
	mem.Clear(SLOT1, 0, 0x71)
	SLOT2 = make([]byte, slot_roms)
	mem.Clear(SLOT2, 0, 0x71)
	SLOT3 = make([]byte, slot_roms)
	mem.Clear(SLOT3, 0, 0x71)
	SLOT4 = make([]byte, slot_roms)
	mem.Clear(SLOT4, 0, 0x71)
	SLOT5 = make([]byte, slot_roms)
	mem.Clear(SLOT5, 0, 0x71)
	SLOT6 = mem.LoadROM(slot_roms, "assets/roms/16SectorP5.bin")
	SLOT7 = make([]byte, slot_roms)
	mem.Clear(SLOT7, 0, 0x71)

	woz.SetupLib()

	Disk1 := disk.Attach(&cpu)
	Disk2 := disk.Attach(&cpu)

	Disk1.LoadDiskImage("imgTest/DOS33.woz")
	// Disk1.LoadDiskImage("imgTest/demo.woz")
	// Disk1.LoadDiskImage("imgTest/Locksmith.woz")
	// Disk1.LoadDiskImage("imgTest/Wolf.woz")
	// Disk1.LoadDiskImage("imgTest/HERO.woz")

	// Disk1.LoadDiskImage("imgTest/Choplifter.woz")
	// DiskDrive.LoadDiskImage("imgTest/POP_A.woz")
	// DiskDrive.LoadDiskImage("imgTest/Karateka.woz")

	// Disk1.LoadDiskImage("imgTest/anti-m.woz")
	// Disk2.LoadDiskImage("imgTest/Choplifter.woz")

	// Disk1.LoadDiskImage("imgTest/Wizardry_boot.woz")
	// Disk1.LoadDiskImage("imgTest/CompInsp.woz")

	IOAccess = &io_access{Disks: [2]*disk.DRIVE{Disk1, Disk2}, Video: &CRTC}

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
	cpu.Init(conf.Mhz, &MEM, conf.Debug || conf.Disassamble)
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
			fmt.Printf("Bank: %d\n", BankSel)
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
		case 'k':
			CRTC.ToggleMonitorColor()
		case 'p':
			if crtc.Is_PAGE2 {
				crtc.Is_PAGE2 = false
				log.Println("PAGE 1")
			} else {
				crtc.Is_PAGE2 = true
				log.Println("PAGE 2")
			}
			CRTC.UpdateVideoRam()
			CRTC.DumpMode()
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
	var speed float64

	// defer timeTrack(time.Now(), "RunEmulation")
	for {
		CRTC.Run(!run)
		if cpu.CycleCount == 1 && !run {
			execInst.Lock()
		}

		if InputLine.KeyCode != 0 && !is_Keypressed {
			key = keyMap[InputLine.KeyCode]
			if InputLine.Mode == 1073742048 {
				key -= 0x60
			}
			IO[0] = key | 0b10000000
			is_Keypressed = true
			InputLine.KeyCode = 0
			InputLine.Mode = 0
		}

		speed = cpu.NextCycle()

		// if cpu.State == mos6510.ReadInstruction {
		// 	outputDriver.DumpCode(cpu.FullInst)
		// 	if conf.Breakpoint == cpu.InstStart {
		// 		conf.Disassamble = true
		// 		run = false
		// 	}
		// }

		if cpu.CycleCount == 1 {
			outputDriver.DumpCode(cpu.FullInst)
			outputDriver.SetSpeed(speed)
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
	outputDriver.ShowCode = false
	outputDriver.ShowFps = true

	go RunEmulation()
	// }()

	outputDriver.Run(true)

	// cpu.DumpStats()
	// <-exit
}