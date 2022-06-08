package disk

import (
	woz "newApple/goWoz"

	"github.com/Djoulzy/emutools/mos6510"
)

type DRIVE struct {
	motorPhases      [4]bool
	IsWriteProtected bool
	IsRunning        bool
	ReadMode         bool
	wozImage         *woz.Disk
	wozTrack         *woz.Track

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
	drive.ReadMode = true
	drive.wozTrack = nil

	return &drive
}

func (D *DRIVE) LoadDiskImage(fileName string) {
	D.wozImage = woz.NewWozDisk(fileName)
}

func (D *DRIVE) StartMotor() {
	D.IsRunning = true
}

func (D *DRIVE) StopMotor() {
	D.IsRunning = false
}

func (D *DRIVE) moveHead(offset int) {
	if offset < 0 {
		D.halftrack -= 0.5
		if D.halftrack < 0 {
			D.halftrack = 0
		}
	} else {
		D.halftrack += 0.5
		if D.halftrack > 40 {
			D.halftrack = 40
		}
	}
	// clog.Test("Drive", "moveHead", "HalfTrack: %0.1f", D.halftrack)
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
