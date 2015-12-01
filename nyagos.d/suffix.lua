share._suffixes={}

share._setsuffix = function(suffix,cmdline)
    local suffix=string.lower(suffix)
    if string.sub(suffix,1,1)=='.' then
        suffix = string.sub(suffix,2)
    end
    if not share._suffixes[suffix] then
        local orgpathext = nyagos.getenv("PATHEXT")
        local newext="."..suffix
        if not string.find(";"..orgpathext..";",";"..newext..";",1,true) then
            nyagos.setenv("PATHEXT",orgpathext..";"..newext)
        end
    end
    local table = share._suffixes
    table[suffix]=cmdline
    share._suffixes = table
end

suffix = setmetatable({},{
    __call = function(t,k,v) share._setsuffix(k,v) return end,
    __newindex = function(t,k,v) share._setsuffix(k,v) return end,
    __index = function(t,k) return share._suffixes[k] end
})

share._org_suffix_argsfilter=nyagos.argsfilter
nyagos.argsfilter = function(args)
    if share._org_suffix_argsfilter then
        local args_ = share._org_suffix_argsfilter(args)
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
    local cmdline = share._suffixes[ string.lower(m) ]
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
suffix.wsf={"cscript","//nologo"}
suffix.ps1={"powershell","-file"}
