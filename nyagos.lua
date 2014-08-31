print("Nihongo Yet Another GOing Shell")
print(string.format("Build at %s with commit %s",nyagos.stamp,nyagos.commit))
print("Copyright (c) 2014 HAYAMA_Kaoru and NYAOS.ORG")

-- This is system-lua files which will be updated.
-- When you want to customize, please edit ~\.nyagos
-- Please do not edit this.

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
        return os.getenv(w)
    end)
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
            local original=os.getenv(left)
            if string.find(right,original) then
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
    open='%COMSPEC% /c for %I in ($*) do @start "%I"',
    rd='%COMSPEC% /c rd $*',
    ren='%COMSPEC% /c ren $*',
    rename='%COMSPEC% /c rename $*',
    rmdir='%COMSPEC% /c rmdir $*',
    start='%COMSPEC% /c start $*',
    ['type']='%COMSPEC% /c type $*',
    ls='ls -oF $*',
    lua_e=function(...)
        local args={...}
        assert(load(args[2]))()
    end,
    which=function(...)
        local args={...}
        for dir1 in string.gmatch(os.getenv('PATH'),"[^;]+") do
            for ext1 in string.gmatch(os.getenv('PATHEXT'),"[^;]+") do
                local path1 = dir1 .. "\\" .. args[2] .. ext1
                if exists(path1) then
                    nyagos.echo(path1)
                end
            end
        end
    end
}

local home = os.getenv("HOME") or os.getenv("USERPROFILE")
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
