package main

import (
	"os"

	"github.com/Djoulzy/emutools/mem"
)

func DumpMem(mem *mem.BANK, file string) error {
	var tmp []byte
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	tmp = make([]byte, ramSize)
	for i := 0; i < ramSize; i++ {
		// tmp[i] = mem.Read(uint16(i))
		tmp[i] = mem.Layouts[0].Layers[0][i]
	}
	f.Write(tmp)

	return nil
}
