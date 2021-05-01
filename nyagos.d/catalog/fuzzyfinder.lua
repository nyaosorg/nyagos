if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

share.fuzzyfinder         = share.fuzzyfinder or {}
share.fuzzyfinder.cmd     = share.fuzzyfinder.cmd   or "fzf.exe"
share.fuzzyfinder.args    = share.fuzzyfinder.args  or {}
share.fuzzyfinder.args.dir     = share.fuzzyfinder.args.dir     or "--preview='dir {1}'"
share.fuzzyfinder.args.cmdhist = share.fuzzyfinder.args.cmdhist or ""
share.fuzzyfinder.args.cdhist  = share.fuzzyfinder.args.cdhist  or "--preview='dir {1}'"
share.fuzzyfinder.args.gitlog  = share.fuzzyfinder.args.gitlog  or "--preview='git show {1} | cat'"

nyagos.alias.dump_temp_out = function()
    for _,val in ipairs(share.dump_temp_out) do
        nyagos.write(val,"\n")
    end
end

nyagos.bindkey("C-o",function(this)
    local word,pos = this:lastword()
    word = string.gsub(word,'"','')
    share.dump_temp_out = nyagos.glob(word.."*")
    local result=nyagos.eval('dump_temp_out | ' .. share.fuzzyfinder.cmd .. " " .. share.fuzzyfinder.args.dir)
    this:call("CLEAR_SCREEN")
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
    local result = nyagos.eval('__dump_history | ' .. share.fuzzyfinder.cmd .. " " .. share.fuzzyfinder.args.cmdhist)
    this:call("CLEAR_SCREEN")
    return result
end)

nyagos.bindkey("M_H" , function(this)
    local result = nyagos.eval('cd --history | ' .. share.fuzzyfinder.cmd .. " " .. share.fuzzyfinder.args.cdhist)
    this:call("CLEAR_SCREEN")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end)

nyagos.bindkey("M_G" , function(this)
    local result = nyagos.eval('git log --pretty="format:%h %s" | ' .. share.fuzzyfinder.cmd .. " " .. share.fuzzyfinder.args.gitlog)
    this:call("CLEAR_SCREEN")
    return string.match(result,"^%S+") or ""
end)
