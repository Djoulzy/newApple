package disk

import (
	"log"
	"time"

	"github.com/Djoulzy/gowoz"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/emutools/mos6510"
)

var MotorIsOn bool = false

type DRIVE struct {
	motorPhases      [4]bool
	IsWriteProtected bool
	IsRunning        bool
	wozImage         *gowoz.WOZFileFormat

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

	D.wozImage, err = gowoz.InitWozFile(fileName)
	if err != nil {
		panic(err)
	}
	D.wozImage.Dump(false)

	D.IsWriteProtected = D.wozImage.INFO.WriteProtected == 1
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
	return D.wozImage.GetNextByte()
}

func (D *DRIVE) SetPhase(phase int, state bool) {
	// fmt.Printf("Set Phase %d - State: %v\n", phase, state)
	if state == false {
		return
	}
	if phase == 3 && D.currentPhase == 0 {
		clog.FileRaw("\nMove Head DOWN")
		D.wozImage.Seek(-0.5)
		D.currentPhase = phase
		return
	}
	if phase == 0 && D.currentPhase == 3 {
		clog.FileRaw("\nMove Head UP")
		D.wozImage.Seek(0.5)
		D.currentPhase = phase
		return
	}
	if phase > D.currentPhase {
		clog.FileRaw("\nMove Head UP")
		D.wozImage.Seek(0.5)
		D.currentPhase = phase
		return
	}
	if phase < D.currentPhase {
		clog.FileRaw("\nMove Head DOWN")
		D.wozImage.Seek(-0.5)
		D.currentPhase = phase
		return
	}
}

func (D *DRIVE) DumpMeta() {
	log.Printf("%s", D.wozImage.META.Metadata)
}
