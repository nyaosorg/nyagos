package readline

const (
	K_BACKSPACE = "BACKSPACE"
	K_CTRL      = "CTRL"
	K_CTRL_A    = "C_A"
	K_CTRL_B    = "C_B"
	K_CTRL_C    = "C_C"
	K_CTRL_D    = "C_D"
	K_CTRL_E    = "C_E"
	K_CTRL_F    = "C_F"
	K_CTRL_G    = "C_G"
	K_CTRL_H    = "C_H"
	K_CTRL_I    = "C_I"
	K_CTRL_J    = "C_J"
	K_CTRL_K    = "C_K"
	K_CTRL_L    = "C_L"
	K_CTRL_M    = "C_M"
	K_CTRL_N    = "C_N"
	K_CTRL_O    = "C_O"
	K_CTRL_P    = "C_P"
	K_CTRL_Q    = "C_Q"
	K_CTRL_R    = "C_R"
	K_CTRL_S    = "C_S"
	K_CTRL_T    = "C_T"
	K_CTRL_U    = "C_U"
	K_CTRL_V    = "C_V"
	K_CTRL_W    = "C_W"
	K_CTRL_X    = "C_X"
	K_CTRL_Y    = "C_Y"
	K_CTRL_Z    = "C_Z"
	K_DELETE    = "DEL"
	K_DOWN      = "DOWN"
	K_END       = "END"
	K_ENTER     = "ENTER"
	K_ESCAPE    = "ESCAPE"
	K_HOME      = "HOME"
	K_LEFT      = "LEFT"
	K_RIGHT     = "RIGHT"
	K_SHIFT     = "SHIFT"
	K_UP        = "UP"
)

const (
	F_INTR                 = "INTR"
	F_BACKWARD_DELETE_CHAR = "BACKWARD_DELETE_CHAR"
	F_BACKWARD_CHAR        = "BACKWARD_CHAR"
	F_CLEAR_SCREEN         = "CLEAR_SCREEN"
	F_DELETE_CHAR          = "DELETE_CHAR"
	F_DELETE_OR_ABORT      = "DELETE_OR_ABORT"
	F_ACCEPT_LINE          = "ACCEPT_LINE"
	F_KILL_LINE            = "KILL_LINE"
	F_UNIX_LINE_DISCARD    = "UNIX_LINE_DISCARD"
	F_FORARD_CHAR          = "FORWARD_CHAR"
	F_BEGINNING_OF_LINE    = "BEGINNING_OF_LINE"
	F_PASS                 = "PASS"
	F_YANK                 = "YANK"
	F_REPAINT              = "REPAINT"
	F_END_OF_LINE          = "END_OF_LINE"
	F_KILL_WHOLE_LINE      = "KILL_WHOLE_LINE"
)

var NAME2CHAR = map[string]rune{
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

var NAME2SCAN = map[string]uint16{
	K_CTRL:   0x11,
	K_DELETE: 0x2E,
	K_DOWN:   0x28,
	K_END:    0x23,
	K_HOME:   0x24,
	K_LEFT:   0x25,
	K_RIGHT:  0x27,
	K_SHIFT:  0x10,
	K_UP:     0x26,
}

var NAME2FUNC = map[string]KeyFuncT{
	F_BACKWARD_DELETE_CHAR: &KeyGoFuncT{KeyFuncBackSpace},
	F_BACKWARD_CHAR:        &KeyGoFuncT{KeyFuncBackword},
	F_KILL_WHOLE_LINE:      &KeyGoFuncT{KeyFuncClear},
	F_DELETE_CHAR:          &KeyGoFuncT{KeyFuncDelete},
	F_DELETE_OR_ABORT:      &KeyGoFuncT{KeyFuncDeleteOrAbort},
	F_ACCEPT_LINE:          &KeyGoFuncT{KeyFuncEnter},
	F_KILL_LINE:            &KeyGoFuncT{KeyFuncClearAfter},
	F_UNIX_LINE_DISCARD:    &KeyGoFuncT{KeyFuncClearBefore},
	F_FORARD_CHAR:          &KeyGoFuncT{KeyFuncForward},
	F_BEGINNING_OF_LINE:    &KeyGoFuncT{KeyFuncHead},
	F_PASS:                 &KeyGoFuncT{KeyFuncPass},
	F_YANK:                 &KeyGoFuncT{KeyFuncPaste},
	F_CLEAR_SCREEN:         &KeyGoFuncT{KeyFuncCLS},
	F_END_OF_LINE:          &KeyGoFuncT{KeyFuncTail},
	F_INTR:                 &KeyGoFuncT{KeyFuncIntr},
}
