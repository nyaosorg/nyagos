Option Explicit
If WScript.Arguments.Count < 2 Then
    WScript.Echo( _
        "Usage: cscript lnk.vbs FILENAME SHORTCUT [WORKINGDIRECTORY]... make shortcut" _
        & vbcrlf & _
        "       cscript lnk.vbs SHORTCUT          ... print shortcut-target")
    WScript.Quit()
End If
Dim src : src=WScript.Arguments.Item(0)
Dim dst : dst=WScript.Arguments.Item(1)
Dim fsObj : Set fsObj=CreateObject("Scripting.FileSystemObject")
src=fsObj.GetAbsolutePathName(src)
dst=fsObj.GetAbsolutePathName(dst)
If fsObj.FolderExists(dst) Then
    dst = dst & "\" & fsObj.getFileName(src)
End If
If Right(dst,4) <> ".lnk" Then
    dst = dst & ".lnk"
End If
Dim shell1 : Set shell1=CreateObject("WScript.Shell")
Dim shortcut1 : Set shortcut1=shell1.CreateShortcut(dst)
If shortcut1 Is Nothing Then
    WScript.Quit()
End If
Dim workDir : workDir = ""
If WScript.Arguments.Count >= 3 Then
    workDir = WScript.Arguments.Item(2)
    workDir = fsObj.GetAbsolutePathName(workDir)
    If fsObj.FolderExists(workDir) Then
        shortcut1.WorkingDirectory = workDir
    Else
        workDir = ""
    End If
End If

shortcut1.TargetPath=src
shortcut1.Save()

WScript.Echo "    " & src & vbcrlf & "--> " & dst
If workDir <> "" Then
    WScript.Echo "    on " & workDir
End If
