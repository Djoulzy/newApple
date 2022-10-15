package main

import (
	mem "github.com/Djoulzy/emutools/mem2"
)

const (
	DISABLED = true
	ENABLED  = false
)

func memLayouts(model int) {

	// Apple 2
	if model == 1 {
		MEM.Attach(0, "RAM", 0x0000, RAM, mem.READWRITE, ENABLED, nil)
		MEM.Attach(0, "BANK2", 0xD000, BANK2, mem.READWRITE, ENABLED, nil)
		MEM.Attach(0, "BANK1", 0xD000, BANK1, mem.READWRITE, ENABLED, nil)

		MEM.Attach(0, "IO", 0xC000, IO, mem.READWRITE, ENABLED, IOAccess)
		MEM.Attach(0, "SLOT1", 0xC100, SLOTS[1], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT2", 0xC200, SLOTS[2], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT3", 0xC300, SLOTS[3], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT4", 0xC400, SLOTS[4], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT5", 0xC500, SLOTS[5], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT6", 0xC600, SLOTS[6], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT7", 0xC700, SLOTS[7], mem.READONLY, ENABLED, nil)

		MEM.Attach(0, "ROM_D", 0xD000, ROM_D, mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "ROM_EF", 0xE000, ROM_EF, mem.READONLY, ENABLED, nil)
	} else {
		// Apple 2e et 2+
		MEM.Attach(0, "RAM", 0x0000, RAM, mem.READWRITE, ENABLED, nil)
		MEM.Attach(0, "AUX", 0x0000, AUX, mem.READWRITE, DISABLED, nil)

		MEM.Attach(0, "BANK2", 0xD000, BANK2, mem.READWRITE, ENABLED, nil)
		MEM.Attach(0, "BANK1", 0xD000, BANK1, mem.READWRITE, ENABLED, nil)

		MEM.Attach(0, "ZP", 0x0000, ZP, mem.READWRITE, ENABLED, nil)

		MEM.Attach(0, "AUX_BANK2", 0xD000, AUX_BANK2, mem.READWRITE, DISABLED, nil)
		MEM.Attach(0, "AUX_BANK1", 0xD000, AUX_BANK1, mem.READWRITE, DISABLED, nil)
		MEM.Attach(0, "ALT_ZP", 0x0000, ALT_ZP, mem.READWRITE, DISABLED, nil)
		MEM.Attach(0, "ROM_C", 0xC000, ROM_C, mem.READONLY, ENABLED, nil)

		MEM.Attach(0, "ROM_D", 0xD000, ROM_D, mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "ROM_EF", 0xE000, ROM_EF, mem.READONLY, ENABLED, nil)

		MEM.Attach(0, "IO", 0xC000, IO, mem.READWRITE, ENABLED, IOAccess)
		MEM.Attach(0, "SLOT1", 0xC100, SLOTS[1], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT2", 0xC200, SLOTS[2], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT3", 0xC300, SLOTS[3], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT4", 0xC400, SLOTS[4], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT5", 0xC500, SLOTS[5], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT6", 0xC600, SLOTS[6], mem.READONLY, ENABLED, nil)
		MEM.Attach(0, "SLOT7", 0xC700, SLOTS[7], mem.READONLY, ENABLED, nil)

		// MEM.Layouts[1].Show()
	}
}
