package main

import (
	"fmt"
	"newApple/config"
	"newApple/crtc"
	"os"
	"runtime"
	"strconv"
	"time"

	PROC "github.com/Djoulzy/emutools/mos6510"
	"github.com/Djoulzy/mmu"
	"github.com/mattn/go-tty"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/Djoulzy/emutools/render"
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

	cpu       PROC.CPU
	MODEL     int
	LayoutSel byte

	RAM    *mmu.RAM
	AUX    *mmu.RAM
	ROM_D  *mmu.ROM
	ROM_EF *mmu.ROM

	IO      *SoftSwitch
	Disks   *DiskInterface
	SLOTS   [8]*mmu.ROM
	CHARGEN *mmu.ROM

	MEM *mmu.MMU

	InputLine    render.KEYPressed
	outputDriver render.SDL2Driver
	CRTC         crtc.CRTC
	trace        bool
	stepper      bool
	lastPC       uint16
	timeGap      time.Duration // 1Mhz = 1 000 000/s = 1000/ms
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func apple2_Roms() {
	ROM_D = mmu.NewROM("ROM_D", romSize, "assets/roms/II/D.bin")
	MEM.Attach(ROM_D, 0xD0, 16)
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "assets/roms/II/EF.bin")
	MEM.Attach(ROM_EF, 0xE0, 32)
	CHARGEN = mmu.NewROM("CHARGEN", chargenSize, "assets/roms/II/3410036.bin")
	// MEM.Attach(ROM_D, 0xD0, 8)
}

// func apple2e_Roms() {
// 	ROM_C = MEM.LoadROM(romSize, "assets/roms/IIe/C.bin")
// 	ROM_D = MEM.LoadROM(romSize, "assets/roms/IIe/D.bin")
// 	ROM_EF = MEM.LoadROM(romSize*2, "assets/roms/IIe/EF.bin")
// 	CHARGEN = MEM.LoadROM(chargenSize*2, "assets/roms/IIe/Video_US.bin")
// }

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
			// SLOTS[i] = MEM.LoadROM(slot_roms, conf.Slots.Catalog[i])
			SLOTS[i] = mmu.NewROM("SLOT_"+strconv.Itoa(i), slot_roms, conf.Slots.Catalog[i])
			MEM.Attach(SLOTS[i], 0xC0+uint(i), 1)
		}
		// else {
		// 	SLOTS[i] = make([]byte, slot_roms)
		// 	MEM.Clear(SLOTS[i], 0, 0x71)
		// }
	}
}

func setup() {
	LayoutSel = 0
	MEM = mmu.Init(256, 256)

	// Common Setup
	RAM = mmu.NewRAM("RAM", ramSize, false)
	RAM.Clear(0x1000, 0xFF)

	AUX = mmu.NewRAM("AUX", ramSize, false)
	AUX.Clear(0x1000, 0xFF)

	Disks = InitDiskInterface(conf)
	IO = InitSoftSwitch("IO", softSwitches, Disks, &CRTC)

	MEM.Attach(RAM, 0x00, 256)
	MEM.Attach(IO, 0xC0, 1)
	loadSlots()
	// IOAccess = InitIO(Disk1, nil, &CRTC)

	// if MODEL == 1 {
	// 	AUX = nil
	apple2_Roms()
	// } else {
	// 	ZP = make([]byte, 0x0200)
	// 	MEM.Clear(ZP, 0x1000, 0xFF)
	// 	ALT_ZP = make([]byte, 0x0200)
	// 	MEM.Clear(ALT_ZP, 0x1000, 0xFF)

	// 	AUX = make([]byte, ramSize)
	// 	MEM.Clear(AUX, 0x1000, 0xFF)
	// 	AUX_BANK1 = make([]byte, romSize)
	// 	MEM.Clear(AUX_BANK1, 0x1000, 0xFF)
	// 	AUX_BANK2 = make([]byte, romSize*3)
	// 	MEM.Clear(AUX_BANK2, 0x1000, 0xFF)
	// 	apple2e_Roms()
	// }

	// memLayouts(MODEL)

	outputDriver = render.SDL2Driver{}
	initKeyboard()
	CRTC.Init(RAM.Buff, AUX.Buff, IO.Buff, CHARGEN.Buff, &outputDriver, conf)
	outputDriver.SetKeyboardLine(&InputLine)

	// Throttle setup
	if conf.Mhz == 0 {
		timeGap = time.Microsecond * 1
	} else {
		timeGap = (time.Duration(conf.ThrottleInterval/conf.Mhz) * time.Microsecond)
	}

	MEM.DumpMap()
	os.Exit(0)
	// CPU Setup
	cpu.Init(conf.CPUModel, MEM, conf.Globals.DebugMode)
}

func RunEmulation() {
	var key byte
	var interCycles int64 = 0
	var cycles uint64 = 0
	var start = time.Now()
	var elapsed time.Duration
	var throttled = false

	keyb, _ = tty.Open()
	defer keyb.Close()

	go input()

	// defer timeTrack(time.Now(), "RunEmulation")
	for {
		CRTC.Run()

		if !throttled {
			if InputLine.KeyCode != 0 && !is_Keypressed {
				key = keyMap[InputLine.KeyCode][InputLine.Mode]
				// log.Printf("KEY DOWN - Code: %d  Mode: %d  -> %d", InputLine.KeyCode, InputLine.Mode, key)
				IO.Buff[0] = key | 0b10000000
				is_Keypressed = true
				InputLine.KeyCode = 0
				InputLine.Mode = 0
			}

			cpu.NextCycle()
			interCycles++
			cycles++
		}

		if interCycles >= conf.ThrottleInterval {
			elapsed = time.Now().Sub(start)
			if elapsed < timeGap {
				throttled = true
			} else {
				outputDriver.SetSpeed(float64(interCycles / elapsed.Microseconds()))
				interCycles = 0
				throttled = false
				start = time.Now()
			}
		}

		if conf.Breakpoint == cpu.InstStart {
			trace = true
			stepper = true
		}

		if cpu.CycleCount == 1 && trace {
			fmt.Printf("%d -- %s\n", cycles, cpu.Trace())
			if stepper {
				if InterractiveMode() {
					go input()
				}
			}
		}
	}
}

func main() {
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

	setup()

	trace = conf.Trace
	stepper = false
	outputDriver.ShowCode = false
	outputDriver.ShowFps = true

	fmt.Printf("Trace : %t - Breakpoint: %04X - Breakcycle: %d\n", trace, conf.Breakpoint, conf.BreakCycle)
	go RunEmulation()
	outputDriver.Run(true)
}
