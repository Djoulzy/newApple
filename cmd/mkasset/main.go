package main

import (
	"newApple/graphic"
	"newApple/mem"
	"newApple/mos6510"
)

const ramSize = 65536

var cpu mos6510.CPU
var MEM mem.BANK
var RAM []byte

func main() {
	var test byte = 0
	var inst []string

	inst = make([]string, 256)
	RAM = make([]byte, ramSize)
	MEM = mem.InitBanks(1, &test)
	memLayouts()

	cpu.Init(&MEM, nil)

	for index, value := range cpu.Mnemonic {
		inst[index] = value.Name
	}
	graphic.MakeAsset(inst)
}
