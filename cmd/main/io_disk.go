package main

import (
	"fmt"
)

var (
	MotorOn       bool = false
	SelectedDrive int  = 0
)

func (C *io_access) diskMotorsON() byte {
	MotorOn = true
	C.Disks[SelectedDrive].StartMotor()
	// C.drivesStatus()
	return 0
}

func (C *io_access) diskMotorsOFF() byte {
	MotorOn = false
	C.Disks[0].StopMotor()
	C.Disks[1].StopMotor()
	// C.drivesStatus()
	return 0
}

func (C *io_access) driveSelect(driveNum int) byte {
	if MotorOn {
		C.Disks[SelectedDrive].StopMotor()
		C.Disks[driveNum].StartMotor()
	}
	SelectedDrive = driveNum
	// C.drivesStatus()
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

	if C.Disks[0].ReadMode {
		D1 = "RD"
	} else {
		D1 = "WR"
	}
	if C.Disks[1].ReadMode {
		D2 = "RD"
	} else {
		D2 = "WR"
	}
	fmt.Printf("Mode:               %3s                 %3s\n", D1, D2)

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
