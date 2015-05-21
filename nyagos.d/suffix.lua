nyagos.suffixes={}

function suffix(suffix,cmdline)
    local suffix=string.lower(suffix)
    if string.sub(suffix,1,1)=='.' then
        suffix = string.sub(suffix,2)
    end
    if not nyagos.suffixes[suffix] then
        local orgpathext = nyagos.getenv("PATHEXT")
        local newext="."..suffix
        if not string.find(";"..orgpathext..";",";"..newext..";",1,true) then
            nyagos.setenv("PATHEXT",orgpathext..";"..newext)
        end
    end
    nyagos.suffixes[suffix]=cmdline
end

local org_filter=nyagos.argsfilter
nyagos.argsfilter = function(args)
    if org_filter then
        local args_ = org_filter(args)
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
        newargs[#newargs+1] = args[i]
    end
    return newargs
end

alias{
    suffix=function(args)
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
}

suffix(".pl",{"perl"})
suffix(".py",{"ipy"})
suffix(".rb",{"ruby"})
suffix(".lua",{"lua"})
suffix(".awk",{"awk","-f"})
suffix(".js",{"cscript","//nologo"})
suffix(".vbs",{"cscript","//nologo"})
suffix(".ps1",{"powershell","-file"})
