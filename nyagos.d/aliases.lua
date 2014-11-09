alias{
    assoc='%COMSPEC% /c assoc $*',
    attrib='%COMSPEC% /c attrib $*',
    echo=function(args) nyagos.write(table.concat(args,' ')..'\n') end,
    copy='%COMSPEC% /c copy $*',
    del='%COMSPEC% /c del $*',
    dir='%COMSPEC% /c dir $*',
    mklink='%COMSPEC% /c mklink $*',
    move='%COMSPEC% /c move $*',
    ren='%COMSPEC% /c ren $*',
    rename='%COMSPEC% /c rename $*',
    rem=function() end,
    start='%COMSPEC% /c start $*',
    ls='ls -oF $*',
    ['type']='%COMSPEC% /c type $*',
    ['for']='%COMSPEC% /c for $*',
    lua_e=function(args) assert(load(args[1]))() end,
    which=function(args) 
        local result=nyagos.which(args[1])
        if result then
            nyagos.write( result..'\n')
        end
    end,
    open=function(args)
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
    end,
    sudo=function(args)
        if #args <= 0 then
            nyagos.shellexecute("runas",nyagos.exe)
            return
        end
        local prog = args[1]
        table.remove(args,1)
        assert(nyagos.shellexecute("runas",prog,table.concat(args," "),nyagos.getwd()))
    end,
}
