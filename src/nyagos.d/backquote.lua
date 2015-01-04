nyagos.filter = function(cmdline)
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
