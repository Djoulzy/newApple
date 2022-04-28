package mos6510

func (C *CPU) initLanguage() {
	C.Mnemonic = map[byte]Instruction{

		0x69: {Name: "ADC", bytes: 2, Cycles: 2, action: C.adc, addr: immediate},
		0x65: {Name: "ADC", bytes: 2, Cycles: 3, action: C.adc, addr: zeropage},
		0x75: {Name: "ADC", bytes: 2, Cycles: 4, action: C.adc, addr: zeropageX},
		0x6D: {Name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absolute},
		0x7D: {Name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absoluteX},
		0x79: {Name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absoluteY},
		0x61: {Name: "ADC", bytes: 2, Cycles: 6, action: C.adc, addr: indirectX},
		0x71: {Name: "ADC", bytes: 2, Cycles: 5, action: C.adc, addr: indirectY},

		// 0x0B: {Name: "ANC", bytes: 2, Cycles: 2, action: C.anc, addr: immediate},

		0x29: {Name: "AND", bytes: 2, Cycles: 2, action: C.and, addr: immediate},
		0x25: {Name: "AND", bytes: 2, Cycles: 3, action: C.and, addr: zeropage},
		0x35: {Name: "AND", bytes: 2, Cycles: 4, action: C.and, addr: zeropageX},
		0x2D: {Name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absolute},
		0x3D: {Name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absoluteX},
		0x39: {Name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absoluteY},
		0x21: {Name: "AND", bytes: 2, Cycles: 6, action: C.and, addr: indirectX},
		0x31: {Name: "AND", bytes: 2, Cycles: 5, action: C.and, addr: indirectY},

		// 0x4B: {Name: "ALR", bytes: 2, Cycles: 2, action: C.alr, addr: immediate},

		0x0A: {Name: "ASL", bytes: 1, Cycles: 2, action: C.asl, addr: implied},
		0x06: {Name: "ASL", bytes: 2, Cycles: 5, action: C.asl, addr: zeropage},
		0x16: {Name: "ASL", bytes: 2, Cycles: 6, action: C.asl, addr: zeropageX},
		0x0E: {Name: "ASL", bytes: 3, Cycles: 6, action: C.asl, addr: absolute},
		0x1E: {Name: "ASL", bytes: 3, Cycles: 7, action: C.asl, addr: absoluteX},

		0x90: {Name: "BCC", bytes: 2, Cycles: 2, action: C.bcc, addr: relative},

		0xB0: {Name: "BCS", bytes: 2, Cycles: 2, action: C.bcs, addr: relative},

		0xF0: {Name: "BEQ", bytes: 2, Cycles: 2, action: C.beq, addr: relative},

		0x24: {Name: "BIT", bytes: 2, Cycles: 3, action: C.bit, addr: zeropage},
		0x2C: {Name: "BIT", bytes: 3, Cycles: 4, action: C.bit, addr: absolute},

		0x30: {Name: "BMI", bytes: 2, Cycles: 2, action: C.bmi, addr: relative},

		0xD0: {Name: "BNE", bytes: 2, Cycles: 2, action: C.bne, addr: relative},

		0x10: {Name: "BPL", bytes: 2, Cycles: 2, action: C.bpl, addr: relative},

		0x00: {Name: "BRK", bytes: 1, Cycles: 7, action: C.brk, addr: implied},

		0x50: {Name: "BVC", bytes: 2, Cycles: 2, action: C.bvc, addr: relative},

		0x70: {Name: "BVS", bytes: 2, Cycles: 2, action: C.bvs, addr: relative},

		0x18: {Name: "CLC", bytes: 1, Cycles: 2, action: C.clc, addr: implied},

		0xD8: {Name: "CLD", bytes: 1, Cycles: 2, action: C.cld, addr: implied},

		0x58: {Name: "CLI", bytes: 1, Cycles: 2, action: C.cli, addr: implied},

		0xB8: {Name: "CLV", bytes: 1, Cycles: 2, action: C.clv, addr: implied},

		0xC9: {Name: "CMP", bytes: 2, Cycles: 2, action: C.cmp, addr: immediate},
		0xC5: {Name: "CMP", bytes: 2, Cycles: 3, action: C.cmp, addr: zeropage},
		0xD5: {Name: "CMP", bytes: 2, Cycles: 4, action: C.cmp, addr: zeropageX},
		0xCD: {Name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absolute},
		0xDD: {Name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absoluteX},
		0xD9: {Name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absoluteY},
		0xC1: {Name: "CMP", bytes: 2, Cycles: 6, action: C.cmp, addr: indirectX},
		0xD1: {Name: "CMP", bytes: 2, Cycles: 5, action: C.cmp, addr: indirectY},

		0xE0: {Name: "CPX", bytes: 2, Cycles: 2, action: C.cpx, addr: immediate},
		0xE4: {Name: "CPX", bytes: 2, Cycles: 3, action: C.cpx, addr: zeropage},
		0xEC: {Name: "CPX", bytes: 3, Cycles: 4, action: C.cpx, addr: absolute},

		0xC0: {Name: "CPY", bytes: 2, Cycles: 2, action: C.cpy, addr: immediate},
		0xC4: {Name: "CPY", bytes: 2, Cycles: 3, action: C.cpy, addr: zeropage},
		0xCC: {Name: "CPY", bytes: 3, Cycles: 4, action: C.cpy, addr: absolute},

		// 0xC7: {Name: "DCP", bytes: 2, Cycles: 5, action: C.dcp, addr: zeropage},
		// 0xD7: {Name: "DCP", bytes: 2, Cycles: 6, action: C.dcp, addr: zeropageX},
		// 0xC3: {Name: "DCP", bytes: 2, Cycles: 8, action: C.dcp, addr: indirectX},
		// 0xD3: {Name: "DCP", bytes: 2, Cycles: 8, action: C.dcp, addr: indirectY},
		// 0xCF: {Name: "DCP", bytes: 3, Cycles: 6, action: C.dcp, addr: absolute},
		// 0xDF: {Name: "DCP", bytes: 3, Cycles: 7, action: C.dcp, addr: absoluteX},
		// 0xDB: {Name: "DCP", bytes: 3, Cycles: 7, action: C.dcp, addr: absoluteY},

		0xC6: {Name: "DEC", bytes: 2, Cycles: 5, action: C.dec, addr: zeropage},
		0xD6: {Name: "DEC", bytes: 2, Cycles: 6, action: C.dec, addr: zeropageX},
		0xCE: {Name: "DEC", bytes: 3, Cycles: 6, action: C.dec, addr: absolute},
		0xDE: {Name: "DEC", bytes: 3, Cycles: 7, action: C.dec, addr: absoluteX},

		0xCA: {Name: "DEX", bytes: 1, Cycles: 2, action: C.dex, addr: implied},

		0x88: {Name: "DEY", bytes: 1, Cycles: 2, action: C.dey, addr: implied},

		0x49: {Name: "EOR", bytes: 2, Cycles: 2, action: C.eor, addr: immediate},
		0x45: {Name: "EOR", bytes: 2, Cycles: 3, action: C.eor, addr: zeropage},
		0x55: {Name: "EOR", bytes: 2, Cycles: 4, action: C.eor, addr: zeropageX},
		0x4D: {Name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absolute},
		0x5D: {Name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absoluteX},
		0x59: {Name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absoluteY},
		0x41: {Name: "EOR", bytes: 2, Cycles: 6, action: C.eor, addr: indirectX},
		0x51: {Name: "EOR", bytes: 2, Cycles: 5, action: C.eor, addr: indirectY},

		0xE6: {Name: "INC", bytes: 2, Cycles: 5, action: C.inc, addr: zeropage},
		0xF6: {Name: "INC", bytes: 2, Cycles: 6, action: C.inc, addr: zeropageX},
		0xEE: {Name: "INC", bytes: 3, Cycles: 6, action: C.inc, addr: absolute},
		0xFE: {Name: "INC", bytes: 3, Cycles: 7, action: C.inc, addr: absoluteX},

		0xE8: {Name: "INX", bytes: 1, Cycles: 2, action: C.inx, addr: implied},

		0xC8: {Name: "INY", bytes: 1, Cycles: 2, action: C.iny, addr: implied},

		// 0xE7: {Name: "ISC", bytes: 2, Cycles: 5, action: C.isc, addr: zeropage},
		// 0xF7: {Name: "ISC", bytes: 2, Cycles: 6, action: C.isc, addr: zeropageX},
		// 0xE3: {Name: "ISC", bytes: 2, Cycles: 8, action: C.isc, addr: indirectX},
		// 0xF3: {Name: "ISC", bytes: 2, Cycles: 8, action: C.isc, addr: indirectY},
		// 0xEF: {Name: "ISC", bytes: 3, Cycles: 6, action: C.isc, addr: absolute},
		// 0xFF: {Name: "ISC", bytes: 3, Cycles: 7, action: C.isc, addr: absoluteX},
		// 0xFB: {Name: "ISC", bytes: 3, Cycles: 7, action: C.isc, addr: absoluteY},

		0x4C: {Name: "JMP", bytes: 3, Cycles: 3, action: C.jmp, addr: absolute},
		0x6C: {Name: "JMP", bytes: 3, Cycles: 5, action: C.jmp, addr: indirect},

		0x20: {Name: "JSR", bytes: 3, Cycles: 6, action: C.jsr, addr: absolute},

		// 0x02: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x12: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x22: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x32: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x42: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x52: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x62: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x72: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0x92: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0xB2: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0xD2: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		// 0xF2: {Name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},

		0xA9: {Name: "LDA", bytes: 2, Cycles: 2, action: C.lda, addr: immediate},
		0xA5: {Name: "LDA", bytes: 2, Cycles: 3, action: C.lda, addr: zeropage},
		0xB5: {Name: "LDA", bytes: 2, Cycles: 4, action: C.lda, addr: zeropageX},
		0xAD: {Name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absolute},
		0xBD: {Name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absoluteX},
		0xB9: {Name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absoluteY},
		0xA1: {Name: "LDA", bytes: 2, Cycles: 6, action: C.lda, addr: indirectX},
		0xB1: {Name: "LDA", bytes: 2, Cycles: 5, action: C.lda, addr: indirectY},

		0xA2: {Name: "LDX", bytes: 2, Cycles: 2, action: C.ldx, addr: immediate},
		0xA6: {Name: "LDX", bytes: 2, Cycles: 3, action: C.ldx, addr: zeropage},
		0xB6: {Name: "LDX", bytes: 2, Cycles: 4, action: C.ldx, addr: zeropageY},
		0xAE: {Name: "LDX", bytes: 3, Cycles: 4, action: C.ldx, addr: absolute},
		0xBE: {Name: "LDX", bytes: 3, Cycles: 4, action: C.ldx, addr: absoluteY},

		0xA0: {Name: "LDY", bytes: 2, Cycles: 2, action: C.ldy, addr: immediate},
		0xA4: {Name: "LDY", bytes: 2, Cycles: 3, action: C.ldy, addr: zeropage},
		0xB4: {Name: "LDY", bytes: 2, Cycles: 4, action: C.ldy, addr: zeropageX},
		0xAC: {Name: "LDY", bytes: 3, Cycles: 4, action: C.ldy, addr: absolute},
		0xBC: {Name: "LDY", bytes: 3, Cycles: 4, action: C.ldy, addr: absoluteX},

		0x4A: {Name: "LSR", bytes: 1, Cycles: 2, action: C.lsr, addr: implied},
		0x46: {Name: "LSR", bytes: 2, Cycles: 5, action: C.lsr, addr: zeropage},
		0x56: {Name: "LSR", bytes: 2, Cycles: 6, action: C.lsr, addr: zeropageX},
		0x4E: {Name: "LSR", bytes: 3, Cycles: 6, action: C.lsr, addr: absolute},
		0x5E: {Name: "LSR", bytes: 3, Cycles: 7, action: C.lsr, addr: absoluteX},

		0xEA: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0x1A: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0x3A: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0x5A: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0x7A: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0xDA: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0xFA: {Name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		// 0x80: {Name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		// 0x82: {Name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		// 0xC2: {Name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		// 0xE2: {Name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		// 0x89: {Name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		// 0x04: {Name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		// 0x44: {Name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		// 0x64: {Name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		// 0x14: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0x34: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0x54: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0x74: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0xD4: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0xF4: {Name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		// 0x0C: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absolute},
		// 0x1C: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		// 0x3C: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		// 0x5C: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		// 0x7C: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		// 0xDC: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		// 0xFC: {Name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},

		0x09: {Name: "ORA", bytes: 2, Cycles: 2, action: C.ora, addr: immediate},
		0x05: {Name: "ORA", bytes: 2, Cycles: 3, action: C.ora, addr: zeropage},
		0x15: {Name: "ORA", bytes: 2, Cycles: 4, action: C.ora, addr: zeropageX},
		0x0D: {Name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absolute},
		0x1D: {Name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absoluteX},
		0x19: {Name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absoluteY},
		0x01: {Name: "ORA", bytes: 2, Cycles: 6, action: C.ora, addr: indirectX},
		0x11: {Name: "ORA", bytes: 2, Cycles: 5, action: C.ora, addr: indirectY},

		0x48: {Name: "PHA", bytes: 1, Cycles: 3, action: C.pha, addr: implied},

		0x08: {Name: "PHP", bytes: 1, Cycles: 3, action: C.php, addr: implied},

		0x68: {Name: "PLA", bytes: 1, Cycles: 4, action: C.pla, addr: implied},

		0x28: {Name: "PLP", bytes: 1, Cycles: 4, action: C.plp, addr: implied},

		// 0x27: {Name: "RLA", bytes: 2, Cycles: 5, action: C.rla, addr: zeropage},
		// 0x37: {Name: "RLA", bytes: 2, Cycles: 6, action: C.rla, addr: zeropageX},
		// 0x23: {Name: "RLA", bytes: 2, Cycles: 8, action: C.rla, addr: indirectX},
		// 0x33: {Name: "RLA", bytes: 2, Cycles: 8, action: C.rla, addr: indirectY},
		// 0x2F: {Name: "RLA", bytes: 3, Cycles: 6, action: C.rla, addr: absolute},
		// 0x3F: {Name: "RLA", bytes: 3, Cycles: 7, action: C.rla, addr: absoluteX},
		// 0x3B: {Name: "RLA", bytes: 3, Cycles: 7, action: C.rla, addr: absoluteY},

		0x2A: {Name: "ROL", bytes: 1, Cycles: 2, action: C.rol, addr: implied},
		0x26: {Name: "ROL", bytes: 2, Cycles: 5, action: C.rol, addr: zeropage},
		0x36: {Name: "ROL", bytes: 2, Cycles: 6, action: C.rol, addr: zeropageX},
		0x2E: {Name: "ROL", bytes: 3, Cycles: 6, action: C.rol, addr: absolute},
		0x3E: {Name: "ROL", bytes: 3, Cycles: 7, action: C.rol, addr: absoluteX},

		0x6A: {Name: "ROR", bytes: 1, Cycles: 2, action: C.ror, addr: implied},
		0x66: {Name: "ROR", bytes: 2, Cycles: 5, action: C.ror, addr: zeropage},
		0x76: {Name: "ROR", bytes: 2, Cycles: 6, action: C.ror, addr: zeropageX},
		0x6E: {Name: "ROR", bytes: 3, Cycles: 6, action: C.ror, addr: absolute},
		0x7E: {Name: "ROR", bytes: 3, Cycles: 7, action: C.ror, addr: absoluteX},

		0x40: {Name: "RTI", bytes: 1, Cycles: 6, action: C.rti, addr: implied},

		0x60: {Name: "RTS", bytes: 1, Cycles: 6, action: C.rts, addr: implied},

		// 0x87: {Name: "SAX", bytes: 2, Cycles: 3, action: C.sax, addr: zeropage},
		// 0x97: {Name: "SAX", bytes: 2, Cycles: 4, action: C.sax, addr: zeropageY},
		// 0x83: {Name: "SAX", bytes: 2, Cycles: 6, action: C.sax, addr: zeropageX},
		// 0x8F: {Name: "SAX", bytes: 3, Cycles: 4, action: C.sax, addr: absolute},

		0xE9: {Name: "SBC", bytes: 2, Cycles: 2, action: C.sbc, addr: immediate},
		0xE5: {Name: "SBC", bytes: 2, Cycles: 3, action: C.sbc, addr: zeropage},
		0xF5: {Name: "SBC", bytes: 2, Cycles: 4, action: C.sbc, addr: zeropageX},
		0xED: {Name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absolute},
		0xFD: {Name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absoluteX},
		0xF9: {Name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absoluteY},
		0xE1: {Name: "SBC", bytes: 2, Cycles: 6, action: C.sbc, addr: indirectX},
		0xF1: {Name: "SBC", bytes: 2, Cycles: 5, action: C.sbc, addr: indirectY},

		// 0xCB: {Name: "SBX", bytes: 2, Cycles: 2, action: C.sbx, addr: immediate},

		// 0x07: {Name: "SLO", bytes: 2, Cycles: 5, action: C.slo, addr: zeropage},
		// 0x17: {Name: "SLO", bytes: 2, Cycles: 6, action: C.slo, addr: zeropageX},
		// 0x03: {Name: "SLO", bytes: 2, Cycles: 8, action: C.slo, addr: indirectX},
		// 0x13: {Name: "SLO", bytes: 2, Cycles: 8, action: C.slo, addr: indirectY},
		// 0x0F: {Name: "SLO", bytes: 3, Cycles: 6, action: C.slo, addr: absolute},
		// 0x1F: {Name: "SLO", bytes: 3, Cycles: 7, action: C.slo, addr: absoluteX},
		// 0x1B: {Name: "SLO", bytes: 3, Cycles: 7, action: C.slo, addr: absoluteY},

		0x38: {Name: "SEC", bytes: 1, Cycles: 2, action: C.sec, addr: implied},

		0xF8: {Name: "SED", bytes: 1, Cycles: 2, action: C.sed, addr: implied},

		0x78: {Name: "SEI", bytes: 1, Cycles: 2, action: C.sei, addr: implied},

		// 0x47: {Name: "SRE", bytes: 2, Cycles: 5, action: C.sre, addr: zeropage},
		// 0x57: {Name: "SRE", bytes: 2, Cycles: 6, action: C.sre, addr: zeropageX},
		// 0x43: {Name: "SRE", bytes: 2, Cycles: 8, action: C.sre, addr: indirectX},
		// 0x53: {Name: "SRE", bytes: 2, Cycles: 8, action: C.sre, addr: indirectY},
		// 0x4F: {Name: "SRE", bytes: 3, Cycles: 6, action: C.sre, addr: absolute},
		// 0x5F: {Name: "SRE", bytes: 3, Cycles: 7, action: C.sre, addr: absoluteX},
		// 0x5B: {Name: "SRE", bytes: 3, Cycles: 7, action: C.sre, addr: absoluteY},

		0x85: {Name: "STA", bytes: 2, Cycles: 3, action: C.sta, addr: zeropage},
		0x95: {Name: "STA", bytes: 2, Cycles: 4, action: C.sta, addr: zeropageX},
		0x8D: {Name: "STA", bytes: 3, Cycles: 4, action: C.sta, addr: absolute},
		0x9D: {Name: "STA", bytes: 3, Cycles: 5, action: C.sta, addr: absoluteX},
		0x99: {Name: "STA", bytes: 3, Cycles: 5, action: C.sta, addr: absoluteY},
		0x81: {Name: "STA", bytes: 2, Cycles: 6, action: C.sta, addr: indirectX},
		0x91: {Name: "STA", bytes: 2, Cycles: 6, action: C.sta, addr: indirectY},

		0x86: {Name: "STX", bytes: 2, Cycles: 3, action: C.stx, addr: zeropage},
		0x96: {Name: "STX", bytes: 2, Cycles: 4, action: C.stx, addr: zeropageY},
		0x8E: {Name: "STX", bytes: 3, Cycles: 4, action: C.stx, addr: absolute},

		0x84: {Name: "STY", bytes: 2, Cycles: 3, action: C.sty, addr: zeropage},
		0x94: {Name: "STY", bytes: 2, Cycles: 4, action: C.sty, addr: zeropageX},
		0x8C: {Name: "STY", bytes: 3, Cycles: 4, action: C.sty, addr: absolute},

		0xAA: {Name: "TAX", bytes: 1, Cycles: 2, action: C.tax, addr: implied},

		0xA8: {Name: "TAY", bytes: 1, Cycles: 2, action: C.tay, addr: implied},

		0xBA: {Name: "TSX", bytes: 1, Cycles: 2, action: C.tsx, addr: implied},

		0x8A: {Name: "TXA", bytes: 1, Cycles: 2, action: C.txa, addr: implied},

		0x9A: {Name: "TXS", bytes: 1, Cycles: 2, action: C.txs, addr: implied},

		0x98: {Name: "TYA", bytes: 1, Cycles: 2, action: C.tya, addr: implied},
	}
}
