package disk

import (
	"encoding/json"
	"fmt"
	"log"
	woz "newApple/goWoz"
	"strconv"
	"time"

	"github.com/Djoulzy/emutools/mos6510"
)

var MotorIsOn bool = false

type WozInfo struct {
	Version            int      `json:"version"`
	WriteProtected     bool     `json:"write_protected"`
	RequireRam         int      `json:"requires_ram"`
	CompatibleHardware []string `json:"compatible_hardware"`
}

type WozMeta struct {
	Side            string   `json:"side"`
	Version         string   `json:"version"`
	Subtitle        string   `json:"subtitle"`
	Language        string   `json:"language"`
	Title           string   `json:"title"`
	Publisher       string   `json:"publisher"`
	RequireRam      string   `json:"requires_ram"`
	Copyright       string   `json:"copyright"`
	RequiresMachine []string `json:"requires_machine"`
	Genre           string   `json:"genre"`
	SideName        string   `json:"side_name"`
}

type WozDisk struct {
	Info WozInfo `json:"info"`
	Meta WozMeta `json:"meta"`
}

type WOZ struct {
	Disk WozDisk `json:"woz"`
}

type DRIVE struct {
	motorPhases      [4]bool
	IsWriteProtected bool
	IsRunning        bool
	wozImage         *woz.Disk
	wozTrack         *woz.Track
	WOZ              WOZ

	halftrack      float64
	currentPhase   int
	diskHasChanges bool

	cpu *mos6510.CPU
}

func Attach(cpu *mos6510.CPU) *DRIVE {
	drive := DRIVE{}
	drive.cpu = cpu

	drive.currentPhase = 0
	drive.motorPhases = [4]bool{false, false, false, false}
	drive.halftrack = 0
	drive.IsWriteProtected = false
	drive.wozTrack = nil

	return &drive
}

func (D *DRIVE) LoadDiskImage(fileName string) {
	D.wozImage = woz.NewWozDisk(fileName)
	if err := json.Unmarshal(D.wozImage.Dump(), &D.WOZ); err != nil {
		panic(err)
	}
	D.DumpMeta()
	D.IsWriteProtected = D.WOZ.Disk.Info.WriteProtected
}

func (D *DRIVE) StartMotor() {
	D.IsRunning = true
}

func (D *DRIVE) motorStopDelay() {
	time.Sleep(time.Millisecond * 1000)
	if !MotorIsOn {
		D.IsRunning = false
	}
}

func (D *DRIVE) StopMotor() {
	go D.motorStopDelay()
}

func (D *DRIVE) moveHead(offset int) {
	if offset < 0 {
		D.halftrack -= 0.5
		if D.halftrack < 0 {
			D.halftrack = 0
		}
	} else {
		D.halftrack += 0.5
		if D.halftrack > 68 {
			D.halftrack = 68
		}
	}
	// clog.FileRaw("\nHalfTrack: %0.1f", D.halftrack)
	fmt.Printf("HalfTrack: %v\n", D.halftrack)
	D.wozTrack = D.wozImage.Seek(D.halftrack)
}

func (D *DRIVE) GetNextByte() byte {
	if D.wozTrack != nil {
		return byte(D.wozTrack.Nibble())
	} else {
		return 0
	}
}

func (D *DRIVE) SetPhase(phase int, state bool) {
	if state == false {
		return
	}
	if phase == 3 && D.currentPhase == 0 {
		D.moveHead(-1)
		D.currentPhase = phase
		return
	}
	if phase == 0 && D.currentPhase == 3 {
		D.moveHead(1)
		D.currentPhase = phase
		return
	}
	if phase > D.currentPhase {
		D.moveHead(1)
		D.currentPhase = phase
		return
	}
	if phase < D.currentPhase {
		D.moveHead(-1)
		D.currentPhase = phase
		return
	}
}

func (D *DRIVE) DumpMeta() {
	log.Printf("WOZ Disk Meta:\n-- Title: %s (%s)\n-- Write Protected: %s", D.WOZ.Disk.Meta.Title, D.WOZ.Disk.Meta.Subtitle, strconv.FormatBool(D.WOZ.Disk.Info.WriteProtected))
}
