local orgfilter = nyagos.filter
nyagos.filter = function(cmdline)
    if orgfilter then
        cmdline = orgfilter(cmdline)
    end
    local last
    repeat
        last = true
        cmdline = cmdline:gsub("(%S*)(%b{})(%S*)", function(left,mid,right)
            last = false
            local contents = string.sub(mid,2,-2)
            local result
            local is_replace = false
            for s in string.gmatch(contents,"[^,]+") do
                local value=left .. s .. right
                if result then
                    result = result .. " " .. value
                    is_replace = true
                else
                    result = value
                end
            end
            if is_replace then
                return result
            else
                return left .. "%u+007B%" .. contents .. "%u+007D%" ..right
            end
        end)
    until last
    return cmdline:gsub("%%u%+007B%%","{"):gsub("%%u%+007D%%","}")
end
