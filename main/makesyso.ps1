Set-PSDebug -strict

# If windres.exe is not found, quit immediately.
try{
    Get-Command windres -ErrorAction Stop | %{}
}catch{
    Write-Host "makesyso.ps1: windres.exe not found"
    exit
}

if ( [IO.File]::Exists("nyagos.syso") -and -not $env:version ){
    Write-Host "makesyso.ps1: nothing to do"
    exit
}

$text = @"
#include <winver.h>

GOPHER ICON NYAGOS.ICO
"@

$version = $env:version
if ($version ){
    Write-Host "makesyso.ps1: set icon and version $version"
    $vcomma = ($version -replace "[\._]",",")
    $text = $text + @"

VS_VERSION_INFO    VERSIONINFO
FILEVERSION        $vcomma
PRODUCTVERSION     $vcomma
FILEFLAGSMASK      VS_FFI_FILEFLAGSMASK
FILEFLAGS          0x00000000L
FILEOS             VOS__WINDOWS32
FILETYPE           VFT_APP
FILESUBTYPE        VFT2_UNKNOWN
BEGIN
    BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x0411, 0x04E4
    END
    BLOCK "StringFileInfo"
    BEGIN
        BLOCK "041104E4"
        BEGIN
            VALUE "CompanyName",      "NYAOS.ORG"
            VALUE "FileDescription",  "Extended Commandline Shell"
            VALUE "FileVersion",      "$version\0"
            VALUE "LegalCopyright",   "Copyright (C) 2014-2016 HAYAMA_Kaoru\0"
            VALUE "OriginalFilename", "NYAGOS.EXE\0"
            VALUE "ProductName",      "Nihongo Yet Another GOing Shell\0"
            VALUE "ProductVersion",   "$version\0"
        END
    END
END
"@
}else{
    Write-Host "makesyso.ps1: set icon."
}

$text | windres.exe --output-format=coff -o nyagos.syso
