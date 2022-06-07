package disk

import (
	"log"
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

	prevHalfTrack  int
	halftrack      float64
	trackLocation  uint32
	trackStart     []uint32
	trackNbits     []uint32
	diskData       []byte
	currentPhase   int
	direction      int
	diskHasChanges bool

	cpu *mos6510.CPU
}

func Attach(cpu *mos6510.CPU) *DRIVE {
	drive := DRIVE{}
	drive.cpu = cpu

	drive.currentPhase = 0
	drive.motorPhases = [4]bool{false, false, false, false}
	drive.direction = 0
	drive.trackStart = make([]uint32, 80)
	drive.trackNbits = make([]uint32, 80)
	drive.prevHalfTrack = 0
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
	/*
		if D.trackStart[D.halftrack] > 0 {
			D.prevHalfTrack = D.halftrack
		}
		D.halftrack += offset
		if D.halftrack < 0 || D.halftrack > 68 {
			if D.halftrack < 0 {
				D.halftrack = 0
			} else if D.halftrack > 68 {
				D.halftrack = 68
			}
		}
		// log.Printf("track=%0.1f\n", float64(D.halftrack)/2)
		// Adjust new track location based on arm position relative to old track loc.
		if D.trackStart[D.halftrack] > 0 && D.prevHalfTrack != D.halftrack {
			// oldloc := D.trackLocation
			D.trackLocation = uint32(math.Floor(float64(D.trackLocation * (D.trackNbits[D.halftrack] / D.trackNbits[D.prevHalfTrack]))))
			if D.trackLocation > 3 {
				D.trackLocation -= 4
			}
			// log.Printf("track=%d %d %d %d %d", D.halftrack, oldloc, D.trackLocation, D.trackNbits[D.halftrack], D.trackNbits[D.prevHalfTrack])
		}
		// if D.wozTrack != nil {
		// 	D.wozImage.Close()
		// }
		D.wozTrack = D.wozImage.Seek(float64(D.halftrack) / 2)
	*/
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
	log.Printf("HalfTrack: %0.1f", D.halftrack)
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
	/*
		// var debug string
		D.motorPhases[phase] = state

		ascend := D.motorPhases[(D.currentPhase+1)%4]
		descend := D.motorPhases[(D.currentPhase+3)%4]
		if !D.motorPhases[D.currentPhase] {
			if D.IsRunning && ascend {
				D.moveHead(1)
				D.currentPhase = (D.currentPhase + 1) % 4
				// debug = fmt.Sprintf(" currPhase= %d track= %0.1f", D.currentPhase, float64(D.halftrack)/2)

			} else if D.IsRunning && descend {
				D.moveHead(-1)
				D.currentPhase = (D.currentPhase + 3) % 4
				// debug = fmt.Sprintf(" currPhase= %d track= %0.1f", D.currentPhase, float64(D.halftrack)/2)
			}
			// log.Printf("***** %s", debug)
		}
	*/
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
