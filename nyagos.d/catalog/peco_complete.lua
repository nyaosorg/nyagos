nyagos.alias.peco_temp_out = function()
    for _,val in ipairs(share.peco_temp_out) do
        nyagos.write(val,"\n")
    end
end

nyagos.bindkey("C-o",function(this)
    local word = this:lastword()
    share.peco_temp_out = nyagos.glob(word.."*")
    local result=nyagos.eval('peco_temp_out | peco')
    this:call("CLEAR_SCREEN")
    if string.find(result," ",1,true) then
        result = '"'..result..'"'
    end
    for i=1,utf8.len(word) do
        this:call("BACKWARD_DELETE_CHAR")
    end
    return result
end)

