package mos6510

import (
	"fmt"
	"time"

	"github.com/Djoulzy/Tools/clog"
)

func (C *CPU) registers() string {
	var i, mask byte
	res := ""
	for i = 0; i < 8; i++ {
		mask = 1 << i
		if C.S&mask == mask {
			res = regString[i] + res
		} else {
			res = "-" + res
		}
	}
	return res
}

func Disassemble(inst instruction, oper uint16) string {
	var token string

	switch inst.addr {
	case implied:
		token = fmt.Sprintf("")
	case immediate:
		token = fmt.Sprintf(" #$%02X", oper)
	case relative:
		token = fmt.Sprintf(" $%02X", oper)
	case zeropage:
		token = fmt.Sprintf(" $%02X", oper)
	case zeropageX:
		token = fmt.Sprintf(" $%02X,X", oper)
	case zeropageY:
		token = fmt.Sprintf(" $%02X,Y", oper)
	case Branching:
		fallthrough
	case CrossPage:
		fallthrough
	case absolute:
		token = fmt.Sprintf(" $%04X", oper)
	case absoluteX:
		token = fmt.Sprintf(" $%04X,X", oper)
	case absoluteY:
		token = fmt.Sprintf(" $%04X,Y", oper)
	case indirect:
		token = fmt.Sprintf(" ($%04X)", oper)
	case indirectX:
		token = fmt.Sprintf(" ($%02X,X)", oper)
	case indirectY:
		token = fmt.Sprintf(" ($%02X),Y", oper)
	}
	return inst.Name + token
}

func (C *CPU) Trace() string {
	return fmt.Sprintf("%s   A:%c[1;33m%02X%c[0m X:%c[1;33m%02X%c[0m Y:%c[1;33m%02X%c[0m SP:%c[1;33m%02X%c[0m   %c[1;31m%04X%c[0m: %-8s %c[1;30m(%d)%c[0m %c[1;37m%-10s%c[0m",
		C.registers(), 27, C.A, 27, 27, C.X, 27, 27, C.Y, 27, 27, C.SP, 27, 27, C.InstStart, 27, C.instDump, 27, C.Inst.Cycles, 27, 27, Disassemble(C.Inst, C.oper), 27)
}

func (C *CPU) DumpRom(start int) [][]byte {
	var code byte
	var listing [][]byte
	var inst instruction
	var ok bool

	listing = make([][]byte, C.ramSize)
	pc := start
	for pc < C.ramSize {
		code = C.ram.Read(uint16(pc))
		if inst, ok = C.Mnemonic[code]; ok {
			// listing[pc] = make([]byte, 5)
			switch inst.addr {
			case implied:
				listing[pc] = []byte{code}
			case immediate:
				// token = fmt.Sprintf(" #$%02X", oper)
				listing[pc] = []byte{code, 1, C.ram.Read(uint16(pc) + 1)}
			case relative:
			case zeropage:
				// token = fmt.Sprintf(" $%02X", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 1)}
			case zeropageX:
				// token = fmt.Sprintf(" $%02X,X", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 1), 0, 1}
			case zeropageY:
				// token = fmt.Sprintf(" $%02X,Y", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 1), 0, 2}
			case absolute:
				// token = fmt.Sprintf(" $%04X", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 2), C.ram.Read(uint16(pc) + 1)}
			case absoluteX:
				// token = fmt.Sprintf(" $%04X,X", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 2), C.ram.Read(uint16(pc) + 1), 4}
			case absoluteY:
				// token = fmt.Sprintf(" $%04X,Y", oper)
				listing[pc] = []byte{code, 2, C.ram.Read(uint16(pc) + 2), C.ram.Read(uint16(pc) + 1), 5}
			case indirect:
				// token = fmt.Sprintf(" ($%04X)", oper)
				listing[pc] = []byte{code, 3, C.ram.Read(uint16(pc) + 2), C.ram.Read(uint16(pc) + 1), 6}
			case indirectX:
				// token = fmt.Sprintf(" ($%02X,X)", oper)
				listing[pc] = []byte{code, 3, C.ram.Read(uint16(pc) + 1), 0, 7}
			case indirectY:
				// token = fmt.Sprintf(" ($%02X),Y", oper)
				listing[pc] = []byte{code, 3, C.ram.Read(uint16(pc) + 1), 0, 8}
			}
			pc += int(inst.bytes)
		} else {
			pc++
		}
	}
	return listing
}

func ColVal(val time.Duration) string {
	if val > time.Microsecond {
		return clog.CSprintf("white", "red", "%10s", val)
	} else {
		return fmt.Sprintf("%10s", val)
	}
}

func (C *CPU) DumpStats() {
	var min time.Duration
	var max time.Duration

	for index, val := range perfStats {
		total := 0
		cpt := 0
		hicount := 0
		min = time.Minute
		max = 0
		for _, duree := range val {
			cpt++
			total += int(duree)
			if duree > time.Microsecond {
				hicount++
			}
			if duree > max {
				max = duree
			}
			if duree < min {
				min = duree
			}
		}
		if cpt > 0 {
			moy := time.Duration(total / cpt)
			hiPercent := float32(hicount) / float32(cpt) * 100
			fmt.Printf("$%02X: (%s) Moy: %s - Max: %s - Min: %s - NbHi: %5d = %6.2f%% - Nb Samples: %d \n", index, C.Mnemonic[index].Name, ColVal(moy), ColVal(max), ColVal(min), hicount, hiPercent, cpt)
		}
	}
}
