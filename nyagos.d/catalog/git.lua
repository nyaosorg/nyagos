if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

-- hub exists, replace git command
local hubpath=nyagos.which("hub.exe")
if hubpath then
  nyagos.alias.git = "hub.exe"
end

share.git = {}

local getcommits = function(args)
    local fd=io.popen("git log --format=\"%h\" -n 20 2>nul","r")
    if not fd then
        return {}
    end
    local result={}
    for line in fd:lines() do
        result[#result+1] = line
    end
    fd:close()
    return result
end

-- setup local branch listup
local branchlist = function(args)
  if string.find(args[#args],"[/\\\\]") then
      return nil
  end
  local gitbranches = getcommits()
  local gitbranch_tmp = nyagos.eval('git for-each-ref  --format="%(refname:short)" refs/heads/ 2> nul')
  for line in gitbranch_tmp:gmatch('[^\n]+') do
    table.insert(gitbranches,line)
  end
  return gitbranches
end

local unquote = function(s)
    s = string.gsub(s,'"','')
    return string.gsub(s,'\\[0-7][0-7][0-7]',function(t)
        return string.char(tonumber(string.sub(t,2),8))
    end)
end

local isUnderUntrackedDir = function(arg,files)
    local matched_count=0
    local last_matched
    local upper_arg = string.upper(arg)
    local upper_arg_len = string.len(upper_arg)
    for i=1,#files do
        if string.upper(string.sub(files[i],1,upper_arg_len)) == upper_arg then
            matched_count = matched_count + 1
            last_matched = files[i]
        end
    end
    if matched_count == 1 and string.match(last_matched,"/$") then
        return true
    elseif matched_count < 1 then
        return true
    end
    return false
end

local addlist = function(args)
    local fd = io.popen("git status -s 2>nul","r")
    if not fd then
        return nil
    end
    local files = {}
    for line in fd:lines() do
        local arrowStart,arrowEnd = string.find(line," -> ",1,true)
        if arrowStart then
            files[#files+1] = unquote(string.sub(line,4,arrowStart-1))
            files[#files+1] = unquote(string.sub(line,arrowEnd+1))
        else
            files[#files+1] = unquote(string.sub(line,4))
        end
    end
    fd:close()
    if isUnderUntrackedDir(args[#args],files) then
        return nil
    end
    return files
end

local checkoutlist = function(args)
    local result = branchlist(args) or {}
    local fd = io.popen("git status -s 2>nul","r")
    if fd then
        for line in fd:lines() do
            if string.sub(line,1,2) == " M" then
                result[1+#result] = unquote(string.sub(line,4))
            end
        end
        fd:close()
    end
    return result
end


--setup current branch string
local currentbranch = function()
  return nyagos.eval('git rev-parse --abbrev-ref HEAD 2> nul')
end

-- subcommands
local gitsubcommands={}

-- see https://github.com/git/git/blob/master/command-list.txt
-- list-up `$git --list-cmds=list-history` for history group subcommands.
-- TODO: feature working in the future
-- 1. List up subcommands for each groups.
-- 2. Assign function for these subcommands, that function is completion function worked as supported multi groups type.

-- keyword (sort alphabetical, currently manual pickup)
gitsubcommands["bisect"]={"bad", "good", "help", "log", "new", "old", "replay", "reset", "run", "skip", "start", "terms", "view", "visualize"}
gitsubcommands["notes"]={"add", "append", "copy", "edit", "get-ref", "list", "merge", "merge", "merge", "prune", "remove", "show"}
gitsubcommands["reflog"]={"delete", "exists", "expire", "show"}
gitsubcommands["rerere"]={"clear", "diff", "forget", "gc", "remaining", "status"}
gitsubcommands["stash"]={"apply", "branch", "clear", "create", "drop", "list", "pop", "push", "show", "store"}
gitsubcommands["submodule"]={"absorbgitdirs", "add", "deinit", "foreach", "init", "set-branch", "set-url", "status", "summary", "sync", "update"}
gitsubcommands["svn"]={"blame", "branch", "clone", "commit-diff", "create-ignore", "dcommit", "fetch", "find-rev", "gc", "info", "init", "log", "mkdirs", "propget", "proplist", "propset", "rebase", "reset", "set-tree", "show-externals", "show-ignore", "tag"}
gitsubcommands["worktree"]={"add", "list", "lock", "move", "prune", "remove", "repair", "unlock"}

-- completion function apply
-- checkout
gitsubcommands["checkout"]=checkoutlist
-- branch select only
gitsubcommands["branch"]=branchlist
gitsubcommands["switch"]=branchlist
gitsubcommands["reset"]=branchlist
gitsubcommands["merge"]=branchlist
gitsubcommands["rebase"]=branchlist
gitsubcommands["revert"]=branchlist
-- current branch's commit
gitsubcommands["show"]=getcommits
-- add unstage file
gitsubcommands["add"]=addlist

local gitvar=share.git
gitvar.subcommand=gitsubcommands
gitvar.commit=getcommits
gitvar.branch=branchlist
gitvar.currentbranch=currentbranch
share.git=gitvar

if not share.maincmds then
    use "subcomplete.lua"
end

if share.maincmds and share.maincmds["git"] then
    -- git command complementation exists.
    nyagos.complete_for.git = function(args)
        while #args > 2 and args[2]:sub(1,1) == "-" do
            table.remove(args,2)
        end
        if #args == 2 then
            return share.maincmds.git
        end
        local subcmd = table.remove(args,2)
        while #args > 2 and args[2]:sub(1,1) == "-" do
            table.remove(args,2)
        end
        local t = share.git.subcommand[subcmd]
        if type(t) == "function" then
            return t(args)
        elseif type(t) == "table" and #args == 2 then
            return t
        end
    end
end

-- EOF
