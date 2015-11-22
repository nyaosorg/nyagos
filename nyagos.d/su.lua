nyagos.alias.sudo = function(args)
    if #args <= 0 then
        nyagos.shellexecute("runas",nyagos.exe)
        return
    end
    local prog = args[1]
    table.remove(args,1)
    assert(nyagos.shellexecute("runas",prog,table.concat(args," "),nyagos.getwd()))
end

share._clone = function(action)
    local status,err = nyagos.shellexecute(action,nyagos.exe)
    if not status and string.match(err,"^Error%(5%)") then
	status,err = nyagos.shellexecute(action,nyagos.getenv("COMSPEC"),'/c "'..nyagos.exe..'"')
    end
    return status,err
end

nyagos.alias.su = function()
    assert(share._clone("runas"))
end
nyagos.alias.clone = function()
    assert(share._clone("open"))
end
