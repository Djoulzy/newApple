package io

import (
	"log"
	"strconv"
)

func (C *SoftSwitch) Write(addr uint16, val byte) {
	switch addr {
	case CLR80VID:
		C.Video.Set40Cols()
	case SET80VID:
		C.Video.Set80Cols()
	case STOREOFF:
		is_80Store = false
	case STOREON:
		is_80Store = true
	case KBDSTRB:
		Is_Keypressed = false
		C.Buff[KBD] = 0

	case TXTCLR:
		C.Video.SetGraphMode()
	case TXTSET:
		C.Video.SetTexMode()
	case MIXCLR:
		C.Video.SetFullMode()
	case MIXSET:
		C.Video.SetMixedMode()
	case LORES:
		C.Video.SetLoResMode()
	case HIRES:
		C.Video.SetHiResMode()
	case TXTPAGE1:
		C.Video.SetPage1()
	case TXTPAGE2:
		C.Video.SetPage2()

	case RDMAINRAM:
		C.Mmu.MountReader("MN_ZPS")
		C.Mmu.MountReader("MN___1")
		C.Mmu.MountReader("MN_TXT")
		C.Mmu.MountReader("MN___2")
		C.Mmu.MountReader("MN_HGR")
		C.Mmu.MountReader("MN___3")
	case RDCARDRAM:
		C.Mmu.MountReader("AX_ZPS")
		C.Mmu.MountReader("AX___1")
		C.Mmu.MountReader("AX_TXT")
		C.Mmu.MountReader("AX___2")
		C.Mmu.MountReader("AX_HGR")
		C.Mmu.MountReader("AX___3")
	case WRMAINRAM:
		C.Mmu.MountWriter("MN_ZPS")
		C.Mmu.MountWriter("MN___1")
		C.Mmu.MountWriter("MN_TXT")
		C.Mmu.MountWriter("MN___2")
		C.Mmu.MountWriter("MN_HGR")
		C.Mmu.MountWriter("MN___3")
	case WRCARDRAM:
		C.Mmu.MountWriter("AX_ZPS")
		C.Mmu.MountWriter("AX___1")
		C.Mmu.MountWriter("AX_TXT")
		C.Mmu.MountWriter("AX___2")
		C.Mmu.MountWriter("AX_HGR")
		C.Mmu.MountWriter("AX___3")

	case SETSTDZP:
		fallthrough
	case SETALTZP:
		log.Printf("ZP Management: %04X\n", addr+0xC000)

	case SETSLOTCXROM:
		for i := 1; i < 8; i++ {
			C.Mmu.Mount("SLOT_"+strconv.Itoa(i), "")
		}
		is_CX_INT = false
	case SETINTCXROM:
		C.Mmu.Mount("ROM_C", "")
		C.Mmu.Mount("IO", "IO")
		is_CX_INT = true
	case SETINTC3ROM:
		C.Mmu.SwapRom("SLOT_3", "ROM_C")
		is_C3_INT = true
	case SETSLOTC3ROM:
		C.Mmu.Mount("SLOT_3", "")
		is_C3_INT = false

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
