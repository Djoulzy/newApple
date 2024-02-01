package io

import (
	"log"
	"newApple/crtc"
)

func (C *SoftSwitch) Read(addr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch addr {
	case RD80VID:
		// PRINT (PEEK(49183))
		if crtc.Set_80COL == 1 {
			return 0x8D
		}
		return 0x00
	case KBD:
		return C.Buff[KBD]
	case STOREON:
		is_80Store = true
		return 0x8D
	case RD80STORE:
		if is_80Store {
			return 0x8D
		}
		return 0x00
	case KBDSTRB:
		if Is_Keypressed {
			Is_Keypressed = false
			C.Buff[KBD] = 0
			return 0x8D
		}
		Is_Keypressed = false
		C.Buff[KBD] = 0
		return 0x00
	case RDLCBNK2:
		if is_BANK2 {
			return 0x8D
		}
		return 0x00
	case RDLCRAM:
		if is_BS_RAM {
			return 0x8D
		}
		return 0x00
	case RDALTZP:
		if is_ALT_ZP {
			return 0x8D
		}
		return 0x00
	case RSTXINT:
		if is_CX_INT {
			return 0x8D
		} else {
			return 0x00
		}
	case RDC3ROM:
		if is_C3_INT {
			return 0x8D
		} else {
			return 0x00
		}
	case RDRAMRD:
		if is_RAMRD {
			return 0x8D
		} else {
			return 0x00
		}
	case RDRAMWRT:
		if is_RAMWRT {
			return 0x8D
		} else {
			return 0x00
		}

	case CLR80VID:
		C.Video.Set40Cols()
		return 0
	case TXTCLR:
		C.Video.SetGraphMode()
		return 0
	case TXTSET:
		C.Video.SetTexMode()
		return 0
	case MIXCLR:
		C.Video.SetFullMode()
		return 0
	case MIXSET:
		C.Video.SetMixedMode()
		return 0
	case LORES:
		is_HIRES = false
		C.Video.SetLoResMode()
		return 0
	case HIRES:
		is_HIRES = true
		C.Video.SetHiResMode()
		return 0
	case TXTPAGE1:
		is_PAGE2 = false
		C.Video.SetPage1()
		return 0
	case TXTPAGE2:
		if is_80Store {
			// TODO
		} else {
			is_PAGE2 = true
			C.Video.SetPage2()
		}
		return 0

	case RDTEXT:
		if crtc.Set_MODE == 0 {
			return 0x80
		}
		return 0x00
	case RDMIXED:
		if crtc.Set_MIXED == 1 {
			return 0x80
		}
		return 0x00
	case RDPAGE2:
		if is_PAGE2 {
			return 0x80
		}
		return 0x00
	case RDHIRES:
		if is_HIRES {
			return 0x80
		}
		return 0x00
	case RDALTCHAR:
		return 13
	case RSTVBL:
		return 0

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
		return 0x80
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
		return 0x80
	case ROMROB2:
		C.Mmu.Mount("ROM_D", "")
		C.Mmu.Mount("ROM_EF", "")
		is_BANK2 = true
		is_BS_RAM = false
		return 0x80
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
		return 0x80

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
		return 0x80
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
		return 0x80
	case ROMROB1:
		C.Mmu.Mount("ROM_D", "")
		C.Mmu.Mount("ROM_EF", "")
		is_BANK2 = false
		is_BS_RAM = false
		return 0x80
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
		return 0x80

	case 0x0084, 0x0085, 0x0086, 0x0087, 0x008C, 0x008D, 0x008E, 0x008F:
		log.Println("Saturn Card not implemented")
		return 0
	case SPKR:
		return 0
	case SLOT6_OFFSET + DRVSM0:
		C.Disks.SetPhase(0, false)
		return 0
	case SLOT6_OFFSET + DRVSM0 + 1:
		C.Disks.SetPhase(0, true)
		return 0
	case SLOT6_OFFSET + DRVSM1:
		C.Disks.SetPhase(1, false)
		return 0
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disks.SetPhase(1, true)
		return 0
	case SLOT6_OFFSET + DRVSM2:
		C.Disks.SetPhase(2, false)
		return 0
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disks.SetPhase(2, true)
		return 0
	case SLOT6_OFFSET + DRVSM3:
		C.Disks.SetPhase(3, false)
		return 0
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disks.SetPhase(3, true)
		return 0

		/*
			PRINT (PEEK(49386))
			PRINT (PEEK(49385))

			Check protect
			PRINT (PEEK(49389))
			PRINT (PEEK(49390))


			Read data
			PRINT (PEEK(49390))
			PRINT (PEEK(49388))


		*/

	case SLOT6_OFFSET + DRIVE: // $C0E8
		// PRINT (PEEK(49384))
		return C.Disks.diskMotorsOFF()
	case SLOT6_OFFSET + DRIVE + 1:
		// PRINT (PEEK(49385))
		return C.Disks.diskMotorsON()

	case SLOT6_OFFSET + DRVSEL: // $C0EA
		// PRINT (PEEK(49386))
		return C.Disks.driveSelect(0)
	case SLOT6_OFFSET + DRVSEL + 1:
		// PRINT (PEEK(49387))
		return C.Disks.driveSelect(1)

	case SLOT6_OFFSET + DRVDATA: // Q6 $C0EC
		// PRINT (PEEK(49388))
		return C.Disks.ShiftOrRead()
	case SLOT6_OFFSET + DRVDATA + 1: // Q6 $C0ED
		// PRINT (PEEK(49389))
		return C.Disks.LoadOrCheck()

	case SLOT6_OFFSET + DRVWRITE: // Q7 $C0EE
		// PRINT (PEEK(49390))
		return C.Disks.SetSequencerMode(SEQ_READ_MODE)
	case SLOT6_OFFSET + DRVWRITE + 1: // Q7 $C0EF
		// PRINT (PEEK(49391))
		return C.Disks.SetSequencerMode(SEQ_WRITE_MODE)

	case 0x0078, 0x0079, 0x007E:
		// log.Printf("Mouse not supported %04X\n", addr+0xC000)
		return 0x00

	case 0x0060, 0x0061, 0x0062, 0x0063:
		// log.Printf("Keyboard switch not supported %04X\n", addr+0xC000)
		return 0x00

	default:
		log.Printf("IO Read Unknown: %04X\n", addr+0xC000)
		return C.Buff[addr]
	}
}
