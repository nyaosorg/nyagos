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

if nyagos.env.OS == "Windows_NT" then
    nyagos.alias.ls="__ls__ -oFh $*"
    nyagos.alias.ll="__ls__ -olFh $*"
    nyagos.alias.wordpad=
        '"' .. nyagos.pathjoin( nyagos.env.PROGRAMFILES ,
        'Windows NT\\Accessories\\wordpad.exe' ) .. '"'
end
