print "Nihongo Yet Another GOing Shell"
print "Copyright (c) 2014 HAYAMA_Kaoru and NYAOS.ORG"

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

set "PROMPT=$e[36;40;1m$L$P$G$_$$ $e[37;1m"
alias "ls=ls -oF"
local home = os.getenv("HOME") or os.getenv("USERPROFILE")
if home then
    local rcfname = home .. [[\.nyagos]]
    if exists(rcfname) then
        loadfile(rcfname)()
    end
end
exec "cd"
