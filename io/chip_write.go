package io

import (
	"log"
	"newApple/crtc"
)

func (C *SoftSwitch) Write(addr uint16, val byte) {
	switch addr {
	case CLR80VID:
		crtc.Is_80COL = false
		C.Video.UpdateGraphMode()
	case SET80VID:
		crtc.Is_80COL = true
		C.Video.UpdateGraphMode()
	case STOREOFF:
		is_80Store = false
	case STOREON:
		is_80Store = true
	case KBDSTRB:
		Is_Keypressed = false
		C.Buff[KBD] = 0
	case TXTCLR:
		crtc.Is_TEXTMODE = false
		C.Video.UpdateGraphMode()
	case TXTSET:
		crtc.Is_TEXTMODE = true
		C.Video.UpdateGraphMode()
	case MIXCLR:
		crtc.Is_MIXEDMODE = false
		C.Video.UpdateGraphMode()
	case MIXSET:
		crtc.Is_MIXEDMODE = true
		C.Video.UpdateGraphMode()
	case LORES:
		log.Printf("HIRES OFF")
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
	case HIRES:
		log.Printf("HIRES ON")
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
	case TXTPAGE1:
		crtc.Is_PAGE2 = false
		if is_80Store {
			// C.Mmu.Disable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}
	case TXTPAGE2:
		crtc.Is_PAGE2 = true
		if is_80Store {
			// C.Mmu.Enable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}

	case SETSTDZP:
		fallthrough
	case SETALTZP:
		log.Printf("ZP Management: %04X\n", addr+0xC000)

	case SETSLOTCXROM:
		fallthrough
	case SETINTCXROM:
		fallthrough
	case SETINTC3ROM:
		fallthrough
	case SETSLOTC3ROM:
		log.Printf("Slot Management: %04X\n", addr+0xC000)

	case SLOT6_OFFSET + DRVSM0:
		log.Println("[WRITE] SetPhase Off")
		C.Disks.SetPhase(0, false)
	case SLOT6_OFFSET + DRVSM0 + 1:
		log.Println("[WRITE] SetPhase On")
		C.Disks.SetPhase(0, true)
	case SLOT6_OFFSET + DRVSM1:
		C.Disks.SetPhase(1, false)
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disks.SetPhase(1, true)
	case SLOT6_OFFSET + DRVSM2:
		C.Disks.SetPhase(2, false)
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disks.SetPhase(2, true)
	case SLOT6_OFFSET + DRVSM3:
		C.Disks.SetPhase(3, false)
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disks.SetPhase(3, true)

	case SLOT6_OFFSET + DRIVE:
		// PRINT (PEEK(49384))
		C.Disks.diskMotorsOFF()
	case SLOT6_OFFSET + DRIVE + 1:
		// PRINT (PEEK(49385))
		C.Disks.diskMotorsON()

	case SLOT6_OFFSET + DRVSEL:
		// PRINT (PEEK(49290))
		C.Disks.driveSelect(0)
	case SLOT6_OFFSET + DRVSEL + 1:
		// PRINT (PEEK(49291))
		C.Disks.driveSelect(1)

	case SLOT6_OFFSET + DRVDATA: // Q6
		// PRINT (PEEK(49292))
		C.Disks.ShiftOrRead()
	case SLOT6_OFFSET + DRVDATA + 1: // Q6
		// PRINT (PEEK(49293))
		C.Disks.LoadOrCheck()

	case SLOT6_OFFSET + DRVWRITE: // Q7
		// PRINT (PEEK(49390))
		C.Disks.SetSequencerMode(SEQ_READ_MODE)
	case SLOT6_OFFSET + DRVWRITE + 1: // Q7
		C.Disks.SetSequencerMode(SEQ_WRITE_MODE)

	default:
		// log.Printf("IO Write Unknown: %04X\n", addr+0xC000)
		C.Buff[addr] = val
	}
	// mem[translatedAddr] = val
}
