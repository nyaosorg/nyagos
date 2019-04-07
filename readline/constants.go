package readline

import "context"

const (
	K_BACKSPACE     = "BACKSPACE"
	K_CAPSLOCK      = "CAPSLOCK"
	K_CLEAR         = "CLEAR"
	K_CTRL          = "CTRL"
	K_CTRL_A        = "C_A"
	K_CTRL_B        = "C_B"
	K_CTRL_BREAK    = "C_BREAK"
	K_CTRL_C        = "C_C"
	K_CTRL_D        = "C_D"
	K_CTRL_E        = "C_E"
	K_CTRL_F        = "C_F"
	K_CTRL_G        = "C_G"
	K_CTRL_H        = "C_H"
	K_CTRL_I        = "C_I"
	K_CTRL_J        = "C_J"
	K_CTRL_K        = "C_K"
	K_CTRL_L        = "C_L"
	K_CTRL_M        = "C_M"
	K_CTRL_N        = "C_N"
	K_CTRL_O        = "C_O"
	K_CTRL_P        = "C_P"
	K_CTRL_Q        = "C_Q"
	K_CTRL_R        = "C_R"
	K_CTRL_S        = "C_S"
	K_CTRL_T        = "C_T"
	K_CTRL_U        = "C_U"
	K_CTRL_V        = "C_V"
	K_CTRL_W        = "C_W"
	K_CTRL_X        = "C_X"
	K_CTRL_Y        = "C_Y"
	K_CTRL_Z        = "C_Z"
	K_DELETE        = "DEL"
	K_DOWN          = "DOWN"
	K_CTRL_DOWN     = "C_DOWN"
	K_END           = "END"
	K_ENTER         = "ENTER"
	K_ESCAPE        = "ESCAPE"
	K_F1            = "F1"
	K_F10           = "F10"
	K_F11           = "F11"
	K_F12           = "F12"
	K_F13           = "F13"
	K_F14           = "F14"
	K_F15           = "F15"
	K_F16           = "F16"
	K_F17           = "F17"
	K_F18           = "F18"
	K_F19           = "F19"
	K_F2            = "F2"
	K_F20           = "F20"
	K_F21           = "F21"
	K_F22           = "F22"
	K_F23           = "F23"
	K_F24           = "F24"
	K_F3            = "F3"
	K_F4            = "F4"
	K_F5            = "F5"
	K_F6            = "F6"
	K_F7            = "F7"
	K_F8            = "F8"
	K_F9            = "F9"
	K_HOME          = "HOME"
	K_LEFT          = "LEFT"
	K_CTRL_LEFT     = "C_LEFT"
	K_PAGEDOWN      = "PAGEDOWN"
	K_PAGEUP        = "PAGEUP"
	K_PAUSE         = "PAUSE"
	K_RIGHT         = "RIGHT"
	K_CTRL_RIGHT    = "C_RIGHT"
	K_SHIFT         = "SHIFT"
	K_UP            = "UP"
	K_CTRL_UP       = "C_UP"
	K_ALT_A         = "M_A"
	K_ALT_B         = "M_B"
	K_ALT_BACKSPACE = "M_BACKSPACE"
	K_ALT_BREAK     = "M_BREAK"
	K_ALT_C         = "M_C"
	K_ALT_D         = "M_D"
	K_ALT_E         = "M_E"
	K_ALT_F         = "M_F"
	K_ALT_G         = "M_G"
	K_ALT_H         = "M_H"
	K_ALT_I         = "M_I"
	K_ALT_J         = "M_J"
	K_ALT_K         = "M_K"
	K_ALT_L         = "M_L"
	K_ALT_M         = "M_M"
	K_ALT_N         = "M_N"
	K_ALT_O         = "M_O"
	K_ALT_P         = "M_P"
	K_ALT_Q         = "M_Q"
	K_ALT_R         = "M_R"
	K_ALT_S         = "M_S"
	K_ALT_T         = "M_T"
	K_ALT_U         = "M_U"
	K_ALT_V         = "M_V"
	K_ALT_W         = "M_W"
	K_ALT_X         = "M_X"
	K_ALT_Y         = "M_Y"
	K_ALT_Z         = "M_Z"
	K_ALT_OEM_2     = "M_OEM_2"
)

const (
	F_ACCEPT_LINE          = "ACCEPT_LINE"
	F_BACKWARD_CHAR        = "BACKWARD_CHAR"
	F_BACKWARD_WORD        = "BACKWARD_WORD"
	F_BACKWARD_DELETE_CHAR = "BACKWARD_DELETE_CHAR"
	F_BEGINNING_OF_LINE    = "BEGINNING_OF_LINE"
	F_CLEAR_SCREEN         = "CLEAR_SCREEN"
	F_DELETE_CHAR          = "DELETE_CHAR"
	F_DELETE_OR_ABORT      = "DELETE_OR_ABORT"
	F_END_OF_LINE          = "END_OF_LINE"
	F_FORWARD_CHAR         = "FORWARD_CHAR"
	F_FORWARD_WORD         = "FORWARD_WORD"
	F_HISTORY_DOWN         = "HISTORY_DOWN" // for compatible
	F_HISTORY_UP           = "HISTORY_UP"   // for compatible
	F_NEXT_HISTORY         = "NEXT_HISTORY"
	F_PREVIOUS_HISTORY     = "PREVIOUS_HISTORY"
	F_INTR                 = "INTR"
	F_ISEARCH_BACKWARD     = "ISEARCH_BACKWARD"
	F_KILL_LINE            = "KILL_LINE"
	F_KILL_WHOLE_LINE      = "KILL_WHOLE_LINE"
	F_PASS                 = "PASS"
	F_QUOTED_INSERT        = "QUOTED_INSERT"
	F_REPAINT_ON_NEWLINE   = "REPAINT_ON_NEWLINE"
	F_SWAPCHAR             = "SWAPCHAR"
	F_UNIX_LINE_DISCARD    = "UNIX_LINE_DISCARD"
	F_UNIX_WORD_RUBOUT     = "UNIX_WORD_RUBOUT"
	F_YANK                 = "YANK"
	F_YANK_WITH_QUOTE      = "YANK_WITH_QUOTE"
)

var name2char = map[string]string{
	K_BACKSPACE: "\x7F",
	K_CTRL_A:    "\x01",
	K_CTRL_B:    "\x02",
	K_CTRL_C:    "\x03",
	K_CTRL_D:    "\x04",
	K_CTRL_E:    "\x05",
	K_CTRL_F:    "\x06",
	K_CTRL_G:    "\x07",
	K_CTRL_H:    "\x08",
	K_CTRL_I:    "\x09",
	K_CTRL_J:    "\x0A",
	K_CTRL_K:    "\x0B",
	K_CTRL_L:    "\x0C",
	K_CTRL_M:    "\x0D",
	K_CTRL_N:    "\x0E",
	K_CTRL_O:    "\x0F",
	K_CTRL_P:    "\x10",
	K_CTRL_Q:    "\x11",
	K_CTRL_R:    "\x12",
	K_CTRL_S:    "\x13",
	K_CTRL_T:    "\x14",
	K_CTRL_U:    "\x15",
	K_CTRL_V:    "\x16",
	K_CTRL_W:    "\x17",
	K_CTRL_X:    "\x18",
	K_CTRL_Y:    "\x19",
	K_CTRL_Z:    "\x20",
	// K_DELETE:        "\x7F",
	K_ENTER:         "\r",
	K_ESCAPE:        "\x1B",
	K_ALT_A:         "\x1Ba",
	K_ALT_B:         "\x1Bb",
	K_ALT_BACKSPACE: "\x1B\b",
	K_ALT_C:         "\x1Bc",
	K_ALT_D:         "\x1Bd",
	K_ALT_E:         "\x1Be",
	K_ALT_F:         "\x1Bf",
	K_ALT_G:         "\x1Bg",
	K_ALT_H:         "\x1Bh",
	K_ALT_I:         "\x1Bi",
	K_ALT_J:         "\x1Bj",
	K_ALT_K:         "\x1Bk",
	K_ALT_L:         "\x1Bl",
	K_ALT_M:         "\x1Bm",
	K_ALT_N:         "\x1Bn",
	K_ALT_O:         "\x1Bo",
	K_ALT_P:         "\x1Bp",
	K_ALT_Q:         "\x1Bq",
	K_ALT_R:         "\x1Br",
	K_ALT_S:         "\x1Bs",
	K_ALT_T:         "\x1Bt",
	K_ALT_U:         "\x1Bu",
	K_ALT_V:         "\x1Bv",
	K_ALT_W:         "\x1Bw",
	K_ALT_X:         "\x1Bx",
	K_ALT_Y:         "\x1By",
	K_ALT_Z:         "\x1Bz",
	K_CLEAR:         "0x0C",
	K_CTRL:          "0x11",
	K_CTRL_BREAK:    "0x03",
	K_DELETE:        "\x1B[3~",
	K_DOWN:          "\x1B[B",
	K_CTRL_DOWN:     "\x1B[1;5B",
	K_END:           "\x1B[F",
	K_F10:           "\x1B[21~",
	K_F11:           "\x1B[23~",
	K_F12:           "\x1B[24~",
	K_F13:           "0x7C",
	K_F14:           "0x7D",
	K_F15:           "0x7E",
	K_F16:           "0x7F",
	K_F17:           "0x80",
	K_F18:           "0x81",
	K_F19:           "0x82",
	K_F1:            "\x1B[OP",
	K_F20:           "0x83",
	K_F21:           "0x84",
	K_F22:           "0x85",
	K_F23:           "0x86",
	K_F24:           "0x87",
	K_F2:            "\x1B[OQ",
	K_F3:            "\x1B[OR",
	K_F4:            "\x1B[OS",
	K_F5:            "\x1B[15~",
	K_F6:            "\x1B[16~",
	K_F7:            "\x1B[17~",
	K_F8:            "\x1B[18~",
	K_F9:            "\x1B[20~",
	K_HOME:          "\x1B[H",
	K_LEFT:          "\x1B[D",
	K_CTRL_LEFT:     "\x1B[1;5D",
	K_PAGEDOWN:      "\x1B[6~",
	K_PAGEUP:        "\x1B[5~",
	K_PAUSE:         "0x13",
	K_RIGHT:         "\x1B[C",
	K_CTRL_RIGHT:    "\x1B[1;5C",
	K_UP:            "\x1B[A",
	K_CTRL_UP:       "\x1B[1;5A",
}

// KeyCode from
// http://msdn.microsoft.com/ja-jp/library/windows/desktop/dd375731(v=vs.85).aspx

var NAME2FUNC = map[string]func(context.Context, *Buffer) Result{
	F_ACCEPT_LINE:          keyFuncEnter,
	F_BACKWARD_CHAR:        keyFuncBackward,
	F_BACKWARD_WORD:        keyFuncBackwardWord,
	F_BACKWARD_DELETE_CHAR: keyFuncBackSpace,
	F_BEGINNING_OF_LINE:    keyFuncHead,
	F_CLEAR_SCREEN:         keyFuncCLS,
	F_DELETE_CHAR:          keyFuncDelete,
	F_DELETE_OR_ABORT:      keyFuncDeleteOrAbort,
	F_END_OF_LINE:          keyFuncTail,
	F_FORWARD_CHAR:         keyFuncForward,
	F_FORWARD_WORD:         keyFuncForwardWord,
	F_HISTORY_DOWN:         keyFuncHistoryDown, // for compatible
	F_HISTORY_UP:           keyFuncHistoryUp,   // for compatible
	F_NEXT_HISTORY:         keyFuncHistoryDown,
	F_PREVIOUS_HISTORY:     keyFuncHistoryUp,
	F_INTR:                 keyFuncIntr,
	F_ISEARCH_BACKWARD:     keyFuncIncSearch,
	F_KILL_LINE:            keyFuncClearAfter,
	F_KILL_WHOLE_LINE:      keyFuncClear,
	F_PASS:                 nil,
	F_QUOTED_INSERT:        keyFuncQuotedInsert,
	F_UNIX_LINE_DISCARD:    keyFuncClearBefore,
	F_UNIX_WORD_RUBOUT:     keyFuncWordRubout,
	F_YANK:                 keyFuncPaste,
	F_YANK_WITH_QUOTE:      keyFuncPasteQuote,
	F_SWAPCHAR:             keyFuncSwapChar,
	F_REPAINT_ON_NEWLINE:   keyFuncRepaintOnNewline,
}

func name2func(keyName string) KeyFuncT {
	if p, ok := NAME2FUNC[keyName]; ok {
		return &KeyGoFuncT{
			Func: p,
			Name: keyName,
		}
	} else {
		return nil
	}
}
