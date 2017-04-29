share.maincmds = {}

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
        local maincmds = share.maincmds
        maincmds["git"] = gitcmds
        share.maincmds = maincmds
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
        local maincmds = share.maincmds
        maincmds["svn"] = svncmds
        share.maincmds = maincmds
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
        local maincmds=share.maincmds
        maincmds["hg"] = hgcmds
        share.maincmds = maincmds
    end
end

if next(share.maincmds) then
    nyagos.completion_hook = function(c)
        if c.pos <= 1 then
            return nil
        end
        local cmdname = string.match(c.text,"^%S+")
        if not cmdname then
            return nil
        end
        --[[
          2nd command completion like :git bisect go[od]
          user define-able

          local subcommands={"good", "bad"}
          local maincmds=share.maincmds
          maincmds["git bisect"] = subcommands
          share.maincmds = maincmds
        --]]
        local cmd2nd = string.match(c.text,"^%S+%s+%S+")
        if share.maincmds[cmd2nd] then
          cmdname = cmd2nd
        end
        local subcmds = share.maincmds[cmdname]
        if not subcmds then
            return nil
        end
        for i=1,#subcmds do
            table.insert(c.list,subcmds[i])
        end
        return c.list
    end
end
