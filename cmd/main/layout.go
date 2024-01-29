package main

import (
	"newApple/io"
	"strconv"

	"github.com/Djoulzy/mmu"
)

const (
	lowRamSize   = 53248
	hiRamSize    = 8192
	bankSize     = 4096
	romSize      = 4096
	softSwitches = 256
	chargenSize  = 2048
	keyboardSize = 2048
	slot_roms    = 256
)

var (
	MAIN_LOW *mmu.RAM
	MAIN_B1  *mmu.RAM
	MAIN_B2  *mmu.RAM
	MAIN_HI  *mmu.RAM

	AUX_LOW *mmu.RAM
	AUX_B1  *mmu.RAM
	AUX_B2  *mmu.RAM
	AUX_HI  *mmu.RAM

	ROM_C  *mmu.ROM
	ROM_D  *mmu.ROM
	ROM_EF *mmu.ROM

	IO      *io.SoftSwitch
	Disks   *io.DiskInterface
	SLOTS   [8]*mmu.ROM
	CHARGEN *mmu.ROM
)

// func memLayouts(model int) {

// 	// Apple 2
// 	if model == 1 {
// 		MEM.Attach(0, "RAM", 0x0000, RAM, READWRITE, ENABLED, nil)
// 		MEM.Attach(0, "BANK2", 0xD000, BANK2, READWRITE, ENABLED, nil)
// 		MEM.Attach(0, "BANK1", 0xD000, BANK1, READWRITE, ENABLED, nil)

// 		MEM.Attach(0, "IO", 0xC000, IO, READWRITE, ENABLED, IOAccess)
// 		MEM.Attach(0, "SLOT1", 0xC100, SLOTS[1], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT2", 0xC200, SLOTS[2], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT3", 0xC300, SLOTS[3], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT4", 0xC400, SLOTS[4], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT5", 0xC500, SLOTS[5], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT6", 0xC600, SLOTS[6], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT7", 0xC700, SLOTS[7], READONLY, ENABLED, nil)

// 		MEM.Attach(0, "ROM_D", 0xD000, ROM_D, READONLY, ENABLED, nil)
// 		MEM.Attach(0, "ROM_EF", 0xE000, ROM_EF, READONLY, ENABLED, nil)
// 	} else {
// 		// Apple 2e et 2+
// 		MEM.Attach(0, "RAM", 0x0000, RAM, READWRITE, ENABLED, nil)
// 		MEM.Attach(0, "AUX", 0x0000, AUX, READWRITE, DISABLED, nil)

// 		MEM.Attach(0, "BANK2", 0xD000, BANK2, READWRITE, ENABLED, nil)
// 		MEM.Attach(0, "BANK1", 0xD000, BANK1, READWRITE, ENABLED, nil)

// 		MEM.Attach(0, "ZP", 0x0000, ZP, READWRITE, ENABLED, nil)

// 		MEM.Attach(0, "AUX_BANK2", 0xD000, AUX_BANK2, READWRITE, DISABLED, nil)
// 		MEM.Attach(0, "AUX_BANK1", 0xD000, AUX_BANK1, READWRITE, DISABLED, nil)
// 		MEM.Attach(0, "ALT_ZP", 0x0000, ALT_ZP, READWRITE, DISABLED, nil)
// 		MEM.Attach(0, "ROM_C", 0xC000, ROM_C, READONLY, ENABLED, nil)

// 		MEM.Attach(0, "ROM_D", 0xD000, ROM_D, READONLY, ENABLED, nil)
// 		MEM.Attach(0, "ROM_EF", 0xE000, ROM_EF, READONLY, ENABLED, nil)

// 		MEM.Attach(0, "IO", 0xC000, IO, READWRITE, ENABLED, IOAccess)
// 		MEM.Attach(0, "SLOT1", 0xC100, SLOTS[1], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT2", 0xC200, SLOTS[2], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT3", 0xC300, SLOTS[3], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT4", 0xC400, SLOTS[4], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT5", 0xC500, SLOTS[5], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT6", 0xC600, SLOTS[6], READONLY, ENABLED, nil)
// 		MEM.Attach(0, "SLOT7", 0xC700, SLOTS[7], READONLY, ENABLED, nil)

// 		// MEM.Layouts[1].Show()
// 	}
// }

func loadSlots() {
	conf.Slots.Catalog[1] = conf.Slots.Slot1
	conf.Slots.Catalog[2] = conf.Slots.Slot2
	conf.Slots.Catalog[3] = conf.Slots.Slot3
	conf.Slots.Catalog[4] = conf.Slots.Slot4
	conf.Slots.Catalog[5] = conf.Slots.Slot5
	conf.Slots.Catalog[6] = conf.Slots.Slot6
	conf.Slots.Catalog[7] = conf.Slots.Slot7

	for i := 1; i < 8; i++ {
		if conf.Slots.Catalog[i] != "" {
			SLOTS[i] = mmu.NewROM("SLOT_"+strconv.Itoa(i), slot_roms, conf.Slots.Catalog[i])
			MEM.Attach(SLOTS[i], 0xC0+uint(i))
			MEM.Mount("SLOT_"+strconv.Itoa(i), "SLOT_"+strconv.Itoa(i))
		}
	}
}

func apple2_Roms() {
	ROM_D = mmu.NewROM("ROM_D", romSize, "assets/roms/II/D.bin")
	MEM.Attach(ROM_D, 0xD0)
	ROM_EF = mmu.NewROM("ROM_EF", romSize*2, "assets/roms/II/EF.bin")
	MEM.Attach(ROM_EF, 0xE0)

	MEM.Mount("ROM_D", "MAIN_B1")
	MEM.Mount("ROM_EF", "MAIN_HI")

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

	MEM.Mount("ROM_C", "MAIN_LOW")
	MEM.Mount("ROM_D", "MAIN_B1")
	MEM.Mount("ROM_EF", "MAIN_HI")

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
	MAIN_LOW = mmu.NewRAM("MAIN_LOW", lowRamSize)
	MAIN_B1 = mmu.NewRAM("MAIN_B1", bankSize)
	MAIN_B2 = mmu.NewRAM("MAIN_B2", bankSize)
	MAIN_HI = mmu.NewRAM("MAIN_HI", hiRamSize)

	AUX_LOW = mmu.NewRAM("AUX_LOW", lowRamSize)
	AUX_B1 = mmu.NewRAM("AUX_B1", bankSize)
	AUX_B2 = mmu.NewRAM("AUX_B2", bankSize)
	AUX_HI = mmu.NewRAM("AUX_HI", hiRamSize)

	MAIN_LOW.Clear(0x1000, 0xFF)
	MAIN_B1.Clear(0x1000, 0xFF)
	MAIN_B2.Clear(0x1000, 0xFF)
	AUX_LOW.Clear(0x1000, 0xFF)
	AUX_B1.Clear(0x1000, 0xFF)
	AUX_B2.Clear(0x1000, 0xFF)

	MEM.Attach(MAIN_LOW, 0x00)
	MEM.Attach(MAIN_B1, 0xD0)
	MEM.Attach(MAIN_B2, 0xD0)
	MEM.Attach(MAIN_HI, 0xE0)

	MEM.Attach(AUX_LOW, 0x00)
	MEM.Attach(AUX_B1, 0xD0)
	MEM.Attach(AUX_B2, 0xD0)
	MEM.Attach(AUX_HI, 0xE0)

	MEM.Mount("MAIN_LOW", "MAIN_LOW")
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
