package main

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

var keyMap map[uint](map[uint]byte)

func initKeyboard() {
	var i uint
	keyMap = make(map[uint]map[uint]byte)

	for i = 0; i < 256; i++ {
		keyMap[i] = make(map[uint]byte)
		keyMap[i][KB_NORMAL] = byte(i)
		keyMap[i][KB_L_SHIFT] = byte(i) - 0x20
		keyMap[i][KB_R_SHIFT] = byte(i) - 0x20
		keyMap[i][KB_L_CTRL] = byte(i) - 0x60
		keyMap[i][KB_R_CTRL] = byte(i) - 0x60
	}

	keyMap[10][KB_NORMAL] = 0x8D // RETURN

	keyMap[51][KB_NORMAL] = 0x22  // "
	keyMap[51][KB_L_SHIFT] = 0x33 // 3
	keyMap[51][KB_R_SHIFT] = 0x33 // 3

	keyMap[53][KB_NORMAL] = 0x28  // (
	keyMap[53][KB_L_SHIFT] = 0x35 // 5
	keyMap[53][KB_R_SHIFT] = 0x35 // 5

	keyMap[41][KB_NORMAL] = 0x29  // )

	keyMap[60][KB_NORMAL] = 0x40  // @
	keyMap[60][KB_L_SHIFT] = 0x23 // #
	keyMap[60][KB_R_SHIFT] = 0x23 // #
}
