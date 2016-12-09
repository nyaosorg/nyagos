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
    if #args < 1 then
        for key,val in pairs(share._suffixes) do
            local right=val
            if type(val) == "table" then
                right = table.concat(val," ")
            end
            print(key .. "=" .. right)
        end
        return
    end
    for i=1,#args do
        local left,right=string.match(args[i],"^%.?([^=]+)%=(.+)$")
        if right then
            local args={}
            for m in string.gmatch(right,"%S+") do
                args[#args+1] = m
            end
            share._setsuffix(left,args)
        else
            print(args[i].."="..(share._suffixes[args[i]] or ""))
        end
    end
end
