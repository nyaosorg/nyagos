nyagos.alias.ls='ls -oF $*'
nyagos.alias.lua_e=function(args) assert(load(args[1]))() end
nyagos.alias["for"]='%COMSPEC% /c "@set PROMPT=$G & @for $*"'
