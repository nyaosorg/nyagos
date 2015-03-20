Option Explicit
If WScript.Arguments.Count < 2 Then
    WScript.Echo( _
        "Usage: cscript lnk.vbs FILENAME SHORTCUT {Option=Value}... make shortcut" _
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
    dst = fsObj.BuildPath(dst,fsObj.GetFileName(src))
End If
If Right(dst,4) <> ".lnk" Then
    dst = dst & ".lnk"
End If
Dim shell1 : Set shell1=CreateObject("WScript.Shell")
Dim shortcut1 : Set shortcut1=shell1.CreateShortcut(dst)
If shortcut1 Is Nothing Then
    WScript.Quit()
End If
shortcut1.TargetPath=src
If WScript.Arguments.Count >= 3 Then
    Dim i
    For i=2 to WScript.Arguments.Count-1
        Dim equation : equation = WScript.Arguments.Item(i)
        Dim pos : pos = InStr(equation,"=")
        If pos >= 0 Then
            equation = Left(equation,pos-1) & "=""" & Mid(equation,pos+1) & """"
            WScript.Echo equation
            Execute "shortcut1." & equation
        End If
    Next
End If
shortcut1.Save()
If Err.Number <> 0 Then
    WScript.Echo "Error: " & Err.Description
Else
    WScript.Echo "    " & src & vbcrlf & "--> " & dst
End If
