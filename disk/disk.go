package disk

import (
	"log"
	"path/filepath"
	"time"

	"github.com/Djoulzy/godsk"
	"github.com/Djoulzy/gowoz"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/emutools/mos6510"
)

var MotorIsOn bool = false

type DiskImage interface {
	IsWriteProtected() bool
	GetNextByte() byte
	Seek(float32)
	GetMeta() map[string]string
	Dump(bool)
	DumpTrack(float32)
	DumpTrackRaw(float32)
}

type DRIVE struct {
	motorPhases      [4]bool
	IsWriteProtected bool
	IsRunning        bool
	diskImage        DiskImage

	currentPhase   int
	diskHasChanges bool

	cpu *mos6510.CPU
}

func Attach(cpu *mos6510.CPU) *DRIVE {
	drive := DRIVE{}
	drive.cpu = cpu

	drive.currentPhase = 0
	drive.motorPhases = [4]bool{false, false, false, false}
	drive.IsWriteProtected = false

	return &drive
}

func (D *DRIVE) LoadDiskImage(fileName string) {
	var err error

	ext := filepath.Ext(fileName)
	switch ext {
	case ".woz":
		D.diskImage, err = gowoz.InitContainer(fileName)
	case ".dsk":
		D.diskImage, err = godsk.InitContainer(fileName)
	default:
		panic("Unknown image disk format")
	}

	if err != nil {
		panic(err)
	}
	D.diskImage.Dump(false)

	D.IsWriteProtected = D.diskImage.IsWriteProtected()
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

func (D *DRIVE) GetNextByte() byte {
	return D.diskImage.GetNextByte()
}

func (D *DRIVE) SetPhase(phase int, state bool) {
	// fmt.Printf("Set Phase %d - State: %v\n", phase, state)
	if state == false {
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

func (D *DRIVE) DumpMeta() {
	log.Printf("%s", D.diskImage.GetMeta())
}

func (D *DRIVE) DumpTrack(trk float32) {
	D.diskImage.DumpTrack(trk)
}
