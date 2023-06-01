package io

import (
	"log"
	"newApple/config"
	"newApple/disk"
)

const (
	SEQ_READ_MODE  = true
	SEQ_WRITE_MODE = false
)

type DiskInterface struct {
	Disks          []*disk.DRIVE
	connectedDrive int

	SelectedDrive int
	SequencerMode bool
	ProtectCheck  byte
}

func InitDiskInterface(conf *config.ConfigData) *DiskInterface {
	tmp := DiskInterface{
		Disks:          make([]*disk.DRIVE, 2),
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
		if conf.Disks.Disk1 != "" {
			C.Disks[0] = disk.Attach(conf.Globals.DebugMode)
			C.Disks[0].LoadDiskImage(conf.Disks.Disk1)
			C.connectedDrive++
		}
		if conf.Disks.Disk2 != "" {
			C.Disks[1] = disk.Attach(conf.Globals.DebugMode)
			C.Disks[1].LoadDiskImage(conf.Disks.Disk2)
			C.connectedDrive++
		}
		if C.connectedDrive == 0 {
			conf.Slots.Slot6 = ""
		}
	}
}

func (C *DiskInterface) diskMotorsON() byte {
	disk.MotorIsOn = true
	// log.Printf("Start motor drive %d", SelectedDrive)
	C.Disks[C.SelectedDrive].StartMotor()
	// C.drivesStatus()
	// clog.FileRaw("\n%s : Start Motor: %04X", time.Now().Format("15:04:05"), cpu.InstStart)
	return 0
}

func (C *DiskInterface) diskMotorsOFF() byte {
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
func (C *DiskInterface) driveSelect(driveNum int) byte {
	var retVal byte = 0

	if C.connectedDrive == 0 {
		retVal = 0x80
	}
	if driveNum == C.SelectedDrive {
		return retVal
	}
	if disk.MotorIsOn {
		C.Disks[C.SelectedDrive].StopMotor()
		if driveNum+1 <= C.connectedDrive {
			C.Disks[driveNum].StartMotor()
			C.SelectedDrive = driveNum
		} else {
			C.SelectedDrive = 0
		}
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
		if C.SelectedDrive != -1 {
			if C.Disks[C.SelectedDrive].IsRunning {
				tmp := C.Disks[C.SelectedDrive].GetNextByte()
				// clog.Debug("IO", "disk", "Read : %02X\n", tmp)
				// fmt.Printf("%02X\n", tmp)
				// clog.FileRaw("\n%s : => READ DATA => %02X [%04X]", time.Now().Format("15:04:05"), tmp, cpu.InstStart)
				// clog.FileRaw("\n%s", cpu.FullDebug)
				return tmp
			}
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
