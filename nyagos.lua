print("Nihongo Yet Another GOing Shell")
print(string.format("Build at %s with commit %s",nyagos.stamp,nyagos.commit))
print("Copyright (c) 2014 HAYAMA_Kaoru and NYAOS.ORG")

-- This is system-lua files which will be updated.
-- When you want to customize, please edit ~\.nyagos
-- Please do not edit this.

io.getenv = nyagos.getenv
io.setenv = nyagos.setenv

local function split(equation)
    local pos=string.find(equation,"=",1,true)
    if pos then
        local left=string.sub(equation,1,pos-1)
        local right=string.sub(equation,pos+1)
        return left,right,pos
    else
        return nil,nil,nil
    end
end

local function expand(text)
    return string.gsub(text,"%%(%w+)%%",function(w)
        return nyagos.getenv(w)
    end)
end

function hasList(list,target)
    local LIST=";"..string.upper(list)..";"
    local TARGET=";"..string.upper(target)..";"
    return string.find(LIST,TARGET,1,true)
end

function addpath(...)
    for _,dir in pairs{...} do
        dir = expand(dir)
        local list=nyagos.getenv("PATH")
        if not hasList(list,dir) then
            nyagos.setenv("PATH",dir..";"..list)
        end
    end
end

function set(equation)
    if type(equation) == 'table' then
        for left,right in pairs(equation) do
            nyagos.setenv(left,expand(right))
        end
    else
        local left,right,pos = split(equation)
        if pos and string.sub(left,-1) == "+" then
            left = string.sub(left,1,-2)
            local original=nyagos.getenv(left)
            if string.find(right,original,1,true) then
                right = original
            else
                right = right .. ";" .. original
            end
        end
        if right then
            nyagos.setenv(left,expand(right))
        end
    end
end

function alias(equation)
    if type(equation) == 'table' then
        for left,right in pairs(equation) do
            nyagos.alias(left,right)
        end
    else
        local left,right,pos = split(equation)
        if right then
            nyagos.alias(left,right)
        end
    end
end

function exists(path)
    local fd=io.open(path)
    if fd then
        fd:close()
        return true
    else
        return false
    end
end

x = nyagos.exec
original_print = print
print = nyagos.echo

nyagos.suffixes={}
function suffix(suffix,cmdline)
    local suffix=string.lower(suffix)
    if string.sub(suffix,1,1)=='.' then
        suffix = string.sub(suffix,2)
    end
    if not nyagos.suffixes[suffix] then
        local orgpathext = nyagos.getenv("PATHEXT")
        local newext="."..suffix
        if not hasList(orgpathext,newext) then
            nyagos.setenv("PATHEXT",orgpathext..";"..newext)
        end
    end
    nyagos.suffixes[suffix]=cmdline
end
suffix(".pl",{"perl"})
suffix(".py",{"ipy"})
suffix(".rb",{"ruby"})
suffix(".lua",{"lua"})
suffix(".awk",{"awk","-f"})
suffix(".js",{"cscript","//nologo"})
suffix(".vbs",{"cscript","//nologo"})

nyagos.argsfilter = function(args)
    local path=nyagos.which(args[0])
    if not path then
        return
    end
    local m = string.match(path,"%.(%w+)$")
    if not m then 
        return
    end
    local cmdline = nyagos.suffixes[ string.lower(m) ]
    if not cmdline then
        return
    end
    local newargs={}
    for i=1,#cmdline do
        newargs[i-1]=cmdline[i]
    end
    newargs[#cmdline] = path
    for i=1,#args do
        newargs[#cmdline+i] = args[i]
    end
    return newargs
end

nyagos.filter = function(cmdline)
    return string.gsub(cmdline,'`([^`]*)`',function(m)
        local r = nyagos.eval(m)
        if not r then
            return false
        end
        r = nyagos.atou(r)
        r = string.gsub(r,'[|&<>!]',function(m)
            return string.format('%%u+%04X%%',string.byte(m,1,1))
        end)
        return string.gsub(r,'%s+$','')
    end)
end

alias{
    assoc='%COMSPEC% /c assoc $*',
    attrib='%COMSPEC% /c attrib $*',
    copy='%COMSPEC% /c copy $*',
    del='%COMSPEC% /c del $*',
    dir='%COMSPEC% /c dir $*',
    ['for']='%COMSPEC% /c for $*',
    md='%COMSPEC% /c md $*',
    mkdir='%COMSPEC% /c mkdir $*',
    mklink='%COMSPEC% /c mklink $*',
    move='%COMSPEC% /c move $*',
    open=function(args)
        local count=0
        for i=1,#args do
            nyagos.exec(string.format('%s /c start "%s"',nyagos.getenv('COMSPEC'),args[i]))
            count = count +1
        end
        if count <= 0 then
            nyagos.exec(string.format('%s /c start .',nyagos.getenv('COMSPEC')))
        end
    end,
    rd='%COMSPEC% /c rd $*',
    ren='%COMSPEC% /c ren $*',
    rename='%COMSPEC% /c rename $*',
    rmdir='%COMSPEC% /c rmdir $*',
    start='%COMSPEC% /c start $*',
    ['type']='%COMSPEC% /c type $*',
    suffix=function(args)
        if #args < 2 then
            print "Usage: suffix SUFFIX COMMAND"
        else
            suffix(args[1],args[2])
        end
    end,
    ls='ls -oF $*',
    lua_e=function(args)
        assert(load(args[1]))()
    end,
    which=function(args)
        nyagos.echo( nyagos.which(args[1]) )
    end
}

local home = nyagos.getenv("HOME") or nyagos.getenv("USERPROFILE")
if home then
    x'cd'
    local rcfname = '.nyagos'
    if exists(rcfname) then
        local chank,err=loadfile(rcfname)
        if chank then
            chank()
        elseif err then
            print(err)
        end
    end
end
