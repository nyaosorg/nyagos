nyagos.alias.ls='ls -oF $*'
nyagos.alias.lua_e=function(args) assert(load(args[1]))() end
nyagos.alias.lua_f=function(args)
    local path=table.remove(args,1)
    assert(loadfile(path))(args)
end
nyagos.alias["for"]='%COMSPEC% /c "@set PROMPT=$G & @for $*"'
