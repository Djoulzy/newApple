package io

import (
	"newApple/crtc"

	"github.com/Djoulzy/mmu"
)

// PRINT PEEK(49173)

var (
	is_BS_RAM  bool = false
	is_RAMRD   bool = false
	is_RAMWRT  bool = false
	is_BANK2   bool = false
	is_C3_INT  bool = true
	is_CX_INT  bool = false
	is_80Store bool = false
	is_ALT_ZP  bool = false
	is_HIRES   bool = false
	is_PAGE2   bool = false
	is_ALTCHAR bool = false
)

type SoftSwitch struct {
	mmu.IC

	Disks *DiskInterface
	Video *crtc.CRTC
}

func InitSoftSwitch(name string, size uint, disk *DiskInterface, vid *crtc.CRTC) *SoftSwitch {
	tmp := SoftSwitch{
		Disks: disk,
		Video: vid,
	}

	tmp.Size = size
	tmp.Name = name
	tmp.Buff = make([]byte, size)

	return &tmp
}

func (C *SoftSwitch) ReadOnly() bool {
	return false
}
