alias{
    assoc='%COMSPEC% /c assoc $*',
    attrib='%COMSPEC% /c attrib $*',
    echo=function(args) nyagos.write(table.concat(args,' ')..'\n') end,
    copy='%COMSPEC% /c copy $*',
    del='%COMSPEC% /c del $*',
    dir='%COMSPEC% /c dir $*',
    md='%COMSPEC% /c md $*',
    mkdir='%COMSPEC% /c mkdir $*',
    mklink='%COMSPEC% /c mklink $*',
    move='%COMSPEC% /c move $*',
    rd='%COMSPEC% /c rd $*',
    ren='%COMSPEC% /c ren $*',
    rename='%COMSPEC% /c rename $*',
    rem=function() end,
    rmdir='%COMSPEC% /c rmdir $*',
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
                    nyagos.exec(string.format('%s /c start "%s"',nyagos.getenv('COMSPEC'),list[i]))
                end
            else
                nyagos.exec(string.format('%s /c start "%s"',nyagos.getenv('COMSPEC'),args[i]))
            end
            count = count +1
        end
        if count <= 0 then
            nyagos.exec(string.format('%s /c start .',nyagos.getenv('COMSPEC')))
        end
    end,
}
