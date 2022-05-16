package main

import (
	"github.com/Djoulzy/emutools/mem"
)

func memLayouts(model int) {

	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	if model == 1 {
		MEM.Layouts[0].Attach("ROM_D0", 0xD000, ROM_D0, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_D8", 0xD800, ROM_D8, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_E0", 0xE000, ROM_E0, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_E8", 0xE800, ROM_E8, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_F0", 0xF000, ROM_F0, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_F8", 0xF800, ROM_F8, mem.READONLY)
	} else {
		MEM.Layouts[0].Attach("ROM_CD", 0xC000, ROM_CD, mem.READONLY)
		MEM.Layouts[0].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY)
	}
	MEM.Layouts[0].Attach("IO", 0xC000, IO, mem.READWRITE)
	MEM.Layouts[0].Accessor("IO", IOAccess)
	// MEM.Layouts[0].Show()

	///////////////////////////////////////////////////////////////////////

	MEM.Layouts[1] = mem.InitConfig(ramSize)
	MEM.Layouts[1].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	if model == 1 {
		MEM.Layouts[1].Attach("ROM_D0", 0xD000, ROM_D0, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_D8", 0xD800, ROM_D8, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_E0", 0xE000, ROM_E0, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_E8", 0xE800, ROM_E8, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_F0", 0xF000, ROM_F0, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_F8", 0xF800, ROM_F8, mem.READONLY)
	} else {
		MEM.Layouts[1].Attach("ROM_CD", 0xC000, ROM_CD, mem.READONLY)
		MEM.Layouts[1].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY)
	}
	MEM.Layouts[1].Attach("IO", 0xC000, IO, mem.READWRITE)
	MEM.Layouts[1].Attach("SLOT1", 0xC100, SLOT1, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT2", 0xC200, SLOT2, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT3", 0xC300, SLOT3, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT4", 0xC400, SLOT4, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT5", 0xC500, SLOT5, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT6", 0xC600, SLOT5, mem.READONLY)
	MEM.Layouts[1].Attach("SLOT7", 0xC700, SLOT7, mem.READONLY)
	MEM.Layouts[1].Accessor("IO", IOAccess)
	// MEM.Layouts[1].Show()

	///////////////////////////////////////////////////////////////////////

	MEM.Layouts[2] = mem.InitConfig(ramSize)
	MEM.Layouts[2].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	if model == 1 {
		MEM.Layouts[2].Attach("ROM_D0", 0xD000, ROM_D0, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_D8", 0xD800, ROM_D8, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_E0", 0xE000, ROM_E0, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_E8", 0xE800, ROM_E8, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_F0", 0xF000, ROM_F0, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_F8", 0xF800, ROM_F8, mem.READONLY)
	} else {
		MEM.Layouts[2].Attach("ROM_CD", 0xC000, ROM_CD, mem.READONLY)
		MEM.Layouts[2].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY)
		MEM.Layouts[2].Attach("SLOT3", 0xC300, SLOT3, mem.READONLY)
	}
	MEM.Layouts[2].Attach("IO", 0xC000, IO, mem.READWRITE)
	MEM.Layouts[2].Accessor("IO", IOAccess)
	// MEM.Layouts[0].Show()

	///////////////////////////////////////////////////////////////////////

	MEM.Layouts[3] = mem.InitConfig(ramSize)
	MEM.Layouts[3].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	if model == 1 {
		MEM.Layouts[3].Attach("ROM_D0", 0xD000, ROM_D0, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_D8", 0xD800, ROM_D8, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_E0", 0xE000, ROM_E0, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_E8", 0xE800, ROM_E8, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_F0", 0xF000, ROM_F0, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_F8", 0xF800, ROM_F8, mem.READONLY)
	} else {
		MEM.Layouts[3].Attach("ROM_CD", 0xC000, ROM_CD, mem.READONLY)
		MEM.Layouts[3].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY)
	}
	MEM.Layouts[3].Attach("IO", 0xC000, IO, mem.READWRITE)
	MEM.Layouts[3].Attach("SLOT1", 0xC100, SLOT1, mem.READONLY)
	MEM.Layouts[3].Attach("SLOT2", 0xC200, SLOT2, mem.READONLY)
	MEM.Layouts[3].Attach("SLOT4", 0xC400, SLOT4, mem.READONLY)
	MEM.Layouts[3].Attach("SLOT5", 0xC500, SLOT5, mem.READONLY)
	MEM.Layouts[3].Attach("SLOT6", 0xC600, SLOT5, mem.READONLY)
	MEM.Layouts[3].Attach("SLOT7", 0xC700, SLOT7, mem.READONLY)
	MEM.Layouts[3].Accessor("IO", IOAccess)
	// MEM.Layouts[1].Show()
}
