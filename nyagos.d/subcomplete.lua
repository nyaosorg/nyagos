local maincmds = {}

local githelp=io.popen("git help -a 2>nul","r")
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
    if #gitcmds > 1 then
        maincmds["git"] = gitcmds
    end
end
local svnhelp=io.popen("svn help 2>nul","r")
if svnhelp then
    local svncmds={}
    for line in svnhelp:lines() do
        local m=string.match(line,"^   ([a-z]+)")
        if m then
            svncmds[ #svncmds+1 ] = m
        end
    end
    svnhelp:close()
    if #svncmds > 1 then
        maincmds["svn"] = svncmds
    end
end

local hghelp=io.popen("hg debugcomplete 2>nul","r")
if hghelp then
    local hgcmds={}
    for line in hghelp:lines() do
        for word in string.gmatch(line,"[a-z]+") do
            hgcmds[#hgcmds+1] = word
        end
    end
    hghelp:close()
    if #hgcmds > 1 then
        maincmds["hg"] = hgcmds
    end
end

if next(maincmds) then
    nyagos.bindkey("C_I",function(this)
        local maincmd1 = this:firstword()
        local subcmds = maincmds[maincmd1]
        if not subcmds then
            return this:call("COMPLETE")
        end
        local lastword,lastword_at = this:lastword()
        if lastword_at ~= string.len(maincmd1)+2 then
            return this:call("COMPLETE")
        end
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
    end)
end
