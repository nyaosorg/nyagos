if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

if share.maincmds then
    return
end

share.maincmds = {}

local function get_local()
    local dir=nyagos.pathjoin(nyagos.env.LOCALAPPDATA,"NYAOS_ORG")
    local stat=nyagos.stat(dir)
    if not stat then
        nyagos.exec('mkdir "'..dir..'"')
    end
    return dir
end

local function load_subcommands_cache(fname)
    fname = nyagos.pathjoin(get_local(),fname)
    local fd=io.open(fname,"r")
    if not fd then
        return nil
    end
    local list = {}
    for line in fd:lines() do
        list[#list+1] = line
    end
    fd:close()
    return list
end

local function save_subcommands_cache(fname,list)
    fname = nyagos.pathjoin(get_local(),fname)
    local fd=io.open(fname,"w")
    if not fd then
        return
    end
    for i=1,#list do
        fd:write(list[i].."\n")
    end
    fd:close()
end

-- git

local function update_cache()
    share.maincmds["git"] = load_subcommands_cache("git-subcommands.txt")
    share.maincmds["hub"] = load_subcommands_cache("hub-subcommands.txt")
    if not share.maincmds["git"] then
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
                save_subcommands_cache("git-subcommands.txt",gitcmds)
                if hub then
                  maincmds["hub"] = gitcmds
                  save_subcommands_cache("hub-subcommands.txt",gitcmds)
                end
                share.maincmds = maincmds
            end
        end
    end

    -- Subversion
    share.maincmds["svn"] = load_subcommands_cache("svn-subcommands.txt")
    if not share.maincmds["svn"] then
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
                save_subcommands_cache("svn-subcommands.txt",svncmds)
            end
        end
    end

    -- Mercurial
    share.maincmds["hg"] = load_subcommands_cache("hg-subcommands.txt")
    if not share.maincmds["hg"] then
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
                save_subcommands_cache("hg-subcommands.txt",hgcmds)
            end
        end
    end

    -- Rclone

    share.maincmds["rclone"] = load_subcommands_cache("rclone-subcommands.txt")
    if not share.maincmds["rclone"] then
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
              save_subcommands_cache("rclone-subcommands.txt",rclonecmds)
          end
        end
    end

    share.maincmds["fsutil"] = load_subcommands_cache("fsutil-subcommands.txt")
    if not share.maincmds["fsutil"] then
        local fd=io.popen("fsutil","r")
        if fd then
            local list = {}
            for line in fd:lines() do
                local m=string.match(line,"^(%w+)%s+%S+")
                if m then
                    list[#list+1] = m
                end
            end
            fd:close()
            if #list >= 1 then
                share.maincmds["fsutil"] = list
                save_subcommands_cache("fsutil-subcommands.txt",list)
            end
        end
    end

    share.maincmds["go"] = {
        "bug", "build", "clean", "doc", "env", "fix",
        "fmt", "generate", "get", "install", "list",
        "mod", "run", "test", "tool", "version", "vet"
    }

    for cmd,subcmdData in pairs(share.maincmds or {}) do
        if not nyagos.complete_for[cmd] then
            nyagos.complete_for[cmd] = function(args)
                local subcmdType = type(subcmdData)
                if "table" == subcmdType then
                    while #args > 2 and args[2]:sub(1,1) == "-" do
                        table.remove(args,2)
                    end
                    if #args == 2 then
                        return subcmdData
                    end
                elseif "function" == subcmdType then
                    return subcmdData(args)
                end
                return nil
            end
        end
    end
end

update_cache()

nyagos.alias.clear_subcommands_cache = function()
    local wildcard = nyagos.pathjoin(get_local(),"*-subcommands.txt")
    local files = nyagos.glob(wildcard)
    if #files >= 2 or not string.find(files[1],"*",1,true) then
        for i=1,#files do
            print("remove "..files[i])
            os.remove(files[i])
        end
    end
    update_cache()
end
