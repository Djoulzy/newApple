package main

import (
	"log"
	"newApple/disk"
)

const (
	SLOT0_OFFSET = 0x90
	SLOT1_OFFSET = 0x90
	SLOT2_OFFSET = 0xA0
	SLOT3_OFFSET = 0xB0
	SLOT4_OFFSET = 0xC0
	SLOT5_OFFSET = 0xD0
	SLOT6_OFFSET = 0xE0
	SLOT7_OFFSET = 0xF0

	DRVSM0   = 0x00
	DRVSM1   = 0x02
	DRVSM2   = 0x04
	DRVSM3   = 0x06
	DRIVE    = 0x08
	DRVSEL   = 0x0A
	DRVDATA  = 0x0C
	DRVWRITE = 0x0E
)

type io_access struct {
	Disk *disk.DRIVE
}

func (C *io_access) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch translatedAddr {
	case SLOT6_OFFSET + DRVSM0:
		log.Printf("Motor Switch 0 off\n")
		C.Disk.NewPhase(0)
		return 0
	case SLOT6_OFFSET + DRVSM0 + 1:
		log.Printf("Motor Switch 0 on\n")
		return 0
	case SLOT6_OFFSET + DRVSM1:
		log.Printf("Motor Switch 1 off\n")
		C.Disk.NewPhase(1)
		return 0
	case SLOT6_OFFSET + DRVSM1 + 1:
		log.Printf("Motor Switch 1 on\n")
		return 0
	case SLOT6_OFFSET + DRVSM2:
		log.Printf("Motor Switch 2 off\n")
		C.Disk.NewPhase(2)
		return 0
	case SLOT6_OFFSET + DRVSM2 + 1:
		log.Printf("Motor Switch 2 on\n")
		return 0
	case SLOT6_OFFSET + DRVSM3:
		log.Printf("Motor Switch 3 off\n")
		C.Disk.NewPhase(3)
		return 0
	case SLOT6_OFFSET + DRVSM3 + 1:
		log.Printf("Motor Switch 3 on\n")
		return 0
	case SLOT6_OFFSET + DRIVE + 1:
		log.Printf("Read - Start Motor\n")
		C.Disk.StartMotor()
		return mem[translatedAddr]
	case SLOT6_OFFSET + DRIVE:
		log.Printf("Read - Stop Motor\n")
		C.Disk.StopMotor()
		return mem[translatedAddr]
	case SLOT6_OFFSET + DRVDATA:
		// log.Printf("Read - Q6 off (read)\n")
		if C.Disk.IsRunning() {
			tmp := C.Disk.GetNextByte()
			// log.Printf("Read: %02X\n", tmp)
			return tmp
		}
		return 0x00
	case SLOT6_OFFSET + DRVDATA + 1:
		log.Printf("Read - Q6 on (WP sense)\n")
		return mem[translatedAddr]
	case 0x10: // Clear keyboard strobe
		mem[0] = 0
		fallthrough
	default:
		return mem[translatedAddr]
	}
}

func (C *io_access) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("Accessor", "MWrite", "Addr: %04X -> %02X", 0xE800+translatedAddr, val)
	switch translatedAddr {
	case SLOT6_OFFSET + DRVSM0:
		fallthrough
	case SLOT6_OFFSET + DRVSM0 + 1:
		fallthrough
	case SLOT6_OFFSET + DRVSM1:
		fallthrough
	case SLOT6_OFFSET + DRVSM1 + 1:
		fallthrough
	case SLOT6_OFFSET + DRVSM2:
		fallthrough
	case SLOT6_OFFSET + DRVSM2 + 1:
		fallthrough
	case SLOT6_OFFSET + DRVSM3:
		fallthrough
	case SLOT6_OFFSET + DRVSM3 + 1:
		log.Printf("Motor Switch\n")
	case SLOT6_OFFSET + DRIVE + 1:
		log.Printf("Start Motor\n")
	case SLOT6_OFFSET + DRIVE:
		log.Printf("Stop Motor\n")
	}
	mem[translatedAddr] = val
}
