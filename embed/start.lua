if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

nyagos.alias.start = function(args)
    -- Remove title-parameter --
    if not args[1] then
        print('start ["title"] [/D directory] PROGNAME ARGS...')
        return
    end
    if string.sub(args.rawargs[1],1,1) == '"' then
        table.remove(args.rawargs,1)
        table.remove(args,1)
    end
    local dir=""
    if args[1] == '/D' or args[1] == '/d' then
        if not args[2] then
            print('start ["title"] [/D directory] PROGNAME ARGS...')
            return
        end
        dir = args[2]
        table.remove(args.rawargs,1)
        table.remove(args.rawargs,1)
        table.remove(args,1)
        table.remove(args,1)
    end
    if not args[1] then
        print('start ["title"] [/D directory] PROGNAME ARGS...')
        return
    end
    local progname = ((nyagos.which(args[1])) or args[1])
    local param = ""
    if #args >= 2 then
        param = table.concat(args.rawargs," ",2)
    end
    if string.match(progname,"^https?:") then
        os.execute("rundll32.exe url.dll,FileProtocolHandler " .. progname)
    else
        assert(nyagos.shellexecute("open",progname,param,dir))
    end
end
