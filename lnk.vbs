Option Explicit
if WScript.Arguments.count < 2 then
    WScript.Echo( _
        "Usage: cscript lnk.vbs FILENAME SHORTCUT ... make shortcut" _
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
shortcut1.TargetPath=src
shortcut1.Save()

WScript.Echo "    " & src & vbcrlf & "--> " & dst
