package io

import (
	"fmt"
	"log"
	"newApple/config"
	"newApple/diskdrive"
)

const (
	SEQ_READ_MODE  = true
	SEQ_WRITE_MODE = false
)

type DiskInterface struct {
	Disks          []*diskdrive.DRIVE
	connectedDrive int

	MotorIsOn     bool
	SelectedDrive int
	SequencerMode bool
	ProtectCheck  byte
}

func InitDiskInterface(conf *config.ConfigData) *DiskInterface {
	tmp := DiskInterface{
		Disks:          make([]*diskdrive.DRIVE, 2),
		connectedDrive: 0,
		SelectedDrive:  0,
		SequencerMode:  false,
		ProtectCheck:   0,
	}
	tmp.loadDisks(conf)
	return &tmp
}

func (C *DiskInterface) loadDisks(conf *config.ConfigData) {
	if conf.Slots.Slot6 != "" {
		C.connectedDrive = 0
		if conf.Disks.Disk1 != "" {
			C.Disks[0] = diskdrive.Attach(conf.Globals.DebugMode)
			if conf.Disks.Disk1 != "empty" {
				C.Disks[0].LoadDiskImage(conf.Disks.Disk1)
			}
			C.connectedDrive++
		} else {
			C.Disks[0] = nil
		}
		if conf.Disks.Disk2 != "" {
			C.Disks[1] = diskdrive.Attach(conf.Globals.DebugMode)
			if conf.Disks.Disk2 != "empty" {
				C.Disks[1].LoadDiskImage(conf.Disks.Disk2)
			}
			C.connectedDrive++
		} else {
			C.Disks[1] = nil
		}
	}
}

func (C *DiskInterface) diskMotorsON() byte {
	// log.Printf("Start motor drive %d", SelectedDrive)
	C.MotorIsOn = true
	C.Disks[C.SelectedDrive].StartMotor()
	// C.drivesStatus()
	// clog.FileRaw("\n%s : Start Motor: %04X", time.Now().Format("15:04:05"), cpu.InstStart)
	return 0
}

func (C *DiskInterface) diskMotorsOFF() byte {
	if !C.MotorIsOn {
		return 0
	}
	if C.Disks[0] != nil {
		C.Disks[0].StopMotor()
	}
	if C.Disks[1] != nil {
		C.Disks[1].StopMotor()
	}
	C.MotorIsOn = false
	// clog.FileRaw("\n%s : Stop Motor: %04X", time.Now().Format("15:04:05"), cpu.InstStart)
	return 0
}

// Select drive 0 or 1
func (C *DiskInterface) driveSelect(driveNum int) byte {
	var retVal byte = 0

	if C.connectedDrive == 0 {
		retVal = 0x80
	}
	// On selection le drive déjà selectionné
	if driveNum == C.SelectedDrive {
		return 0x80
	}
	// On switch de drive
	C.Disks[C.SelectedDrive].IsSpinning = false
	C.SelectedDrive = driveNum
	if C.Disks[C.SelectedDrive] != nil && C.MotorIsOn {
		C.Disks[C.SelectedDrive].IsSpinning = true
	}

	return retVal
}

func (C *DiskInterface) SetPhase(phase int, state bool) {
	C.Disks[C.SelectedDrive].SetPhase(phase, state)
}

func (C *DiskInterface) SetSequencerMode(mode bool) byte {
	var retVal byte = 0
	C.SequencerMode = mode
	if C.SequencerMode { // SEQ_READ_MODE
		retVal = C.ProtectCheck
		C.ProtectCheck = 0
		// log.Printf("Sequencer read = %02X (%04X)", retVal, cpu.InstStart)
	}
	return retVal
}

func (C *DiskInterface) ShiftOrRead() byte {
	if C.SequencerMode { // SEQ_READ_MODE
		if C.Disks[C.SelectedDrive] != nil && C.Disks[C.SelectedDrive].IsSpinning {
			tmp := C.Disks[C.SelectedDrive].GetNextByte()
			// clog.Debug("IO", "disk", "Read : %02X\n", tmp)
			// fmt.Printf("%02X\n", tmp)
			// clog.FileRaw("\n%s : => READ DATA => %02X [%04X]", time.Now().Format("15:04:05"), tmp, cpu.InstStart)
			// clog.FileRaw("\n%s", cpu.FullDebug)
			return tmp
		}
		return 0x00
	}
	// else {
	// 	Shift sequencer
	// }
	return 0
}

func (C *DiskInterface) LoadOrCheck() byte {
	// if SequencerMode == SEQ_READ_MODE {
	if C.SelectedDrive != -1 {
		if C.Disks[C.SelectedDrive].IsWriteProtected {
			// log.Printf("Disk is Write Protected: %04X", cpu.InstStart)
			C.ProtectCheck = 0x80
			// log.Printf("Protection Check = %02X (%04X)", C.ProtectCheck, cpu.InstStart)
		} else {
			log.Printf("Disk is Writable")
			C.ProtectCheck = 0
		}
	}
	// Load sequencer

	return 0
}

func (C *DiskInterface) GetStats() []string {
	sel := ""
	mtr := ""
	sts := ""
	stat := make([]string, 2)
	tmp := make([]string, 3)
	tmp[0] = "Motr Drv Stat Selec     " + "    Phy.Trk Dat.Trk Rev  Pos"
	if C.SelectedDrive == 0 {
		sel = "*"
	}
	if C.MotorIsOn {
		mtr = "ON"
	}
	if C.Disks[0] != nil && C.Disks[0].IsSpinning {
		sts = "Spin"
	} else {
		sts = "Stop"
	}
	stat[0] = fmt.Sprintf(" %2s   1  %s   %s       ", mtr, sts, sel)
	sel = ""
	mtr = ""
	sts = ""
	if !C.MotorIsOn {
		mtr = "OF"
	}
	if C.SelectedDrive == 1 {
		sel = "*"
	}
	if C.Disks[1] != nil && C.Disks[1].IsSpinning {
		sts = "Spin"
	} else {
		sts = "Stop"
	}
	stat[1] = fmt.Sprintf(" %2s   2  %s   %s       ", mtr, sts, sel)

	if C.Disks[0] != nil {
		tmp[1] = stat[0] + C.Disks[0].GetStatus()
	} else {
		tmp[1] = stat[0]
	}
	if C.Disks[1] != nil {
		tmp[2] = stat[1] + C.Disks[1].GetStatus()
	} else {
		tmp[2] = stat[1]
	}
	return tmp
}

// func (C *DiskInterface) drivesStatus() {
// 	var D1, D2 string
// 	fmt.Printf("==============================================\n")
// 	fmt.Printf("                 =Drive 1=           =Drive 2=\n")
// 	if C.SelectedDrive == 0 {
// 		fmt.Printf("Selected:           X\n")
// 	} else {
// 		fmt.Printf("Selected:                                X\n")
// 	}

// 	if C.Disks[0].IsRunning {
// 		D1 = "ON"
// 	} else {
// 		D1 = "OFF"
// 	}
// 	if C.Disks[1].IsRunning {
// 		D2 = "ON"
// 	} else {
// 		D2 = "OFF"
// 	}
// 	fmt.Printf("Motors:             %3s                 %3s\n", D1, D2)

// 	// if C.Disks[0].ReadMode {
// 	// 	D1 = "RD"
// 	// } else {
// 	// 	D1 = "WR"
// 	// }
// 	// if C.Disks[1].ReadMode {
// 	// 	D2 = "RD"
// 	// } else {
// 	// 	D2 = "WR"
// 	// }
// 	// fmt.Printf("Mode:               %3s                 %3s\n", D1, D2)

// 	if C.Disks[0].IsWriteProtected {
// 		D1 = "ON"
// 	} else {
// 		D1 = "OFF"
// 	}
// 	if C.Disks[1].IsWriteProtected {
// 		D2 = "ON"
// 	} else {
// 		D2 = "OFF"
// 	}
// 	fmt.Printf("WriteProtect:       %3s                 %3s\n", D1, D2)
// 	fmt.Printf("==============================================\n")
// }
