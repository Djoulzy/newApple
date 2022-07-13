package main

import (
	"github.com/Djoulzy/emutools/mem"
)

const (
	DISABLED = true
	ENABLED  = false
)

func memLayouts(model int) {

	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0x0000, RAM, mem.READWRITE, ENABLED)
	MEM.Layouts[0].Attach("BANK2", 0xD000, BANK2, mem.READWRITE, ENABLED)
	MEM.Layouts[0].Attach("BANK1", 0xD000, BANK1, mem.READWRITE, ENABLED)

	MEM.Layouts[0].Attach("ZP", 0x0000, ZP, mem.READWRITE, ENABLED)

	if conf.Model == "2e" {
		MEM.Layouts[0].Attach("AUX", 0x0000, AUX, mem.READWRITE, DISABLED)
		MEM.Layouts[0].Attach("AUX_BANK2", 0xD000, AUX_BANK2, mem.READWRITE, DISABLED)
		MEM.Layouts[0].Attach("AUX_BANK1", 0xD000, AUX_BANK1, mem.READWRITE, DISABLED)
		MEM.Layouts[0].Attach("ALT_ZP", 0x0000, ALT_ZP, mem.READWRITE, DISABLED)

		MEM.Layouts[0].Attach("ROM_C", 0xC000, ROM_C, mem.READONLY, ENABLED)
	}
	MEM.Layouts[0].Attach("ROM_D", 0xD000, ROM_D, mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY, ENABLED)

	MEM.Layouts[0].Attach("IO", 0xC000, IO, mem.READWRITE, ENABLED)
	MEM.Layouts[0].Attach("SLOT1", 0xC100, SLOTS[1], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT2", 0xC200, SLOTS[2], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT3", 0xC300, SLOTS[3], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT4", 0xC400, SLOTS[4], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT5", 0xC500, SLOTS[5], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT6", 0xC600, SLOTS[6], mem.READONLY, ENABLED)
	MEM.Layouts[0].Attach("SLOT7", 0xC700, SLOTS[7], mem.READONLY, ENABLED)
	MEM.Layouts[0].Accessor("IO", IOAccess)
	// MEM.Layouts[0].Show()
}
