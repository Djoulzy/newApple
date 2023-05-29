package main

import (
	"log"
	"newApple/crtc"

	"github.com/Djoulzy/mmu"
)

const (
	// MEMORY MANAGEMENT SOFT SWITCHES (W)
	_80STOREOFF  = 0xC000
	_80STOREON   = 0xC001
	RAMRDON      = 0xC003
	RAMRDOFF     = 0xC002
	RAMWRTON     = 0xC005
	RAMWRTOFF    = 0xC004
	INTCXROMOFF  = 0xC006
	INTCXROMON   = 0xC007
	ALZTPOFF     = 0xC008
	ALZTPON      = 0xC009
	SLOTC3ROMOFF = 0xC00A
	SLOTC3ROMON  = 0xC00B
	BSRBANK2     = 0xC011
	BSRREADRAM   = 0xC012

	// VIDEO SOFT SWITCHES (W/R)
	_80COLOFF     = 0xC00C
	_80COLON      = 0xC00D
	RAMRD         = 0xC013
	RAMWRT        = 0xC014
	ALTCHARSETOFF = 0xC00E
	ALTCHARSETON  = 0xC00F
	TEXTOFF       = 0xC050
	TEXTON        = 0xC051
	MIXEDOFF      = 0xC052
	MIXEDON       = 0xC053
	PAGE2OFF      = 0xC054
	PAGE2ON       = 0xC055
	HIRESOFF      = 0xC056
	HIRESON       = 0xC057

	// SOFT SWITCH STATUS FLAGS (R bit 7)
	AKD        = 0xC010
	INTCXROM   = 0xC015
	SLOTC3ROM  = 0xC017
	ALTZP      = 0xC016
	TEXT       = 0xC01A
	MIXED      = 0xC01B
	PAGE2      = 0xC01C
	HIRES      = 0xC01D
	ALTCHARSET = 0xC01E
	_80COL     = 0xC01F
	_80STORE   = 0xC018

	// BANK SWITCHING
	RDRAM_B2  = 0xC080
	RDROM_WB2 = 0xC081
	RDROM_2   = 0xC082
	RWRAM_B2  = 0xC083
	RDRAM_B1  = 0xC088
	RDROM_WB1 = 0xC089
	RDROM_1   = 0xC08A
	RWRAM_B1  = 0xC08B

	SATURN_CTRL1 = 0xC084
	SATURN_CTRL2 = 0xC085
	SATURN_CTRL3 = 0xC086
	SATURN_CTRL4 = 0xC087
	SATURN1      = 0xC08C
	SATURN2      = 0xC08D
	SATURN3      = 0xC08E
	SATURN4      = 0xC08F

	// OTHER
	SPKR = 0x30

	// SLOTS
	SLOT0_OFFSET = 0xC090
	SLOT1_OFFSET = 0xC090
	SLOT2_OFFSET = 0xC0A0
	SLOT3_OFFSET = 0xC0B0
	SLOT4_OFFSET = 0xC0C0
	SLOT5_OFFSET = 0xC0D0
	SLOT6_OFFSET = 0xC0E0
	SLOT7_OFFSET = 0xC0F0

	DRVSM0   = 0x00 // Q0
	DRVSM1   = 0x02 // Q1
	DRVSM2   = 0x04 // Q2
	DRVSM3   = 0x06 // Q3
	DRIVE    = 0x08 // Q4
	DRVSEL   = 0x0A // Q5
	DRVDATA  = 0x0C // Q6
	DRVWRITE = 0x0E // Q7
)

// PRINT PEEK(49173)

var (
	is_READ_RAM   bool = false
	is_BANK2      bool = false
	is_C3_INT     bool = true
	is_CX_INT     bool = false
	is_Keypressed bool = false
	is_80Store    bool = false
	is_ALT_ZP     bool = false
)

type SoftSwitch struct {
	mmu.IC

	Disks *DiskInterface
	Video *crtc.CRTC
}

func InitSoftSwitch(name string, size int, disk *DiskInterface, vid *crtc.CRTC) *SoftSwitch {
	tmp := SoftSwitch{
		Disks: disk,
		Video: vid,
	}
	tmp.Name = name
	tmp.Buff = make([]byte, size)
	// if d1 == nil && d2 != nil {
	// 	d1 = d2
	// 	d2 = nil
	// }
	// if d1 != nil {
	// 	tmp.Disks[0] = d1
	// 	tmp.connectedDrive++
	// }
	// if d2 != nil {
	// 	tmp.Disks[1] = d2
	// 	tmp.connectedDrive++
	// }

	return &tmp
}

func (C *SoftSwitch) ReadOnly() bool {
	return false
}

func (C *SoftSwitch) Read(addr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch addr {
	case _80COL:
		// PRINT (PEEK(49183))
		if crtc.Is_80COL {
			return 0x8D
		}
		return 0x00
	case _80STOREOFF:
		return C.Buff[_80STOREOFF]
	case _80STORE:
		if is_80Store {
			return 0x8D
		}
		return 0x00
	case AKD:
		if is_Keypressed {
			is_Keypressed = false
			C.Buff[_80STOREOFF] = 0
			return 0x8D
		}
		is_Keypressed = false
		C.Buff[_80STOREOFF] = 0
		return 0x00
	case BSRBANK2:
		if is_BANK2 {
			return 0x8D
		}
		return 0x00
	case BSRREADRAM:
		if is_READ_RAM {
			return 0x8D
		}
		return 0x00
	case ALTZP:
		if is_ALT_ZP {
			return 0x8D
		}
		return 0x00
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
		// log.Printf("HIRES OFF %04X", cpu.PC)
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
		return 0
	case HIRESON:
		// log.Printf("HIRES ON %04X", cpu.PC)
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
		return 0
	case PAGE2OFF:
		crtc.Is_PAGE2 = false
		C.Video.UpdateVideoRam()
		return 0
	case PAGE2ON:
		crtc.Is_PAGE2 = true
		C.Video.UpdateVideoRam()
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

	case RDRAM_B2:
		log.Println("RDRAM_B2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadOnly("BANK2")
		is_READ_RAM = true
		is_BANK2 = true
		return 0x80
	case RDROM_WB2:
		log.Println("RDROM_WB2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadWrite("BANK2")
		is_READ_RAM = false
		is_BANK2 = true
		return 0x80
	case RDROM_2:
		log.Println("RDROM_2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Disable("BANK2")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RWRAM_B2:
		log.Println("RWRAM_B2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadWrite("BANK2")
		is_READ_RAM = true
		is_BANK2 = true
		return 0x80

	case RDROM_1:
		log.Println("RDROM_1")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Disable("BANK2")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RDRAM_B1:
		log.Println("RDRAM_B1")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadOnly("BANK1")
		is_READ_RAM = true
		is_BANK2 = false
		return 0x80
	case RDROM_WB1:
		log.Println("RDROM_WB1")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadWrite("BANK1")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RWRAM_B1:
		log.Println("RWRAM_B1")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadWrite("BANK1")
		is_READ_RAM = true
		is_BANK2 = false
		return 0x80

	case SATURN_CTRL1:
		fallthrough
	case SATURN_CTRL2:
		fallthrough
	case SATURN_CTRL3:
		fallthrough
	case SATURN_CTRL4:
		fallthrough
	case SATURN1:
		fallthrough
	case SATURN2:
		fallthrough
	case SATURN3:
		fallthrough
	case SATURN4:
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

	default:
		log.Printf("Read Unknown: %02X\n", addr)
		return 0x00
	}
}

func (C *SoftSwitch) Write(addr uint16, val byte) {
	switch addr {
	case _80COLOFF:
		crtc.Is_80COL = false
		C.Video.UpdateGraphMode()
	case _80COLON:
		crtc.Is_80COL = true
		C.Video.UpdateGraphMode()
	case _80STOREOFF:
		is_80Store = false
	case _80STOREON:
		is_80Store = true
	case AKD:
		is_Keypressed = false
		C.Buff[_80STOREOFF] = 0
	case ALZTPOFF:
		log.Printf("ALT_ZP Off")
		is_ALT_ZP = false
		C.Mmu.Enable("ZP")
		C.Mmu.Disable("ALT_ZP")
	case ALZTPON:
		log.Printf("ALT_ZP On")
		is_ALT_ZP = true
		C.Mmu.Enable("ALT_ZP")
		C.Mmu.Disable("ZP")
	case INTCXROMOFF:
		log.Println("INTCXROMOFF")
		is_CX_INT = false
		C.Mmu.Enable("SLOT1")
		C.Mmu.Enable("SLOT2")
		C.Mmu.Enable("SLOT3")
		C.Mmu.Enable("SLOT4")
		C.Mmu.Enable("SLOT5")
		C.Mmu.Enable("SLOT6")
		C.Mmu.Enable("SLOT7")
	case INTCXROMON:
		log.Println("INTCXROMON")
		is_CX_INT = true
		C.Mmu.Disable("SLOT1")
		C.Mmu.Disable("SLOT2")
		C.Mmu.Disable("SLOT3")
		C.Mmu.Disable("SLOT4")
		C.Mmu.Disable("SLOT5")
		C.Mmu.Disable("SLOT6")
		C.Mmu.Disable("SLOT7")
	case SLOTC3ROMON:
		C.Mmu.Enable("SLOT3")
		is_C3_INT = false
	case SLOTC3ROMOFF:
		C.Mmu.Disable("SLOT3")
		is_C3_INT = true
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
		log.Printf("HIRES OFF")
		crtc.Is_HIRESMODE = false
		C.Video.UpdateGraphMode()
	case HIRESON:
		log.Printf("HIRES ON")
		crtc.Is_HIRESMODE = true
		C.Video.UpdateGraphMode()
	case PAGE2OFF:
		crtc.Is_PAGE2 = false
		if is_80Store {
			// C.Mmu.Disable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}
	case PAGE2ON:
		crtc.Is_PAGE2 = true
		if is_80Store {
			// C.Mmu.Enable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}

	case RDRAM_B2:
		log.Println("RDROM_WB2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadOnly("BANK2")
	case RDROM_WB2:
		log.Println("RDROM_WB2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadWrite("BANK2")
	case RDROM_2:
		log.Println("RDROM_WB2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Disable("BANK2")
	case RWRAM_B2:
		log.Println("RDROM_WB2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Enable("BANK2")
		C.Mmu.ReadWrite("BANK2")

	case RDROM_1:
		log.Println("RDROM_WB2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Disable("BANK1")
		C.Mmu.Disable("BANK2")
	case RDRAM_B1:
		log.Println("RDROM_WB2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadOnly("BANK1")
	case RDROM_WB1:
		log.Println("RDROM_WB2")
		C.Mmu.Enable("ROM_D")
		C.Mmu.Enable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadWrite("BANK1")
	case RWRAM_B1:
		log.Println("RDROM_WB2")
		C.Mmu.Disable("ROM_D")
		C.Mmu.Disable("ROM_EF")
		C.Mmu.Enable("BANK1")
		C.Mmu.Disable("BANK2")
		C.Mmu.ReadWrite("BANK1")

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
		// log.Printf("Write Unknown: %02X\n", translatedAddr)
	}
	// mem[translatedAddr] = val
}

// func (C *SoftSwitch) MWriteUnder(mem []C.Mmu.MEMCell, addr uint16, value byte) {
// 	log.Printf("WRITE UNDER")
// 	C.Buff[addr].Under = value
// }
