package diskdrive

import (
	"log"
	"path/filepath"
	"time"

	"github.com/Djoulzy/godsk"
	"github.com/Djoulzy/gowoz"

	"github.com/Djoulzy/Tools/clog"
)

type DiskImage interface {
	IsWriteProtected() bool
	GetNextByte() byte
	Seek(float32)
	GetMeta() map[string]string
	Dump(bool)
	DumpTrack(float32)
	DumpTracksRaw()
	GetStatus() string
	GetCurrentTrack() float32
}

type DRIVE struct {
	IsEmpty          bool
	motorPhases      [4]bool
	IsWriteProtected bool
	MotorWillStop    bool
	IsSpinning       bool
	diskImage        DiskImage

	currentPhase   int
	diskHasChanges bool
}

var debug bool

func Attach(debugMode bool) *DRIVE {
	drive := DRIVE{}

	drive.currentPhase = 0
	drive.motorPhases = [4]bool{false, false, false, false}
	drive.IsWriteProtected = false
	drive.IsEmpty = true
	drive.IsSpinning = false
	drive.MotorWillStop = false

	debug = debugMode
	return &drive
}

func (D *DRIVE) LoadDiskImage(fileName string) {
	var err error

	ext := filepath.Ext(fileName)
	switch ext {
	case ".woz":
		D.diskImage, err = gowoz.InitContainer(fileName, debug)
	case ".dsk":
		D.diskImage, err = godsk.InitContainer(fileName, debug)
	default:
		panic("Unknown image disk format")
	}

	if err != nil {
		panic(err)
	}

	D.IsEmpty = false
	D.IsWriteProtected = D.diskImage.IsWriteProtected()
}

func (D *DRIVE) StartMotor() {
	D.IsSpinning = true
	D.MotorWillStop = false
}

func (D *DRIVE) motorStopDelay() {
	time.Sleep(time.Millisecond * 1000)
	if D.MotorWillStop {
		D.IsSpinning = false
		D.MotorWillStop = false
	}
}

func (D *DRIVE) StopMotor() {
	// fmt.Printf("Stop Motor\n")
	if D.IsSpinning && !D.MotorWillStop {
		D.MotorWillStop = true
		go D.motorStopDelay()
	}
}

func (D *DRIVE) GetNextByte() byte {
	if D.IsEmpty {
		return 0
	}
	return D.diskImage.GetNextByte()
}

func (D *DRIVE) SetPhase(phase int, state bool) {
	// fmt.Printf("Set Phase %d - State: %v\n", phase, state)
	if D.IsEmpty {
		return
	}
	if !state {
		return
	}
	if phase == 3 && D.currentPhase == 0 {
		clog.FileRaw("\nMove Head DOWN")
		D.diskImage.Seek(-0.5)
		D.currentPhase = phase
		return
	}
	if phase == 0 && D.currentPhase == 3 {
		clog.FileRaw("\nMove Head UP")
		D.diskImage.Seek(0.5)
		D.currentPhase = phase
		return
	}
	if phase > D.currentPhase {
		clog.FileRaw("\nMove Head UP")
		D.diskImage.Seek(0.5)
		D.currentPhase = phase
		return
	}
	if phase < D.currentPhase {
		clog.FileRaw("\nMove Head DOWN")
		D.diskImage.Seek(-0.5)
		D.currentPhase = phase
		return
	}
}

func (D *DRIVE) Dump(full bool) {
	D.diskImage.Dump(full)
}

func (D *DRIVE) DumpMeta() {
	log.Printf("%s", D.diskImage.GetMeta())
}

func (D *DRIVE) DumpTrack(trk float32) {
	D.diskImage.DumpTrack(trk)
}

func (D *DRIVE) DumpTracksRaw() {
	D.diskImage.DumpTracksRaw()
}

func (D *DRIVE) GetStatus() string {
	if D.IsEmpty {
		return "Empty"
	}
	return D.diskImage.GetStatus()
}
