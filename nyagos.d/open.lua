nyagos.alias("open",function(args)
    local count=0
    for i=1,#args do
        local list=nyagos.glob(args[i])
        if list and #list >= 1 then
            for i=1,#list do
                local fd = io.open(list[i])
                if fd then
                    fd:close()
                    assert(nyagos.shellexecute("open",list[i]))
                else
                    nyagos.writerr(list[i]..": not found.\n")
                end
            end
        else
            local fd = io.open(args[i])
            if fd then
                fd:close()
                assert(nyagos.shellexecute("open",args[i]))
            else
                print(args[i] .. ": not found.\n")
            end
        end
        count = count +1
    end
    if count <= 0 then
        local fd = io.open("open.cmd")
        if fd then
            fd:close()
            nyagos.exec("open.cmd")
        else
            assert(nyagos.shellexecute("open","."))
        end
    end
end)
