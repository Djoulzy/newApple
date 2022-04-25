package main

var keyMap map[uint]byte //{
// 0:          Keyb_NULL,
// 8:          Keyb_DEL,
// 13:         Keyb_RETURN,
// 27:         Keyb_RUNSTOP,
// 32:         Keyb_SPACE,
// 44:         Keyb_COMMA,
// 45:         Keyb_MINUS,
// 48:         Keyb_0,
// 49:         Keyb_1,
// 50:         Keyb_2,
// 51:         Keyb_3,
// 52:         Keyb_4,
// 53:         Keyb_5,
// 54:         Keyb_6,
// 55:         Keyb_7,
// 56:         Keyb_8,
// 57:         Keyb_9,
// 58:         Keyb_COLON,    // :
// 59:         Keyb_SEMICOLN, // ;
// 60:         Keyb_AROBASE,
// 61:         Keyb_EQUAL,
// 94:         Keyb_ARROW_UP,
// 97: 0xC1,
// 98: 0xC2,
// 99:         Keyb_C,
// 100:        Keyb_D,
// 101:        Keyb_E,
// 102:        Keyb_F,
// 103:        Keyb_G,
// 104:        Keyb_H,
// 105:        Keyb_I,
// 106:        Keyb_J,
// 107:        Keyb_K,
// 108:        Keyb_L,
// 109:        Keyb_M,
// 110:        Keyb_N,
// 111:        Keyb_O,
// 112:        Keyb_P,
// 113:        Keyb_Q,
// 114:        Keyb_R,
// 115:        Keyb_S,
// 116:        Keyb_T,
// 117:        Keyb_U,
// 118:        Keyb_V,
// 119:        Keyb_W,
// 120:        Keyb_X,
// 121:        Keyb_Y,
// 122:        Keyb_Z,
// 1073742048: Keyb_CTRL,
// 1073742049: Keyb_LSHIFT,
// 1073742051: Keyb_CBM,
// 1073742053: Keyb_RSHIFT,
// 1073741905: Keyb_CRSR_DOWN,
// 1073741903: Keyb_CRSR_RIGHT,
//}

func initKeyboard() {
	var i uint
	keyMap = make(map[uint]byte)
	cpt := 0xB0
	for i = 48; i < 58; i++ {
		keyMap[i] = byte(cpt)
		cpt++
	}

	cpt = 0xC1
	for i = 97; i < 123; i++ {
		keyMap[i] = byte(cpt)
		cpt++
	}

	keyMap[13] = 0x8D // Return
	keyMap[32] = 0xA0 // Space
	keyMap[45] = 0xAD // Minus
}
