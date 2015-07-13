local orgfilter = nyagos.filter
nyagos.filter = function(cmdline)
    if orgfilter then
        local cmdline_ = orgfilter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    local save={}
    cmdline = cmdline:gsub("((['\"])[^%2]*%2)", function(s,_)
        local i=#save+1
        save[i] = s
        return '\a('..i..')'
    end)
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
