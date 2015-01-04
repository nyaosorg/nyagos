nyagos.alias("sudo",function(args)
    if #args <= 0 then
        nyagos.shellexecute("runas",nyagos.exe)
        return
    end
    local prog = args[1]
    table.remove(args,1)
    assert(nyagos.shellexecute("runas",prog,table.concat(args," "),nyagos.getwd()))
end)

nyagos.alias("su",function() nyagos.shellexecute("runas",nyagos.exe) end)
nyagos.alias("clone",function() nyagos.shellexecute("open",nyagos.exe) end)
