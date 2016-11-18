nyagos.alias.open = function(args)
    local count=0
    for i=1,#args do
        local list=nyagos.glob(args[i])
        if list and #list >= 1 then
            for j=1,#list do
                if nyagos.access(list[j],0) then
                    if #list >= 2 then
                        nyagos.write("open: "..list[j].." ? [y/n/q] ")
                        local ch = nyagos.getkey()
                        ch = string.lower(utf8.char(ch))
                        if string.match(ch,"^%g+$") then
                            nyagos.write(ch)
                        end
                        nyagos.write("\n")
                        if ch == 'q' then
                            break
                        elseif ch ~= 'y' then
                            goto continue
                        end
                    end
                    assert(nyagos.shellexecute("open",list[j]))
                else
                    print(args[i] .. ": can not get status")
                end
                ::continue::
            end
        else
            if nyagos.access(args[i],0) then
                assert(nyagos.shellexecute("open",args[i]))
            else
                print(args[i] .. ": can not get status")
            end
        end
        count = count +1
    end
    if count <= 0 then
        if nyagos.access(".\\open.cmd",0) then
            nyagos.exec("open.cmd")
        else
            assert(nyagos.shellexecute("open","."))
        end
    end
end
