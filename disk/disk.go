package disk

import (
	"hash/crc32"
	"io/ioutil"
	"log"
	"math"
)

type DRIVE struct {
	prevHalfTrack byte
	halftrack     byte
	trackLocation int

	trackStart []int
	trackNbits []int

	diskData []byte

	motorIsRunning      bool
	diskImageHasChanges bool
	isWriteProtected    bool
}

var pickbit = []byte{128, 64, 32, 16, 8, 4, 2, 1}
var crcTable *crc32.Table

func Attach() *DRIVE {
	drive := DRIVE{}

	drive.trackStart = make([]int, 80)
	drive.trackNbits = make([]int, 80)
	drive.prevHalfTrack = 0
	drive.halftrack = 0

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

func (D *DRIVE) moveHead(offset byte) {
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

	// Adjust new track location based on arm position relative to old track loc.
	if D.trackStart[D.halftrack] > 0 && D.prevHalfTrack != D.halftrack {
		D.trackLocation = int(math.Floor(float64(D.trackLocation * (D.trackNbits[D.halftrack] / D.trackNbits[D.prevHalfTrack]))))
		if D.trackLocation > 3 {
			D.trackLocation -= 4
		}
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
	// log.Printf(" trackLocation= %d byte= %02X\n", D.trackLocation, result)

	return result
}

func Mycrc32(data []uint8, offset int) int {
	if crcTable[255] == 0 {
		crcTable = crc32.MakeTable(0xEDB88320)
	}
	crc := 0 ^ (-1)
	for i := offset; i < len(data); i++ {
		crc = int(crc>>8) ^ int(crcTable[(uint8(crc)^data[i])&0xFF])
	}

	return crc ^ (-1)
}

func (D *DRIVE) destectFormat(header []byte) bool {
	for i := 0; i < len(header); i++ {
		if D.diskData[i] != header[i] {
			return false
		}
	}
	return true
}

func (D *DRIVE) decodeDiskData(fileName string) {
	woz2 := []byte{0x57, 0x4F, 0x5A, 0x32, 0xFF, 0x0A, 0x0D, 0x0A}
	woz1 := []byte{0x57, 0x4F, 0x5A, 0x31, 0xFF, 0x0A, 0x0D, 0x0A}

	D.diskImageHasChanges = false
	if D.destectFormat(woz2) {
		D.isWriteProtected = D.diskData[22] == 1
		crc := D.diskData[8:12]
		storedCRC := int(crc[0]) + (int(crc[1]) << 8) + (int(crc[2]) << 16) + int(crc[3])*int(math.Pow(2, 24))
		actualCRC := crc32.Checksum(D.diskData, crcTable)
		if (storedCRC != 0) && (uint32(storedCRC) != actualCRC) {
			log.Printf("CRC checksum error: %s\n", fileName)
		}
		for htrack := 0; htrack < 80; htrack++ {
			tmap_index := uint16(D.diskData[88+htrack*2])
			if tmap_index < 255 {
				tmap_offset := 256 + 8*tmap_index
				trk := D.diskData[tmap_offset : tmap_offset+8]
				D.trackStart[htrack] = 512*int(trk[0]) + (int(trk[1]) << 8)
				// const nBlocks = trk[2] + (trk[3] << 8)
				D.trackNbits[htrack] = int(trk[4]) + int(trk[5])<<8 + int(trk[6])<<16 + int(trk[7])*int(math.Pow(2, 24))
			} else {
				D.trackStart[htrack] = 0
				D.trackNbits[htrack] = 51200
				log.Printf("empty woz2 track %d\n", htrack/2)
			}
		}
		return
	}

	if D.destectFormat(woz1) {
		D.isWriteProtected = D.diskData[22] == 1
		for htrack := 0; htrack < 80; htrack++ {
			tmap_index := int(D.diskData[88+htrack*2])
			if tmap_index < 255 {
				D.trackStart[htrack] = 256 + tmap_index*6656
				trk := D.diskData[D.trackStart[htrack]+6646 : D.trackStart[htrack]+6656]
				D.trackNbits[htrack] = int(trk[2]) + (int(trk[3]) << 8)
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
