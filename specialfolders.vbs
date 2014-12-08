Option Explicit
dim objShell,arg1
set objShell = WScript.CreateObject("WScript.Shell")
for each arg1 in WScript.Arguments
    WScript.Echo( objShell.SpecialFolders(arg1) )
next
