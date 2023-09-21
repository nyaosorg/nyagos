package textwidth

import (
	"os"
)

var (
	isVsCodeTerminal = os.Getenv("VSCODE_PID") != ""

	isWindowsTerminal = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" && !isVsCodeTerminal

	ambiguousIsWide = !isWindowsTerminal
)

var RuneWidth = newRuneWidth(ambiguousIsWide)
