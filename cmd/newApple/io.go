package main

type accessor struct {
}

func (C *accessor) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", 0xE800+translatedAddr)
	return mem[translatedAddr]
}

func (C *accessor) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("Accessor", "MWrite", "Addr: %04X -> %02X", 0xE800+translatedAddr, val)
	mem[translatedAddr] = val
}
