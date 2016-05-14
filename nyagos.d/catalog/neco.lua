nyagos.alias.neco_temp_out = function()
    for _,val in ipairs(share.neco_temp_out) do
        nyagos.write(val,"\n")
    end
end

nyagos.bindkey("C-o",function(this)
    local word,pos = this:lastword()
    word = string.gsub(word,'"','')
    share.neco_temp_out = nyagos.glob(word.."*")
    local result=nyagos.eval('neco_temp_out | neco')
    if string.find(result," ",1,true) then
        result = '"'..result..'"'
    end
    assert( this:replacefrom(pos,result) )
end)

nyagos.bindkey("C_R", function(this)
    local path = nyagos.pathjoin(nyagos.env.appdata,'NYAOS_ORG\\nyagos.history')
    local stat = nyagos.stat(path)
    if stat and stat.size > 0 then
        return nyagos.eval('neco < '..path)
    end
end)

nyagos.bindkey("M_H" , function(this)
    local result = nyagos.eval('cd --history | neco')
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end)

nyagos.bindkey("M_G" , function(this)
    local result = nyagos.eval('git log --pretty="format:%h %s" | neco')
    return string.match(result,"^%S+") or ""
end)
