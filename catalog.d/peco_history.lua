nyagos.bindkey("C_R", function(this)
    local path = nyagos.pathjoin(nyagos.env.appdata,'NYAOS_ORG\\nyagos.history')
    local result = nyagos.eval('peco < '..path)
    this:call("CLEAR_SCREEN")
    return result
end)

