local version = os.getenv("VERSION")
if not version or not string.match(version,"^%d+%.%d+%.%d+%_%d+$") then
    print [[
#include <winver.h>

GOPHER ICON        NYAGOS.ICO
]]
    os.exit()
end

local dict = {
    ["%VERSION%"] = version,
    ["%VCOMMA%"] = string.gsub(version,"[%.%_]",","),
}
local text,count = string.gsub([[
#include <winver.h>

GOPHER ICON        NYAGOS.ICO

VS_VERSION_INFO    VERSIONINFO 
FILEVERSION        %VCOMMA%
PRODUCTVERSION     %VCOMMA%
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
            VALUE "FileVersion",      "%VERSION%\0"
            VALUE "LegalCopyright",   "Copyright (C) 2014-2015 HAYAMA_Kaoru\0"
            VALUE "OriginalFilename", "NYAGOS.EXE\0"
            VALUE "ProductName",      "Nihongo Yet Another GOing Shell\0"
            VALUE "ProductVersion",   "%VERSION%\0"
        END 
    END
END
]],"%%[^%%]+%%",function(m)
    if dict[m] then
        return dict[m]
    else
        return m
    end
end)

print(text)
