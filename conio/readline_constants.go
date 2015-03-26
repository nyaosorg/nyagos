package conio

const (
	K_BACKSPACE  = "BACKSPACE"
	K_CAPSLOCK   = "CAPSLOCK"
	K_CLEAR      = "CLEAR"
	K_CTRL       = "CTRL"
	K_CTRL_A     = "C_A"
	K_CTRL_B     = "C_B"
	K_CTRL_BREAK = "C_BREAK"
	K_CTRL_C     = "C_C"
	K_CTRL_D     = "C_D"
	K_CTRL_E     = "C_E"
	K_CTRL_F     = "C_F"
	K_CTRL_G     = "C_G"
	K_CTRL_H     = "C_H"
	K_CTRL_I     = "C_I"
	K_CTRL_J     = "C_J"
	K_CTRL_K     = "C_K"
	K_CTRL_L     = "C_L"
	K_CTRL_M     = "C_M"
	K_CTRL_N     = "C_N"
	K_CTRL_O     = "C_O"
	K_CTRL_P     = "C_P"
	K_CTRL_Q     = "C_Q"
	K_CTRL_R     = "C_R"
	K_CTRL_S     = "C_S"
	K_CTRL_T     = "C_T"
	K_CTRL_U     = "C_U"
	K_CTRL_V     = "C_V"
	K_CTRL_W     = "C_W"
	K_CTRL_X     = "C_X"
	K_CTRL_Y     = "C_Y"
	K_CTRL_Z     = "C_Z"
	K_DELETE     = "DEL"
	K_DOWN       = "DOWN"
	K_END        = "END"
	K_ENTER      = "ENTER"
	K_ESCAPE     = "ESCAPE"
	K_F1         = "F1"
	K_F10        = "F10"
	K_F11        = "F11"
	K_F12        = "F12"
	K_F13        = "F13"
	K_F14        = "F14"
	K_F15        = "F15"
	K_F16        = "F16"
	K_F17        = "F17"
	K_F18        = "F18"
	K_F19        = "F19"
	K_F2         = "F2"
	K_F20        = "F20"
	K_F21        = "F21"
	K_F22        = "F22"
	K_F23        = "F23"
	K_F24        = "F24"
	K_F3         = "F3"
	K_F4         = "F4"
	K_F5         = "F5"
	K_F6         = "F6"
	K_F7         = "F7"
	K_F8         = "F8"
	K_F9         = "F9"
	K_HOME       = "HOME"
	K_LEFT       = "LEFT"
	K_PAGEDOWN   = "PAGEDOWN"
	K_PAGEUP     = "PAGEUP"
	K_PAUSE      = "PAUSE"
	K_RIGHT      = "RIGHT"
	K_SHIFT      = "SHIFT"
	K_UP         = "UP"
	K_ALT_A      = "M_A"
	K_ALT_B      = "M_B"
	K_ALT_BREAK  = "M_BREAK"
	K_ALT_C      = "M_C"
	K_ALT_D      = "M_D"
	K_ALT_E      = "M_E"
	K_ALT_F      = "M_F"
	K_ALT_G      = "M_G"
	K_ALT_H      = "M_H"
	K_ALT_I      = "M_I"
	K_ALT_J      = "M_J"
	K_ALT_K      = "M_K"
	K_ALT_L      = "M_L"
	K_ALT_M      = "M_M"
	K_ALT_N      = "M_N"
	K_ALT_O      = "M_O"
	K_ALT_P      = "M_P"
	K_ALT_Q      = "M_Q"
	K_ALT_R      = "M_R"
	K_ALT_S      = "M_S"
	K_ALT_T      = "M_T"
	K_ALT_U      = "M_U"
	K_ALT_V      = "M_V"
	K_ALT_W      = "M_W"
	K_ALT_X      = "M_X"
	K_ALT_Y      = "M_Y"
	K_ALT_Z      = "M_Z"
)

const (
	F_ACCEPT_LINE          = "ACCEPT_LINE"
	F_BACKWARD_CHAR        = "BACKWARD_CHAR"
	F_BACKWARD_DELETE_CHAR = "BACKWARD_DELETE_CHAR"
	F_BEGINNING_OF_LINE    = "BEGINNING_OF_LINE"
	F_CLEAR_SCREEN         = "CLEAR_SCREEN"
	F_DELETE_CHAR          = "DELETE_CHAR"
	F_DELETE_OR_ABORT      = "DELETE_OR_ABORT"
	F_END_OF_LINE          = "END_OF_LINE"
	F_FORARD_CHAR          = "FORWARD_CHAR"
	F_HISTORY_DOWN         = "HISTORY_DOWN"
	F_HISTORY_UP           = "HISTORY_UP"
	F_INTR                 = "INTR"
	F_ISEARCH_BACKWARD     = "ISEARCH_BACKWARD"
	F_KILL_LINE            = "KILL_LINE"
	F_KILL_WHOLE_LINE      = "KILL_WHOLE_LINE"
	F_PASS                 = "PASS"
	F_REPAINT              = "REPAINT"
	F_UNIX_LINE_DISCARD    = "UNIX_LINE_DISCARD"
	F_YANK                 = "YANK"
)

var name2char = map[string]rune{
	K_BACKSPACE: '\b',
	K_CTRL_A:    rune('a' & 0x1F),
	K_CTRL_B:    rune('b' & 0x1F),
	K_CTRL_C:    rune('c' & 0x1F),
	K_CTRL_D:    rune('d' & 0x1F),
	K_CTRL_E:    rune('e' & 0x1F),
	K_CTRL_F:    rune('f' & 0x1F),
	K_CTRL_G:    rune('g' & 0x1F),
	K_CTRL_H:    rune('h' & 0x1F),
	K_CTRL_I:    rune('i' & 0x1F),
	K_CTRL_J:    rune('j' & 0x1F),
	K_CTRL_K:    rune('k' & 0x1F),
	K_CTRL_L:    rune('l' & 0x1F),
	K_CTRL_M:    rune('m' & 0x1F),
	K_CTRL_N:    rune('n' & 0x1F),
	K_CTRL_O:    rune('o' & 0x1F),
	K_CTRL_P:    rune('p' & 0x1F),
	K_CTRL_Q:    rune('q' & 0x1F),
	K_CTRL_R:    rune('r' & 0x1F),
	K_CTRL_S:    rune('s' & 0x1F),
	K_CTRL_T:    rune('t' & 0x1F),
	K_CTRL_U:    rune('u' & 0x1F),
	K_CTRL_V:    rune('v' & 0x1F),
	K_CTRL_W:    rune('w' & 0x1F),
	K_CTRL_X:    rune('x' & 0x1F),
	K_CTRL_Y:    rune('y' & 0x1F),
	K_CTRL_Z:    rune('z' & 0x1F),
	K_DELETE:    '\x7F',
	K_ENTER:     '\r',
	K_ESCAPE:    rune('[' & 0x1F),
}

// KeyCode from
// http://msdn.microsoft.com/ja-jp/library/windows/desktop/dd375731(v=vs.85).aspx

var name2scan = map[string]uint16{
	K_CAPSLOCK:   0x14,
	K_CLEAR:      0x0C,
	K_CTRL:       0x11,
	K_CTRL_BREAK: 0x03,
	K_DELETE:     0x2E,
	K_DOWN:       0x28,
	K_END:        0x23,
	K_F10:        0x79,
	K_F11:        0x7A,
	K_F12:        0x7B,
	K_F13:        0x7C,
	K_F14:        0x7D,
	K_F15:        0x7E,
	K_F16:        0x7F,
	K_F17:        0x80,
	K_F18:        0x81,
	K_F19:        0x82,
	K_F1:         0x70,
	K_F20:        0x83,
	K_F21:        0x84,
	K_F22:        0x85,
	K_F23:        0x86,
	K_F24:        0x87,
	K_F2:         0x71,
	K_F3:         0x72,
	K_F4:         0x73,
	K_F5:         0x74,
	K_F6:         0x75,
	K_F7:         0x76,
	K_F8:         0x77,
	K_F9:         0x78,
	K_HOME:       0x24,
	K_LEFT:       0x25,
	K_PAGEDOWN:   0x22,
	K_PAGEUP:     0x21,
	K_PAUSE:      0x13,
	K_RIGHT:      0x27,
	K_SHIFT:      0x10,
	K_UP:         0x26,
}

var name2alt = map[string]uint16{
	K_ALT_A: 0x41,
	K_ALT_B: 0x42,
	K_ALT_C: 0x43,
	K_ALT_D: 0x44,
	K_ALT_E: 0x45,
	K_ALT_F: 0x46,
	K_ALT_G: 0x47,
	K_ALT_H: 0x48,
	K_ALT_I: 0x49,
	K_ALT_J: 0x4A,
	K_ALT_K: 0x4B,
	K_ALT_L: 0x4C,
	K_ALT_M: 0x4D,
	K_ALT_N: 0x4E,
	K_ALT_O: 0x4F,
	K_ALT_P: 0x50,
	K_ALT_Q: 0x51,
	K_ALT_R: 0x52,
	K_ALT_S: 0x53,
	K_ALT_T: 0x54,
	K_ALT_U: 0x55,
	K_ALT_V: 0x56,
	K_ALT_W: 0x57,
	K_ALT_X: 0x58,
	K_ALT_Y: 0x59,
	K_ALT_Z: 0x5A,
}

var name2func = map[string]KeyFuncT{
	F_ACCEPT_LINE:          &KeyGoFuncT{KeyFuncEnter},
	F_BACKWARD_CHAR:        &KeyGoFuncT{KeyFuncBackword},
	F_BACKWARD_DELETE_CHAR: &KeyGoFuncT{KeyFuncBackSpace},
	F_BEGINNING_OF_LINE:    &KeyGoFuncT{KeyFuncHead},
	F_CLEAR_SCREEN:         &KeyGoFuncT{KeyFuncCLS},
	F_DELETE_CHAR:          &KeyGoFuncT{KeyFuncDelete},
	F_DELETE_OR_ABORT:      &KeyGoFuncT{KeyFuncDeleteOrAbort},
	F_END_OF_LINE:          &KeyGoFuncT{KeyFuncTail},
	F_FORARD_CHAR:          &KeyGoFuncT{KeyFuncForward},
	F_HISTORY_DOWN:         &KeyGoFuncT{KeyFuncHistoryDown},
	F_HISTORY_UP:           &KeyGoFuncT{KeyFuncHistoryUp},
	F_INTR:                 &KeyGoFuncT{KeyFuncIntr},
	F_ISEARCH_BACKWARD:     &KeyGoFuncT{KeyFuncIncSearch},
	F_KILL_LINE:            &KeyGoFuncT{KeyFuncClearAfter},
	F_KILL_WHOLE_LINE:      &KeyGoFuncT{KeyFuncClear},
	F_PASS:                 &KeyGoFuncT{KeyFuncPass},
	F_UNIX_LINE_DISCARD:    &KeyGoFuncT{KeyFuncClearBefore},
	F_YANK:                 &KeyGoFuncT{KeyFuncPaste},
}
