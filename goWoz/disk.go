package woz

import (
	"log"

	"github.com/DataDog/go-python3"
)

type Disk struct {
	pyRef       *python3.PyObject
	actualTrack float64
	Tracks      map[float64]*Track
}

func NewWozDisk(fileName string) *Disk {
	WozDiskImage := python3.PyDict_GetItemString(oDict, "WozDiskImage")
	log.Printf("%v", WozDiskImage)
	if !(WozDiskImage != nil && python3.PyCallable_Check(WozDiskImage)) {
		panic("Can't instantiate WozDiskImage")
	}

	args := python3.PyTuple_New(1)
	defer args.DecRef()
	python3.PyTuple_SetItem(args, 0, python3.PyBytes_FromString(fileName))

	woz := Disk{}
	woz.Tracks = make(map[float64]*Track)
	woz.pyRef = WozDiskImage.CallObject(args)
	return &woz
}

func (W *Disk) Close() {
	W.pyRef.DecRef()
}

func (W *Disk) Seek(num float64) *Track {
	if val, ok := W.Tracks[num]; ok {
		// val.Find("D5 AA 96")
		return val
	}
	ret := W.pyRef.CallMethodArgs("seek", python3.PyFloat_FromDouble(num))
	if ret == python3.Py_None {
		// log.Printf("None detected")
		W.Tracks[num] = nil
		return nil
	}
	W.Tracks[num] = NewWozTrack(ret)
	return W.Tracks[num]
}

func (W *Disk) Dump() []byte {
	return []byte(python3.PyBytes_AsString(W.pyRef.CallMethodArgs("to_json", nil)))
}
