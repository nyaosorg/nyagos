Option Explicit
Dim objShell,arg1
Set objShell = WScript.CreateObject("WScript.Shell")
For Each arg1 In WScript.Arguments
    WScript.Echo( objShell.SpecialFolders(arg1) )
Next
