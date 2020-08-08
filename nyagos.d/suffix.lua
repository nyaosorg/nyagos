if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

share._suffixes={}

share._setsuffix = function(suffix,cmdline)
    suffix=string.gsub(string.lower(suffix),"^%.","")
    if not share._suffixes[suffix] then
        local newext="."..suffix
        local orgpathext = nyagos.env.PATHEXT
        if orgpathext then
            if not string.find(";"..orgpathext..";",";"..newext..";",1,true) then
                nyagos.env.PATHEXT = orgpathext..";"..newext
            end
        else
            nyagos.env.PATHEXT = newext
        end
    end
    share._suffixes[suffix]=cmdline
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
    local m = string.match(args[0],"%.(%w+)$")
    if not m then
        return
    end
    local cmdline = share._suffixes[ string.lower(m) ]
    if not cmdline then
        return
    end
    local path=nyagos.which(args[0])
    if not path then
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
            local val = share._suffixes[args[i]]
            if not val then
                val = ""
            elseif type(val) == "table" then
                val = table.concat(val," ")
            end
            print(args[i].."="..val)
        end
    end
end

for key,val in pairs{
    awk={"gawk","-f"},
    js={"cscript","//nologo"},
    lua={"nyagos.exe","--norc","--lua-file"},
    pl={"perl"},
    ps1={"powershell","-ExecutionPolicy","RemoteSigned","-file"},
    rb={"ruby"},
    vbs={"cscript","//nologo"},
    wsf={"cscript","//nologo"},
    py={"python"},
} do
    share._setsuffix( key , val )
end
