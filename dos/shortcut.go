package dos

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// WScriptShell create OLE Object for "WScript.Shell"
func WScriptShell() (*ole.IUnknown, *ole.IDispatch, error) {
	agent, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return nil, nil, err
	}
	agentDis, err := agent.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		agent.Release()
		return nil, nil, err
	}
	return agent, agentDis, nil
}

// ReadShortcut reads *.lnk file and returns targetpath and working-directory.
func ReadShortcut(path string) (string, string, error) {
	agent, agentDis, err := WScriptShell()
	if err != nil {
		return "", "", err
	}
	defer agent.Release()
	defer agentDis.Release()
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

// MakeShortcut makes *.lnk file
func MakeShortcut(from, to, dir string) error {
	agent, agentDis, err := WScriptShell()
	if err != nil {
		return err
	}
	defer agent.Release()
	defer agentDis.Release()
	shortcut, err := oleutil.CallMethod(agentDis, "CreateShortCut", to)
	if err != nil {
		return err
	}
	shortcutDis := shortcut.ToIDispatch()
	defer shortcutDis.Release()
	_, err = oleutil.PutProperty(shortcutDis, "TargetPath", from)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(shortcutDis, "WorkingDirectory", dir)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(shortcutDis, "Save")
	return err
}
