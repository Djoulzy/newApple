package io

import (
	"log"
	"strconv"
)

func (C *SoftSwitch) Write(addr uint16, val byte) {
	switch addr {
	case CLR80VID:
		C.Video.Set40Cols()
		is_80COL = false
	case SET80VID:
		C.Video.Set80Cols()
		is_80COL = true
	case STOREOFF:
		is_STORE = false
	case STOREON:
		is_STORE = true
	case KBDSTRB:
		Is_Keypressed = false
		C.Buff[KBD] = 0

	case CLRALTCHAR:
		is_ALTCHAR = false
	case SETALTCHAR:
		is_ALTCHAR = true
	case TXTCLR:
		C.Video.SetGraphMode()
	case TXTSET:
		C.Video.SetTexMode()
	case MIXCLR:
		C.Video.SetFullMode()
	case MIXSET:
		C.Video.SetMixedMode()
	case LORES:
		is_HIRES = false
		C.Video.SetLoResMode()
	case HIRES:
		is_HIRES = true
		C.Video.SetHiResMode()
	case TXTPAGE1:
		if !is_STORE {
			C.Video.SetPage1()
		} else {
			C.Mmu.Mount("MN_TXT", "MN_TXT")
			if is_HIRES {
				C.Mmu.Mount("MN_HGR", "MN_HGR")
			}
		}
		is_PAGE2 = false
	case TXTPAGE2:
		if !is_STORE {
			C.Video.SetPage2()
		} else {
			C.Mmu.Mount("AX_TXT", "AX_TXT")
			if is_HIRES {
				C.Mmu.Mount("AX_HGR", "AX_HGR")
			}
		}
		is_PAGE2 = true

	case RDMAINRAM:
		C.Mmu.MountReader("MN___1")
		C.Mmu.MountReader("MN___2")
		C.Mmu.MountReader("MN___3")
		if is_STORE {
			if !is_HIRES {
				C.Mmu.MountReader("MN_HGR")
			}
		} else {
			C.Mmu.MountReader("MN_TXT")
		}
		is_RAMRD = false
	case RDCARDRAM:
		C.Mmu.MountReader("AX___1")
		C.Mmu.MountReader("AX___2")
		C.Mmu.MountReader("AX___3")
		if is_STORE {
			if !is_HIRES {
				C.Mmu.MountReader("AX_HGR")
			}
		} else {
			C.Mmu.MountReader("AX_TXT")
		}
		is_RAMRD = true
	case WRMAINRAM:
		C.Mmu.MountWriter("MN___1")
		C.Mmu.MountWriter("MN___2")
		C.Mmu.MountWriter("MN___3")
		if is_STORE {
			if !is_HIRES {
				C.Mmu.MountWriter("MN_HGR")
			}
		} else {
			C.Mmu.MountWriter("MN_TXT")
		}
		is_RAMWRT = false
	case WRCARDRAM:
		C.Mmu.MountWriter("AX___1")
		C.Mmu.MountWriter("AX___2")
		C.Mmu.MountWriter("AX___3")
		if is_STORE {
			if !is_HIRES {
				C.Mmu.MountWriter("AX_HGR")
			}
		} else {
			C.Mmu.MountWriter("AX_TXT")
		}
		is_RAMWRT = true

	case RAMROB2:
		if is_ALT_ZP {
			C.Mmu.Mount("AX_BK2", "")
			C.Mmu.Mount("AX___4", "")
		} else {
			C.Mmu.Mount("MN_BK2", "")
			C.Mmu.Mount("MN___4", "")
		}
		is_BANK2 = true
		is_BS_RAM = true
	case ROMWB2:
		if is_ALT_ZP {
			C.Mmu.Mount("ROM_D", "AX_BK2")
			C.Mmu.Mount("ROM_EF", "AX___4")
		} else {
			C.Mmu.Mount("ROM_D", "MN_BK2")
			C.Mmu.Mount("ROM_EF", "MN___4")
		}
		is_BANK2 = true
		is_BS_RAM = false
	case ROMROB2:
		C.Mmu.Mount("ROM_D", "")
		C.Mmu.Mount("ROM_EF", "")
		is_BANK2 = true
		is_BS_RAM = false
	case RAMRWB2:
		if is_ALT_ZP {
			C.Mmu.Mount("AX_BK2", "AX_BK2")
			C.Mmu.Mount("AX___4", "AX___4")
		} else {
			C.Mmu.Mount("MN_BK2", "MN_BK2")
			C.Mmu.Mount("MN___4", "MN___4")
		}
		is_BS_RAM = true
		is_BANK2 = true

	case RAMROB1:
		if is_ALT_ZP {
			C.Mmu.Mount("AX_BK1", "")
			C.Mmu.Mount("AX___4", "")
		} else {
			C.Mmu.Mount("MN_BK1", "")
			C.Mmu.Mount("MN___4", "")
		}
		is_BANK2 = false
		is_BS_RAM = true
	case ROMWB1:
		if is_ALT_ZP {
			C.Mmu.Mount("ROM_D", "AX_BK1")
			C.Mmu.Mount("ROM_EF", "AX___4")
		} else {
			C.Mmu.Mount("ROM_D", "MN_BK1")
			C.Mmu.Mount("ROM_EF", "MN___4")
		}
		is_BANK2 = false
		is_BS_RAM = false
	case ROMROB1:
		C.Mmu.Mount("ROM_D", "")
		C.Mmu.Mount("ROM_EF", "")
		is_BANK2 = false
		is_BS_RAM = false
	case RAMRWB1:
		if is_ALT_ZP {
			C.Mmu.Mount("AX_BK1", "AX_BK1")
			C.Mmu.Mount("AX___4", "AX___4")
		} else {
			C.Mmu.Mount("MN_BK1", "MN_BK1")
			C.Mmu.Mount("MN___4", "MN___4")
		}
		is_BS_RAM = true
		is_BANK2 = false

	case SETSTDZP:
		is_ALT_ZP = false
		C.Mmu.SwapChip("AX_ZPS", "MN_ZPS")
		C.Mmu.SwapChip("AX_BK1", "MN_BK1")
		C.Mmu.SwapChip("AX_BK2", "MN_BK2")
		C.Mmu.SwapChip("AX___4", "MN___4")
	case SETALTZP:
		is_ALT_ZP = true
		C.Mmu.SwapChip("MN_ZPS", "AX_ZPS")
		C.Mmu.SwapChip("MN_BK1", "AX_BK1")
		C.Mmu.SwapChip("MN_BK2", "AX_BK2")
		C.Mmu.SwapChip("MN___4", "AX___4")

	case SETSLOTCXROM:
		for i := 1; i < 8; i++ {
			if i == 3 {
				continue
			}
			C.Mmu.Mount("SLOT_"+strconv.Itoa(i), "")
		}
		is_CX_INT = false
	case SETINTCXROM:
		C.Mmu.Mount("ROM_C", "")
		C.Mmu.Mount("IO", "IO")
		is_CX_INT = true
	case SETINTC3ROM:
		C.Mmu.SwapChip("SLOT_3", "ROM_C")
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
		log.Printf("IO Write Unknown: %04X (value: %02X)\n", addr+0xC000, val)
		C.Buff[addr] = val
	}
	// mem[translatedAddr] = val
}
