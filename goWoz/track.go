package woz

import (
	"log"

	"github.com/DataDog/go-python3"
	"github.com/Djoulzy/Tools/clog"
)

type Track struct {
	pyRef *python3.PyObject
	iter  *python3.PyObject
}

func NewWozTrack(trk *python3.PyObject) *Track {
	track := Track{}
	track.pyRef = trk
	// track.Find("D5 AA 96")
	// track.iter = track.pyRef.CallMethodArgs("nibble")
	return &track
}

func (T *Track) Find(pattern string) bool {
	ret := T.pyRef.CallMethodArgs("find", python3.PyBytes_FromString(pattern))
	if python3.PyBool_Check(ret) {
		log.Printf("NO DATA")
		return false
	}
	return true
}

func (T *Track) Nibble() int {
	iter := T.pyRef.CallMethodArgs("nibble")
	num := iter.CallMethodArgs("__next__")
	bytes := python3.PyLong_AsLong(num)
	clog.FileRaw("%02X ", bytes)
	// fmt.Printf("%02X ", bytes)
	return bytes
}

func (T *Track) Close() {
	T.pyRef.DecRef()
}
