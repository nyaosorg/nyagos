// +build !windows

package frame

func coInitialize() {}

func coUnInitialize() {}

func enableVirtualTerminalProcessing() {}

func isEscapeSequenceAvailable() bool {
	return true
}
