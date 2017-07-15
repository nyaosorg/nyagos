if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

backquote = {
    org = nyagos.filter,
    replace = function(m)
        m = string.sub(m,2,string.len(m)-1)
        local r = nyagos.eval(m)
        if not r then
            return false
        end
        r = nyagos.atou(r)
        r = string.gsub(r,'[|&<>!]',function(m)
            return string.format('%%u+%04X%%',string.byte(m,1,1))
        end)
        return string.gsub(r,'%s+$','')
    end
}

nyagos.filter = function(cmdline)
    if backquote.org then
        local cmdline_ = backquote.org(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    cmdline = cmdline:gsub('`[^`]*`',backquote.replace)
    cmdline = cmdline:gsub('%$(%b())',backquote.replace)
    return cmdline
end
