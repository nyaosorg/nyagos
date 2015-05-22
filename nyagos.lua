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

function x(s)
    for line in string.gmatch(s,'[^\r\n]+') do
        nyagos.exec(line)
    end
end
print( "Nihongo Yet Another GOing Shell " .. (nyagos.version or "") ..
    " Powered by " ..  _VERSION .. " & nyole.dll ".. nyoleVer)
if not nyagos.version or string.len(nyagos.version) <= 0 then
    print("Build at "..(nyagos.stamp or "").." with commit "..(nyagos.commit or ""))
end
print("Copyright (c) 2014,2015 HAYAMA_Kaoru and NYAOS.ORG")

local function expand(text)
    return string.gsub(text,"%%(%w+)%%",function(w)
        return nyagos.getenv(w)
    end)
end

local function set_(f,equation,expand)
    if type(equation) == 'table' then
        for left,right in pairs(equation) do
            f(left,expand(right))
        end
        return true
    end
    local pluspos=string.find(equation,"+=",1,true)
    if pluspos and pluspos > 0 then
        local left=string.sub(equation,1,pluspos-1)
        equation = string.format("%s=%s;%%%s%%",
                        left,string.sub(equation,pluspos+2),left)
    end
    local pos=string.find(equation,"=",1,true)
    if pos then
        local left=string.sub(equation,1,pos-1)
        local right=string.sub(equation,pos+1)
        f( left , expand(right) )
        return true
    end
    return false,(equation .. ': invalid format')
end

function set(equation) 
    set_(nyagos.setenv,equation,expand)
end
function alias(equation)
    set_(nyagos.alias,equation,function(x) return x end)
end
function addpath(...)
    for _,dir in pairs{...} do
        dir = expand(dir)
        local list=nyagos.getenv("PATH")
        if not string.find(";"..list..";",";"..dir..";",1,true) then
            nyagos.setenv("PATH",dir..";"..list)
        end
    end
end
function nyagos.echo(s)
    nyagos.write((s or '<nil>')..'\n')
end
io.getenv = nyagos.getenv
io.setenv = nyagos.setenv
original_print = print
print = nyagos.echo

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
