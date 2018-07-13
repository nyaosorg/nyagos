if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

nyagos.alias.lua_e=function(args)
    if #args >= 1 then
        local f,err =loadstring(args[1])
        if f then
            f()
        else
            io.stderr:write(err,"\n")
        end
    end
end
nyagos.alias.lua_f=function(args)
    local save=_G["arg"]
    local script = args[1]
    local param = {}
    for i=1,#args do
        param[i-1] = args[i]
    end
    local f, err = loadfile(script)
    _G["arg"] = param
    if f then
        f()
    else
        io.stderr:write(err,"\n")
    end
    _G["arg"] = save
end

nyagos.alias["for"]=function(args)
    local batchpathu = nyagos.env.temp .. os.tmpname() .. ".cmd"
    local batchpatha = nyagos.utoa(batchpathu)
    local fd,fd_err = nyagos.open(batchpathu,"w")
    if not fd then
        nyagos.writerr(fd_err.."\n")
        return
    end
    local cmdline = "@for "..table.concat(args.rawargs," ").."\n"
    fd:write("@set prompt=$G\n")
    fd:write(cmdline)
    fd:close()
    nyagos.rawexec(nyagos.env.comspec,"/c",batchpathu)
    os.remove(batchpatha)
end
nyagos.alias.kill = function(args)
    local command="taskkill.exe"
    for i=1,#args do
        if args[i] == "-f" then
            command="taskkill.exe /F"
        else
            nyagos.exec(command .. " /PID " .. args[i])
        end
    end
end
nyagos.alias.killall = function(args)
    local command="taskkill.exe"
    for i=1,#args do
        if args[i] == "-f" then
            command="taskkill.exe /F"
        else
            nyagos.exec(command .. " /IM " .. args[i])
        end
    end
end

-- on chcp, font-width is changed.
nyagos.alias.chcp = function(args)
    nyagos.resetcharwidth()
    nyagos.rawexec(args[0],(table.unpack or unpack)(args))
end

-- wildcard expand for external commands
nyagos.alias.wildcard = function(args)
    local newargs = {}
    for i=1,#args do
        local tmp = nyagos.glob(args[i])
        for j=1,#tmp do
            newargs[ #newargs+1] = tmp[j]
        end
    end
    -- for i=1,#newargs do print(newargs[i]) end
    nyagos.exec( newargs )
end

-- print the absolute path
function nyagos.alias.abspath(args)
    local cwd = nyagos.getwd()
    for i = 1,#args do
        local path1 = nyagos.pathjoin(cwd,args[i])
        nyagos.write(path1)
        nyagos.write("\n")
    end
end

-- chompf - cut last CRLF of the files and output them to STDOUT
function nyagos.alias.chompf(args)
    local lf=""
    if #args <= 0 then
        for line in nyagos.lines() do
            nyagos.write(lf)
            nyagos.write(line)
            lf = "\n"
        end
    else
        for i=1,#args do
            for line in nyagos.lines(args[i]) do
                nyagos.write(lf)
                nyagos.write(line)
                lf = "\n"
            end
        end
    end
end
