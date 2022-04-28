package mos6510

import (
	"log"
)

func (C *CPU) bcc() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if !C.issetC() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bcs() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if C.issetC() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) beq() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if C.issetZ() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bmi() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if C.issetN() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bne() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if !C.issetZ() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bpl() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if !C.issetN() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) brk() {
	switch C.Inst.addr {
	case implied:
		C.pushWordStack(C.PC + 1)
		C.setB(true)
		C.pushByteStack(C.S)
		C.PC = C.readWord(IRQBRK_Vector)
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bvc() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if !C.issetV() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bvs() {
	switch C.Inst.addr {
	case relative:
		C.Oper = C.getRelativeAddr(C.Oper)
		if C.issetV() {
			C.Inst.addr = Branching
			C.State = Compute
			C.Inst.Cycles++
		}
	case Branching:
		if C.PC&0xFF00 == C.Oper&0xFF00 {
			C.PC = C.Oper
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case CrossPage:
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) jmp() {
	switch C.Inst.addr {
	case absolute:
		C.PC = C.Oper
	case indirect:
		C.PC = C.readWord(C.Oper)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) jsr() {
	switch C.Inst.addr {
	case absolute:
		C.pushWordStack(C.InstStart + 2)
		C.PC = C.Oper
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) rti() {
	switch C.Inst.addr {
	case implied:
		C.S = C.pullByteStack()
		C.setB(false)
		C.setU(false)
		C.PC = C.pullWordStack()
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) rts() {
	switch C.Inst.addr {
	case implied:
		C.PC = C.pullWordStack() + 1
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) nop() {
	switch C.Inst.addr {
	case implied:
		fallthrough
	case immediate:
		fallthrough
	case zeropage:
		fallthrough
	case zeropageX:
		fallthrough
	case absolute:
		fallthrough
	case absoluteX:
	default:
		log.Fatal("Bad addressing mode")
	}

}
