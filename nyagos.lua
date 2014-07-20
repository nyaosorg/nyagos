print "Nihongo Yet Another GOing Shell (c) 2014 K.Hayama"

function set(equation)
    local pos=string.find(equation,"=",1,true)
    local left=string.sub(equation,1,pos-1)
    local right=string.sub(equation,pos+1)
    if pos and string.sub(left,-1) == "+" then
        left = string.sub(left,1,-2)
        local original=os.getenv(left)
        if string.find(right,original) then
            right = right .. ";" .. original
        else
            right = original
        end
    end
    right = string.gsub(right,"%%(%w+)%%",function(w)
        return os.getenv(w)
    end)
    setenv(left,right)
end

set "PROMPT=$e[36;40;1m$L$P$G$_$$ $e[37;1m"
alias("ls","ls -oF")
local home = os.getenv("HOME") or os.getenv("USERPROFILE")
if home then
    local rcfname = home .. [[\.nyagos]]
    fd = io.open(rcfname)
    if fd then
        fd:close()
        loadfile(rcfname)()
    end
end
system "cd"
