package main

import (
	"fmt"
	"log"
	"newApple/config"
	"newApple/crtc"
	"newApple/io"
	"runtime"
	"time"

	PROC "github.com/Djoulzy/emutools/mos6510"
	"github.com/Djoulzy/mmu"
	"github.com/mattn/go-tty"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/Djoulzy/emutools/render"
)

const (
	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &config.ConfigData{}

	cpu       PROC.CPU
	MODEL     int
	LayoutSel byte

	MEM *mmu.MMU

	InputLine    render.KEYPressed
	outputDriver render.SDL2Driver
	CRTC         crtc.CRTC
	trace        bool
	stepper      bool
	timeGap      time.Duration // 1Mhz = 1 000 000/s = 1000/ms
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
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
			if InputLine.KeyCode != 0 && !io.Is_Keypressed {
				key = io.KeyMap[InputLine.KeyCode][InputLine.Mode]
				// log.Printf("KEY DOWN - Code: %d  Mode: %d  -> %d", InputLine.KeyCode, InputLine.Mode, key)
				IO.Buff[0] = key | 0b10000000
				io.Is_Keypressed = true
				InputLine.KeyCode = 0
				InputLine.Mode = 0
			}

			cpu.NextCycle()
			interCycles++
			cycles++
		}

		if interCycles >= conf.ThrottleInterval {
			elapsed = time.Since(start)
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

		if cpu.CycleCount == 1 {
			outputDriver.SetDriveStat(Disks.GetStats())
			if trace {
				fmt.Printf("%d -- %s\n", cycles, cpu.Trace())
				if stepper {
					if InterractiveMode() {
						go input()
					}
				}
			}
		}
	}
}

func setup() {
	MEM = mmu.Init(256, 256)

	setupMemoryLayout()

	outputDriver = render.SDL2Driver{}
	io.InitKeyboard()
	CRTC.Init(MEM, CHARGEN, &outputDriver, conf)
	outputDriver.SetKeyboardLine(&InputLine)

	// Throttle setup
	if conf.Mhz == 0 {
		timeGap = time.Microsecond * 1
	} else {
		timeGap = (time.Duration(conf.ThrottleInterval/conf.Mhz) * time.Microsecond)
	}

	// CPU Setup
	cpu.Init(conf.CPUModel, MEM, conf.Globals.DebugMode)
}

func main() {
	confload.Load("config.ini", conf)

	clog.LogLevel = conf.LogLevel
	clog.StartLogging = conf.StartLogging
	if conf.FileLog != "" {
		clog.EnableFileLog(conf.FileLog)
	}
	switch conf.Model {
	case "Apple2":
		MODEL = 1
	case "Apple2e":
		fallthrough
	default:
		MODEL = 2
	}

	log.Printf("-%v-\n", conf.Disks.Disk1)
	log.Printf("-%v-\n", conf.Disks.Disk2)
	setup()

	trace = conf.Trace
	stepper = false
	outputDriver.ShowCode = false
	outputDriver.ShowFps = true

	fmt.Printf("Trace : %t - Breakpoint: %04X - Breakcycle: %d\n", trace, conf.Breakpoint, conf.BreakCycle)
	go RunEmulation()
	outputDriver.Run(true)
}
