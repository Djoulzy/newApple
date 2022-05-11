package main

import (
	"log"
	"newApple/disk"
)

const (
	KBD          = 0x00
	KBDSTRB      = 0x10
	SETSLOTCXROM = 0x06
	SETINTCXROM  = 0x07
	SETINTC3ROM  = 0x0A
	SETSLOTC3ROM = 0x0B
	RDCXROM      = 0x15
	RDC3ROM      = 0x17

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

var (
	C3_INT bool = true
	CX_INT bool = true
)

type io_access struct {
	Disk *disk.DRIVE
}

func (C *io_access) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch translatedAddr {
	case KBD:
		return mem[translatedAddr]
	case KBDSTRB:
		mem[0] = 0
		return mem[translatedAddr]
	case RDCXROM:
		if CX_INT {
			return 0x00
		} else {
			return 0xFF
		}
	case RDC3ROM:
		if C3_INT {
			return 0x00
		} else {
			return 0xFF
		}
	case SLOT6_OFFSET + DRVSM0:
		C.Disk.SetPhase(0, false)
		return 0
	case SLOT6_OFFSET + DRVSM0 + 1:
		C.Disk.SetPhase(0, true)
		return 0
	case SLOT6_OFFSET + DRVSM1:
		C.Disk.SetPhase(1, false)
		return 0
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disk.SetPhase(1, true)
		return 0
	case SLOT6_OFFSET + DRVSM2:
		C.Disk.SetPhase(2, false)
		return 0
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disk.SetPhase(2, true)
		return 0
	case SLOT6_OFFSET + DRVSM3:
		C.Disk.SetPhase(3, false)
		return 0
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disk.SetPhase(3, true)
		return 0

	case SLOT6_OFFSET + DRIVE + 1:
		C.Disk.StartMotor()
		return mem[translatedAddr]
	case SLOT6_OFFSET + DRIVE:
		C.Disk.StopMotor()
		return mem[translatedAddr]

	case SLOT6_OFFSET + DRVWRITE:
		C.Disk.ReadMode = true
		if C.Disk.IsWriteProtected {
			return 0xFF
		} else {
			return 0
		}
	case SLOT6_OFFSET + DRVWRITE + 1:
		C.Disk.ReadMode = false
		return 0

	case SLOT6_OFFSET + DRVDATA:
		if C.Disk.IsRunning() && C.Disk.ReadMode {
			return C.Disk.GetNextByte()
		}
		return 0x00
	case SLOT6_OFFSET + DRVDATA + 1:
		return 0x00

	default:
		// log.Printf("Read Unknown: %02X\n", translatedAddr)
		return mem[translatedAddr]
	}
}

func (C *io_access) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("Accessor", "MWrite", "Addr: %04X -> %02X", 0xE800+translatedAddr, val)
	switch translatedAddr {
	case KBD:
	case KBDSTRB:
		mem[0] = 0
	case SETSLOTCXROM:
		CX_INT = false
		BankSel = 1
	case SETINTCXROM:
		CX_INT = true
		BankSel = 0
	case SETINTC3ROM:
		C3_INT = true
	case SETSLOTC3ROM:
		C3_INT = false
	case SLOT6_OFFSET + DRVSM0:
		C.Disk.SetPhase(0, false)
	case SLOT6_OFFSET + DRVSM0 + 1:
		C.Disk.SetPhase(0, true)
	case SLOT6_OFFSET + DRVSM1:
		C.Disk.SetPhase(1, false)
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disk.SetPhase(1, true)
	case SLOT6_OFFSET + DRVSM2:
		C.Disk.SetPhase(2, false)
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disk.SetPhase(2, true)
	case SLOT6_OFFSET + DRVSM3:
		C.Disk.SetPhase(3, false)
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disk.SetPhase(3, true)
	case SLOT6_OFFSET + DRIVE + 1:
		log.Printf("Write Start Motor\n")
	case SLOT6_OFFSET + DRIVE:
		log.Printf("Write Stop Motor\n")
	case SLOT6_OFFSET + DRVWRITE:
		log.Printf("Write DRVWRITE\n")
	case SLOT6_OFFSET + DRVWRITE + 1:
		log.Printf("Write DRVWRITE+1\n")
	case SLOT6_OFFSET + DRVDATA:
		log.Printf("Write DRVDATA\n")
	case SLOT6_OFFSET + DRVDATA + 1:
		log.Printf("Write DRVDATA+1\n")
	}
	mem[translatedAddr] = val
}
