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

share.org_dollar_complete = nyagos.completion_hook
nyagos.completion_hook = function(c)
    local text = string.gsub(c.word,"$(%w+)",function(m)
        return nyagos.env[m]
    end)
    text = string.gsub(text,"$%b{}",function(m)
        return nyagos.env[string.sub(m,3,string.len(m)-1)]
    end)
    if text == c.word then
        if share.org_dollar_complete then
            return share.org_dollar_complete(c)
        else
            return nil
        end
    end
    local pattern = text.."*"
    local result = nyagos.glob(pattern)
    if #result < 1 or (#result == 1 and result[1] == pattern )then
        return nil
    end
    for i =1,#result do
        if string.len(result[i]) >= string.len(text)+1 then
            result[i] = c.word .. string.sub(result[i],string.len(text)+1)
        end
    end
    return result
end
