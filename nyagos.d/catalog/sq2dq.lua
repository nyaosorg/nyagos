local function sq2dq(source)
    local sq = false
    local dq = false
    return (string.gsub(source,"[\"']",function(c)
        if c == "'" then
            if not dq then
                sq = not sq
                return '"'
            end
        end
        if c == '"' then
            if sq then
                return [[\"]]
            end
            dq = not dq
        end
        return c
    end))
end

if nyagos then
    local oldfilter = nyagos.filter
    nyagos.filter = function(s)
        s = sq2dq(s)
        if oldfilter then
            s = oldfilter(s)
        end
        return s
    end
else
    assert(sq2dq([['(print "ahaha")']])   == [["(print \"ahaha\")"]])
    assert(sq2dq([["(print \"ahaha\")"]]) == [["(print \"ahaha\")"]])
    assert(sq2dq([["(print 'ahaha')"]])   == [["(print 'ahaha')"]]  )
end
