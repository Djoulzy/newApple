package woz

import (
	"os"
	"path/filepath"

	"github.com/DataDog/go-python3"
)

var oImport, oModule, oDict *python3.PyObject

func SetupLib() {
	python3.Py_Initialize()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + dir + "\")")
	oImport = python3.PyImport_ImportModule("goWoz.wozardry") //ret val: new ref

	oModule := python3.PyImport_AddModule("goWoz.wozardry") //ret val: borrowed ref (from oImport)
	oDict = python3.PyModule_GetDict(oModule)
}

func CloseLib() {
	// oDict.DecRef()
	oImport.DecRef()
	python3.Py_Finalize()
}
