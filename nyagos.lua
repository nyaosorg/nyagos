print "Nihongo Yet Another GOing Shell (c) 2014 K.Hayama"
setenv("PROMPT","$e[36;40;1m$L$P$G$_$$ $e[37;1m")
local home = os.getenv("USERPROFILE") or os.getenv("HOME")
if home then
    local rcfname = home .. [[\.nyagos]]
    fd = io.open(rcfname)
    if fd then
        fd:close()
        loadfile(fcfname)
    end
end

function addpath(path1)
    setenv(path1..";"..os.getenv("PATH"))
end

