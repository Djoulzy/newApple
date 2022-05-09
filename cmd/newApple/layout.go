package main

import (
	"github.com/Djoulzy/emutools/mem"
)

func memLayouts() {
	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	MEM.Layouts[0].Attach("IO", 0xC000, IO, mem.READWRITE)
	// MEM.Layouts[0].Attach("SLOT1", 0xC100, SLOT6, mem.READONLY)
	// MEM.Layouts[0].Attach("SLOT2", 0xC200, SLOT6, mem.READONLY)
	// MEM.Layouts[0].Attach("SLOT3", 0xC300, SLOT6, mem.READONLY)
	// MEM.Layouts[0].Attach("SLOT4", 0xC400, SLOT6, mem.READONLY)
	// MEM.Layouts[0].Attach("SLOT5", 0xC500, SLOT6, mem.READONLY)
	MEM.Layouts[0].Attach("SLOT6", 0xC600, SLOT6, mem.READONLY)
	// MEM.Layouts[0].Attach("SLOT7", 0xC700, SLOT6, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_D0", 0xD000, ROM_D0, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_D8", 0xD800, ROM_D8, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_E0", 0xE000, ROM_E0, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_E8", 0xE800, ROM_E8, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_F0", 0xF000, ROM_F0, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_F8", 0xF800, ROM_F8, mem.READONLY)
	// MEM.Layouts[0].Attach("KERNAL", 0xF000, KERNAL, mem.READONLY)
	MEM.Layouts[0].Accessor("IO", IOAccess)
	MEM.Layouts[0].Show()
}
