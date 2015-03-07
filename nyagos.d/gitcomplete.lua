local githelp=io.popen("git help -a","r")
local gitcmds=nil
if githelp then
    gitcmds={}
    for line in githelp:lines() do
        if string.match(line,"^  %S") then
            for word in string.gmatch(line,"%S+") do
                gitcmds[ #gitcmds+1 ] = word
            end
        end
    end
    githelp:close()

    nyagos.bindkey("C_I",function(this)
        if this:firstword() ~= "git" then
            this:call("COMPLETE")
            return
        end
        local lastword,lastword_at = this:lastword()
        if lastword_at ~= 5 then
            this:call("COMPLETE")
            return
        end
        local list={}
        local lastword_len = string.len(lastword)
        for i=1,#gitcmds do
            if string.sub(gitcmds[i],1,lastword_len) == lastword then
                list[#list+1] = gitcmds[i]
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
    end)
end
