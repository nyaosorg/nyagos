nyagos._suffixes={}

nyagos._setsuffix = function(suffix,cmdline)
    local suffix=string.lower(suffix)
    if string.sub(suffix,1,1)=='.' then
        suffix = string.sub(suffix,2)
    end
    if not nyagos._suffixes[suffix] then
        local orgpathext = nyagos.getenv("PATHEXT")
        local newext="."..suffix
        if not string.find(";"..orgpathext..";",";"..newext..";",1,true) then
            nyagos.setenv("PATHEXT",orgpathext..";"..newext)
        end
    end
    local table = nyagos._suffixes
    table[suffix]=cmdline
    nyagos._suffixes = table
end

suffix = setmetatable({},{
    __call = function(t,k,v) nyagos._setsuffix(k,v) return end,
    __newindex = function(t,k,v) nyagos._setsuffix(k,v) return end,
    __index = function(t,k) return nyagos._suffixes[k] end
})

nyagos._org_suffix_argsfilter=nyagos.argsfilter
nyagos.argsfilter = function(args)
    if nyagos._org_suffix_argsfilter then
        local args_ = nyagos._org_suffix_argsfilter(args)
        if args_ then
            args = args_
        end
    end
    local path=nyagos.which(args[0])
    if not path then
        return
    end
    local m = string.match(path,"%.(%w+)$")
    if not m then
        return
    end
    local cmdline = nyagos._suffixes[ string.lower(m) ]
    if not cmdline then
        return
    end
    local newargs={}
    if type(cmdline) == 'table' then
        for i=1,#cmdline do
            newargs[i-1]=cmdline[i]
        end
    elseif type(cmdline) == 'string' then
        newargs[0] = cmdline
    end
    newargs[#newargs+1] = path
    for i=1,#args do
        newargs[#newargs+1] = args[i]
    end
    return newargs
end

nyagos.alias.suffix = function(args)
    if #args < 2 then
        print "Usage: suffix SUFFIX COMMAND"
    else
        local cmdline={}
        for i=2,#args do
            cmdline[#cmdline+1]=args[i]
        end
        suffix(args[1],cmdline)
    end
end

suffix.pl="perl"
if nyagos.which("ipy") then
  suffix.py="ipy"
elseif nyagos.which("py") then
  suffix.py="py"
else
  suffix.py="python"
end
suffix.rb="ruby"
suffix.lua="lua"
suffix.awk={"awk","-f"}
suffix.js={"cscript","//nologo"}
suffix.vbs={"cscript","//nologo"}
suffix.ps1={"powershell","-file"}
