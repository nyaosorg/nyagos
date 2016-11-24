share.org_brace_filter = nyagos.filter
nyagos.filter = function(cmdline)
    if share.org_brace_filter then
        local cmdline_ = share.org_brace_filter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    local save={}
    local masking = function(s)
        local i=#save+1
        save[i] = s
        return '\a('..i..')'
    end
    cmdline = cmdline:gsub('"[^"]*"', masking)
    cmdline = cmdline:gsub("'[^']*'", masking)
    repeat
        local last = true
        cmdline = cmdline:gsub("(%S*)(%b{})(%S*)", function(left,mid,right)
            local contents = string.sub(mid,2,-2)
            local result
            local init = 1
            while true do
                local comma=string.find(contents,",",init,true)
                local value=left .. string.sub(contents,init,(comma or 0)-1) .. right
                if result then
                    result = result .. " " .. value
                else
                    result = value
                end
                if not comma then
                    break
                end
                init = comma + 1
            end
            if init > 1 then
                last = false
                return result
            else
                return nil
            end
        end)
    until last
    cmdline = cmdline:gsub("\a%((%d+)%)",function(s)
        return save[s+0]
    end)
    return cmdline
end
