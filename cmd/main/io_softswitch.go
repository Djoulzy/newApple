package main

import (
	"log"
	"newApple/crtc"

	"github.com/Djoulzy/mmu"
)

const (
	// _80STOREOFF   = 0x0000 // MEMORY MANAGEMENT SOFT SWITCHES (W)
	// _80STOREON    = 0x0001
	// RAMRDON       = 0x0003
	// RAMRDOFF      = 0x0002
	// RAMWRTON      = 0x0005
	// RAMWRTOFF     = 0x0004
	// INTCXROMOFF   = 0x0006
	// INTCXROMON    = 0x0007
	// ALZTPOFF      = 0x0008
	// ALZTPON       = 0x0009
	// SLOTC3ROMOFF  = 0x000A
	// SLOTC3ROMON   = 0x000B
	// BSRBANK2      = 0x0011
	// BSRREADRAM    = 0x0012
	// _80COLOFF     = 0x000C // VIDEO SOFT SWITCHES (W/R)
	// _80COLON      = 0x000D
	// RAMRD         = 0x0013
	// RAMWRT        = 0x0014
	// ALTCHARSETOFF = 0x000E
	// ALTCHARSETON  = 0x000F
	// TEXTOFF       = 0x0050
	// TEXTON        = 0x0051
	// MIXEDOFF      = 0x0052
	// MIXEDON       = 0x0053
	// PAGE2OFF      = 0x0054
	// PAGE2ON       = 0x0055
	// HIRESOFF      = 0x0056
	// HIRESON       = 0x0057
	// DISXY         = 0x0058 // Mouse
	// ENBXY         = 0x0059
	// DISVBL        = 0x005A
	// ENVBL         = 0x005B
	// X0EDGEON      = 0x005C
	// X0EDGEOFF     = 0x005D
	// Y0EDGEON      = 0x005E
	// Y0EDGEOFF     = 0x005F
	// RDMOUX1       = 0x0066
	// RDMOUY1       = 0x0067
	// IOUDISABLE    = 0x0078
	// IOUENABLE     = 0x0079
	// IOUDISON      = 0x007E
	// IOUDISOFF     = 0x007F
	// KEY4080       = 0x0060 // Keyboard Switch 40/80
	// RDBTN0        = 0x0061
	// BUTN1         = 0x0062
	// AKD           = 0x0010 // SOFT SWITCH STATUS FLAGS (R bit 7)
	// INTCXROM      = 0x0015
	// SLOTC3ROM     = 0x0017
	// ALTZP         = 0x0016
	// TEXT          = 0x001A
	// MIXED         = 0x001B
	// PAGE2         = 0x001C
	// HIRES         = 0x001D
	// ALTCHARSET    = 0x001E
	// _80COL        = 0x001F
	// _80STORE      = 0x0018
	// RDRAM_B2      = 0x0080 // BANK SWITCHING
	// RDROM_WB2     = 0x0081
	// RDROM_2       = 0x0082
	// RWRAM_B2      = 0x0083
	// RDRAM_B1      = 0x0088
	// RDROM_WB1     = 0x0089
	// RDROM_1       = 0x008A
	// RWRAM_B1      = 0x008B
	// SATURN_CTRL1  = 0x0084 // SATURN CARD
	// SATURN_CTRL2  = 0x0085
	// SATURN_CTRL3  = 0x0086
	// SATURN_CTRL4  = 0x0087
	// SATURN1       = 0x008C
	// SATURN2       = 0x008D
	// SATURN3       = 0x008E
	// SATURN4       = 0x008F
	// SPKR          = 0x0030 // OTHER
	// SLOT0_OFFSET  = 0x0090 // SLOTS
	// SLOT1_OFFSET  = 0x0090
	// SLOT2_OFFSET  = 0x00A0
	// SLOT3_OFFSET  = 0x00B0
	// SLOT4_OFFSET  = 0x00C0
	// SLOT5_OFFSET  = 0x00D0
	// SLOT6_OFFSET  = 0x00E0
	// SLOT7_OFFSET  = 0x00F0

	// DRIVE OPERATIONS
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

	case DISXY:
		fallthrough
	case ENBXY:
		fallthrough
	case DISVBL:
		fallthrough
	case ENVBL:
		fallthrough
	case X0EDGEON:
		fallthrough
	case X0EDGEOFF:
		fallthrough
	case Y0EDGEON:
		fallthrough
	case Y0EDGEOFF:
		fallthrough
	case RDMOUX1:
		fallthrough
	case RDMOUY1:
		fallthrough
	case IOUDISABLE:
		fallthrough
	case IOUENABLE:
		fallthrough
	case IOUDISON:
		fallthrough
	case IOUDISOFF:
		log.Printf("Mouse not supported %04X\n", addr)
		return 0x00

	case KEY4080:
		fallthrough
	case RDBTN0:
		fallthrough
	case BUTN1:
		log.Printf("Keyboard switch not supported %04X\n", addr)
		return 0x00

	default:
		log.Printf("Read Unknown: %04X\n", addr)
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
		log.Printf("Write Unknown: %04X\n", addr)
	}
	// mem[translatedAddr] = val
}

// func (C *SoftSwitch) MWriteUnder(mem []C.Mmu.MEMCell, addr uint16, value byte) {
// 	log.Printf("WRITE UNDER")
// 	C.Buff[addr].Under = value
// }
