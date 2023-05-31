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
		if crtc.Is_80COL {
			return 0x8D
		}
		return 0x00
	case KBD:
		return C.Buff[KBD]
	case STOREON:
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
		if is_READ_RAM {
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
	case TXTCLR:
		crtc.Is_TEXTMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case TXTSET:
		crtc.Is_TEXTMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case MIXCLR:
		crtc.Is_MIXEDMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case MIXSET:
		crtc.Is_MIXEDMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case LORES:
		// log.Printf("HIRES OFF %04X", cpu.PC)
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case HIRES:
		// log.Printf("HIRES ON %04X", cpu.PC)
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case TXTPAGE1:
		crtc.Is_PAGE2 = false
		C.Video.UpdateVideoRam()
		return 0
	case TXTPAGE2:
		crtc.Is_PAGE2 = true
		C.Video.UpdateVideoRam()
		return 0

	case RDTEXT:
		if crtc.Is_TEXTMODE {
			return 0x80
		}
		return 0x00
	// case MIXED:
	// 	if crtc.Is_MIXEDMODE {
	// 		return 0x80
	// 	}
	// 	return 0x00
	// case PAGE2:
	// 	if crtc.Is_PAGE2 {
	// 		return 0x80
	// 	}
	// 	return 0x00
	// case HIRES:
	// 	if crtc.Is_HIRESMODE {
	// 		return 0x80
	// 	}
	// 	return 0x00

	// case RDRAM_B2:
	// 	log.Println("RDRAM_B2")
	// 	C.Mmu.Disable("ROM_D")
	// 	C.Mmu.Disable("ROM_EF")
	// 	C.Mmu.Disable("BANK1")
	// 	C.Mmu.Enable("BANK2")
	// 	C.Mmu.ReadOnly("BANK2")
	// 	is_READ_RAM = true
	// 	is_BANK2 = true
	// 	return 0x80
	// case RDROM_WB2:
	// 	log.Println("RDROM_WB2")
	// 	C.Mmu.Enable("ROM_D")
	// 	C.Mmu.Enable("ROM_EF")
	// 	C.Mmu.Disable("BANK1")
	// 	C.Mmu.Enable("BANK2")
	// 	C.Mmu.ReadWrite("BANK2")
	// 	is_READ_RAM = false
	// 	is_BANK2 = true
	// 	return 0x80
	// case RDROM_2:
	// 	log.Println("RDROM_2")
	// 	C.Mmu.Enable("ROM_D")
	// 	C.Mmu.Enable("ROM_EF")
	// 	C.Mmu.Disable("BANK1")
	// 	C.Mmu.Disable("BANK2")
	// 	is_READ_RAM = false
	// 	is_BANK2 = false
	// 	return 0x80
	// case RWRAM_B2:
	// 	log.Println("RWRAM_B2")
	// 	C.Mmu.Disable("ROM_D")
	// 	C.Mmu.Disable("ROM_EF")
	// 	C.Mmu.Disable("BANK1")
	// 	C.Mmu.Enable("BANK2")
	// 	C.Mmu.ReadWrite("BANK2")
	// 	is_READ_RAM = true
	// 	is_BANK2 = true
	// 	return 0x80

	// case RDROM_1:
	// 	log.Println("RDROM_1")
	// 	C.Mmu.Enable("ROM_D")
	// 	C.Mmu.Enable("ROM_EF")
	// 	C.Mmu.Disable("BANK1")
	// 	C.Mmu.Disable("BANK2")
	// 	is_READ_RAM = false
	// 	is_BANK2 = false
	// 	return 0x80
	// case RDRAM_B1:
	// 	log.Println("RDRAM_B1")
	// 	C.Mmu.Disable("ROM_D")
	// 	C.Mmu.Disable("ROM_EF")
	// 	C.Mmu.Enable("BANK1")
	// 	C.Mmu.Disable("BANK2")
	// 	C.Mmu.ReadOnly("BANK1")
	// 	is_READ_RAM = true
	// 	is_BANK2 = false
	// 	return 0x80
	// case RDROM_WB1:
	// 	log.Println("RDROM_WB1")
	// 	C.Mmu.Enable("ROM_D")
	// 	C.Mmu.Enable("ROM_EF")
	// 	C.Mmu.Enable("BANK1")
	// 	C.Mmu.Disable("BANK2")
	// 	C.Mmu.ReadWrite("BANK1")
	// 	is_READ_RAM = false
	// 	is_BANK2 = false
	// 	return 0x80
	// case RWRAM_B1:
	// 	log.Println("RWRAM_B1")
	// 	C.Mmu.Disable("ROM_D")
	// 	C.Mmu.Disable("ROM_EF")
	// 	C.Mmu.Enable("BANK1")
	// 	C.Mmu.Disable("BANK2")
	// 	C.Mmu.ReadWrite("BANK1")
	// 	is_READ_RAM = true
	// 	is_BANK2 = false
	// 	return 0x80

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
		log.Printf("Mouse not supported %04X\n", addr)
		return 0x00

	case 0x0060, 0x0061, 0x0062, 0x0063:
		log.Printf("Keyboard switch not supported %04X\n", addr)
		return 0x00

	default:
		log.Printf("Read Unknown: %04X\n", addr)
		return 0x00
	}
}
