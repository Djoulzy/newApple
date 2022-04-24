package crtc

func (C *CRTC) Read(addr uint16) byte {
	// reg := addr - ((addr >> 6) << 6)
	// clog.Trace("VIC", "Read", "addr: %04X - Reg: %02X (%d)", addr, reg, reg)
	// switch reg {
	// default:
	// 	return C.Reg[reg]
	// }
	return 0
}
