local orgfilter = nyagos.filter
nyagos.filter = function(cmdline)
    if orgfilter then
        local cmdline_ = orgfilter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
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
    return cmdline
end
