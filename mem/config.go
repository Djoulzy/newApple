package mem

import (
	"fmt"

	"github.com/Djoulzy/Tools/clog"
)

const (
	PAGE_DIVIDER = 8
	READWRITE    = false
	READONLY     = true
)

type MEMAccess interface {
	MRead([]byte, uint16) byte
	MWrite([]byte, uint16, byte)
}

type CONFIG struct {
	Layers       [][]byte    // Liste des couches de memoire
	LayersName   []string    // Nom de la couche
	Start        []uint16    // Addresse de début de la couche
	PagesUsed    [][]bool    // Pages utilisées par la couche
	ReadOnly     []bool      // Mode d'accès à la couche
	LayerByPages []int       // Couche active pour la page
	Accessors    []MEMAccess // Reader/Writer de la couche
	TotalPages   int         // Nb total de pages
}

func InitConfig(size int) CONFIG {
	C := CONFIG{}

	C.Layers = make([][]byte, 0, 20)
	C.LayersName = make([]string, 0, 20)
	C.Start = make([]uint16, 0, 20)
	C.PagesUsed = make([][]bool, 0, 20)
	C.ReadOnly = make([]bool, 0, 20)
	C.Accessors = make([]MEMAccess, 0, 20)

	C.TotalPages = int(size >> PAGE_DIVIDER)
	C.LayerByPages = make([]int, C.TotalPages)
	return C
}

func (C *CONFIG) Attach(name string, start uint16, content []byte, mode bool) {
	nbPages := len(content) >> PAGE_DIVIDER
	startPage := int(start >> PAGE_DIVIDER)

	C.Layers = append(C.Layers, content)
	layerNum := len(C.Layers) - 1

	C.LayersName = append(C.LayersName, name)
	C.Start = append(C.Start, start)
	C.ReadOnly = append(C.ReadOnly, mode)
	C.Accessors = append(C.Accessors, C)

	C.PagesUsed = append(C.PagesUsed, make([]bool, C.TotalPages))
	for i := 0; i < C.TotalPages; i++ {
		C.PagesUsed[layerNum][i] = false
	}

	for i := 0; i < nbPages; i++ {
		C.LayerByPages[startPage+i] = layerNum
		C.PagesUsed[layerNum][startPage+i] = true
	}
}

func (C *CONFIG) Accessor(layerName string, access MEMAccess) {
	for index, name := range C.LayersName {
		if name == layerName {
			C.Accessors[index] = access
		}
	}
}

func (C *CONFIG) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("MEM", "MRead", "Addr: %04X -> %02X", addr, mem[addr])
	return mem[translatedAddr]
}

func (C *CONFIG) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("MEM", "MWrite", "Addr: %04X -> %02X", addr, val)
	mem[translatedAddr] = val
}

func (C *CONFIG) Show() {
	var line [128]string
	var pos int
	var indice float32 = 128 / float32(C.TotalPages)

	for layerNum := range C.Layers {
		for page := 0; page < C.TotalPages; page++ {
			pos = int(indice * float32(page))
			if C.PagesUsed[layerNum][page] {
				line[pos] = clog.CSprintf("black", "green", " ")
			} else {
				line[pos] = clog.CSprintf("black", "dark_gray", " ")
			}
		}
		fmt.Printf("%10s:", C.LayersName[layerNum])
		for i := range line {
			fmt.Printf("%s", line[i])
		}
		fmt.Println()
	}
	clog.CPrintf("dark_gray", "black", "\n%12s", " ")
	clog.CPrintf("black", "green", "%s", " Read/Write ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "yellow", "%s", " Read Only ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "red", "%s", " Write Only ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "light_gray", "%s", " Masked ")
	clog.CPrintf("black", "black", "%s", " ")
	fmt.Printf("\n\n")
}
