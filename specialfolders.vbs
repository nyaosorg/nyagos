Option Explicit
Dim objShell,arg1
Set objShell = WScript.CreateObject("WScript.Shell")
for each arg1 in WScript.Arguments
    WScript.Echo( objShell.SpecialFolders(arg1) )
next
