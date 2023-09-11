package io

const (
	KB_NORMAL  = 0
	KB_L_SHIFT = 1073742049
	KB_R_SHIFT = 1073742053
	KB_L_CTRL  = 1073742048
	KB_R_CTRL  = 1073742052
	KB_L_META  = 1073742051
	KB_R_META  = 1073742055
	KB_L_ALT   = 1073742050
	KB_R_ALT   = 1073742054
)

var (
	KeyMap        map[uint](map[uint]byte)
	Is_Keypressed bool = false
)

func InitKeyboard() {
	var i uint
	KeyMap = make(map[uint]map[uint]byte)

	for i = 0; i < 256; i++ {
		KeyMap[i] = make(map[uint]byte)
		KeyMap[i][KB_NORMAL] = byte(i)
		KeyMap[i][KB_L_SHIFT] = byte(i) - 0x20
		KeyMap[i][KB_R_SHIFT] = byte(i) - 0x20
		KeyMap[i][KB_L_CTRL] = byte(i) - 0x60
		KeyMap[i][KB_R_CTRL] = byte(i) - 0x60
	}

	KeyMap[10][KB_NORMAL] = 0x8D // RETURN

	KeyMap[51][KB_NORMAL] = 0x33  // "
	KeyMap[51][KB_L_SHIFT] = 0x22 // 3
	KeyMap[51][KB_R_SHIFT] = 0x22 // 3

	KeyMap[53][KB_NORMAL] = 0x35  // (
	KeyMap[53][KB_L_SHIFT] = 0x28 // 5
	KeyMap[53][KB_R_SHIFT] = 0x28 // 5

	KeyMap[41][KB_NORMAL] = 0x29 // )

	KeyMap[60][KB_NORMAL] = 0x40  // @
	KeyMap[60][KB_L_SHIFT] = 0x23 // #
	KeyMap[60][KB_R_SHIFT] = 0x23 // #
}
