nyagos.org_backquote_filter = nyagos.filter
nyagos.filter = function(cmdline)
    if nyagos.org_backquote_filter then
        local cmdline_ = nyagos.org_backquote_filter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    return string.gsub(cmdline,'`([^`]*)`',function(m)
        local r = nyagos.eval(m)
        if not r then
            return false
        end
        r = nyagos.atou(r)
        r = string.gsub(r,'[|&<>!]',function(m)
            return string.format('%%u+%04X%%',string.byte(m,1,1))
        end)
        return string.gsub(r,'%s+$','')
    end)
end
