--------------------------------------------------------------------------
-- DO NOT EDIT THIS. PLEASE EDIT ~\.nyagos OR ADD SCRIPT INTO nyagos.d\ --
--------------------------------------------------------------------------

if nyagos == nil then
    print("This is the startup script for NYAGOS")
    print("Do not run this with lua.exe")
    os.exit(0)
end

nyagos.ole = require('nyole')
local fsObj = nyagos.ole.create_object_utf8('Scripting.FileSystemObject')
nyagos.fsObj = fsObj
local nyoleVer = fsObj:GetFileVersion(nyagos.ole.dll_utf8)

print( "Nihongo Yet Another GOing Shell " .. (nyagos.version or "") ..
    " Powered by " ..  _VERSION .. " & nyole.dll ".. nyoleVer)
if not nyagos.version or string.len(nyagos.version) <= 0 then
    print("Build at "..(nyagos.stamp or "").." with commit "..(nyagos.commit or ""))
end
print("Copyright (c) 2014,2015 HAYAMA_Kaoru and NYAOS.ORG")

local function include(fname)
    local chank,err=loadfile(fname)
    if err then
        print(err)
    elseif chank then
        local ok,err=pcall(chank)
        if not ok then
            print(fname .. ": " ..err)
        end
    else
        print(fname .. ":fail to load")
    end
end

local dotFolderPath = string.gsub(nyagos.exe,"%.exe$",".d")
local dotFolder = fsObj:GetFolder(dotFolderPath)
local files = dotFolder:files()
for p in files:__iter__() do
    if string.match(p.Name,"%.[lL][uU][aA]$") then
        local path = nyagos.fsObj:BuildPath(dotFolderPath,p.Name)
        include(path)
    end
end

local home = nyagos.getenv("HOME") or nyagos.getenv("USERPROFILE")
if home then 
    local dotfile = fsObj:BuildPath(home,'.nyagos')
    local fd=io.open(dotfile)
    if fd then
        fd:close()
        include(dotfile)
    end
end
