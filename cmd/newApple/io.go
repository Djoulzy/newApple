package main

import (
	"log"
	"newApple/crtc"
	"newApple/disk"
)

const (
	// MEMORY MANAGEMENT SOFT SWITCHES (W)
	_80STOREOFF  = 0x00
	INTCXROMOFF  = 0x06
	INTCXROMON   = 0x07
	SLOTC3ROMOFF = 0x0A
	SLOTC3ROMON  = 0x0B

	// VIDEO SOFT SWITCHES (W/R)
	_80COLOFF     = 0x0C
	_80COLON      = 0x0D
	ALTCHARSETOFF = 0x0E
	ALTCHARSETON  = 0x0F
	TEXTOFF       = 0x50
	TEXTON        = 0x51
	MIXEDOFF      = 0x52
	MIXEDON       = 0x53
	PAGE2OFF      = 0x54
	PAGE2ON       = 0x55
	HIRESOFF      = 0x56
	HIRESON       = 0x57

	// SOFT SWITCH STATUS FLAGS (R bit 7)
	AKD        = 0x10
	INTCXROM   = 0x15
	SLOTC3ROM  = 0x17
	TEXT       = 0x1A
	MIXED      = 0x1B
	PAGE2      = 0x1C
	HIRES      = 0x1D
	ALTCHARSET = 0x1E
	_80COL     = 0x1F

	// OTHER
	SPKR = 0x30

	// SLOTS
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

// PRINT PEEK(49173)

var (
	is_C3_INT bool = true
	is_CX_INT bool = false
)

type io_access struct {
	Disk  *disk.DRIVE
	Video *crtc.CRTC
}

func (C *io_access) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch translatedAddr {
	case _80STOREOFF:
		return mem[translatedAddr]
	case AKD:
		mem[_80STOREOFF] = 0
		return mem[translatedAddr]
	case INTCXROM:
		if is_CX_INT {
			return 0x8D
		} else {
			return 0x00
		}
	case SLOTC3ROM:
		if is_C3_INT {
			return 0x8D
		} else {
			return 0x00
		}
	case TEXTOFF:
		crtc.Is_TEXTMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case TEXTON:
		crtc.Is_TEXTMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case MIXEDOFF:
		crtc.Is_MIXEDMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case MIXEDON:
		crtc.Is_MIXEDMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case HIRESOFF:
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case HIRESON:
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case PAGE2OFF:
		crtc.Is_PAGE2 = false
		C.Video.UpdateVideoRam(crtc.TEXTPAGE1)
		return 0
	case PAGE2ON:
		crtc.Is_PAGE2 = true
		C.Video.UpdateVideoRam(crtc.TEXTPAGE2)
		return 0

	case TEXT:
		if crtc.Is_TEXTMODE {
			return 0x80
		}
		return 0x00
	case MIXED:
		if crtc.Is_MIXEDMODE {
			return 0x80
		}
		return 0x00
	case PAGE2:
		if crtc.Is_PAGE2 {
			return 0x80
		}
		return 0x00
	case HIRES:
		if crtc.Is_HIRESMODE {
			return 0x80
		}
		return 0x00
	case SPKR:
		return 0
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
			return 0x8D
		} else {
			return 0x00
		}
	case SLOT6_OFFSET + DRVWRITE + 1:
		C.Disk.ReadMode = false
		return 0

	case SLOT6_OFFSET + DRVDATA:
		if C.Disk.IsRunning() && C.Disk.ReadMode {
			tmp := C.Disk.GetNextByte()
			// log.Printf("Read : %02X\n", tmp)
			return tmp
		}
		return 0x00
	case SLOT6_OFFSET + DRVDATA + 1:
		return 0x00
	case SLOT6_OFFSET + DRVSEL:
		return 0x00
	case SLOT6_OFFSET + DRVSEL + 1:
		return 0x00
	default:
		// log.Printf("Read Unknown: %02X\n", translatedAddr)
		return 0x00
	}
}

func (C *io_access) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("Accessor", "MWrite", "Addr: %04X -> %02X", 0xE800+translatedAddr, val)
	switch translatedAddr {
	case _80STOREOFF:
		mem[_80STOREOFF] = val
	case AKD:
		mem[_80STOREOFF] = 0
	case INTCXROMOFF:
		log.Printf("WRITE MemConf = 1")
		is_CX_INT = false
		BankSel = 1
	case INTCXROMON:
		log.Printf("WRITE MemConf = 0")
		is_CX_INT = true
		BankSel = 0
	case SLOTC3ROMON:
		if !is_CX_INT {
			BankSel = 3
		}
		is_C3_INT = true
	case SLOTC3ROMOFF:
		if is_CX_INT {
			BankSel = 2
		}
		is_C3_INT = false
	case TEXTOFF:
		crtc.Is_TEXTMODE = false
		C.Video.UpdateGraphMode()
	case TEXTON:
		crtc.Is_TEXTMODE = true
		C.Video.UpdateGraphMode()
	case MIXEDOFF:
		crtc.Is_MIXEDMODE = false
		C.Video.UpdateGraphMode()
	case MIXEDON:
		crtc.Is_MIXEDMODE = true
		C.Video.UpdateGraphMode()
	case HIRESOFF:
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
	case HIRESON:
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
	case PAGE2OFF:
		crtc.Is_PAGE2 = false
		C.Video.UpdateVideoRam(crtc.TEXTPAGE1)
	case PAGE2ON:
		crtc.Is_PAGE2 = true
		C.Video.UpdateVideoRam(crtc.TEXTPAGE2)
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
	default:
		// log.Printf("Write Unknown: %02X\n", translatedAddr)
	}
	// mem[translatedAddr] = val
}
