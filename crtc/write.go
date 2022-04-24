package crtc

func (C *CRTC) Write(addr uint16, val byte) {
	// reg := addr - ((addr >> 6) << 6)
	// // clog.Trace("VIC", "Write", "addr: %04X - Reg: %02X (%d)", addr, reg, reg)
	// switch reg {
	// default:
	// 	C.Reg[reg] = val
	// }
}

func (C *CRTC) testBit(reg uint16, mask byte) bool {
	if C.Reg[reg]&mask == mask {
		return true
	}
	return false
}
