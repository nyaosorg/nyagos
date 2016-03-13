share.org_dollar_filter = nyagos.filter
nyagos.filter = function(cmdline)
    if share.org_dollar_filter then
        local cmdline_ = share.org_dollar_filter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    return cmdline:gsub("%$(%w+)",function(m)
        return nyagos.env[m]
    end):gsub("$%b{}", function(m)
        if string.len(m) >= 3 then
            return nyagos.env[string.sub(m,3,-2)]
        else
            return nil
        end
    end)
end
