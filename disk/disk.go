package disk

import (
	"hash/crc32"
	"io/ioutil"
	"log"
	"math"
)

type DRIVE struct {
	motorPhases      [4]bool
	IsWriteProtected bool
	motorIsRunning   bool
	ReadMode         bool

	prevHalfTrack  int
	halftrack      int
	trackLocation  uint32
	trackStart     []uint32
	trackNbits     []uint32
	diskData       []byte
	currentPhase   int
	direction      int
	diskHasChanges bool
}

var pickbit = []byte{128, 64, 32, 16, 8, 4, 2, 1}
var crcTable *crc32.Table

func Attach() *DRIVE {
	drive := DRIVE{}

	drive.currentPhase = 0
	drive.motorPhases = [4]bool{false, false, false, false}
	drive.direction = 0
	drive.trackStart = make([]uint32, 80)
	drive.trackNbits = make([]uint32, 80)
	drive.prevHalfTrack = 0
	drive.halftrack = 0
	drive.IsWriteProtected = false

	crcTable = crc32.MakeTable(0xEDB88320)
	return &drive
}

func (D *DRIVE) LoadDiskImage(fileName string) {
	D.diskData = make([]byte, 0x2A518)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 0x2A518; i++ {
		D.diskData[i] = byte(data[i])
	}
	D.decodeDiskData(fileName)
}

func (D *DRIVE) StartMotor() {
	D.motorIsRunning = true
}

func (D *DRIVE) IsRunning() bool {
	return D.motorIsRunning
}

func (D *DRIVE) StopMotor() {
	D.motorIsRunning = false
}

func (D *DRIVE) moveHead(offset int) {
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
		oldloc := D.trackLocation
		D.trackLocation = uint32(math.Floor(float64(D.trackLocation * (D.trackNbits[D.halftrack] / D.trackNbits[D.prevHalfTrack]))))
		if D.trackLocation > 3 {
			D.trackLocation -= 4
		}
		log.Printf("track=%d %d %d %d %d", D.halftrack, oldloc, D.trackLocation, D.trackNbits[D.halftrack], D.trackNbits[D.prevHalfTrack])
	}
}

func (D *DRIVE) getNextBit() byte {
	var bit byte
	D.trackLocation = D.trackLocation % D.trackNbits[D.halftrack]
	if D.trackStart[D.halftrack] > 0 {
		fileOffset := D.trackStart[D.halftrack] + (D.trackLocation >> 3)
		byteRead := D.diskData[fileOffset]
		b := D.trackLocation & 7
		bit = (byteRead & pickbit[b]) >> (7 - b)
	} else {
		// TODO: Freak out like a MC3470 and return random bits
		bit = 1
	}
	D.trackLocation++
	return bit
}

var JulesCpt int = 0
var JulesTmp int = 0

func (D *DRIVE) GetNextByte() byte {
	var bit, result byte

	if len(D.diskData) == 0 {
		return 0
	}
	result = 0
	for bit = 0; bit == 0; bit = D.getNextBit() {
	}
	result = 0x80 // the bit we just retrieved is the high bit
	for i := 6; i >= 0; i-- {
		result |= D.getNextBit() << i
	}
	// fmt.Printf("Track: %d Location: %d byte= %02X\n", D.halftrack, D.trackLocation, result)
	return result
}

func (D *DRIVE) SetPhase(phase int, state bool) {
	// var debug string
	D.motorPhases[phase] = state

	ascend := D.motorPhases[(D.currentPhase+1)%4]
	descend := D.motorPhases[(D.currentPhase+3)%4]
	if !D.motorPhases[D.currentPhase] {
		if D.motorIsRunning && ascend {
			D.moveHead(1)
			D.currentPhase = (D.currentPhase + 1) % 4
			// debug = fmt.Sprintf(" currPhase= %d track= %0.1f", D.currentPhase, float64(D.halftrack)/2)

		} else if D.motorIsRunning && descend {
			D.moveHead(-1)
			D.currentPhase = (D.currentPhase + 3) % 4
			// debug = fmt.Sprintf(" currPhase= %d track= %0.1f", D.currentPhase, float64(D.halftrack)/2)
		}
		// log.Printf("***** %s", debug)
	}
}

func (D *DRIVE) destectFormat(header []byte) bool {
	for i := 0; i < len(header); i++ {
		if D.diskData[i] != header[i] {
			return false
		}
	}
	return true
}

func get_crc32(data []byte, offset int) uint32 {
	crc := 0 ^ ^uint32(0)
	for i := offset; i < len(data); i++ {
		crc = crcTable[(crc^uint32(data[i]))&0xFF] ^ (crc >> 8)
	}
	return crc ^ ^uint32(0)
}

func (D *DRIVE) decodeDiskData(fileName string) {
	woz2 := []byte{0x57, 0x4F, 0x5A, 0x32, 0xFF, 0x0A, 0x0D, 0x0A}
	woz1 := []byte{0x57, 0x4F, 0x5A, 0x31, 0xFF, 0x0A, 0x0D, 0x0A}

	D.diskHasChanges = false
	if D.destectFormat(woz2) {
		// D.IsWriteProtected = D.diskData[22] == 1
		crc := D.diskData[8:12]
		storedCRC := uint32(crc[0]) + (uint32(crc[1]) << 8) + (uint32(crc[2]) << 16) + uint32(crc[3])*uint32(math.Pow(2, 24))
		actualCRC := get_crc32(D.diskData, 12)
		if (storedCRC != 0) && (storedCRC != actualCRC) {
			log.Printf("CRC checksum error: %s (stored: %X - calculated: %X)\n", fileName, storedCRC, actualCRC)
		}
		for htrack := 0; htrack < 80; htrack++ {
			tmap_index := uint32(D.diskData[88+htrack*2])
			if tmap_index < 255 {
				tmap_offset := 256 + 8*tmap_index
				trk := D.diskData[tmap_offset : tmap_offset+8]
				D.trackStart[htrack] = 512*uint32(trk[0]) + (uint32(trk[1]) << 8)
				// const nBlocks = trk[2] + (trk[3] << 8)
				D.trackNbits[htrack] = uint32(trk[4]) + uint32(trk[5])<<8 + uint32(trk[6])<<16 + uint32(trk[7])*uint32(math.Pow(2, 24))
			} else {
				D.trackStart[htrack] = 0
				D.trackNbits[htrack] = 51200
				log.Printf("empty woz2 track %d\n", htrack/2)
			}
		}
		return
	}

	if D.destectFormat(woz1) {
		// D.IsWriteProtected = D.diskData[22] == 1
		for htrack := 0; htrack < 80; htrack++ {
			tmap_index := int(D.diskData[88+htrack*2])
			if tmap_index < 255 {
				D.trackStart[htrack] = 256 + uint32(tmap_index)*6656
				trk := D.diskData[D.trackStart[htrack]+6646 : D.trackStart[htrack]+6656]
				D.trackNbits[htrack] = uint32(trk[2]) + (uint32(trk[3]) << 8)
			} else {
				D.trackStart[htrack] = 0
				D.trackNbits[htrack] = 51200
				log.Printf("empty woz2 track %d\n", htrack/2)
			}
		}
		return
	}
	log.Printf("Unknown disk format.\n")
}
