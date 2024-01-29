package main

import (
	"newApple/io"
	"strconv"

	"github.com/Djoulzy/mmu"
)

const (
	romSize      = 4096
	softSwitches = 256
	slot_roms    = 256
	chargenSize  = 2048
)

var (
	MN_ZPS = mmu.NewRAM("MN_ZPS", 0x0200)
	MN___1 = mmu.NewRAM("MN___1", 0x0200)
	MN_TXT = mmu.NewRAM("MN_TXT", 0x0400)
	MN___2 = mmu.NewRAM("MN___2", 0x1800)
	MN_HGR = mmu.NewRAM("MN_HGR", 0x2000)
	MN___3 = mmu.NewRAM("MN___3", 0x9000)
	MN_SLT = mmu.NewRAM("MN_SLT", 0x0800)

	MN_BK1 = mmu.NewRAM("MN_BK1", 0x1000)
	MN_BK2 = mmu.NewRAM("MN_BK2", 0x1000)
	MN___4 = mmu.NewRAM("MN___4", 0x2000)

	// AUX_ZP = mmu.NewRAM("AX_ZP", zpStack)
	// AUX_LO = mmu.NewRAM("AX_LO", lowRamSize)
	// AUX_B1 = mmu.NewRAM("AX_B1", bankSize)
	// AUX_B2 = mmu.NewRAM("AX_B2", bankSize)
	// AUX_HI = mmu.NewRAM("AX_HI", hiRamSize)

	ROM_C  *mmu.ROM
	ROM_D  *mmu.ROM
	ROM_EF *mmu.ROM

	IO      *io.SoftSwitch
	Disks   *io.DiskInterface
	SLOTS   [8]*mmu.ROM
	CHARGEN *mmu.ROM
)

func loadSlots() {
	conf.Slots.Catalog[1] = conf.Slots.Slot1
	conf.Slots.Catalog[2] = conf.Slots.Slot2
	conf.Slots.Catalog[3] = conf.Slots.Slot3
	conf.Slots.Catalog[4] = conf.Slots.Slot4
	conf.Slots.Catalog[5] = conf.Slots.Slot5
	conf.Slots.Catalog[6] = conf.Slots.Slot6
	conf.Slots.Catalog[7] = conf.Slots.Slot7

	for i := 1; i < 8; i++ {
		// if conf.Slots.Catalog[i] != "" {
		SLOTS[i] = mmu.NewROM("SLOT_"+strconv.Itoa(i), slot_roms, conf.Slots.Catalog[i])
		MEM.Attach(SLOTS[i], 0xC0+uint(i))
		MEM.Mount("SLOT_"+strconv.Itoa(i), "")
		// }
	}
}

func apple2_Roms() {
	ROM_D = mmu.NewROM("ROM_D", romSize, "assets/roms/II/D.bin")
	MEM.Attach(ROM_D, 0xD0)
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "assets/roms/II/EF.bin")
	MEM.Attach(ROM_EF, 0xE0)

	MEM.Mount("ROM_D", "MN_BK1")
	MEM.Mount("ROM_EF", "MN___4")

	CHARGEN = mmu.NewROM("CHARGEN", chargenSize, "assets/roms/II/3410036.bin")
	// MEM.Attach(ROM_D, 0xD0, 8)
}

func apple2e_Roms() {
	ROM_C = mmu.NewROM("ROM_C", romSize, "assets/roms/IIe/C.bin")
	MEM.Attach(ROM_C, 0xC0)
	ROM_D = mmu.NewROM("ROM_D", romSize, "assets/roms/IIe/D.bin")
	MEM.Attach(ROM_D, 0xD0)
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "assets/roms/IIe/EF.bin")
	MEM.Attach(ROM_EF, 0xE0)

	MEM.Mount("ROM_C", "MN_SLT")
	MEM.Mount("ROM_D", "MN_BK1")
	MEM.Mount("ROM_EF", "MN___4")

	CHARGEN = mmu.NewROM("CHARGEN", chargenSize*2, "assets/roms/IIe/Video_US.bin")
}

// func apple2_Roms() {
// 	ROM_D = mmu.NewROM("ROM_D", romSize, conf.Rom_Apple2.Rom_D)
// 	MEM.Attach(ROM_D, 0xD0, mmu.READONLY)
// 	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, conf.Rom_Apple2.Rom_EF)
// 	MEM.Attach(ROM_EF, 0xE0, mmu.READONLY)

// 	MEM.Mount("ROM_D", mmu.READONLY)
// 	MEM.Mount("ROM_EF", mmu.READONLY)

// 	CHARGEN = mmu.NewROM("CHARGEN", chargenSize, conf.Rom_Apple2.Chargen)
// 	// MEM.Attach(ROM_D, 0xD0, 8)
// }

// func apple2e_Roms() {
// 	ROM_C = mmu.NewROM("ROM_C", romSize, conf.Rom_Apple2e.Rom_C)
// 	MEM.Attach(ROM_C, 0xC0, mmu.READONLY)
// 	ROM_D = mmu.NewROM("ROM_D", romSize, conf.Rom_Apple2e.Rom_D)
// 	MEM.Attach(ROM_D, 0xD0, mmu.READONLY)
// 	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, conf.Rom_Apple2e.Rom_EF)
// 	MEM.Attach(ROM_EF, 0xE0, mmu.READONLY)

// 	MEM.Mount("ROM_C", mmu.READONLY)
// 	MEM.Mount("ROM_D", mmu.READONLY)
// 	MEM.Mount("ROM_EF", mmu.READONLY)

// 	CHARGEN = mmu.NewROM("CHARGEN", chargenSize*2, conf.Rom_Apple2e.Chargen)
// }

func initRam() {
	MN_ZPS.Clear(0x1000, 0xFF)
	MN___1.Clear(0x1000, 0xFF)
	MN_TXT.Clear(0x1000, 0xFF)
	MN___2.Clear(0x1000, 0xFF)
	MN_HGR.Clear(0x1000, 0xFF)
	MN___3.Clear(0x1000, 0xFF)
	MN_SLT.Clear(0x1000, 0xFF)
	MN_BK1.Clear(0x1000, 0xFF)
	MN_BK2.Clear(0x1000, 0xFF)
	MN___4.Clear(0x1000, 0xFF)

	// AUX_LO.Clear(0x1000, 0xFF)
	// AUX_B1.Clear(0x1000, 0xFF)
	// AUX_B2.Clear(0x1000, 0xFF)

	MEM.Attach(MN_ZPS, 0x00)
	MEM.Attach(MN___1, 0x02)
	MEM.Attach(MN_TXT, 0x04)
	MEM.Attach(MN___2, 0x08)
	MEM.Attach(MN_HGR, 0x20)
	MEM.Attach(MN___3, 0x40)
	MEM.Attach(MN_SLT, 0xC8)
	MEM.Attach(MN_BK1, 0xD0)
	MEM.Attach(MN_BK2, 0xD0)
	MEM.Attach(MN___4, 0xE0)

	// MEM.Attach(AUX_LO, 0x00)
	// MEM.Attach(AUX_B1, 0xD0)
	// MEM.Attach(AUX_B2, 0xD0)
	// MEM.Attach(AUX_HI, 0xE0)

	MEM.Mount("MN_ZPS", "MN_ZPS")
	MEM.Mount("MN___1", "MN___1")
	MEM.Mount("MN_TXT", "MN_TXT")
	MEM.Mount("MN___2", "MN___2")
	MEM.Mount("MN_HGR", "MN_HGR")
	MEM.Mount("MN___3", "MN___3")
	MEM.Mount("MN_SLT", "MN_SLT")
}

func setupMemoryLayout() {
	initRam()

	if MODEL == 1 {
		apple2_Roms()
	} else {
		apple2e_Roms()
	}

	Disks = io.InitDiskInterface(conf)
	IO = io.InitSoftSwitch("IO", softSwitches, Disks, &CRTC)
	IO.SetMMU(MEM)
	MEM.Attach(IO, 0xC0)
	MEM.Mount("IO", "IO")

	loadSlots()

	MEM.CheckMapIntegrity()
	MEM.DumpMap()
	// os.Exit(0)
}
