Option Explicit
if WScript.Arguments.count < 2 then
    WScript.Echo( _
        "Usage: cscript lnk.vbs FILENAME SHORTCUT [WORKINGDIRECTORY]... make shortcut" _
        & vbcrlf & _
        "       cscript lnk.vbs SHORTCUT          ... print shortcut-target")
    WScript.Quit()
end if 
dim src : src=WScript.Arguments.Item(0)
dim dst : dst=WScript.Arguments.Item(1)
dim fsObj : set fsObj=CreateObject("Scripting.FileSystemObject")
src=fsObj.GetAbsolutePathName(src)
dst=fsObj.GetAbsolutePathName(dst)
if fsObj.FolderExists(dst) then
    dst = dst & "\" & fsObj.getFileName(src)
end if
if right(dst,4) <> ".lnk" then
    dst = dst & ".lnk"
end if
dim shell1 : set shell1=CreateObject("WScript.Shell")
dim shortcut1 : set shortcut1=shell1.CreateShortcut(dst)
if shortcut1 is nothing then
    WScript.Quit()
end if
dim workDir : workDir = ""
if WScript.Arguments.Count >= 3 then
    workDir = WScript.Arguments.Item(2)
    workDir = fsObj.GetAbsolutePathName(workDir)
    if fsObj.FolderExists(workDir) then
        shortcut1.WorkingDirectory = workDir
    else
        workDir = ""
    end if
end if

shortcut1.TargetPath=src
shortcut1.Save()

WScript.Echo "    " & src & vbcrlf & "--> " & dst
if workDir <> "" then
    WScript.Echo "    on " & workDir
end if
