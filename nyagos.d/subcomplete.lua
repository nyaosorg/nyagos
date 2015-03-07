local maincmds = {}

local githelp=io.popen("git help -a","r")
if githelp then
    local gitcmds={}
    for line in githelp:lines() do
        if string.match(line,"^  %S") then
            for word in string.gmatch(line,"%S+") do
                gitcmds[ #gitcmds+1 ] = word
            end
        end
    end
    githelp:close()
    if next(gitcmds) then
        maincmds["git"] = gitcmds
    end
end

if next(maincmds) then
    nyagos.bindkey("C_I",function(this)
        local firstword = this:firstword()
        for maincmd1,subcmds in pairs(maincmds) do
            if maincmd1 == firstword then
                local lastword,lastword_at = this:lastword()
                if lastword_at == string.len(maincmd1)+2 then
                    local list={}
                    local lastword_len = string.len(lastword)
                    for i=1,#subcmds do
                        if string.sub(subcmds[i],1,lastword_len) == lastword then
                            list[#list+1] = subcmds[i]
                        end
                    end
                    local commonprefix = nyagos.commonprefix(list)
                    if string.len(commonprefix) > lastword_len then
                        local result = string.sub(commonprefix,lastword_len+1)
                        if #list == 1 then
                            result = result .. " "
                        end
                        return result
                    elseif #list >= 2 then
                        this:boxprint(list)
                    end
                    return nil
                end
            end
        end
        return this:call("COMPLETE")
    end)
end
