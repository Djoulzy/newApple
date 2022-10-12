package main

import (
	"fmt"
	"newApple/config"
	"newApple/crtc"
	"newApple/disk"
	"runtime"
	"time"

	"github.com/Djoulzy/emutools/mos6510"
	"github.com/mattn/go-tty"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/Djoulzy/emutools/mem2"
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

	cpu       mos6510.CPU
	MODEL     int
	LayoutSel byte

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

	MEM      mem2.BANK
	IOAccess mem2.MEMAccess

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
	ROM_D = mem2.LoadROM(romSize, "assets/roms/II/D.bin")
	ROM_EF = mem2.LoadROM(romSize*2, "assets/roms/II/EF.bin")
	CHARGEN = mem2.LoadROM(chargenSize, "assets/roms/II/3410036.bin")
}

func apple2e_Roms() {
	ROM_C = mem2.LoadROM(romSize, "assets/roms/IIe/C.bin")
	ROM_D = mem2.LoadROM(romSize, "assets/roms/IIe/D.bin")
	ROM_EF = mem2.LoadROM(romSize*2, "assets/roms/IIe/EF.bin")
	CHARGEN = mem2.LoadROM(chargenSize*2, "assets/roms/IIe/Video_US.bin")
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
			SLOTS[i] = mem2.LoadROM(slot_roms, conf.Slots.Catalog[i])
		} else {
			SLOTS[i] = make([]byte, slot_roms)
			mem2.Clear(SLOTS[i], 0, 0x71)
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
	LayoutSel = 0
	MEM = mem2.InitBanks(nbMemLayout, &LayoutSel)

	// Common Setup
	RAM = make([]byte, ramSize)
	mem2.Clear(RAM, 0x1000, 0xFF)
	BANK1 = make([]byte, romSize)
	mem2.Clear(BANK1, 0x1000, 0xFF)
	BANK2 = make([]byte, romSize*3)
	mem2.Clear(BANK2, 0x1000, 0xFF)

	IO = make([]byte, softSwitches)
	mem2.Clear(IO, 0, 0x00)
	Disk1, _ := loadDisks()
	loadSlots()
	IOAccess = InitIO(Disk1, nil, &CRTC)

	if MODEL == 1 {
		AUX = nil
		apple2_Roms()
	} else {
		ZP = make([]byte, 0x0200)
		mem2.Clear(ZP, 0x1000, 0xFF)
		ALT_ZP = make([]byte, 0x0200)
		mem2.Clear(ALT_ZP, 0x1000, 0xFF)

		AUX = make([]byte, ramSize)
		mem2.Clear(AUX, 0x1000, 0xFF)
		AUX_BANK1 = make([]byte, romSize)
		mem2.Clear(AUX_BANK1, 0x1000, 0xFF)
		AUX_BANK2 = make([]byte, romSize*3)
		mem2.Clear(AUX_BANK2, 0x1000, 0xFF)
		apple2e_Roms()
	}

	memLayouts(MODEL)

	outputDriver = render.SDL2Driver{}
	initKeyboard()
	CRTC.Init(RAM, AUX, IO, CHARGEN, &outputDriver, conf)
	outputDriver.SetKeyboardLine(&InputLine)

	// Throttle setup
	if conf.Mhz == 0 {
		timeGap = time.Microsecond * 1
	} else {
		timeGap = (time.Duration(conf.ThrottleInterval/conf.Mhz) * time.Microsecond)
	}

	// CPU Setup
	cpu.Init(conf.CPUModel, &MEM, conf.Debug || conf.Disassamble)
}

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Now().Sub(start)
// 	log.Printf("%s took %s", name, elapsed)
// }

func RunEmulation() {
	var key byte
	var cycles int64 = 0
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
				IO[0] = key | 0b10000000
				is_Keypressed = true
				InputLine.KeyCode = 0
				InputLine.Mode = 0
			}

			cpu.NextCycle()
			cycles++
		}

		if cycles >= conf.ThrottleInterval {
			elapsed = time.Now().Sub(start)
			if elapsed < timeGap {
				throttled = true
			} else {
				outputDriver.SetSpeed(float64(cycles / elapsed.Microseconds()))
				cycles = 0
				throttled = false
				start = time.Now()
			}
		}

		if cpu.CycleCount == 1 && trace {
			fmt.Println(cpu.Trace())
			if stepper {
				if InterractiveMode() {
					go input()
				}
			}
		}

		// if cpu.CycleCount == 1 {
		// 	if trace {
		// 		run = false
		// 	}
		// 	// if cpu.InstStart > 0x0300 && cpu.InstStart < 0xC000 {
		// 	// 	clog.FileRaw("\n%d: %s", cpu.Cycles, cpu.FullDebug)
		// 	// }
		// 	// outputDriver.DumpCode(cpu.FullInst)
		// 	outputDriver.SetSpeed(speed)
		// 	if conf.Breakpoint == cpu.InstStart {
		// 		fmt.Printf("%s\n", cpu.FullDebug)
		// 		trace = true
		// 	}
		// }
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

	trace = false
	stepper = false
	outputDriver.ShowCode = false
	outputDriver.ShowFps = true

	go RunEmulation()

	outputDriver.Run(true)
}
