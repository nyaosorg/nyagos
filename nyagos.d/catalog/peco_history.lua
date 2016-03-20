nyagos.bindkey("C_R", function(this)
    local path = nyagos.pathjoin(nyagos.env.appdata,'NYAOS_ORG\\nyagos.history')
    local result = nyagos.eval('peco < '..path)
    this:call("CLEAR_SCREEN")
    return result
end)

nyagos.bindkey("M_H" , function(this)
    local result = nyagos.eval('cd --history | peco')
    this:call("CLEAR_SCREEN")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end)
