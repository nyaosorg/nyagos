nyagos.alias.sudo = function(args)
    if #args <= 0 then
        nyagos.shellexecute("runas",nyagos.exe)
        return
    end
    local prog = args[1]
    table.remove(args,1)
    local cwd = nyagos.netdrivetounc(nyagos.getwd())
    assert(nyagos.shellexecute("runas",prog,table.concat(args," "),cwd))
end

share._clone = function(action)
    local cwd = nyagos.netdrivetounc(nyagos.getwd())
    print(cwd)
    local status,err = nyagos.shellexecute(action,nyagos.exe,"",cwd)
    if status then
        return status,err
    end
    if string.match(err,"^Error%(5%)") or string.match(err,"winapi error") then
	status,err = nyagos.shellexecute(action,nyagos.getenv("COMSPEC"),'/c "'..nyagos.exe,"",cwd)
    end
    return status,err
end

nyagos.alias.su = function()
    assert(share._clone("runas"))
end
nyagos.alias.clone = function()
    assert(share._clone("open"))
end
