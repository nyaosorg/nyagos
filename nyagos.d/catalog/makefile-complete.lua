if not nyagos then
    os.exit(0)
end

nyagos.complete_for["make"] = function(args)
    if #args >= 2 and args[#args-1] == "-f" then
        return nil
    end
    if #args >= 1 and string.find(args[#args],"=",1,true) then
        return nil
    end
    local target
    local fd=io.open("Makefile","r")
    if fd then
        target = {}
        for line in fd:lines() do
            local m =string.match(line,'^([^:\t%s]+)%s*:[^=]')
            if m then
                target[1+#target] = m
            end
        end
        fd:close()
    end
    return target
end

nyagos.complete_for["nmake"] = nyagos.complete_for["make"]

