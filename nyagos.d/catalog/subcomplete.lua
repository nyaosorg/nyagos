if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

share.maincmds = {}

-- git
local githelp=io.popen("git help -a 2>nul","r")
local hubhelp=io.popen("hub help -a 2>nul","r")
if githelp then
    local gitcmds={}
    local hub=false
    if hubhelp then
      hub=true
      local startflag = false
      local found=false
      for line in hubhelp:lines() do
        if not found then
          if startflag then
            -- skip blank line
            if string.match(line,"%S") then
              -- found commands
              for word in string.gmatch(line, "%S+") do
                gitcmds[ #gitcmds+1 ] = word
              end
              found = true
            end
          end
          if string.match(line,"hub custom") then
            startflag = true
          end
        end
      end
      hubhelp:close()
    end
    for line in githelp:lines() do
        local word = string.match(line,"^ +(%S+)")
        if nil ~= word then
          gitcmds[ #gitcmds+1 ] = word
        end
    end
    githelp:close()
    if #gitcmds > 1 then
        local maincmds = share.maincmds
        maincmds["git"] = gitcmds
        if hub then
          maincmds["hub"] = gitcmds
        end
        share.maincmds = maincmds
    end
end

-- Subversion
local svnhelp=nyagos.eval("svn help 2>nul","r")
if string.len(svnhelp) > 5 then
    local svncmds={}
    for line in string.gmatch(svnhelp,"[^\n]+") do
        local m=string.match(line,"^ +([a-z]+)")
        if m then
            svncmds[ #svncmds+1 ] = m
        end
    end
    if #svncmds > 1 then
        local maincmds = share.maincmds
        maincmds["svn"] = svncmds
        share.maincmds = maincmds
    end
end

-- Mercurial
local hghelp=nyagos.eval("hg debugcomplete 2>nul","r")
if string.len(hghelp) > 5 then
    local hgcmds={}
    for line in string.gmatch(hghelp,"[^\n]+") do
        for word in string.gmatch(line,"[a-z]+") do
            hgcmds[#hgcmds+1] = word
        end
    end
    if #hgcmds > 1 then
        local maincmds=share.maincmds
        maincmds["hg"] = hgcmds
        share.maincmds = maincmds
    end
end

-- Rclone
local rclonehelp=io.popen("rclone --help 2>nul","r")
if rclonehelp then
  local rclonecmds={}
  local startflag = false
  local flagsfound=false
  for line in rclonehelp:lines() do
      if not flagsfound then
          if string.match(line,"Available Commands:") then
            startflag = true
          end
          if string.match(line,"Flags:") then
              flagsfound = true
          end
          if startflag then
            local m=string.match(line,"^%s+([a-z]+)")
            if m then
                rclonecmds[ #rclonecmds+1 ] = m
            end
          end
      end
  end
  rclonehelp:close()
  if #rclonecmds > 1 then
      local maincmds=share.maincmds
      maincmds["rclone"] = rclonecmds
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
        cmdname = string.lower(string.gsub(cmdname,"%.%w+$",""))
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
        local subcmds = {}
        local subcmdData = share.maincmds[cmdname]
        if not subcmdData then
            return nil
        end
        local subcmdType = type(subcmdData)
        if "table" == subcmdType then
          subcmds = subcmdData
        elseif "function" == subcmdType then
          subcmds = subcmdData()
        end
        for i=1,#subcmds do
            table.insert(c.list,subcmds[i])
        end
        return c.list
    end
end
