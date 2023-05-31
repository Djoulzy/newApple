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
	case SETSTDZP:
		log.Printf("ALT_ZP Off")
		is_ALT_ZP = false
		C.Mmu.Enable("ZP")
		C.Mmu.Disable("ALT_ZP")
	case SETALTZP:
		log.Printf("ALT_ZP On")
		is_ALT_ZP = true
		C.Mmu.Enable("ALT_ZP")
		C.Mmu.Disable("ZP")
	case SETSLOTCXROM:
		log.Println("SETSLOTCXROM")
		is_CX_INT = false
		C.Mmu.Enable("SLOT1")
		C.Mmu.Enable("SLOT2")
		C.Mmu.Enable("SLOT3")
		C.Mmu.Enable("SLOT4")
		C.Mmu.Enable("SLOT5")
		C.Mmu.Enable("SLOT6")
		C.Mmu.Enable("SLOT7")
	case SETINTCXROM:
		log.Println("SETINTCXROM")
		is_CX_INT = true
		C.Mmu.Disable("SLOT1")
		C.Mmu.Disable("SLOT2")
		C.Mmu.Disable("SLOT3")
		C.Mmu.Disable("SLOT4")
		C.Mmu.Disable("SLOT5")
		C.Mmu.Disable("SLOT6")
		C.Mmu.Disable("SLOT7")
	case SETSLOTC3ROM:
		C.Mmu.Enable("SLOT3")
		is_C3_INT = false
	case SETINTC3ROM:
		C.Mmu.Disable("SLOT3")
		is_C3_INT = true
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

	case SLOT6_OFFSET + DRVSM0:
		C.Disks.SetPhase(0, false)
	case SLOT6_OFFSET + DRVSM0 + 1:
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
		log.Printf("Write Unknown: %04X\n", addr)
	}
	// mem[translatedAddr] = val
}
