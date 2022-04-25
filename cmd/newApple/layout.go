package main

import (
	"newApple/mem"
)

func memLayouts() {
	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0x0000, RAM, mem.READWRITE)
	MEM.Layouts[0].Attach("IO", 0xC000, IO, mem.READWRITE)
	MEM.Layouts[0].Attach("ROM_CD", 0xC000, ROM_CD, mem.READONLY)
	MEM.Layouts[0].Attach("ROM_EF", 0xE000, ROM_EF, mem.READONLY)
	// MEM.Layouts[0].Attach("KERNAL", 0xF000, KERNAL, mem.READONLY)
	// MEM.Layouts[0].Accessor("IO", IOAccess)
	MEM.Layouts[0].Show()
}
