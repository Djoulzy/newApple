package main

import (
	"github.com/Djoulzy/emutools/mem"
)

func memLayouts(model int) {

	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0x0000, RAM, mem.READWRITE, false)
	if model == 2 {
		MEM.Layouts[0].Attach("ROM_C", 0xC000, ROM_C, mem.READONLY, false)
	}
	MEM.Layouts[0].Attach("ROM_D", 0xD000, ROM_D, mem.READONLY, false)
	MEM.Layouts[0].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY, false)

	MEM.Layouts[0].Attach("IO", 0xC000, IO, mem.READWRITE, false)
	MEM.Layouts[0].Attach("SLOT1", 0xC100, SLOT1, mem.READONLY, false)
	MEM.Layouts[0].Attach("SLOT2", 0xC200, SLOT2, mem.READONLY, false)
	MEM.Layouts[0].Attach("SLOT3", 0xC300, SLOT3, mem.READONLY, true)
	MEM.Layouts[0].Attach("SLOT4", 0xC400, SLOT4, mem.READONLY, false)
	MEM.Layouts[0].Attach("SLOT5", 0xC500, SLOT5, mem.READONLY, false)
	MEM.Layouts[0].Attach("SLOT6", 0xC600, SLOT6, mem.READONLY, false)
	MEM.Layouts[0].Attach("SLOT7", 0xC700, SLOT7, mem.READONLY, false)
	MEM.Layouts[0].Accessor("IO", IOAccess)
	// MEM.Layouts[0].Show()
}
