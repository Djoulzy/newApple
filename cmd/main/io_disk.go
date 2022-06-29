package main

import (
	"fmt"
	"log"
	"newApple/disk"

	"github.com/Djoulzy/Tools/clog"
)

const (
	SEQ_READ_MODE  = true
	SEQ_WRITE_MODE = false
)

var (
	SelectedDrive int  = 0
	SequencerMode bool = false
	ProtectCheck  byte = 0
)

func (C *io_access) diskMotorsON() byte {
	disk.MotorIsOn = true
	C.Disks[SelectedDrive].StartMotor()
	// C.drivesStatus()
	// clog.FileRaw("\n%s : Start Motor: %04X", time.Now().Format("15:04:05"), cpu.InstStart)
	return 0
}

func (C *io_access) diskMotorsOFF() byte {
	if C.connectedDrive == 0 || !disk.MotorIsOn {
		return 0
	}

	C.Disks[0].StopMotor()
	if C.connectedDrive > 1 {
		C.Disks[1].StopMotor()
	}

	disk.MotorIsOn = false
	// clog.FileRaw("\n%s : Stop Motor: %04X", time.Now().Format("15:04:05"), cpu.InstStart)
	return 171
}

// Select drive 0 or 1
func (C *io_access) driveSelect(driveNum int) byte {
	var retVal byte = 0

	if C.connectedDrive == 0 {
		retVal = 0x80
	}
	if driveNum == SelectedDrive {
		return retVal
	}
	if disk.MotorIsOn {
		if SelectedDrive != driveNum {
			C.Disks[SelectedDrive].StopMotor()
		}
		if driveNum+1 <= C.connectedDrive {
			C.Disks[driveNum].StartMotor()
			SelectedDrive = driveNum
		} else {
			SelectedDrive = -1
		}
	}

	return retVal
}

func (C *io_access) SetSequencerMode(mode bool) byte {
	var retVal byte = 0
	SequencerMode = mode
	if SequencerMode == SEQ_READ_MODE {
		retVal = ProtectCheck
		ProtectCheck = 0
		// log.Printf("Sequencer read = %02X (%04X)", retVal, cpu.InstStart)
	}
	return retVal
}

func (C *io_access) ShiftOrRead() byte {
	if SequencerMode == SEQ_READ_MODE {
		if SelectedDrive != -1 {
			if C.Disks[SelectedDrive].IsRunning {
				tmp := C.Disks[SelectedDrive].GetNextByte()
				// clog.Debug("IO", "disk", "Read : %02X\n", tmp)
				// fmt.Printf("%02X\n", tmp)
				// clog.FileRaw("\n%s : => READ DATA => %02X [%04X]", time.Now().Format("15:04:05"), tmp, cpu.InstStart)
				clog.FileRaw("\n%s", cpu.FullDebug)
				return tmp
			}
		}
		return 0x00
	} else {
		// Shift sequencer
	}
	return 0
}

func (C *io_access) LoadOrCheck() byte {
	// if SequencerMode == SEQ_READ_MODE {
	if SelectedDrive != -1 {
		if C.Disks[SelectedDrive].IsWriteProtected {
			// log.Printf("Disk is Write Protected: %04X", cpu.InstStart)
			ProtectCheck = 0x80
			log.Printf("Protection Check = %02X (%04X)", ProtectCheck, cpu.InstStart)
		} else {
			log.Printf("Disk is Writable")
			ProtectCheck = 0
		}
	}
	// Load sequencer

	return 0
}

func (C *io_access) drivesStatus() {
	var D1, D2 string
	fmt.Printf("==============================================\n")
	fmt.Printf("                 =Drive 1=           =Drive 2=\n")
	if SelectedDrive == 0 {
		fmt.Printf("Selected:           X\n")
	} else {
		fmt.Printf("Selected:                                X\n")
	}

	if C.Disks[0].IsRunning {
		D1 = "ON"
	} else {
		D1 = "OFF"
	}
	if C.Disks[1].IsRunning {
		D2 = "ON"
	} else {
		D2 = "OFF"
	}
	fmt.Printf("Motors:             %3s                 %3s\n", D1, D2)

	// if C.Disks[0].ReadMode {
	// 	D1 = "RD"
	// } else {
	// 	D1 = "WR"
	// }
	// if C.Disks[1].ReadMode {
	// 	D2 = "RD"
	// } else {
	// 	D2 = "WR"
	// }
	// fmt.Printf("Mode:               %3s                 %3s\n", D1, D2)

	if C.Disks[0].IsWriteProtected {
		D1 = "ON"
	} else {
		D1 = "OFF"
	}
	if C.Disks[1].IsWriteProtected {
		D2 = "ON"
	} else {
		D2 = "OFF"
	}
	fmt.Printf("WriteProtect:       %3s                 %3s\n", D1, D2)
	fmt.Printf("==============================================\n")
}
