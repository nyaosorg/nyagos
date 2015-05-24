WScript.Echo("#include <winver.h>");
WScript.Echo("");
WScript.Echo("GOPHER ICON NYAGOS.ICO");

var objWsh = new ActiveXObject("WScript.Shell");
var version = objWsh.ExpandEnvironmentStrings("%VERSION%");
if( ! version || ! version.match(/^\d+\.\d+\.\d+\_\d+$/) ){
    WScript.Quit(1);
}
var vcomma = version.replace(/[%.%_]/g,",");

WScript.Echo('VS_VERSION_INFO    VERSIONINFO ')
WScript.Echo('FILEVERSION        ' + vcomma)
WScript.Echo('PRODUCTVERSION     ' + vcomma)
WScript.Echo('FILEFLAGSMASK      VS_FFI_FILEFLAGSMASK ')
WScript.Echo('FILEFLAGS          0x00000000L')
WScript.Echo('FILEOS             VOS__WINDOWS32')
WScript.Echo('FILETYPE           VFT_APP')
WScript.Echo('FILESUBTYPE        VFT2_UNKNOWN')
WScript.Echo('BEGIN')
WScript.Echo('    BLOCK "VarFileInfo"')
WScript.Echo('    BEGIN')
WScript.Echo('        VALUE "Translation", 0x0411, 0x04E4')
WScript.Echo('    END')
WScript.Echo('    BLOCK "StringFileInfo"')
WScript.Echo('    BEGIN')
WScript.Echo('        BLOCK "041104E4"')
WScript.Echo('        BEGIN')
WScript.Echo('            VALUE "CompanyName",      "NYAOS.ORG"')
WScript.Echo('            VALUE "FileDescription",  "Extended Commandline Shell"')
WScript.Echo('            VALUE "FileVersion",      "'+version+'\\0"')
WScript.Echo('            VALUE "LegalCopyright",   "Copyright (C) 2014-2015 HAYAMA_Kaoru\\0"')
WScript.Echo('            VALUE "OriginalFilename", "NYAGOS.EXE\\0"')
WScript.Echo('            VALUE "ProductName",      "Nihongo Yet Another GOing Shell\\0"')
WScript.Echo('            VALUE "ProductVersion",   "'+version+'\\0"')
WScript.Echo('        END ')
WScript.Echo('    END')
WScript.Echo('END')
