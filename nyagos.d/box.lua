nyagos.alias.dump_temp_out = function()
    for _,val in ipairs(share.dump_temp_out) do
        nyagos.write(val,"\n")
    end
end

nyagos.bindkey("C-o",function(this)
    local word,pos = this:lastword()
    word = string.gsub(word,'"','')
    local wildcard = word.."*"
    local list = nyagos.glob(wildcard)
    if #list == 1 and list[1] == wildcard then
        return
    end
    share.dump_temp_out = list
    nyagos.write("\n")
    local result=nyagos.eval('dump_temp_out | box')
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result," ",1,true) then
        result = '"'..result..'"'
    end
    assert( this:replacefrom(pos,result) )
end)

nyagos.alias.__dump_history = function()
    local uniq={}
    for i=nyagos.gethistory()-1,1,-1 do
        local line = nyagos.gethistory(i)
        if line ~= "" and not uniq[line] then
            nyagos.write(line,"\n")
            uniq[line] = true
        end
    end
end

nyagos.bindkey("C_R", function(this)
    nyagos.write("\n")
    local result = nyagos.eval('__dump_history | box')
    this:call("REPAINT_ON_NEWLINE")
    return result
end)

nyagos.bindkey("M_H" , function(this)
    nyagos.write("\n")
    local result = nyagos.eval('cd --history | box')
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end)

nyagos.bindkey("M_G" , function(this)
    nyagos.write("\n")
    local result = nyagos.eval('git log --pretty="format:%h %s" | box')
    this:call("REPAINT_ON_NEWLINE")
    return string.match(result,"^%S+") or ""
end)
