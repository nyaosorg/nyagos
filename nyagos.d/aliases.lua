alias{
    ls='ls -oF $*',
    lua_e=function(args) assert(load(args[1]))() end
}
