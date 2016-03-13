nyagos.bindkey("C-o",function(this)
    local word = this:lastword()
    local word_noquote = string.gsub(word,[['"']],"")
    local result=nyagos.eval('ls -1 -a -d "'..word_noquote..'*" | peco')
    this:call("CLEAR_SCREEN")
    if string.find(result," ",1,true) then
        result = '"'..result..'"'
    end
    for i=1,utf8.len(word) do
        this:call("BACKWARD_DELETE_CHAR")
    end
    return result
end)

