alias{
    ls='ls -oF $*',
    echo=function(args) nyagos.write(table.concat(args,' ')..'\n') end,
    lua_e=function(args) assert(load(args[1]))() end
}
