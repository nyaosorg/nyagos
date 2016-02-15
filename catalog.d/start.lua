nyagos.alias.start = function(args)
    local dir=""
    if args[1] == '/D' or args[1] == '/d' then
        dir = args[2]
        table.remove(args,1)
        table.remove(args,1)
    end
    local progname = args[1] -- nyagos.pathjoin(nyagos.getwd(),args[1])
    local param = ""
    if #args > 2 then
        param = '"' .. table.concat(args,'" "',2) .. '"'
    end
    print("open "..progname.." "..param.." /D "..dir)
    assert(nyagos.shellexecute("open",progname,param,dir))
end
