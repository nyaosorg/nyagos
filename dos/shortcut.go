package dos

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func ReadShortcut(path string) (string, string, error) {
	agent, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return "", "", err
	}
	defer agent.Release()
	agentDis, err := agent.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return "", "", err
	}
	agentDis.Release()
	shortcut, err := oleutil.CallMethod(agentDis, "CreateShortCut", path)
	if err != nil {
		return "", "", err
	}
	shortcutDis := shortcut.ToIDispatch()
	defer shortcutDis.Release()
	targetPath, err := oleutil.GetProperty(shortcutDis, "TargetPath")
	if err != nil {
		return "", "", err
	}
	workingDir, err := oleutil.GetProperty(shortcutDis, "WorkingDirectory")
	if err != nil {
		return "", "", err
	}
	return targetPath.ToString(), workingDir.ToString(), err
}
