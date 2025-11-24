function getUsage(command)
    print("$ " .. command)
    local subcommand = {}
    local fd = assert(io.popen(command))
    for line in fd:lines() do
        local m = string.match(line,"^  ([a-z][-a-z]+)")
        if m then
            subcommand[m] = {}
        end
    end
    fd:close()
    return subcommand
end

function dump(fd,obj,indent)
    local t = type(obj)
    if t == "string" then
        fd:write('"'..obj..'"')
    elseif t == "number" then
        fd:write(obj)
    elseif t == "table" then
        fd:write("{")
        for key,val in pairs(obj) do
            fd:write("\r\n"..string.rep("    ",indent+1).."[")
            dump(fd,key,indent+1)
            fd:write("]=")
            dump(fd,val,indent+1)
            fd:write(",")
        end
        if next(obj) then
            fd:write("\r\n"..string.rep("    ",indent).."}")
        else
            fd:write("}")
        end
    elseif t == "boolean" then
        if obj then
            fd:write("true")
        else
            fd:write("false")
        end
    else
        fd:write("nil")
    end
end

local jj = getUsage("jj -h")
for name,_ in pairs(jj) do
    if name ~= "help" then
        if string.sub(name,1,1) ~= "-" then
            jj[name] = getUsage("jj ".. name .. " -h")
        end
    end
end

local fd = assert(io.open("complete-jj.lua","w+"))
fd:write("share.jj=")
dump(fd,jj,0)

local script = string.gsub([[

nyagos.complete_for["jj"] = function(args)
    if not string.match(args[#args],"^[-a-z]+") then
        return nil
    end

    local j = share.jj
    local last = nil
    while true do
        repeat
            table.remove(args,1)
            if #args <= 0 then
                return last
            end
            last = args[1]
        until string.sub(last,1,1) ~= "-"

        local nextj = j[ last ]
        if not nextj then
            local result = {}
            for key,val in pairs(j) do
                result[#result+1] = key
            end
            if next(result) then
                return result
            else
                return nil
            end
        end
        j = nextj
    end
end
]],"\r?\n","\r\n")

fd:write(script)
fd:close()
