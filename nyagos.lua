print "Nihongo Yet Another GOing Shell"
print "Copyright (c) 2014 HAYAMA_Kaoru and NYAOS.ORG"

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

function set(equation)
    local left,right,pos = split(equation)
    if pos and string.sub(left,-1) == "+" then
        left = string.sub(left,1,-2)
        local original=os.getenv(left)
        if string.find(right,original) then
            right = right .. ";" .. original
        else
            right = original
        end
    end
    if right then
        right = string.gsub(right,"%%(%w+)%%",function(w)
            return os.getenv(w)
        end)
        nyagos.setenv(left,right)
    end
end

function alias(equation)
    local left,right,pos = split(equation)
    if right then
        nyagos.alias(left,right)
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

exec = nyagos.exec

alias 'assoc=%COMSPEC% /c assoc $*'
alias 'attrib=%COMSPEC% /c attrib $*'
alias 'copy=%COMSPEC% /c copy $*'
alias 'del=%COMSPEC% /c del $*'
alias 'dir=%COMSPEC% /c dir $*'
alias 'for=%COMSPEC% /c for $*'
alias 'md=%COMSPEC% /c md $*'
alias 'mkdir=%COMSPEC% /c mkdir $*'
alias 'mklink=%COMSPEC% /c mklink $*'
alias 'move=%COMSPEC% /c move $*'
alias 'open=%COMSPEC% /c for %I in ($*) do @start "%I"'
alias 'rd=%COMSPEC% /c rd $*'
alias 'ren=%COMSPEC% /c ren $*'
alias 'rename=%COMSPEC% /c rename $*'
alias 'rmdir=%COMSPEC% /c rmdir $*'
alias 'start=%COMSPEC% /c start $*'
alias 'type=%COMSPEC% /c type $*'
alias 'ls=ls -oF $*'

local home = os.getenv("HOME") or os.getenv("USERPROFILE")
if home then
    exec "cd"
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
