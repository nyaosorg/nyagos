nyagos.alias("open",function(args)
    local count=0
    for i=1,#args do
        local list=nyagos.glob(args[i])
        if list and #list >= 1 then
            for i=1,#list do
                nyagos.shellexecute("open",list[i])
            end
        else
            nyagos.shellexecute("open",args[i])
        end
        count = count +1
    end
    if count <= 0 then
        nyagos.shellexecute("open",".")
    end
end)
