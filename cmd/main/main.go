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
	"runtime"
	"strconv"

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

	RAM   []byte
	BANK1 []byte
	BANK2 []byte

	ZP        []byte
	ALT_ZP    []byte
	AUX       []byte
	AUX_BANK1 []byte
	AUX_BANK2 []byte

	ROM_C  []byte
	ROM_D  []byte
	ROM_EF []byte

	IO      []byte
	SLOTS   [8][]byte
	CHARGEN []byte

	MEM      mem.BANK
	IOAccess mem.MEMAccess

	InputLine    render.KEYPressed
	outputDriver render.SDL2Driver
	CRTC         crtc.CRTC
	run          bool
	trace        bool
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

func loadSlots() {
	conf.Slots.Catalog[1] = conf.Slots.Slot1
	conf.Slots.Catalog[2] = conf.Slots.Slot2
	conf.Slots.Catalog[3] = conf.Slots.Slot3
	conf.Slots.Catalog[4] = conf.Slots.Slot4
	conf.Slots.Catalog[5] = conf.Slots.Slot5
	conf.Slots.Catalog[6] = conf.Slots.Slot6
	conf.Slots.Catalog[7] = conf.Slots.Slot7

	for i := 1; i < 8; i++ {
		if conf.Slots.Catalog[i] != "" {
			SLOTS[i] = mem.LoadROM(slot_roms, conf.Slots.Catalog[i])
		} else {
			SLOTS[i] = make([]byte, slot_roms)
			mem.Clear(SLOTS[i], 0, 0x71)
		}
	}
}

func loadDisks() (*disk.DRIVE, *disk.DRIVE) {
	var dsk1, dsk2 *disk.DRIVE

	dsk1 = nil
	dsk2 = nil
	if conf.Slots.Slot6 != "" {
		if conf.Disks.Disk1 != "" {
			dsk1 = disk.Attach(&cpu)
			dsk1.LoadDiskImage(conf.Disks.Disk1)
		}
		if conf.Disks.Disk2 != "" {
			dsk2 = disk.Attach(&cpu)
			dsk2.LoadDiskImage(conf.Disks.Disk2)
		}
		if dsk1 == nil && dsk2 == nil {
			conf.Slots.Slot6 = ""
		}
	}
	return dsk1, dsk2
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

	ZP = make([]byte, 0x0200)
	mem.Clear(ZP, 0x1000, 0xFF)
	ALT_ZP = make([]byte, 0x0200)
	mem.Clear(ALT_ZP, 0x1000, 0xFF)

	AUX = make([]byte, ramSize)
	mem.Clear(RAM, 0x1000, 0xFF)
	AUX_BANK1 = make([]byte, romSize)
	mem.Clear(BANK1, 0x1000, 0xFF)
	AUX_BANK2 = make([]byte, romSize*3)

	mem.Clear(BANK2, 0x1000, 0xFF)
	IO = make([]byte, softSwitches)
	mem.Clear(IO, 0, 0x00)

	Disk1, Disk2 := loadDisks()
	loadSlots()

	IOAccess = InitIO(Disk1, Disk2, &CRTC)

	// Disk1.DumpTrack(1)
	// Disk1.ReadTrackRaw(0, 53404)
	// Disk1.Dump(true)

	// panic(1)
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
	CRTC.Init(RAM, AUX, IO, CHARGEN, &outputDriver, conf)
	outputDriver.SetKeyboardLine(&InputLine)

	// CPU Setup
	cpu.Init(conf.CPUModel, conf.Mhz, &MEM, conf.Debug || conf.Disassamble)

	MEM.CheckLayoutForAddr(0x0020)
	MEM.CheckLayoutForAddr(0xC610)
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
			run = true
			trace = false
		case 'l':
		case ' ':
			fmt.Printf("%s\n", cpu.FullDebug)
			trace = true
			run = true
		case 'w':
			fmt.Printf("\nFill Screen")
			cpt := 0
			for i := 0x0400; i < 0x0800; i++ {
				RAM[uint16(i)] = byte(cpt)
				AUX[uint16(i)] = byte(cpt)
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
			fmt.Printf("%s\n", cpu.FullDebug)
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

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Now().Sub(start)
// 	log.Printf("%s took %s", name, elapsed)
// }

func RunEmulation() {
	var key byte
	var speed float64

	// defer timeTrack(time.Now(), "RunEmulation")
	for {
		CRTC.Run(!run)

		if run {
			if InputLine.KeyCode != 0 && !is_Keypressed {
				key = keyMap[InputLine.KeyCode][InputLine.Mode]
				// log.Printf("KEY DOWN - Code: %d  Mode: %d  -> %d", InputLine.KeyCode, InputLine.Mode, key)
				IO[0] = key | 0b10000000
				is_Keypressed = true
				InputLine.KeyCode = 0
				InputLine.Mode = 0
			}

			speed = cpu.NextCycle()
		}

		if cpu.CycleCount == 1 {
			if trace {
				run = false
			}
			// if cpu.InstStart > 0x0300 && cpu.InstStart < 0xC000 {
			// 	clog.FileRaw("\n%d: %s", cpu.Cycles, cpu.FullDebug)
			// }
			outputDriver.DumpCode(cpu.FullInst)
			outputDriver.SetSpeed(speed)
			if conf.Breakpoint == cpu.InstStart {
				fmt.Printf("%s\n", cpu.FullDebug)
				trace = true
			}
		}
	}
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
	trace = false
	outputDriver.ShowCode = false
	outputDriver.ShowFps = true

	go RunEmulation()
	// }()

	outputDriver.Run(true)

	// cpu.DumpStats()
	// <-exit
}
