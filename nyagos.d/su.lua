nyagos.alias("sudo",function(args)
    if #args <= 0 then
        nyagos.shellexecute("runas",nyagos.exe)
        return
    end
    local prog = args[1]
    table.remove(args,1)
    assert(nyagos.shellexecute("runas",prog,table.concat(args," "),nyagos.getwd()))
end)

nyagos.alias("su",function()
    assert(nyagos.shellexecute("runas",nyagos.exe))
end)
nyagos.alias("clone",function()
    assert(nyagos.shellexecute("open",nyagos.exe))
end)
