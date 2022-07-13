package main

import (
	"log"
	"newApple/crtc"
	"newApple/disk"
)

const (
	// MEMORY MANAGEMENT SOFT SWITCHES (W)
	_80STOREOFF  = 0x00
	_80STOREON   = 0x01
	RAMRDON      = 0x03
	RAMRDOFF     = 0x02
	RAMWRTON     = 0x05
	RAMWRTOFF    = 0x04
	INTCXROMOFF  = 0x06
	INTCXROMON   = 0x07
	ALZTPOFF     = 0x08
	ALZTPON      = 0x09
	SLOTC3ROMOFF = 0x0A
	SLOTC3ROMON  = 0x0B
	BSRBANK2     = 0x11
	BSRREADRAM   = 0x12

	// VIDEO SOFT SWITCHES (W/R)
	_80COLOFF     = 0x0C
	_80COLON      = 0x0D
	RAMRD         = 0x13
	RAMWRT        = 0x14
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
	ALTZP      = 0x16
	TEXT       = 0x1A
	MIXED      = 0x1B
	PAGE2      = 0x1C
	HIRES      = 0x1D
	ALTCHARSET = 0x1E
	_80COL     = 0x1F
	_80STORE   = 0x18

	// BANK SWITCHING
	RDRAM_B2  = 0x80
	RDROM_WB2 = 0x81
	RDROM_2   = 0x82
	RWRAM_B2  = 0x83
	RDRAM_B1  = 0x88
	RDROM_WB1 = 0x89
	RDROM_1   = 0x8A
	RWRAM_B1  = 0x8B

	SATURN_CTRL1 = 0x84
	SATURN_CTRL2 = 0x85
	SATURN_CTRL3 = 0x86
	SATURN_CTRL4 = 0x87
	SATURN1      = 0x8C
	SATURN2      = 0x8D
	SATURN3      = 0x8E
	SATURN4      = 0x8F

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

type io_access struct {
	Disks [2]*disk.DRIVE
	Video *crtc.CRTC

	connectedDrive int
}

func InitIO(d1 *disk.DRIVE, d2 *disk.DRIVE, vid *crtc.CRTC) *io_access {
	tmp := io_access{}
	tmp.Video = vid
	tmp.connectedDrive = 0
	if d1 == nil && d2 != nil {
		d1 = d2
		d2 = nil
	}
	if d1 != nil {
		tmp.Disks[0] = d1
		tmp.connectedDrive++
	}
	if d2 != nil {
		tmp.Disks[1] = d2
		tmp.connectedDrive++
	}

	log.Printf("DiskII drive connected: %d\n", tmp.connectedDrive)
	return &tmp
}

func (C *io_access) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", translatedAddr)
	switch translatedAddr {
	case _80COL:
		// PRINT (PEEK(49183))
		if crtc.Is_80COL {
			return 0x8D
		}
		return 0x00
	case _80STOREOFF:
		return mem[_80STOREOFF]
	case _80STORE:
		if is_80Store {
			return 0x8D
		}
		return 0x00
	case AKD:
		if is_Keypressed {
			is_Keypressed = false
			mem[_80STOREOFF] = 0
			return 0x8D
		}
		is_Keypressed = false
		mem[_80STOREOFF] = 0
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
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadOnly("BANK2")
		is_READ_RAM = true
		is_BANK2 = true
		return 0x80
	case RDROM_WB2:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadWrite("BANK2")
		is_READ_RAM = false
		is_BANK2 = true
		return 0x80
	case RDROM_2:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Disable("BANK2")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RWRAM_B2:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadWrite("BANK2")
		is_READ_RAM = true
		is_BANK2 = true
		return 0x80

	case RDROM_1:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Disable("BANK2")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RDRAM_B1:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadOnly("BANK1")
		is_READ_RAM = true
		is_BANK2 = false
		return 0x80
	case RDROM_WB1:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadWrite("BANK1")
		is_READ_RAM = false
		is_BANK2 = false
		return 0x80
	case RWRAM_B1:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadWrite("BANK1")
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
		C.Disks[SelectedDrive].SetPhase(0, false)
		return 0
	case SLOT6_OFFSET + DRVSM0 + 1:
		C.Disks[SelectedDrive].SetPhase(0, true)
		return 0
	case SLOT6_OFFSET + DRVSM1:
		C.Disks[SelectedDrive].SetPhase(1, false)
		return 0
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disks[SelectedDrive].SetPhase(1, true)
		return 0
	case SLOT6_OFFSET + DRVSM2:
		C.Disks[SelectedDrive].SetPhase(2, false)
		return 0
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disks[SelectedDrive].SetPhase(2, true)
		return 0
	case SLOT6_OFFSET + DRVSM3:
		C.Disks[SelectedDrive].SetPhase(3, false)
		return 0
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disks[SelectedDrive].SetPhase(3, true)
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
		return C.diskMotorsOFF()
	case SLOT6_OFFSET + DRIVE + 1:
		// PRINT (PEEK(49385))
		return C.diskMotorsON()

	case SLOT6_OFFSET + DRVSEL: // $C0EA
		// PRINT (PEEK(49386))
		return C.driveSelect(0)
	case SLOT6_OFFSET + DRVSEL + 1:
		// PRINT (PEEK(49387))
		return C.driveSelect(1)

	case SLOT6_OFFSET + DRVDATA: // Q6 $C0EC
		// PRINT (PEEK(49388))
		return C.ShiftOrRead()
	case SLOT6_OFFSET + DRVDATA + 1: // Q6 $C0ED
		// PRINT (PEEK(49389))
		return C.LoadOrCheck()

	case SLOT6_OFFSET + DRVWRITE: // Q7 $C0EE
		// PRINT (PEEK(49390))
		return C.SetSequencerMode(SEQ_READ_MODE)
	case SLOT6_OFFSET + DRVWRITE + 1: // Q7 $C0EF
		// PRINT (PEEK(49391))
		return C.SetSequencerMode(SEQ_WRITE_MODE)

	default:
		// log.Printf("Read Unknown: %02X\n", translatedAddr)
		return 0x00
	}
}

func (C *io_access) MWrite(mem []byte, translatedAddr uint16, val byte) {
	switch translatedAddr {
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
		mem[_80STOREOFF] = 0
	case ALZTPOFF:
		log.Printf("ALT_ZP Off")
		is_ALT_ZP = false
		MEM.Enable("ZP")
		MEM.Disable("ALT_ZP")
	case ALZTPON:
		log.Printf("ALT_ZP On")
		is_ALT_ZP = true
		MEM.Enable("ALT_ZP")
		MEM.Disable("ZP")
	case INTCXROMOFF:
		is_CX_INT = false
		MEM.Enable("SLOT1")
		MEM.Enable("SLOT2")
		MEM.Enable("SLOT3")
		MEM.Enable("SLOT4")
		MEM.Enable("SLOT5")
		MEM.Enable("SLOT6")
		MEM.Enable("SLOT7")
	case INTCXROMON:
		is_CX_INT = true
		MEM.Disable("SLOT1")
		MEM.Disable("SLOT2")
		MEM.Disable("SLOT3")
		MEM.Disable("SLOT4")
		MEM.Disable("SLOT5")
		MEM.Disable("SLOT6")
		MEM.Disable("SLOT7")
	case SLOTC3ROMON:
		MEM.Enable("SLOT3")
		is_C3_INT = false
	case SLOTC3ROMOFF:
		MEM.Disable("SLOT3")
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
			// MEM.Disable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}
	case PAGE2ON:
		crtc.Is_PAGE2 = true
		if is_80Store {
			// MEM.Enable("AUX")
		} else {
			C.Video.UpdateVideoRam()
		}

	case RDRAM_B2:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadOnly("BANK2")
	case RDROM_WB2:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadWrite("BANK2")
	case RDROM_2:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Disable("BANK2")
	case RWRAM_B2:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Enable("BANK2")
		MEM.ReadWrite("BANK2")

	case RDROM_1:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Disable("BANK1")
		MEM.Disable("BANK2")
	case RDRAM_B1:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadOnly("BANK1")
	case RDROM_WB1:
		MEM.Enable("ROM_D")
		MEM.Enable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadWrite("BANK1")
	case RWRAM_B1:
		MEM.Disable("ROM_D")
		MEM.Disable("ROM_EF")
		MEM.Enable("BANK1")
		MEM.Disable("BANK2")
		MEM.ReadWrite("BANK1")

	case SLOT6_OFFSET + DRVSM0:
		C.Disks[SelectedDrive].SetPhase(0, false)
	case SLOT6_OFFSET + DRVSM0 + 1:
		C.Disks[SelectedDrive].SetPhase(0, true)
	case SLOT6_OFFSET + DRVSM1:
		C.Disks[SelectedDrive].SetPhase(1, false)
	case SLOT6_OFFSET + DRVSM1 + 1:
		C.Disks[SelectedDrive].SetPhase(1, true)
	case SLOT6_OFFSET + DRVSM2:
		C.Disks[SelectedDrive].SetPhase(2, false)
	case SLOT6_OFFSET + DRVSM2 + 1:
		C.Disks[SelectedDrive].SetPhase(2, true)
	case SLOT6_OFFSET + DRVSM3:
		C.Disks[SelectedDrive].SetPhase(3, false)
	case SLOT6_OFFSET + DRVSM3 + 1:
		C.Disks[SelectedDrive].SetPhase(3, true)

	case SLOT6_OFFSET + DRIVE:
		// PRINT (PEEK(49384))
		C.diskMotorsOFF()
	case SLOT6_OFFSET + DRIVE + 1:
		// PRINT (PEEK(49385))
		C.diskMotorsON()

	case SLOT6_OFFSET + DRVSEL:
		// PRINT (PEEK(49290))
		C.driveSelect(0)
	case SLOT6_OFFSET + DRVSEL + 1:
		// PRINT (PEEK(49291))
		C.driveSelect(1)

	case SLOT6_OFFSET + DRVDATA: // Q6
		// PRINT (PEEK(49292))
		C.ShiftOrRead()
	case SLOT6_OFFSET + DRVDATA + 1: // Q6
		// PRINT (PEEK(49293))
		C.LoadOrCheck()

	case SLOT6_OFFSET + DRVWRITE: // Q7
		// PRINT (PEEK(49390))
		C.SetSequencerMode(SEQ_READ_MODE)
	case SLOT6_OFFSET + DRVWRITE + 1: // Q7
		C.SetSequencerMode(SEQ_WRITE_MODE)

	default:
		// log.Printf("Write Unknown: %02X\n", translatedAddr)
	}
	// mem[translatedAddr] = val
}
