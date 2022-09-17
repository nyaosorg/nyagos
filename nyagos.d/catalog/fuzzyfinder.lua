if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

-- default/non-configured setting: Use fzf(exists)/peco
local ff = {}
if nyagos.which("fzf.exe") then
  ff.cmd        =  "fzf.exe"
else
  ff.cmd        =  "peco.exe"
end
ff.args         =  {}
ff.args.dir     =  ""
ff.args.cmdhist =  ""
ff.args.cdhist  =  ""
if nyagos.which("fzf.exe") then
  ff.args.gitlog  =  "--preview='git show {1}'"
else
  ff.args.gitlog  =  ""
end

share.fuzzyfinder = share.fuzzyfinder or ff

-- Compatibility settings with old settings
share.selector = share.fuzzyfinder.cmd

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
