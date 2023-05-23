package main

// const (
// 	DISABLED  = true
// 	ENABLED   = false
// 	READWRITE = false
// 	READONLY  = true
// )

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
