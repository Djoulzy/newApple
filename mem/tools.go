package mem

import (
	"fmt"
	"io/ioutil"

	"github.com/Djoulzy/Tools/clog"
)

func LoadROM(size int, file string) []byte {
	val := make([]byte, size)
	if len(file) > 0 {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if len(data) != size {
			panic("Bad ROM Size")
		}
		for i := 0; i < size; i++ {
			val[i] = byte(data[i])
		}
	}
	return val
}

func Clear(zone []byte) {
	cpt := 0
	fill := byte(0x00)
	for i := range zone {
		zone[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
}

func Fill(zone []byte, val byte) {
	for i := range zone {
		zone[i] = val
	}
}

func dispBin(bin byte) []byte {
	var mask byte = 0b10000000
	var res string = ""
	for i := 0; i < 8; i++ {
		if (bin & mask) > 0 {
			res += clog.CSprintf("black", "white", " ")
		} else {
			res += clog.CSprintf("black", "dark_gray", " ")
		}
		mask >>= 1
	}
	return []byte(res)
}

func DisplayCharRom(zone []byte, bytePerLine int, nbLines int, nbDispPerLine int) {
	var charStartAddr int
	var charLine []byte

	size := len(zone)
	charSize := bytePerLine * nbLines
	nbChar := size / charSize

	clog.CPrintf("light_gray", "black", "- Nb Char found: %d\n", nbChar)
	for y := 0; y < nbChar; y += nbDispPerLine {
		for l := 0; l < nbLines; l++ {
			for x := 0; x < nbDispPerLine; x++ {
				charStartAddr = (y * charSize) + (x * nbLines) + l
				charLine = dispBin(zone[charStartAddr])
				if l == 0 {
					clog.CPrintf("light_gray", "black", "%03X", y+x)
				} else {
					fmt.Printf("   ")
				}
				clog.CPrintf("light_gray", "black", "%s ", charLine)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
